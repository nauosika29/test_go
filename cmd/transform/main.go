// 使用 sqlx，而非 ORM 套件來操作資料庫。
// 優點：直接使用 SQL。
// 缺點：需要手動撰寫 SQL 語句，可能較不直觀。
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"testretail/pkg/db"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Product 結構體代表資料庫中的商品表。
// 使用 sqlx 的 struct tag 來指定結構體字段與資料庫表字段的對應關係。
// 注意：Go 語言需要明確指定資料型態，並且不會自動進行型態轉換。
// sql.NullString 用於處理可能為空的字段。(如果只給string, 遇到資料庫是null/nil的話會出錯)
type Product struct {
	ID          int            `db:"id"`
	ProductGUID string         `db:"product_guid"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"short_description"`
}

// CharacterProduct 結構體代表商品與角色之間的關聯表。
type CharacterProduct struct {
	ProductID     int    `db:"product_id"`
	CharacterID   int    `db:"character_id"`
	CharacterType string `db:"charter_type"`
}

// Character 結構體代表資料庫中的角色表。
type Character struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// BigQueryProduct 結構體代表轉換後的商品格式，用於 JSON 輸出。
type BigQueryProduct struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
}

// 如果專案只做到把資料庫的資料轉成JSON格式, 那這邊定義main是沒問題的
// 但如果要做到把資料轉成JSON格式並且上傳到Google BigQuery, 等於這個檔案只會是一個函數
// 所以我會改成一個開頭大寫的名稱, 例如Transform
// 上面的package main可能會修改成package transform
// 檔名也會改成transform.go
// 路徑會改為pkg/transform/transform.go
func main() {
	// 連接資料庫
	db := sqlx.NewDb(db.ConnectDB(), "postgres")
	defer db.Close()

	// 查詢並轉換商品
	products, err := QueryAndTransformProducts(db)
	if err != nil {
		log.Fatalf("Error querying and transforming products: %s", err)
	}

	// 將結果轉換為 JSON 格式並輸出
	jsonData, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		log.Fatalf("Error converting to JSON: %s", err)
	}
	fmt.Println(string(jsonData))
}

func QueryAndTransformProducts(db *sqlx.DB) ([]BigQueryProduct, error) {
	// 查詢前10個商品數據
	// 使用明確指定的欄位名稱而不是 SELECT *，以避免可能發生的結構體映射錯誤。
	// 如果使用 SELECT *，當資料庫表的結構發生變化（例如新增或刪除欄位）時，
	// 可能會導致我們定義的 Go 結構體與資料庫表的結構不匹配，從而引發錯誤。
	// 例如，如果資料庫表中有一個欄位 eslite_sn，但在 Go 結構體中沒有對應的字段，
	// 使用 SELECT * 會導致 "missing destination name eslite_sn in *main.Product" 的錯誤。
	// *main.Produc == 這邊定義的type Product struct
	rows, err := db.Queryx("SELECT id, product_guid, name, short_description FROM products LIMIT 10")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.StructScan(&p); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	// 轉換數據
	var result []BigQueryProduct
	for _, product := range products {
		// 查詢關聯的作者，條件為 character_type 包含 "作者" 且不包含 "作者(原文)"
		var characterProduct CharacterProduct
		err := db.Get(&characterProduct, "SELECT product_id, character_id, charter_type FROM character_products WHERE product_id = $1 AND charter_type LIKE '%作者%' AND charter_type NOT LIKE '%作者(原文)%' LIMIT 1", product.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				// 在使用 sqlx 函式庫進行資料庫查詢時，如果沒有找到符合條件的記錄，會返回 sql.ErrNoRows 錯誤。
				// 在這個例子中，這種情況發生在嘗試查詢商品相關聯的作者時。
				// 如果沒有這段程式碼，沒有作者的商品會被忽略，並且不會出現在最終的結果中。
				// 為了確保所有查詢到的商品都被包含在最終的結果中，我們選擇將作者欄位設置為空字串，並繼續處理其他商品。

				// Go 需要你清楚地去處理每一種可能的情況，包括那些可能並不是錯誤，但需要特殊處理的情況。
				result = append(result, BigQueryProduct{
					Type:        "PP",
					ID:          product.ProductGUID,
					Title:       product.Name,
					Description: product.Description.String,
					Author:      "",
				})
				continue
			}
			return nil, err
		}

		var character Character
		err = db.Get(&character, "SELECT id, name FROM characters WHERE id = $1", characterProduct.CharacterID)
		if err != nil {
			return nil, err
		}

		// 轉換並添加到結果切片
		bqProduct := BigQueryProduct{
			Type:        "PP",
			ID:          product.ProductGUID,
			Title:       product.Name,
			Description: product.Description.String,
			Author:      character.Name,
		}
		result = append(result, bqProduct)
	}

	return result, nil
}

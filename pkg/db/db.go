// package db 定義了操作資料庫的函數和類型。
// 這個包提供了連接和操作資料庫所需的工具。
// 他不能單獨運行，而是被其他包所引用。
// 所以package命名為db, 代表database相關的函數
package db

// 導入所需的包
import (
	"database/sql" // 用於操作 SQL 資料庫的包
	"fmt"          // 用於格式化輸出的包
	"log"          // 用於記錄錯誤信息的包
	"os"           // 用於操作環境變數的包

	// 導入 PostgreSQL 的驅動程序包。我們使用 _ 作為包的別名來表示我們只是想要初始化這個包，而不是直接使用其中的函數或變量。
	// 初始化這個包的目的是註冊 PostgreSQL 驅動到 database/sql 包，這樣我們就可以使用 database/sql 包來操作 PostgreSQL 資料庫。
	"github.com/joho/godotenv" // 用於從 .env 文件加載環境變數的包
	_ "github.com/lib/pq"
)

// 定義一個函數用於連接資料庫，返回一個資料庫連接對象
func ConnectDB() *sql.DB { // C大寫代表public, 小寫代表private
	// 加載 .env 文件中的環境變數
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file") // 如果加載失敗，記錄錯誤信息並終止程序
	}

	// 從環境變數中獲取資料庫連接信息
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// 格式化資料庫連接字符串
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)

	// 嘗試連接資料庫
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err) // 如果連接失敗，記錄錯誤信息並終止程序
	}

	// 嘗試 ping 資料庫以確保連接成功
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database: %s", err) // 如果 ping 失敗，記錄錯誤信息並終止程序
	}

	// 如果連接成功，返回資料庫連接對象
	// log.Println("Successfully connected!")
	return db
}

// 到上方為止為設計資料庫連接的函數
// 以下為測試是否連接成功的函數
// 解開log.Println("Successfully connected!")的註解
// 把package db改成package main(最上面那行)
// 解開以下main函數的註解(可能會有一些開發工具套件的錯誤提醒)
// 運行 go run pkg/db/db.go, 應該會看到錯誤訊息(如果沒設好環境變數的話)或是Successfully connected!

// func main() {
// 	db := ConnectDB()
// 	defer db.Close() // 保持defer db.Close()的原因是, 確保在函數結束時關閉資料庫連接
// }

// 相對的, 如果其他程序要使用這個函數, 只要import這個package, 就可以使用ConnectDB()(db.ConnectDB())函數

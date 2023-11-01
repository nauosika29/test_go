# testretail

## 目標
- 撈取指定商品資料，並轉換成json格式，上傳到Google BigQuery
- 2023/11/1 只做到轉換成json

## 安裝與配置

### 前提條件

裝 Go 語言環境。你可以從 [Go 官方網站](https://golang.org/dl/) 下載並安裝。

### 下載專案

```bash
git clone https://github.com/你的GitHub用戶名/testretail.git
cd testretail
```

### 安裝依賴

```bash
go mod tidy
```

### env 設定
- 請見自己的env
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=myuser
DB_PASSWORD=mypassword
DB_NAME=mydatabase
```

### 目前可運行指令
- 印結構體
```bash
go run cmd/dirprinter/main.go
```

- 商品資料轉換成json
```bash
go run cmd/jsonprinter/main.go
```
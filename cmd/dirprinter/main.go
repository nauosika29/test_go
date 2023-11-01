// 產生目錄樹狀結構
// 可以安裝 tree 就好(但我更害怕 brew update)
// 終端執行 $ go run cmd/dirprinter/main.go .
// 那個.代表目前的目錄, 也可以改成其他目錄
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GPT依照我要求產出的code
// 我還看不懂一堆參數與變數的用途, 就不多說明了
func printDir(out *os.File, prefix, dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(out, "failed to read dir %s: %s\n", dir, err)
		return
	}

	for i, file := range files {
		// 過濾一些不需要的檔案
		if strings.HasPrefix(file.Name(), ".") && file.Name() != ".gitignore" {
			continue
		}

		fmt.Fprintf(out, "%s├── %s\n", prefix, file.Name())
		if file.IsDir() {
			newPrefix := prefix + "│   "
			if i == len(files)-1 {
				newPrefix = prefix + "    "
			}
			printDir(out, newPrefix, filepath.Join(dir, file.Name()))
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run main.go <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	fmt.Printf("%s\n", dir)
	printDir(os.Stdout, "", dir)
}

// 執行後應該可以看到以下結果(如果沒有再去更動結構的話)
// .
// ├── .gitignore
// ├── cmd
// │   ├── dirprinter
// │   │   ├── main.go
// │   ├── transform
// │       ├── main.go
// ├── go.mod
// ├── go.sum
// ├── pkg
//     ├── 123.txt
//     ├── db
//         ├── db.go
// cmd在Go專案中(非強制 可用其他命名), 通常是用來放main.go的, 也就是可以執行的程式
// 一個大的專案可以執行的程式可能會很多, 可以在cmd底下建立不同資料夾, 並且放入main.go檔案
// 在cmd裡面的資料夾就可以當作是一個執行檔的名稱, 也就是一個package
// 很多資訊會說一個專案只能有一個main.go, 但其實是可以有很多的, 就是利用cmd資料夾來分開
// 很像是cmd這個大的main package, 裡面有很多小的main package, 結構管理起來會比較清楚
// main.go也不一定要叫main.go, 可以叫其他名稱, 但是這個檔案的package一定要是main

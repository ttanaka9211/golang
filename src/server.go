package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// 修正部分
func handle(w http.ResponseWriter, r *http.Request) {
	// DB接続部分。接続文字列は(MySQLユーザー名:パスワード@tcp(ホスト名:ポート)/DB名)
	db, err := sql.Open("mysql", "test:test@tcp(db:3306)/sample")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("select * from users where id < ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	var (
		id   int
		name string
	)
	rows, err := stmt.Query(2)

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}
		fmt.Println(id, name)
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		panic(err)
	}
	io.WriteString(w, "Hello, "+name)
}

func main() {
	// docker-composeで設定したポートを設定する
	portNumber := "1323"
	http.HandleFunc("/", handle)
	fmt.Println("Server listening on port ", portNumber)
	http.ListenAndServe(":"+portNumber, nil)
}

package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	db, err := sql.Open("mysql", "domino:yAUpZwWnjfrPBsWD@tcp(192.168.1.129:3306)/gold_member?charset=utf8&parseTime=True&loc=Local")
	checkError(err)
	for i := 1000000; i < 1001000; i++ {
		result, err := db.Query(fmt.Sprintf("insert into userinfo%d (mid,mnick,sex,mtime) values (?,?,?,?)", i%10), i, fmt.Sprintf("R%d", i), 0, time.Now().Unix())
		// checkError(err)
		fmt.Println(result, err)
		// for result.Next() {
		// 	var mid int
		// 	var exp int
		// 	var money int
		// 	var wtimes int
		// 	var ltimes int
		// 	var safebox int
		// 	result.Scan(&mid, &exp, &money, &wtimes, &ltimes, &safebox)
		// 	fmt.Println(mid, exp, money, wtimes, ltimes, safebox)
		// }
	}
	db.Close()
}

func robot_add(db *sql.DB) {
	for i := 1000000; i < 1001000; i++ {
		result, err := db.Query("insert into robot (mid) values (?)", i)
		// checkError(err)
		fmt.Println(result, err)
		// for result.Next() {
		// 	var mid int
		// 	var exp int
		// 	var money int
		// 	var wtimes int
		// 	var ltimes int
		// 	var safebox int
		// 	result.Scan(&mid, &exp, &money, &wtimes, &ltimes, &safebox)
		// 	fmt.Println(mid, exp, money, wtimes, ltimes, safebox)
		// }
	}
}

func gameinfo_add(db *sql.DB) {
	for i := 1000000; i < 1001000; i++ {
		result, err := db.Query(fmt.Sprintf("insert into gameinfo%d (mid) values (?)", i%10), i)
		// checkError(err)
		fmt.Println(result, err)
		// for result.Next() {
		// 	var mid int
		// 	var exp int
		// 	var money int
		// 	var wtimes int
		// 	var ltimes int
		// 	var safebox int
		// 	result.Scan(&mid, &exp, &money, &wtimes, &ltimes, &safebox)
		// 	fmt.Println(mid, exp, money, wtimes, ltimes, safebox)
		// }
	}
}

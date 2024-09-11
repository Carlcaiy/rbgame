package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println(sql.Drivers())
	test_mysql()
	test_sqlite()
}

func test_mysql() {
	dsn := "dev:^PwdOfDev2020$@tcp(192.168.6.17:3306)/tongits_activity?charset=utf8mb4"
	db, err := sql.Open("mysql", dsn)
	checkErr(err)
	rows, err := db.Query("SELECT task_id,COUNT(*) FROM `task_info_1` WHERE activity_id = 12 && `status`=2 && update_time > 0  GROUP BY task_id")
	checkErr(err)
	for rows.Next() {
		taskId := 0
		count := 0
		checkErr(rows.Scan(&taskId, &count))
		fmt.Println("taskId", taskId, "count", count)
	}
}

func test_sqlite() {
	db, err := sql.Open("sqlite3", "./foo.db")
	checkErr(err)
	logic(db)
}

func logic(db *sql.DB) {
	db.Exec(`CREATE TABLE IF NOT EXISTS userinfo (
		username INT
		department 
	)`)

	// 插入数据
	stmt, err := db.Prepare("INSERT INTO userinfo(username, department, created) values(?,?,?)")
	checkErr(err)

	res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)
	// 更新数据
	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)

	res, err = stmt.Exec("astaxieupdate", id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	// 查询数据
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)

	for rows.Next() {
		var uid int
		var username string
		var department string
		var created time.Time
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)
	}

	// 删除数据
	stmt, err = db.Prepare("delete from userinfo where uid=?")
	checkErr(err)

	res, err = stmt.Exec(id)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

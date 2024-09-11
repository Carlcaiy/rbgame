package test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"
	// _ "github.com/go-sql-driver/mysql"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func TestSQL(t *testing.T) {
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

func TestLocal(t *testing.T) {
	db, err := sql.Open("mysql", "domino:yAUpZwWnjfrPBsWD@tcp(192.168.1.129:3306)/test?charset=utf8&parseTime=True&loc=Local")
	checkErr(err)

	// 插入数据
	stmt, err := db.Prepare("INSERT userinfo SET username=?,department=?,created=?")
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
		var created string
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

func TestImg(t *testing.T) {
	db, err := sql.Open("mysql", "dev:^PwdOfDev2020$@tcp(192.168.6.17:3306)/tongits_admin?charset=utf8mb4&parseTime=True&loc=Local")
	checkError(err)
	result, err := db.Query("SELECT upload_path FROM `image_library` WHERE id >= 440")
	checkError(err)
	str := "zip myzip"
	for result.Next() {
		var upload_path string
		result.Scan(&upload_path)
		str += " ." + upload_path
	}
	db.Close()
	fmt.Println(str)
}

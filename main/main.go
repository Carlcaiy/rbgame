package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, _ := sql.Open("mysql", "dev:^PwdOfDev2020$@tcp(192.168.6.17:3306)/tongits_admin?charset=utf8mb4&parseTime=True&loc=Local")
	result, _ := db.Query("SELECT cdn_path FROM `image_library` WHERE id >= 440")
	verify := make(map[string][]string)
	for result.Next() {
		var upload_path string
		result.Scan(&upload_path)
		str := strings.Split(upload_path, "/")
		verify[str[2]] = append(verify[str[2]], upload_path)
	}
	for k, vs := range verify {
		if k != "img_library" && k != "sound_library" {
			str := "zip " + k
			for _, v := range vs {
				str += " ." + v + "/*"
			}
			fmt.Println(str)
		}
	}
	db.Close()
}

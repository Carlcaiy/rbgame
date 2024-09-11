package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {

	f, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		panic(err)
	}
	db = f

	// 1.创建路由
	// 默认使用了2个中间件Logger(), Recovery()
	r := gin.Default()

	r.LoadHTMLGlob("tem/*")
	// 路由组1 ，处理GET请求
	v1 := r.Group("/v1")
	// {} 是书写规范
	{
		v1.GET("/info", info)
		v1.GET("/activity", activity)
		v1.GET("/activity/cfg", activity_cfg)
	}
	v2 := r.Group("/v2")
	{
		v2.POST("/regist", regist)
		v2.POST("/login", login)
		v2.POST("/submit", submit)
	}
	r.Run("192.168.10.202:8000")
}

type ActivityReq struct {
	Uid        int32 `form:"uid" binding:"required,gt=0"`
	ActivityId int32 `form:"activity_id" binding:"required,gt=0"`
}

type ActivityRsp struct {
	Finish map[int32]int32
}

func info(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "http://www.baidu.com")
	// c.HTML(http.StatusOK, "index.html", gin.H{"title": "我是测试", "ce": "123456"})
}

func activity(c *gin.Context) {
	req := new(ActivityReq)
	if err := c.BindQuery(req); err != nil {
		c.String(403, err.Error())
		return
	}
	table := `CREATE TABLE IF NOT EXISTS activity_data (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		uid int NOT NULL DEFAULT '0', 
		activity_type smallint DEFAULT '0',
		activity_id int NOT NULL DEFAULT '0',
		data json NOT NULL,
		update_time int NOT NULL DEFAULT '0')`
	if _, err := db.Exec(table); err != nil {
		c.String(403, err.Error())
		return
	}

	data := ""
	if err := db.QueryRow("SELECT data FROM activity_data WHERE uid=? AND activity_id=?", req.Uid, req.ActivityId).Scan(&data); err != nil {
		if err != sql.ErrNoRows {
			c.String(403, err.Error())
			return
		} else {
			if result, err := db.Exec("INSERT into activity_data (uid,activity_type,activity_id,data,update_time) VALUES (?,?,?,?,?)", req.Uid, 1, req.ActivityId, `{}`, time.Now().Unix()); err != nil {
				c.String(403, err.Error())
				return
			} else {
				if lastInsertId, err := result.LastInsertId(); err != nil {
					c.String(403, err.Error())
					return
				} else {
					fmt.Println("lastInsertId", lastInsertId)
				}
				if rowNum, err := result.RowsAffected(); err != nil {
					c.String(403, err.Error())
					return
				} else {
					fmt.Println("rowNum", rowNum)
				}
			}
		}
	}
	c.JSON(200, data)
}

func activity_cfg(c *gin.Context) {

	c.JSON(200, "success")
}

func regist(c *gin.Context) {
	types := c.DefaultPostForm("type", "post")
	name := c.PostForm("name")
	password := c.PostForm("password")
	fmt.Println(name, password)
	if name == "" || password == "" {
		c.String(404, "name=%s password=%s types=%s", name, password, types)
		return
	}
	db.Prepare("insert into `user` name,password VALUES ?,?")
	c.String(200, "name=%s password=%s types=%s\n", name, password, types)
}

func login(c *gin.Context) {
	name := c.PostForm("name")
	c.String(200, fmt.Sprintf("hello %s\n", name))
}

func submit(c *gin.Context) {
	name := c.PostForm("name")
	c.String(200, fmt.Sprintf("hello %s\n", name))
}

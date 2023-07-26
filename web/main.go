package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	_ "rbgame/web/memory"
	"rbgame/web/session"
	"strconv"
	"strings"
	"time"
)

var globalSessions *session.Manager

// 然后在 init 函数中初始化
func init() {
	fmt.Println("main init")
	s, err := session.NewManager("memory", "gosessionid", 3600)
	if err != nil {
		panic(err)
	}
	globalSessions = s
	go globalSessions.GC()
}

func main() {
	http.HandleFunc("/", sayhelloName) // 设置访问的路由
	http.HandleFunc("/login", login)   // 设置访问的路由
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/count", count)
	err := http.ListenAndServe(":9090", nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	r.ParseForm()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, sess.Get("username"))
	} else {
		sess.Set("username", r.Form["username"])
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // 解析 url 传递的参数，对于 POST 则解析响应包的主体（request body）
	// 注意:如果没有调用 ParseForm 方法，下面无法获取表单的数据
	fmt.Println(r.Form) // 这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") // 这个写入到 w 的是输出到客户端的
}

// func login(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("method:", r.Method) // 获取请求的方法
// 	if r.Method == "GET" {
// 		crutime := time.Now().Unix()
// 		h := md5.New()
// 		io.WriteString(w, strconv.FormatInt(crutime, 10))
// 		token := fmt.Sprintf("%x", h.Sum(nil))

// 		t, _ := template.ParseFiles("login.gtpl")
// 		log.Println(t.Execute(w, token))
// 	} else {
// 		err := r.ParseForm() // 解析 url 传递的参数，对于 POST 则解析响应包的主体（request body）
// 		if err != nil {
// 			// handle error http.Error() for example
// 			log.Fatal("ParseForm: ", err)
// 		}
// 		// 请求的是登录数据，那么执行登录的逻辑判断
// 		token := r.Form.Get("token")
// 		fmt.Println("token:", token, template.HTMLEscapeString(r.Form.Get("token")))

// 		fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username"))) // 输出到服务器端
// 		fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
// 		template.HTMLEscape(w, []byte(r.Form.Get("username"))) // 输出到客户端
// 	}
// }

// 处理 /upload  逻辑
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) // 获取请求的方法
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666) // 此处假设当前目录下已存在test目录
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func count(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	token := sess.Get("token")
	fmt.Println("form.token:", r.Form.Get("token"), "mem.token", token)
	// if r.Form.Get("token") != token {
	// 	log.Fatal("token 不一致")
	// 	return
	// }
	h := md5.New()
	salt := "astaxie%^7&8888"
	io.WriteString(h, salt+time.Now().String())
	new_token := fmt.Sprintf("%x", h.Sum(nil))
	sess.Set("token", new_token)

	ct := sess.Get("countnum")
	if ct == nil {
		sess.Set("countnum", 1)
	} else {
		sess.Set("countnum", (ct.(int) + 1))
	}

	createtime := sess.Get("createtime")
	if createtime == nil {
		sess.Set("createtime", time.Now().Unix())
	} else if (createtime.(int64) + 60) < (time.Now().Unix()) {
		globalSessions.SessionDestroy(w, r)
		sess = globalSessions.SessionStart(w, r)
	}

	t, _ := template.ParseFiles("count.gtpl")
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, sess.Get("countnum"))
}

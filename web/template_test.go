package main

import (
	"fmt"
	"html/template"
	"os"
	"testing"
	"time"
)

type Friend struct {
	Fname string
}

type Person struct {
	UserName string
	Emails   []string
	Friends  []*Friend
}

func TestTemplate1(t *testing.T) {
	tl := template.New("some template")
	tl, _ = tl.Parse("hello {{.UserName}}!")
	person := Person{UserName: "Astaxie"}
	tl.Execute(os.Stdout, person)
}

func TestTemplate2(pt *testing.T) {
	f1 := Friend{Fname: "minux.ma"}
	f2 := Friend{Fname: "xushiwei"}
	t := template.New("filedname example")
	t, _ = t.Parse(`hello {{.UserName}}!
		{{range .Emails}}
			an email {{.}}
		{{end}}
		{{with .Friends}}
		{{range .}}
			my friend name is {{.Fname}}
		{{end}}
		{{end}}
	`)
	p := Person{UserName: "Astaxie", Emails: []string{"astaxie@beego.me", "astaxie@gmail.com"}, Friends: []*Friend{&f1, &f2}}
	t.Execute(os.Stdout, p)
}

func TestIf(t *testing.T) {
	tEmpty := template.New("template test")
	tEmpty = template.Must(tEmpty.Parse("空 pipeline if demo: {{if ``}} 不会输出. {{end}}\n"))
	tEmpty.Execute(os.Stdout, nil)

	tWithValue := template.New("template test")
	tWithValue = template.Must(tWithValue.Parse("不为空的 pipeline if demo: {{if `anything`}} 我有内容，我会输出。{{end}}\n"))
	tWithValue.Execute(os.Stdout, nil)

	tIfElse := template.New("template test")
	tIfElse = template.Must(tIfElse.Parse("if-else demo: {{if ``}} if部分 {{else}} else部分。{{end}}\n"))
	tIfElse.Execute(os.Stdout, nil)
}

func TestExecute(t *testing.T) {
	s1, _ := template.ParseFiles("header.tmpl", "content.tmpl", "footer.tmpl")
	s1.ExecuteTemplate(os.Stdout, "header", nil)
	fmt.Println()
	s1.ExecuteTemplate(os.Stdout, "content", nil)
	fmt.Println()
	s1.ExecuteTemplate(os.Stdout, "footer", nil)
	fmt.Println()
	s1.Execute(os.Stdout, nil)
}

func TestDir(t *testing.T) {
	os.Mkdir("astaxie", 0777)
	os.MkdirAll("astaxie/test1/test2", 0777)
	time.Sleep(time.Second * 10)
	err := os.Remove("astaxie")
	if err != nil {
		fmt.Println(err)
	}
	os.RemoveAll("astaxie")
}

func TestFile(t *testing.T) {
	userFile := "astaxie.txt"
	fout, err := os.Create(userFile)
	if err != nil {
		fmt.Println(userFile, err)
		return
	}
	defer fout.Close()
	for i := 0; i < 10; i++ {
		fout.WriteString("Just a test str!\r\n")
		time.Sleep(time.Second * 1)
		fout.Write([]byte("Just a test b!\r\n"))
		time.Sleep(time.Second * 1)
	}
}

func TestRead(t *testing.T) {
	userFile := "astaxie.txt"
	fl, err := os.Open(userFile)
	if err != nil {
		fmt.Println(userFile, err)
		return
	}
	defer fl.Close()
	buf := make([]byte, 1024)
	for {
		n, _ := fl.Read(buf)
		if n == 0 {
			break
		}
		os.Stdout.Write(buf[:n])
	}
}

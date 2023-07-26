package main

import (
	"fmt"
	"net/http"
	"plugin"
	"strconv"
)

func main() {
	http.HandleFunc("/plugin", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			recover()
		}()
		values := r.URL.Query()
		version := values.Get("version")
		if version == "" {
			w.Write([]byte("默认实现,更多精彩敬请期待"))
			return
		}

		version += ".so"
		fmt.Println(version)
		p, err := plugin.Open(version)
		if err != nil {
			w.Write([]byte("功能暂未上线，敬请期待"))
			return
		}

		fmt.Println("open so success")
		s, err := p.Lookup("MyPrin")
		if err != nil {
			w.Write([]byte("功能错误，敬请期待"))
		}
		fmt.Println("look up myprin success")

		f := s.(func(a, b int) int)
		v := f(4, 5)
		w.Write([]byte(strconv.Itoa(v)))
	})
	http.ListenAndServe(":8080", nil)
}

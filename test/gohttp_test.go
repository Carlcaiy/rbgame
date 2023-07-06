package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestHttpServer(t *testing.T) {
	http.ListenAndServe(":8888", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		fmt.Println("get request")
		select {
		case <-ctx.Done():
			fmt.Println("request canceleld")
		case <-time.After(time.Second):
			w.Write([]byte("process finished"))
		}
	}))
}

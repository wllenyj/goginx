package http

import (
	"testing"
	"time"
	"strings"
	"bytes"
	"net/http"
	"fmt"
	"io/ioutil"
)

func httpServer() {
	http.HandleFunc("/timeout", func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(2*time.Second)
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "pong")
	})
	http.HandleFunc("/params", func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("params %s\n", req.URL.Query())
		fmt.Fprintf(w, "%s", req.URL.RawQuery)
	})
	http.HandleFunc("/headers", func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("User-Agent: %s\n", req.Header.Get("User-Agent"))
		fmt.Fprintf(w, "%s", req.Header.Get("User-Agent"))
	})
	http.HandleFunc("/post", func(w http.ResponseWriter, req *http.Request) {
		body, _ := ioutil.ReadAll(req.Body)
		fmt.Printf("post: %s\n", body)
		
		fmt.Fprintf(w, "%s %s", req.Method, body)
	})

	err := http.ListenAndServe(":35241", nil)

	if err != nil {
		fmt.Println("ListenAndServe error: ", err.Error())
	}
}

func init() {
	go httpServer()
}

func TestHttpGet(t *testing.T) {
	body, err := Get("http://www.baidu.com")
	if err != nil {
		t.Errorf("http failed. %s\n", err)
	}
	if len(body) == 0 {
		t.Errorf("http failed. body empty")
	}
}

func TestHttpTimeout(t *testing.T) {
	body, err := Request("http://localhost:35241/timeout", nil, nil, nil, 1*time.Second)
	if !strings.Contains(err.Error(), "Timeout") {
		t.Errorf("http failed. %s\n", err)
	}
	if len(body) != 0 {
		t.Errorf("http failed. body empty")
	}
}

func TestHttpRet(t *testing.T) {
	body, err := Request("http://localhost:35241/ping", nil, nil, nil, 1*time.Second)
	if err != nil {
		t.Errorf("http failed. %s\n", err)
	}
	if bytes.Equal(body, []byte("pong")) != true {
		t.Errorf("http body failed. %s", body)
	}
}

func TestHttpParams(t *testing.T) {
	params := map[string]string{
		"user":"aaa",
		"pwd":"bbb",
	}
	body, err := Request("http://localhost:35241/params", params, nil, nil, 1*time.Second)
	if err != nil {
		t.Errorf("http failed. %s\n", err)
	}
	if bytes.Equal(body, []byte("pwd=bbb&user=aaa")) != true {
		t.Errorf("http body failed. %s", body)
	}
}

func TestHttpHeaders(t *testing.T) {
	headers := map[string]string{
		"User-Agent":"matrix",
	}
	body, err := Request("http://localhost:35241/headers", nil, headers, nil, 1*time.Second)
	if err != nil {
		t.Errorf("http failed. %s\n", err)
	}
	if bytes.Equal(body, []byte("matrix")) != true {
		t.Errorf("http body failed. %s", body)
	}
}

func TestHttpPost(t *testing.T) {
	body, err := Post("http://localhost:35241/post", []byte("body"))
	if err != nil {
		t.Errorf("http failed. %s\n", err)
	}
	if bytes.Equal(body, []byte("POST body")) != true {
		t.Errorf("http body failed. %s", body)
	}
}

func TestHttpMatrix(t *testing.T) {
	body, err := Post("http://11.160.112.123/miner/collection_notify", []byte(`
	{"pid":"4000000000000000000000002115011400581710"}`))
	if err != nil {
		t.Errorf("http failed. %s\n", err)
	}
	if len(body) == 0 {
		t.Errorf("http failed. body empty")
	}
	fmt.Printf("body:%s\n", body)
}

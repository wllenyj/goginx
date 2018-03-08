package main

import (
	"flag"
	"fmt"
	"github.com/wllenyj/goginx"
	"log"
	"net/http"
	"os"
)

var port = flag.String("p", "8883", "port")

func handler(w http.ResponseWriter, r *http.Request) {
	ret := fmt.Sprintf("pid:%d\n", os.Getpid())
	w.Write([]byte(ret))
}

func main() {
	flag.Parse()
	program, err := goginx.Daemon()
	if err != nil {
		log.Printf("Daemon err. %s", err)
		return
	}
	http.HandleFunc("/hello", handler)

	if err = program.ListenAndServe(":"+*port, nil); err != nil {
		log.Printf("listen err. %s", err)
		return
	}
	program.Run()
}

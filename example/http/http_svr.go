package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"goginx"
	"net/http"
	"log"
	"os"
	"flag"
	"time"
)

var config = flag.String("c", "", "")
var port = flag.String("p", "23523", "")

func handler(w http.ResponseWriter, r *http.Request) {
	ret := fmt.Sprintf("pid:%d\n", os.Getpid())
	time.Sleep(7*time.Second)
	w.Write([]byte(ret))
	log.Println(ret)
}
func handler1(w http.ResponseWriter, r *http.Request) {
	ret := fmt.Sprintf("pid:%d\n", os.Getpid())
	w.Write([]byte(ret))
	log.Println(ret)
}

func main() {
	//fmt.Printf("%d start\n", os.Getpid())
	program, err := goginx.Daemon()
	if err != nil {
		log.Printf("Daemon err. %s", err)
		return
	}

	fmt.Printf("port:%s\n", *port)

	fmt.Printf("%d daemon start\n", os.Getpid())

	mux1 := mux.NewRouter()
	mux1.HandleFunc("/hello", handler).
		Methods("GET")
	mux1.HandleFunc("/hello1", handler1).
		Methods("GET")

	if err = program.ListenAndServer(":" + *port, mux1); err != nil {
		log.Printf("listen err. %s", err)
		return
	}

	program.Run()
}

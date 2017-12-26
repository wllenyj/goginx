package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/wllenyj/goginx"
	"time"
	"log"
	"os"
	"flag"
	"./link"
	"./link/codec"
)

var config = flag.String("c", "", "")

type AddReq struct {
	A, B int
}

type AddRsp struct {
	C int
}

func ListenAndServe(prog *goginx.Program) error {
	json := codec.Json()
	json.Register(AddReq{})
	json.Register(AddRsp{})

	listener, err := prog.ListenTCP("0.0.0.0:41231")
	if err != nil {
		fmt.Printf("listen err. %s", err)
		return err
	}
	server := link.NewServer(listener, json, 0, link.HandlerFunc(serverSessionLoop))
	prog.AddService(server)

	return nil
}

func serverSessionLoop(session *link.Session) {
	for {
		req, err := session.Receive()
		if err != nil {
			log.Printf("%d %d receive err.%s", os.Getpid(), session.ID(), err)
			break
		}

		err = session.Send(&AddRsp{
			req.(*AddReq).A + req.(*AddReq).B,
		})

		//log.Printf("%d %d server recv %d", os.Getpid(), session.ID(), req.(*AddReq).A)
		if err != nil {
			log.Printf("%d %d send err.%s", os.Getpid(), session.ID(), err)
			break
		}
	}
	session.Close()
}

func handler(w http.ResponseWriter, r *http.Request) {
	ret := fmt.Sprintf("pid:%d\n", os.Getpid())
	log.Printf("%d http request ", os.Getpid())
	time.Sleep(2 * time.Second)
	w.Write([]byte(ret))
}

func main() {
	fmt.Printf("%d start\n", os.Getpid())
	program, err := goginx.Daemon()
	if err != nil {
		log.Printf("Daemon err. %s", err)
		return
	}

	fmt.Printf("%d daemon start\n", os.Getpid())

	ListenAndServe(program)	

	mux1 := mux.NewRouter()
	mux1.HandleFunc("/hello", handler).Methods("GET")

	if err = program.ListenAndServe(":41232", mux1); err != nil {
		log.Printf("listen err. %s", err)
		return
	}

	program.Run()

	fmt.Println("vim-go")
}

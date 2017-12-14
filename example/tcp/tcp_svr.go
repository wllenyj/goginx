package main

import (
	"fmt"
	"goginx"
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

	//addr := server.Listener().Addr().String()
	//

	//client, err := link.Dial("tcp", addr, json, 0)
	//
	//checkErr(err)
	//go clientSessionLoop(client)

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

		log.Printf("%d %d server recv %d", os.Getpid(), session.ID(), req.(*AddReq).A)
		if err != nil {
			log.Printf("%d %d send err.%s", os.Getpid(), session.ID(), err)
			break
		}
	}
	session.Close()
}

func clientSessionLoop(session *link.Session) {
	for i := 0; i < 600; i++ {
		err := session.Send(&AddReq{
			i, i,
		})
		checkErr(err)
		log.Printf("Send: %d + %d", i, i)

		rsp, err := session.Receive()
		checkErr(err)
		log.Printf("Receive: %d", rsp.(*AddRsp).C)
		time.Sleep(time.Second)	
	}
}


func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
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

	program.Run()

	fmt.Println("vim-go")
}

package main

import (
	//"fmt"
	"./link"
	"./link/codec"
	"log"
	"os"
	"time"
)

type AddReq struct {
	A, B int
}

type AddRsp struct {
	C int
}

func clientSessionLoop(num int, session *link.Session) {
	for i := 0; i < 30; i++ {
		err := session.Send(&AddReq{
			i, i,
		})
		if err != nil {
			log.Printf("send err.%s", err)
			break
		}
		log.Printf("%d:%d Send: %d + %d", num, os.Getpid(), i, i)

		rsp, err := session.Receive()
		if err != nil {
			log.Printf("receive err.%s", err)
			break
		}

		log.Printf("%d:%d Receive: %d", num, os.Getpid(), rsp.(*AddRsp).C)
		time.Sleep(time.Second)
	}
	session.Close()
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func client(i int) {
	json := codec.Json()
	json.Register(AddReq{})
	json.Register(AddRsp{})

	client, err := link.Dial("tcp", ":41231", json, 0)

	checkErr(err)
	clientSessionLoop(i, client)
}

func main() {
	for i := 0; i < 50; i++ {
		go client(i)
		time.Sleep(time.Second)
	}
	client(101)
}

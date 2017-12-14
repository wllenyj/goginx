package main

import (
	"github.com/wllenyj/goginx"
	"github.com/gorilla/mux"
	"fmt"
	"net/http"
	"os"
	"log"
)

//type myhandler struct {
//}
//
//func (h *myhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprintf(w, "Hi, This is an example of http service in golang!\n")
//}
func handler(w http.ResponseWriter, r *http.Request) {
	ret := fmt.Sprintf("pid:%d\n", os.Getpid())
	w.Write([]byte(ret))
}

func main() {
	fmt.Printf("%d start\n", os.Getpid())
	program, err := goginx.Daemon()
	if err != nil {
		log.Printf("Daemon err. %s", err)
		return
	}

	//pool := x509.NewCertPool()
	//caCertPath := "ca.crt"

	//caCrt, err := ioutil.ReadFile(caCertPath)
	//if err != nil {
	//	fmt.Println("ReadFile err:", err)
	//	return
	//}
	//pool.AppendCertsFromPEM(caCrt)

	//s := &http.Server{
	//	Addr:    ":8043",
	//	Handler: &myhandler{},
	//	TLSConfig: &tls.Config{
	//		ClientCAs:  pool,
	//		ClientAuth: tls.RequireAndVerifyClientCert,
	//	},
	//}
	mux1 := mux.NewRouter()
	mux1.HandleFunc("/hello", handler).
		Methods("GET")

	err = program.ListenAndServeTLS(":8043", "server.crt", "server.key", mux1)
	if err != nil {
		fmt.Println("ListenAndServeTLS err:", err)
	}

	program.Run()
}

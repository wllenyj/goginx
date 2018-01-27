# goginx
A library of Go application, it can make the Go application run like Nginx. Implements the daemon, graceful restart and graceful shutdown.

## Usage:

http.go
``` go
package main

import (
    "fmt"
    "github.com/wllenyj/goginx"
    "log"
    "net/http"
    "os"
)

func handler(w http.ResponseWriter, r *http.Request) {
    ret := fmt.Sprintf("pid:%d\n", os.Getpid())
    w.Write([]byte(ret))
}

func main() {
    program, err := goginx.Daemon()
    if err != nil {
        log.Printf("Daemon err. %s", err)
        return
    }   
    http.HandleFunc("/hello", handler)

    if err = program.ListenAndServe(":8882", nil); err != nil {
        log.Printf("listen err. %s", err)
        return
    }   
    program.Run()
}
```
``` shell
$> go get github.com/wllenyj/goginx
$> go build http.go
$> ./http -p 8883
$> curl "localhost:8883/hello"
pid:5442
$> ./http -s restart
$> curl "localhost:8883/hello"
pid:6312
$> ./http -s stop
$> curl "localhost:8883/hello"
curl: (7) couldn't connect to host
```

## References
* [facebook grace](https://github.com/facebookgo/grace)
* [go-daemon](https://github.com/sevlyar/go-daemon)
* [endless](https://github.com/fvbock/endless)
* [overseer](https://github.com/jpillora/overseer)


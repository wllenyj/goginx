package goginx

import (
	"net"
	"os"
	"sync"
	//"syscall"
	"time"
	"log"
)

func newGoginxListener(l net.Listener, wg *sync.WaitGroup) *goginxListener {
	return &goginxListener{
		Listener: l,
		wg:       wg,
	}
}

type goginxListener struct {
	net.Listener
	stopped bool
	wg      *sync.WaitGroup
}

func (this *goginxListener) Accept() (net.Conn, error) {
	conn, err := this.Listener.(*net.TCPListener).AcceptTCP()
	if err != nil {
		return nil, err
	}
	conn.SetKeepAlive(true)                  // see http.tcpKeepAliveListener
	conn.SetKeepAlivePeriod(3 * time.Minute) // see http.tcpKeepAliveListener

	nxconn := goginxConn{
		Conn: conn,
		wg:   this.wg,
	}
	this.wg.Add(1)
	return nxconn, nil
}

func (this *goginxListener) Close() error {
	//if this.stopped {
	//	return syscall.EINVAL
	//}
	//this.stopped = true
	log.Println("listener close")
	return this.Listener.Close()
}

func (this *goginxListener) File() *os.File {
	// returns a dup(2) - FD_CLOEXEC flag *not* set
	//tl := this.Listener.(*net.TCPListener)
	//fl, _ := tl.File()
	tl, _ := this.Listener.(filer).File()
	return tl
}

type filer interface {
	File() (*os.File, error)
}

type goginxTlsListenner struct {
	goginxListener
}

type goginxConn struct {
	net.Conn
	wg *sync.WaitGroup
}

func (nxc goginxConn) Close() error {
	err := nxc.Conn.Close()
	if err == nil {
		nxc.wg.Done()
	}
	return err
}

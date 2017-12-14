package goginx

import (
	"crypto/tls"
	"goginx/daemon"
	"goginx/gracenet"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Service interface {
	//Server(net.Listener) error
	Serve() error

	// graceful quit
	Shutdown()
}

type goginxServer struct {
	svr      *Service
	listener net.Listener
}

type Program struct {
	daemon *daemon.Daemon
	net    *gracenet.Net

	sds     []Service
	wg      sync.WaitGroup
	sigChan chan os.Signal
}

func NewProgram() *Program {
	return &Program{
		net:     &gracenet.Net{},
		sds:     make([]Service, 0),
		sigChan: make(chan os.Signal),
	}
}

func Daemon(options ...daemon.Option) (*Program, error) {
	opts, err := daemon.DefaultOption(options...)
	if err != nil {
		return nil, err
	}
	cntxt := &daemon.Daemon{
		Opts: opts,
	}
	daemon.ConfigSignal("stop", syscall.SIGTERM)
	daemon.ConfigSignal("restart", syscall.SIGHUP)
	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			log.Fatalln("Unable send signal to the daemon:", err)
		}
		daemon.SendCommands(d)
		os.Exit(0)
		return nil, nil
	}
	d, err := cntxt.Reborn()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if d != nil {
		os.Exit(0)
		return nil, nil
	}
	program := &Program{
		daemon: cntxt,
		net: &gracenet.Net{
			ProcAttr: cntxt,
		},
		sds:     make([]Service, 0),
		sigChan: make(chan os.Signal, 1),
	}
	program.net.AddObserver(cntxt)
	return program, nil
}

func (this *Program) AddService(svr Service) {
	this.sds = append(this.sds, svr)
}

func (prog *Program) ListenTCP(addr string) (net.Listener, error) {
	return prog.net.Listen("tcp", addr)
}

func (prog *Program) ListenAndServer(addr string, handler http.Handler) error {
	srv := NewHttpServer(addr, handler)

	l, err := prog.ListenTCP(addr)
	if err != nil {
		return err
	}
	log.Printf("net listen ok")
	srv.goginxListener = newGoginxListener(l, &srv.wg)

	prog.AddService(srv)
	return nil
}

func (prog *Program) ListenAndServeTLS(addr string, certFile string, keyFile string, handler http.Handler) error {
	srv := NewHttpServer(addr, handler)

	config := &tls.Config{}
	if srv.TLSConfig != nil {
		*config = *srv.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	config.Certificates = make([]tls.Certificate, 1)
	var err error
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	l, err := prog.ListenTCP(addr)
	if err != nil {
		return err
	}
	//srv.tlsInnerListener = newEndlessListener(l, srv)
	inner := newGoginxListener(l, &srv.wg)

	srv.goginxListener = tls.NewListener(inner, config)

	prog.AddService(srv)
	return nil
}

func (this *Program) Run() error {
	if this.net.Inherited() {
		log.Println(syscall.Getpid(), " kill parent ", syscall.Getppid())
		if err := syscall.Kill(syscall.Getppid(), syscall.SIGTERM); err != nil {
			log.Printf("failed to close parent: %s", err)
		}
	}
	go this.handleSignals()

	for _, s := range this.sds {
		svr := s
		go func() {
			svr.Serve()
			this.wg.Done()
		}()
		this.wg.Add(1)
	}

	this.wg.Wait()
	if this.daemon != nil {
		this.daemon.Release()
	}
	return nil
}

func (this *Program) handleSignals() {
	var sig os.Signal

	//hookableSignals := []os.Signal{
	//	syscall.SIGHUP,
	//	syscall.SIGUSR1,
	//	syscall.SIGUSR2,
	//	syscall.SIGINT,
	//	syscall.SIGTERM,
	//	syscall.SIGTSTP,
	//}

	signal.Notify(
		this.sigChan,
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGTSTP,
	)

	pid := syscall.Getpid()
	for {
		sig = <-this.sigChan
		//srv.signalHooks(PRE_SIGNAL, sig)
		switch sig {
		case syscall.SIGHUP:
			log.Println(pid, "Received SIGHUP. forking.")
			//for _, srv := range this.sds {
			//}
			if _, err := this.net.StartProcess(); err != nil {
				log.Println("StartProcess err.%s", err)
			}
		case syscall.SIGUSR1:
			log.Println(pid, "Received SIGUSR1.")
		case syscall.SIGUSR2:
			log.Println(pid, "Received SIGUSR2.")
			//srv.hammerTime(0 * time.Second)
		case syscall.SIGINT, syscall.SIGTERM:
			signal.Stop(this.sigChan)	
			log.Printf("%d Received %v.", pid, sig)
			for _, srv := range this.sds {
				srv.Shutdown()
			}
			return
		case syscall.SIGTSTP:
			log.Println(pid, "Received SIGTSTP.")
		default:
			log.Printf("Received %v: nothing i care about...\n", sig)
		}
		//srv.signalHooks(POST_SIGNAL, sig)
	}
}

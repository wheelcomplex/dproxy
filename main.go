package main

import (
	"flag"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
	proxy "github.com/wheelcomplex/fasthttp-reverse-proxy"
)

var (
	addr        = flag.String("l", ":5003", "http server address")
	backendAddr = flag.String("b", "127.0.0.1:5003", "backend http server address")
	showHelp    = flag.Bool("h", false, "show help")
	proxyServer *proxy.ReverseProxy
)

// ProxyHandler ... fasthttp.RequestHandler func
func ProxyHandler(ctx *fasthttp.RequestCtx) {
	// all proxy to localhost
	proxyServer.ServeHTTP(ctx)
}

func main() {
	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	listenHost, listenPort, err := parseAddr(*addr)
	if err != nil {
		log.Printf("Invalid listen address: %s\n", *addr)
		log.Fatal(err)
	}
	backendHost, backendPort, err := parseAddr(*backendAddr)

	if err != nil {
		log.Printf("Invalid backend address: %s\n", *backendAddr)
		log.Fatal(err)
	}

	var be = net.JoinHostPort(backendHost, backendPort)
	var laddr = net.JoinHostPort(listenHost, listenPort)

	proxyServer = proxy.NewReverseProxy(be)

	log.Printf("dproxy listening on %s, proxy to %s ...\n", laddr, be)

	if err := fasthttp.ListenAndServe(laddr, ProxyHandler); err != nil {
		log.Fatal(err)
	}
}

func parseAddr(s string) (host, port string, err error) {
	s = strings.Trim(s, ":")
	_, err = strconv.ParseUint(s, 10, 16)
	if err == nil {
		// this is a port only
		host = ""
		port = s
		return
	}
	// The port starts after the last colon.
	i := strings.LastIndex(s, ":")
	if i < 0 {
		s = s + ":80"
	}
	host, port, err = net.SplitHostPort(s)

	return
}

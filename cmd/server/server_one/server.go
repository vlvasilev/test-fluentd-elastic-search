package main

import (
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/I330716/test-fluentd-elastic-search/pkg/server/server_one"
)

func main() {
	var myIP = flag.String("address", "0.0.0.0", "ip on wich the server will listen")
	var myPort = flag.String("port", "8000", "port on wich the server will listen")
	var kubeconfig = flag.String("kubeconfig", "./conf/kubeconfig.yaml", "The config file for the kubernetes kluster")
	flag.Parse()

	serverOneHandler := server_one.ServerOneHandler{}
	serverOneHandler.Init(*myIP, *myPort, *kubeconfig)

	server := http.Server{
		Addr:    serverOneHandler.GetListeningEndPoint(),
		Handler: &serverOneHandler,
	}
	log.Println("Server Start To Listen On: " + serverOneHandler.GetListeningEndPoint())
	err := server.ListenAndServe()
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("Server Exit Normaly")
	}

}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}

	io.WriteString(w, "My server: "+r.URL.String())
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

var mux map[string]func(http.ResponseWriter, *http.Request)

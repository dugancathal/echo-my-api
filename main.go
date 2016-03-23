package main

import (
	"fmt"
	"net/http"
	"log"
	"github.com/dugancathal/stuffs/controllers"
	"flag"
	"github.com/dugancathal/stuffs/defaultsparser"
)

var port *int64
var defaultsPath *string

func init() {
	port = flag.Int64("port", 9292, "The port to run on")
	defaultsPath = flag.String("defaults-path", "", "A directory of config files to load in")
	flag.Parse()
}

func main() {
	defaultParser := defaultsparser.NewDirectoryDefaultParser(*defaultsPath)
	echoController := controllers.NewEchoControllerFromConfig(defaultParser.RouteMapping())
	http.HandleFunc("/__config", echoController.GetConfiguration)
	http.HandleFunc("/set-response-for/", echoController.HandleSetReq)
	http.HandleFunc("/get-requests-to/", echoController.HandleGetReq)
	http.HandleFunc("/", echoController.HandleReq)

	fmt.Printf("Running on port %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

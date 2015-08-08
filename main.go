package main // import "github.com/Luzifer/locationmaps"

//go:generate go-bindata assets

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Luzifer/rconfig"
	"github.com/gorilla/mux"
	"github.com/lancecarlson/couchgo"
)

var (
	config = struct {
		Listen            string `flag:"listen" default:":3000" description:"Address and port to listen on"`
		CouchDBConnection string `flag:"couchdb" description:"Connection string for CouchDB"`
		Version           bool   `flag:"version" default:"false" description:"Print version and quit"`
	}{}

	couchConn *couch.Client
	version   = "dev"
)

func init() {
	if err := rconfig.Parse(&config); err != nil {
		fmt.Printf("Unable to load config: %s\n", err)
		os.Exit(1)
	}

	if config.Version {
		fmt.Printf("locationmaps %s\n", version)
		os.Exit(0)
	}

	c, err := couch.NewClientURL(config.CouchDBConnection)
	if err != nil {
		fmt.Printf("Unable to connect to CouchDB: %s\n", err)
		os.Exit(1)
	}
	couchConn = c
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/{user:.+}.html", handleMap).
		Methods("GET")
	r.HandleFunc("/{user:.+}.json", handleGetLatest).
		Methods("GET")

	r.HandleFunc("/simple.add", handleSimpleAdd).
		Methods("POST")

	http.ListenAndServe(config.Listen, r)
}

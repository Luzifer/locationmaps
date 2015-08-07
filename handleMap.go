package main

import (
	"fmt"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
)

func handleMap(res http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tplString, err := Asset("assets/map.html")
	if err != nil {
		fmt.Printf("Unable to load asset map.html: %s\n", err)
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tpl, err := pongo2.FromString(string(tplString))
	if err != nil {
		fmt.Printf("Unable to parse map.html: %s\n", err)
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data, err := tpl.Execute(pongo2.Context{
		"user": vars["user"],
	})
	if err != nil {
		fmt.Printf("Unable to execute map.html: %s\n", err)
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res.Write([]byte(data))
}

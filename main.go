package main

import (
	"fmt"
	"github.com/MegaMindInKZ/task-techno.git/cache"
	"github.com/MegaMindInKZ/task-techno.git/config"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

var Router *mux.Router
var Server *http.Server

func main() {
	defer config.End()
	fmt.Println(Server.ListenAndServe())
}

func init() {
	Router = mux.NewRouter().StrictSlash(true)
	routes()
	cache.LocalCache = cache.NewCache(1000)
	cache.SetUp()

	Server = &http.Server{
		Handler:      Router,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func routes() {
	Router.HandleFunc("/admin/redirects", GetRedirectsAdminHandler).Methods("GET")
	Router.HandleFunc("/admin/redirects/{id}", GetRedirectAdminHandler).Methods("GET")
	Router.HandleFunc("/admin/redirects/{id}", PatchRedirectAdminHandler).Methods("PATCH")
	Router.HandleFunc("/admin/redirects", PostRedirectAdminHandler).Methods("POST")
	Router.HandleFunc("/admin/redirects/{id}", DeleteRedirectAdminHandler).Methods("DELETE")
	Router.HandleFunc("/redirects", GetRedirectHandler).Methods("GET")

}

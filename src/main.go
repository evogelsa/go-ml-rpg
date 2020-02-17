package main

import (
	"fmt"
	"log"
	"net/http"

	"github.iu.edu/evogelsa/go-ml-rpg/ui"

	"github.com/gorilla/mux"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<h1>Hello, World!</h1>`)
}

func newRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", home)
	r.HandleFunc("/uitest", ui.Load)

	return r
}

func main() {
	r := newRouter()

	log.Fatal(http.ListenAndServe(":8080", r))
}

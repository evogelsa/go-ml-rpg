package ui

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func Load(w http.ResponseWriter, r *http.Request) {
	f, err := ioutil.ReadFile("ui/assets/gfx.txt")
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, string(f))
}

package main

import (
	"math/rand"
	"time"

	"github.iu.edu/evogelsa/go-ml-rpg/web"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	web.Server()
}

package main

import (
	"net/http"

	"github.com/ironpark/potion/core"
)

func main() {
	p := core.New()
	p.Get("/", Home)
	p.Get("/pub/*", http.FileServer(http.Dir("public/")))
	http.ListenAndServe(":3007", p.Serve())
}

func Home(c *core.Potion) {
	c.HTML(200, "front.tmpl", nil)
}

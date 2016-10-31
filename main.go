package main

import (
	"net/http"

	"github.com/ironpark/potion/core"
)

func main() {
	p := core.New()
	p.Get("/", Home)
	p.Pubilc("public/")
	http.ListenAndServe(":3007", p.Serve())
}

func Home(c *core.Context) {
	c.HTML(200, "front.tmpl", nil)
}

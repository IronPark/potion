package main

import (
	"fmt"
	"net/http"

	"github.com/ironpark/potion/core"
)

func main() {
	fmt.Print("test")
	p := core.New()
	p.Get("/", Home)

	users := p.Group("/users")
	users.Get("", Users)
	user := users.Group("/:id")
	user.Get("", User)
	user.Get("/profile", UserProfile)

	http.ListenAndServe(":3007", p.Serve())
}

// Home ...
func Home(c *core.Potion) {
	var username string
	c.AppContext.Log.Println("Found User")
	c.Response().Write([]byte("Welcome Home " + username))
}

// Users ...
func Users(c *core.Potion) {
	c.AppContext.Log.Println("In Users Function")
	c.Response().Write([]byte("Users"))
}

// User ...
func User(c *core.Potion) {
	id := c.Param("id")
	var username string
	// username = c.AppContext.DB.find(user by id.....)
	c.AppContext.Log.Println("Found User")
	c.Response().Write([]byte("Welcome " + username + " with id " + id))
}

// UserProfile ...
func UserProfile(c *core.Potion) {
	id := c.Param("id")
	var profile string
	c.AppContext.Log.Println("Found User Profile")
	c.Response().Write([]byte("Here's your profile " + profile + " user " + id))
}

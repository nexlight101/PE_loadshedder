package main

import (
	"log"
	"net/http"

	"github.com/nexlight101/PE_loadshedder/conf"
	"github.com/nexlight101/PE_loadshedder/modules"

	"github.com/julienschmidt/httprouter"
)

func main() {
	// Get a new router.
	r := httprouter.New()

	// Get a template controller value.
	c := modules.NewController(conf.TPL)

	// Handle dictionary routes
	r.GET("/", c.Main)              // Landing page
	r.GET("/index", c.Index)        // GET ROUTE displays the date, stage & suburb form
	r.GET("/schedule", c.Schedule)  // Find the schedule
	r.POST("/findArea", c.FindArea) // POST ROUTE reads the date, stage & suburb fiels
	r.GET("/about", c.About)        // GET ROUTE displays the about information

	// Handle icon
	http.Handle("/favicon.ico", http.NotFoundHandler())
	// Serve CSS
	// r.ServeFiles("/public/stylesheets/*filepath", http.Dir("public/stylesheets"))

	// if not found look for a static file
	static := httprouter.New()
	static.ServeFiles("/public/stylesheets/*filepath", http.Dir("public/stylesheets"))
	// r.NotFound = http.FileServer(http.Dir("public/stylesheets"))
	r.NotFound = static
	// Server
	log.Fatal(http.ListenAndServe(":80", r))

}

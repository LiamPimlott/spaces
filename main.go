package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("hello world")
	r := mux.NewRouter()

	// root route will serve the built react app.
	r.Handle("/", http.FileServer(http.Dir("./frontend/build")))
	// serve static assest like images, css from the /static/{file} route
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend/build/static"))))

	http.ListenAndServe(":8080", r)
}

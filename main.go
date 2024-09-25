package main

import (
	"fmt"
	"net/http"
	"portfolio/functions"
)

// port of the server
const port = ":8080"

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { functions.Index(w, r) })

	//load the CSS, the JS and the images
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))
	http.Handle("/script/", http.StripPrefix("/script/", http.FileServer(http.Dir("./script"))))

	//start the local host
	fmt.Println("\n(http://localhost:8080/) - Server started on port", port)
	http.ListenAndServe(port, nil)
}

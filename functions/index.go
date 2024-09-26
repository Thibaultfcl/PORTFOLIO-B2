package functions

import "net/http"

func Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./tmpl/index.html")
}

func Contact(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./tmpl/contact.html")
}
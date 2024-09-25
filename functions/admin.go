package functions

import (
	"database/sql"
	"net/http"
)

func Admin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	http.ServeFile(w, r, "./tmpl/admin.html")
}

package route

import "net/http"

func RedirectToHome(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Add("HX-Redirect", "/")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	http.Redirect(w, r, "/", http.StatusFound)
}

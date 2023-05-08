package middleware

import (
	"fmt"
	"net/http"
)

type BasicAuth struct {
	Username  string
	Password  string
	AppliesTo []string
}

func (m *BasicAuth) Do(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		applies := false
		if m.Username == "" || m.Password == "" {
			next.ServeHTTP(w, r)
			return
		}

		for _, a := range m.AppliesTo {
			if r.URL.Path == a {
				applies = true
				break
			}
		}

		if !applies {
			next.ServeHTTP(w, r)
			return
		}

		user, pass, ok := r.BasicAuth()
		if !ok || user != m.Username || pass != m.Password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "401 Unauthorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}

package middleware

import (
	"fmt"
	"net/http"
)

func Authentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientToken :=	r.Header.Get("token")
		fmt.Println("token::: ",clientToken)
		if clientToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_, returnFlag := helper.ExtractClaims(clientToken)
		fmt.Println(returnFlag)
		if !returnFlag {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		do stuff
		w.WriteHeader(http.StatusUnauthorized)
		h.ServeHTTP(w, r)
	})
}


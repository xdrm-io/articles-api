package services

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/xdrm-io/aicra/api"
)

// AuthMiddleware is used to update request's context to consider permissions
// active as for the current request' user
//
// As this is an example project, there is no authorization mechanism, you just
// have to set a header `User` with the user id as a value to be authenticated
// as the provided user.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// pass request execution anyways
		defer next.ServeHTTP(w, r)

		// get current authentication
		auth := api.GetAuth(r.Context())
		if auth == nil {
			log.Println("cannot access authentication")
			return
		}

		userHeader := r.Header.Get("User")
		userID, err := strconv.ParseInt(userHeader, 10, 64)
		if err != nil {
			log.Println("invalid 'User' header; invalid int")
			return
		}

		// add permissions
		auth.Active = append(auth.Active, fmt.Sprintf("user[%d]", userID))
	})
}

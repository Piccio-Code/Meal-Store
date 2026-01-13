package main

import (
	"context"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"net/http"
)

type CurrentUserID string

const CurrentUserIDKey = CurrentUserID("CurrentUserIDKey")

func (app *application) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			token, err := jwt.ParseRequest(r, jwt.WithVerify(false)) // TODO: per la production toggle on

			if err != nil {
				app.errorLog.Println(err)
				app.UnauthorizedError(w, r)
				return
			}

			id, ok := token.Subject()

			if !ok {
				app.UnauthorizedError(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), CurrentUserIDKey, id)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		},
	)
}

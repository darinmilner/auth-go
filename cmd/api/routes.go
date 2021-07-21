package main

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) wrapMiddleware(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (app *application) routes() http.Handler {
	router := httprouter.New()
	secure := alice.New(app.checkToken)
	router.HandlerFunc(http.MethodGet, "/status", app.StatusHandler)
	router.HandlerFunc(http.MethodGet, "/v1/open-route/:id", app.OpenRoute)

	//Forget password
	router.HandlerFunc(http.MethodPost, "/v1/forgot-password", app.ForgotPassword)
	router.HandlerFunc(http.MethodPost, "/v1/reset-password", app.ResetPassword)

	//signin
	router.HandlerFunc(http.MethodPost, "/v1/signin", app.Signin)
	router.HandlerFunc(http.MethodPost, "/v1/register", app.Register)
	//Secure route
	router.GET("/v1/secure", app.wrapMiddleware(secure.ThenFunc(app.SecuredRoute)))

	return app.enableCORS(router)
}

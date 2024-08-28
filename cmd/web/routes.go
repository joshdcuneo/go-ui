package main

import (
	"net/http"

	"github.com/joshdcuneo/go-ui/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.EFS))

	mux.HandleFunc("GET /healthcheck", healthcheck)

	webMiddleware := newMiddlewareGroup(app.sessionManager.LoadAndSave, app.preventCSRF, app.authenticate)

	mux.Handle("GET /user/signup", webMiddleware.thenFunc(app.userSignup))
	mux.Handle("POST /user/signup", webMiddleware.thenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", webMiddleware.thenFunc(app.userLogin))
	mux.Handle("POST /user/login", webMiddleware.thenFunc(app.userLoginPost))
	mux.Handle("GET /{$}", webMiddleware.thenFunc(app.home))

	authMiddleware := webMiddleware.append(app.requireAuthentication)

	mux.Handle("POST /user/logout", authMiddleware.thenFunc(app.userLogoutPost))

	globalMiddleware := newMiddlewareGroup(app.recoverPanic, app.logRequest, commonHeaders)
	return globalMiddleware.then(mux)
}

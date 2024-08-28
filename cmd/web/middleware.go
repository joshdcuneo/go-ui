package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"slices"

	"github.com/justinas/nosurf"
)

type Middleware func(http.Handler) http.Handler

type middlewareGroup struct {
	middlewares []Middleware
}

func newMiddlewareGroup(middlewares ...Middleware) middlewareGroup {
	slices.Reverse(middlewares)
	return middlewareGroup{middlewares: middlewares}
}

func (m *middlewareGroup) then(handler http.Handler) http.Handler {
	for _, middleware := range m.middlewares {
		handler = middleware(handler)
	}
	return handler
}

func (m *middlewareGroup) thenFunc(handler http.HandlerFunc) http.Handler {
	return m.then(http.HandlerFunc(handler))
}

func (m *middlewareGroup) append(middleware Middleware) *middlewareGroup {
	m.middlewares = append([]Middleware{middleware}, m.middlewares...)
	return m
}

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO update for bunny
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info("received request",
			slog.String("ip", ip),
			slog.String("proto", proto),
			slog.String("method", method),
			slog.String("uri", uri),
		)
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, r, fmt.Errorf("%s", rec))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusFound)
		}

		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

func (app *application) preventCSRF(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := app.sessionManager.GetInt(r.Context(), string(authenticatedUserIDKey))
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		exists, _ := app.users.Exists(id)
		if !exists {
			app.sessionManager.Remove(r.Context(), string(authenticatedUserIDKey))
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

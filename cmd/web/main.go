package main

import (
	"crypto/tls"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"

	"github.com/joshdcuneo/go-ui/internal/dbutils"
	"github.com/joshdcuneo/go-ui/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type config struct {
	addr        string
	staticDir   string
	logFormat   string
	logLevel    string
	databaseDSN string
}

type application struct {
	users          models.UserModelInterface
	templates      map[string]*template.Template
	sessionManager *scs.SessionManager
	logger         *slog.Logger
}

func main() {

	cfg := config{}
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static/", "Path to static assets")
	flag.StringVar(&cfg.logFormat, "log-format", "text", "Log format (text or json)")
	flag.StringVar(&cfg.databaseDSN, "database-dsn", "./db.sqlite", "Database DSN")
	flag.Parse()

	logger := newLogger(cfg)

	db, err := dbutils.NewDB(cfg.databaseDSN)
	if err != nil {
		logger.Error(err.Error())
	}
	defer db.Close()

	templates, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
	}

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode

	app := &application{
		users:          &models.UserModel{DB: db},
		templates:      templates,
		sessionManager: sessionManager,
		logger:         logger,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         cfg.addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("starting server", slog.String("addr", cfg.addr))

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)
}

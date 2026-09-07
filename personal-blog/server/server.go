package server

import (
	"crypto/sha256"
	"crypto/subtle"
	"log"
	"net/http"
	"os"

	"github.com/pbnjk/backend/personal-blog/handler"
)

type Server struct {
	mux      *http.ServeMux
	username string
	password string
}

func NewServer() Server {
	server := Server{
		mux:      http.NewServeMux(),
		username: os.Getenv("ADMINUSER"),
		password: os.Getenv("ADMINPWD"),
	}
	server.setup()

	return server
}

func (s *Server) Serve() {
	certfile := os.Getenv("CERTFILE")
	keyfile := os.Getenv("KEYFILE")
	log.Fatal(http.ListenAndServeTLS(":8080", certfile, keyfile, s.mux))
}

func (s *Server) setup() {
	s.mux.HandleFunc("/", handler.HandleHome)
	s.mux.HandleFunc("GET /article/{id}", handler.HandleArticle)
	s.mux.HandleFunc("GET /admin", s.wrapAuth(handler.HandleAdmin))

	fs := http.FileServer(http.Dir("public/"))
	s.mux.Handle("/static/", http.StripPrefix("/static/", fs))
}

func (s *Server) wrapAuth(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()

		// Unauthorized
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			correctUsernameHash := sha256.Sum256([]byte(s.username))

			passwordHash := sha256.Sum256([]byte(password))
			correctPasswordHash := sha256.Sum256([]byte(s.password))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], correctUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], correctPasswordHash[:]) == 1)
			if usernameMatch && passwordMatch {
				f.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

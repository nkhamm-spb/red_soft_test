package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/nkhamm-spb/red_soft_test/config"
	_ "github.com/nkhamm-spb/red_soft_test/docs"
	"github.com/nkhamm-spb/red_soft_test/httpserver/httphandlers"
	"github.com/nkhamm-spb/red_soft_test/storage"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	config *config.Server

	httpServer *http.Server
	router     *mux.Router
}

func New(ctx context.Context, storage *storage.Storage, config *config.Server) (*Server, error) {
	log.Printf("Creating new HTTP server")

	server := &Server{config: config}
	server.httpServer = &http.Server{}

	server.router = mux.NewRouter()
	server.router.Handle("/api/users/{id:[0-9]+}/get_user", &httphandlers.HandlerGetUser{Storage: storage}).Methods("GET")
	server.router.Handle("/api/users/{id:[0-9]+}/edit_user", &httphandlers.HandlerEditUser{Storage: storage}).Methods("PUT")
	server.router.Handle("/api/users/add_user", &httphandlers.HandlerAddUser{Storage: storage}).Methods("POST")
	server.router.Handle("/api/users/get_by_surname/{surname}", &httphandlers.HandlerGetBySurname{Storage: storage}).Methods("GET")
	server.router.Handle("/api/users/get_all", &httphandlers.HandlerGetAll{Storage: storage}).Methods("GET")

	server.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // URL документации
	))

	log.Println("HTTP server is created")
	return server, nil
}

func (s *Server) Run() error {
	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	log.Printf("Starting HTTP server on address: %s\n", address)
	return http.ListenAndServe(address, s.router)
}

func (s *Server) Shutdown() error {
	log.Printf("Waiting for shutdown HTTP Server")

	cancelCtx, cancel := context.WithTimeout(context.Background(), time.Duration(5*int(time.Second)))
	defer cancel()

	err := s.httpServer.Shutdown(cancelCtx)

	return err
}

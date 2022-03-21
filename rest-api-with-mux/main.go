package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"example.com/with_mux/helpers"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Mux paketini dahil etmek için : go get -u github.com/gorilla/mux

func main() {
	r := mux.NewRouter()

	CORSOptions()
	// tüm response'ları sıkıştırmak için
	r.Use(loggingMiddleware)
	r.Use(authenticationMiddleware)

	s := r.PathPrefix("/products").Subrouter()
	// "/products/{name}/"
	s.HandleFunc("/{name}", ProductNameHandler)
	// "/products/{id:[0-9]+}"
	s.HandleFunc("/id/{id:[0-9]+}", ProductIdHandler)

	p := r.PathPrefix("/person").Subrouter()
	p.HandleFunc("/", personCreate).Methods("POST")

	// r.HandleFunc("/products/{name}", ProductNameHandler).
	// 	Methods("GET").
	// 	Schemes("http")
	// r.HandleFunc("/products/{id:[0-9]+}", ProductIdHandler)

	srv := &http.Server{
		Addr:         "0.0.0.0:8090",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	ShutdownServer(srv, time.Second*10)
}

func ProductIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/json")
	d := ApiResponse{
		Data: vars["id"],
	}
	resp, _ := json.Marshal(d)
	w.Write(resp)
}

func ProductNameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// eğer query string parametresi okumak istiyorsanız r.URL.Query().Get("param")
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-type", "application/json")
	d := ApiResponse{
		Data: vars["name"],
	}
	resp, _ := json.Marshal(d)
	w.Write(resp)
}

func personCreate(w http.ResponseWriter, r *http.Request) {
	var p Person

	err := helpers.DecodeJSONBody(w, r, &p)
	if err != nil {
		var mr *helpers.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintf(w, "Person: %+v", p)
}

type ApiResponse struct {
	Data interface{} `json:"data"`
}

type Person struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.URL.Query())
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		token := r.Header.Get("Authorization")
		if strings.HasPrefix(r.URL.Path, "/products") {
			if token != "" {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Token not found", http.StatusUnauthorized)
			}
		} else {
			next.ServeHTTP(w, r)
		}

	})
}

func ShutdownServer(srv *http.Server, timeout time.Duration) {
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func CORSOptions() {
	handlers.AllowedOrigins([]string{"https://www.example.com"})
	handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	handlers.AllowedMethods([]string{"POST", "GET", "PUT", "PATCH"})
}

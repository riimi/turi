package main

import (
	"context"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"mq/academy/handler"
	"mq/academy/repo"
	"mq/academy/session"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

type GameServer struct {
	Port       string
	IsReady    *atomic.Value
	IsHealthy  *atomic.Value
	Repo       repo.GameRepo
	SessionMgr session.Manager
}

type GameServerOption func(*GameServer) error

func Port(port string) GameServerOption {
	return func(server *GameServer) error {
		server.Port = port
		return nil
	}
}

func NewGameServer(opts ...GameServerOption) *GameServer {
	server := &GameServer{
		IsReady:    &atomic.Value{},
		IsHealthy:  &atomic.Value{},
		Repo:       repo.NewRepo(),
		SessionMgr: session.NewManagerInMem(),
		Port:       "1323",
	}

	for _, opt := range opts {
		if err := opt(server); err != nil {
			log.Fatal(err)
		}
	}

	return server
}

func (server *GameServer) Router() *mux.Router {
	if server.IsReady == nil || server.IsHealthy == nil {
		log.Fatal("server nil")
	}
	server.IsReady.Store(false)
	server.IsHealthy.Store(true)
	go func() {
		log.Printf("Readyz probe is negative by default...")
		time.Sleep(10 * time.Second)
		server.IsReady.Store(true)
		log.Printf("Readyz prove is now positive")
	}()

	r := mux.NewRouter()
	// for health check
	r.HandleFunc("/", handler.Home(BuildTime, Release)).Methods("GET")
	r.HandleFunc("/health", handler.Health(server.IsHealthy))
	r.HandleFunc("/healthy", handler.Healthy(server.IsHealthy))
	r.HandleFunc("/unhealthy", handler.Unhealthy(server.IsHealthy))
	r.HandleFunc("/ready", handler.Readyz(server.IsReady))

	// user apis
	r.HandleFunc("/users/{name}", handler.Users(server.Repo)).Methods("GET", "POST")
	r.HandleFunc("/users/{name}/friends", handler.Friends(server.Repo)).Methods("GET", "POST")
	r.HandleFunc("/users/{name}/maybe-friends", handler.MaybeFriends(server.Repo)).Methods("GET")
	r.HandleFunc("/login", handler.Login(server.Repo, server.SessionMgr)).Methods("POST")

	// for profiling
	debug := r.PathPrefix("/debug/pprof").Subrouter()
	debug.HandleFunc("/", pprof.Index)
	debug.HandleFunc("/cmdline", pprof.Cmdline)
	debug.HandleFunc("/profile", pprof.Profile)
	debug.HandleFunc("/symbol", pprof.Symbol)
	debug.HandleFunc("/trace", pprof.Trace)
	debug.HandleFunc("/{category}", pprof.Index)

	return r
}

func (server *GameServer) Run() {
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	srv := &http.Server{
		Addr:    ":" + server.Port,
		Handler: server.Router(),
	}

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
	log.Printf("The service is ready to listen and serve. port: %s", server.Port)

	<-interrupt
	log.Print("The service is shutting down...")
	srv.Shutdown(context.Background())
	log.Print("Done")
}

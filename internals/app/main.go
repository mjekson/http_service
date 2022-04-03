package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mjekson/http_service/api"
	"github.com/mjekson/http_service/api/middleware"
	"github.com/mjekson/http_service/internals/app/db"
	"github.com/mjekson/http_service/internals/app/handlers"
	"github.com/mjekson/http_service/internals/app/processors"
	"github.com/mjekson/http_service/internals/cfg"
)

type Server struct {
	config cfg.Cfg
	ctx    context.Context
	srv    *http.Server
	db     *pgxpool.Pool
}

func NewServer(config cfg.Cfg, ctx context.Context) *Server {
	server := new(Server)
	server.ctx = ctx
	server.config = config
	return server
}

func (server *Server) Serve() {
	log.Println("----Starting server----")
	var err error
	server.db, err = pgxpool.Connect(server.ctx, server.config.GetDBString())
	if err != nil {
		log.Fatalln(err)
	}
	carsStorage := db.NewCarsStorage(server.db)
	usersStorage := db.NewUsersStorage(server.db)

	carsProcessor := processors.NewCarsProcessor(carsStorage)
	usersProcessor := processors.NewUsersProcessor(usersStorage)

	userHandler := handlers.NewUsersHandler(usersProcessor)
	carsHandler := handlers.NewCarsHandler(carsProcessor)

	routes := api.CreateRoutes(userHandler, carsHandler)
	routes.Use(middleware.RequestLog)

	server.srv = &http.Server{
		Addr:    ":" + server.config.Port,
		Handler: routes,
	}
	log.Println("----Server started----")

	err = server.srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}

func (server *Server) Shutdown() {
	log.Printf("----Server stopped----")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	server.db.Close()
	defer func() {
		cancel()
	}()

	var err error
	if err = server.srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("----Shutdown Failed----:%v", err)
	}

	log.Printf("----Server exited properly----")

	if err == http.ErrServerClosed {
		err = nil
	}
}

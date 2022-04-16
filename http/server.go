package http

import (
	"fmt"
	"github.com/exluap/kit/log"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"net/http"
	"time"
)

const (
	ReadTimeout     = time.Second * 10
	WriteTimeout    = time.Second * 10
	ReadBufferSize  = 1024
	WriteBufferSize = 1024
)

type Cors struct {
	AllowedHeaders []string
	AllowedOrigins []string
	AllowedMethods []string
	Debug          bool
}

type Config struct {
	Port  string
	Cors  *Cors
	Trace bool
}

// Server represents HTTP server
type Server struct {
	Srv        *http.Server        // Srv - internal server
	RootRouter *mux.Router         // RootRouter - root router
	WsUpgrader *websocket.Upgrader // WsUpgrader - websocket upgrader
	logger     log.CLoggerFunc     // logger
}

type RouteSetter interface {
	Set() error
}

type WsUpgrader interface {
	Set(router *mux.Router, upgrader *websocket.Upgrader)
}

// getOptions getting cors options preconfigured
func getOptions(cfg *Config) cors.Options {
	if cfg.Cors == nil {
		return cors.Options{
			AllowCredentials: true,
		}
	}

	return cors.Options{
		AllowedOrigins:   cfg.Cors.AllowedOrigins,
		AllowedMethods:   cfg.Cors.AllowedMethods,
		AllowedHeaders:   cfg.Cors.AllowedHeaders,
		AllowCredentials: true,
		Debug:            cfg.Cors.Debug,
	}
}

func NewHttpServer(cfg *Config, logger log.CLoggerFunc) *Server {
	r := mux.NewRouter()

	corsHandler := cors.New(getOptions(cfg)).Handler(r)

	s := &Server{
		Srv: &http.Server{
			Addr:         fmt.Sprintf(":%s", cfg.Port),
			Handler:      corsHandler,
			WriteTimeout: WriteTimeout,
			ReadTimeout:  ReadTimeout,
		},
		WsUpgrader: &websocket.Upgrader{
			ReadBufferSize:  ReadBufferSize,
			WriteBufferSize: WriteBufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		logger: logger,
	}
	if cfg.Trace {
		r.Use(s.loggingMiddleware)
	}
	s.RootRouter = r

	return s
}

func (s *Server) SetWsUpgrader(upgradeSetter WsUpgrader) {
	upgradeSetter.Set(s.RootRouter, s.WsUpgrader)
}

func (s *Server) Listen() {
	go func() {
		l := s.logger().Pr("http").Cmp("server").Mth("listen").F(log.FF{"url": s.Srv.Addr})
		l.Inf("start listening")

		// if tls parameters are specified, list tls connection
		if err := s.Srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				l.E(ErrHttpSrvListen(err)).St().Err()
			} else {
				l.Dbg("server closed")
			}
			return
		}
	}()
}

func (s *Server) Close() {
	_ = s.Srv.Close()
}

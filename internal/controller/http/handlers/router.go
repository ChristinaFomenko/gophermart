package handler

import (
	"github.com/ChristinaFomenko/gophermart/internal/controller/http/middlewares"
	"github.com/ChristinaFomenko/gophermart/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"go.uber.org/zap"
)

const (
	signingKey = "WAdjK$12KCFjH6905u#uEZQjF349%835hFqzA"
)

type Handler struct {
	Service   *service.Service
	TokenAuth *jwtauth.JWTAuth
	log       *zap.Logger
}

func NewHandler(service *service.Service, log *zap.Logger) *Handler {
	tokenAuth := jwtauth.New("HS256", []byte(signingKey), nil)

	return &Handler{
		Service:   service,
		TokenAuth: tokenAuth,
		log:       log,
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middlewares.GzipHandle)

	router.Group(func(router chi.Router) {
		router.Post("/api/user/register", h.registration)
		router.Post("/api/user/login", h.authentication)
	})

	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(h.TokenAuth))
		router.Use(jwtauth.Authenticator)

		router.Post("/api/user/orders", h.loadOrders)
		router.Get("/api/user/orders", h.getUploadedOrders)
		router.Post("/api/user/balance/withdraw", h.deductionOfPoints)
		router.Get("/api/user/withdrawals", h.getWithdrawalOfPoints)
		router.Get("/api/user/balance", h.getCurrentBalance)
	})

	return router
}

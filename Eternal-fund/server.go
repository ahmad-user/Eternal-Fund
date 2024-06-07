package main

import (
	"database/sql"
	"eternal-fund/config"
	"eternal-fund/controller"
	"eternal-fund/middleware"
	"eternal-fund/repository"
	"eternal-fund/usecase"
	"eternal-fund/usecase/service"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Server struct {
	userUC        usecase.UserUseCase
	campaignsUC   usecase.CampaignsUseCase
	authUc        usecase.AuthUseCase
	transactionUC usecase.TransactionUseCase
	jwtService    service.JwtService
	engine        *gin.Engine
}

func (s *Server) initRoute() {
	rg := s.engine.Group("/api/v1")

	authMiddleware := middleware.NewAuthMiddleware(s.jwtService)
	controller.NewUserController(s.userUC, rg, authMiddleware).Routing()
	controller.NewCampaignsController(s.campaignsUC, rg, authMiddleware).Routing()
	controller.NewAuthController(s.authUc, rg).Route()
	controller.NewTransactionController(s.transactionUC, rg, authMiddleware).Routing()
}

func (s *Server) Run() {
	s.initRoute()

	s.engine.Run(":2000")
}

func NewServer() *Server {

	c, _ := config.NewConfig()

	urlConnect := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.DbPort, c.DbUser, c.DbPassword, c.DbName)

	database, err := sql.Open(c.Driver, urlConnect)
	if err != nil {
		panic("connection Error")
	}
	userRepo := repository.NewUserRepo(database)
	userUC := usecase.NewUserUseCase(userRepo)

	campaignsRepo := repository.NewCampaignsRepo(database)
	campaignsUseCase := usecase.NewCampaignsUseCase(campaignsRepo, userRepo)

	jwtService := service.NewJwtService(c.TokenConfig)
	authUseCase := usecase.NewAuthUseCase(jwtService, userUC)

	transactionRepo := repository.NewTransactionRepo(database)
	paymentService := service.NewPaymentService()
	transactionUC := usecase.NewTransactionUseCase(transactionRepo, campaignsRepo, paymentService)

	return &Server{
		userUC:        userUC,
		campaignsUC:   campaignsUseCase,
		transactionUC: transactionUC,
		engine:        gin.Default(),
		jwtService:    jwtService,
		authUc:        authUseCase,
	}

}

package main

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/ttani03/ouchi-ipam/docs"

	"github.com/ttani03/ouchi-ipam/gen/sqlc"
)

type Server struct {
	db *sqlc.Queries

	echo *echo.Echo
}

func NewPool(url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}

func NewSever(pool *pgxpool.Pool) (*Server, error) {
	s := &Server{
		db: sqlc.New(pool),
	}

	s.init()
	return s, nil
}

func (s *Server) init() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// swagger API
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Subnet API
	e.GET("/subnets", s.getSubnets)
	e.GET("/subnets/:id", s.getSubnet)
	e.POST("/subnets", s.createSubnet)
	e.DELETE("/subnets/:id", s.deleteSubnet)

	// IpAddress API
	e.GET("/subnets/:id/ip", s.getIpAddresses)
	e.POST("/subnets/:id/ip", s.reserveIpAddress)
	e.POST("/subnets/:id/ip/:ip", s.reserveSpecifiedIpAddress)
	e.DELETE("/subnets/:id/ip/:ip", s.freeIpAddress)

	s.echo = e
}

func (s *Server) Start(address string) error {
	return s.echo.Start(address)
}

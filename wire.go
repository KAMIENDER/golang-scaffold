//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/KAMIENDER/golang-scaffold/infra/auth"
	"github.com/KAMIENDER/golang-scaffold/infra/config"
	"github.com/KAMIENDER/golang-scaffold/infra/database/mysql"
	"github.com/KAMIENDER/golang-scaffold/infra/database/nosql"
	"github.com/KAMIENDER/golang-scaffold/server"
)

var BardSet = wire.NewSet(
	config.NewConfig,
	nosql.NewRedis,
	mysql.NewDatabase,
	gin.New,
	wire.Struct(new(server.Handler), "*"),
	auth.NewAuthManager,

	wire.Bind(new(nosql.NoSQLDB), new(*nosql.Redis)),
)

func NewHandler() (*server.Handler, error) {
	wire.Build(BardSet)
	return nil, nil
}

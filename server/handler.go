package server

import (
	"context"
	"net/http"

	"github.com/KAMIENDER/golang-scaffold/infra/auth"
	bizerror "github.com/KAMIENDER/golang-scaffold/infra/biz_error"
	"github.com/KAMIENDER/golang-scaffold/infra/database/nosql"
	"github.com/KAMIENDER/golang-scaffold/view"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handle[Req any, Data any](c *gin.Context, execute func(ctx context.Context, req Req) (Data, *bizerror.BizError)) {
	var req Req
	resp := view.Resp[Data]{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := execute(c, req)
	if err != nil {
		resp.Code = err.Code()
		resp.Msg = err.Msg()
		c.JSON(http.StatusOK, gin.H{"resp": resp})
		return
	}
	resp.Data = data
	resp.Code = "00000"
	c.JSON(http.StatusOK, gin.H{"resp": resp})
}

type Handler struct {
	Router      *gin.Engine
	NoSQLClient nosql.NoSQLDB
	DB          *gorm.DB
	AuthManager *auth.AuthManager
}

func (h *Handler) innerInit() {
	h.Router.Use([]gin.HandlerFunc{
		gin.Logger(),
		gin.Recovery(),
	}...)

	h.AuthManager.SetupAuthBoss(h.Router)

}

func (h *Handler) Run() {
	h.innerInit()
	h.Router.Run()
}

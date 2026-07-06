package http

import (
	"fmt"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	svc port.StoreService
}

func NewStoreHandler(svc port.StoreService) *StoreHandler {
	return &StoreHandler{svc}
}

type registerStoreRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=50" example:"我的寄卖行"`
	Email    string `json:"email" binding:"required,email" example:"admin@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"12345678"`
}

type registerResponse struct {
	Store  *domain.Store `json:"store"`
	User   *domain.User  `json:"user"`
	Token  string        `json:"token"`
	Domain string        `json:"domain"`
}

func (sh *StoreHandler) Register(ctx *gin.Context) {
	var req registerStoreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	store, user, token, err := sh.svc.Register(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := &registerResponse{
		Store:  store,
		User:   user,
		Token:  token,
		Domain: store.Domain,
	}

	handleSuccess(ctx, rsp)
}

func (sh *StoreHandler) GetStoreByDomain(ctx *gin.Context) {
	domain := ctx.Query("domain")
	if domain == "" {
		validationError(ctx, fmt.Errorf("domain is required"))
		return
	}

	store, err := sh.svc.GetStoreByDomain(ctx, domain)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, store)
}
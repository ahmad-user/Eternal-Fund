package controller

import (
	"eternal-fund/model/dto"
	commonresponse "eternal-fund/model/dto/common_response"
	"eternal-fund/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUc      usecase.AuthUseCase
	routerGroup *gin.RouterGroup
}

func (a *AuthController) loginHandler(ctx *gin.Context) {
	var payload dto.AuthReqDto
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	token, err := a.authUc.Login(payload)
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	commonresponse.SendSingleResponse(ctx, token, "Success Login")

}

func (a *AuthController) Route() {
	a.routerGroup.POST("/auth/login", a.loginHandler)
}

func NewAuthController(authUc usecase.AuthUseCase, rg *gin.RouterGroup) *AuthController {
	return &AuthController{authUc: authUc, routerGroup: rg}
}

package controller

import (
	"net/http"
	"strconv"
	"time"

	"eternal-fund/middleware"
	"eternal-fund/model"
	commonresponse "eternal-fund/model/dto/common_response"
	"eternal-fund/usecase"

	"github.com/gin-gonic/gin"
)

type campaignController struct {
	campaignUseCase usecase.CampaignsUseCase
	router          *gin.RouterGroup
	authMiddleware  middleware.AuthMiddleware
}

func (cc *campaignController) createCampaignHandler(ctx *gin.Context) {
	var input model.Campaigns
	if err := ctx.ShouldBindJSON(&input); err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	campaign, err := cc.campaignUseCase.CreateCampaigns(input)
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	campaign.User.PasswordHash = ""
	campaign.User.CreatedAt = time.Time{}
	campaign.User.UpdatedAt = time.Time{}

	commonresponse.SendSingleResponse(ctx, campaign, "Campaign created successfully")
}

func (cc *campaignController) getCampaignsHandler(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid page number")
		return
	}
	size, err := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid size number")
		return
	}

	campaigns, paging, err := cc.campaignUseCase.FindAllCampaigns(page, size)
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var data []interface{}
	for _, campaign := range campaigns {
		data = append(data, campaign)
	}

	commonresponse.SendManyResponse(ctx, data, paging, "Campaigns retrieved successfully")
}

func (cc *campaignController) getCampaignByIdHandler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid campaign ID")
		return
	}

	campaign, err := cc.campaignUseCase.FindByIdCampaigns(id)
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	campaign.User.PasswordHash = ""
	campaign.User.CreatedAt = time.Time{}
	campaign.User.UpdatedAt = time.Time{}

	commonresponse.SendSingleResponse(ctx, campaign, "Campaign retrieved successfully")
}

func (cc *campaignController) updateCampaignHandler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid campaign ID")
		return
	}
	updatedCampaign, err := cc.campaignUseCase.UpdateCampaigns(id)
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	updatedCampaign.User.PasswordHash = ""
	updatedCampaign.User.CreatedAt = time.Time{}
	updatedCampaign.User.UpdatedAt = time.Time{}
	commonresponse.SendSingleResponse(ctx, updatedCampaign, "Campaign updated successfully")
}

func (cc *campaignController) deleteCampaignHandler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid campaign ID")
		return
	}

	err = cc.campaignUseCase.DeleteCampaigns(id)
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	commonresponse.SendSingleResponse(ctx, nil, "Campaign deleted successfully")
}

func (cc *campaignController) Routing() {
	cc.router.POST("/campaigns", cc.authMiddleware.CheckToken("user", "admin"), cc.createCampaignHandler)
	cc.router.GET("/campaigns", cc.getCampaignsHandler)
	cc.router.GET("/campaigns/:campaign_id", cc.authMiddleware.CheckToken("user", "admin"), cc.getCampaignByIdHandler)
	cc.router.PUT("/campaigns/:campaign_id", cc.authMiddleware.CheckToken("user", "admin"), cc.updateCampaignHandler)
	cc.router.DELETE("/campaigns/:campaign_id", cc.authMiddleware.CheckToken("user", "admin"), cc.deleteCampaignHandler)

}

func NewCampaignsController(campaignUseCase usecase.CampaignsUseCase, rg *gin.RouterGroup, authMiddleware middleware.AuthMiddleware) *campaignController {
	return &campaignController{
		campaignUseCase: campaignUseCase,
		router:          rg,
		authMiddleware:  authMiddleware,
	}
}

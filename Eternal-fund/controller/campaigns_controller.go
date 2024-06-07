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

func (cc *campaignController) uploadCampaignImageHandler(ctx *gin.Context) {
	campaignID, err := strconv.Atoi(ctx.Param("campaign_id"))
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid campaign ID")
		return
	}
	file, err := ctx.FormFile("file")
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, "Error in file upload: "+err.Error())
		return
	}

	var input model.CampaignImage
	if err := ctx.ShouldBind(&input); err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	fileLocation := "images/campaigns/" + file.Filename
	if err := ctx.SaveUploadedFile(file, fileLocation); err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusInternalServerError, "Error saving file: "+err.Error())
		return
	}
	input.CampaignID = campaignID
	input.FileLocation = fileLocation
	campaignImage, err := cc.campaignUseCase.SaveCampaignImage(input, fileLocation)
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusInternalServerError, "Error saving campaign image: "+err.Error())
		return
	}

	commonresponse.SendSingleResponse(ctx, campaignImage, "Campaign image uploaded successfully")
}

func (cc *campaignController) Routing() {
	cc.router.POST("/campaigns", cc.authMiddleware.CheckToken("user", "admin"), cc.createCampaignHandler)
	cc.router.GET("/campaigns", cc.getCampaignsHandler)
	cc.router.GET("/campaigns/:campaign_id", cc.authMiddleware.CheckToken("user", "admin"), cc.getCampaignByIdHandler)
	cc.router.PUT("/campaigns/:campaign_id", cc.authMiddleware.CheckToken("user", "admin"), cc.updateCampaignHandler)
	cc.router.DELETE("/campaigns/:campaign_id", cc.authMiddleware.CheckToken("user", "admin"), cc.deleteCampaignHandler)

	cc.router.POST("/campaigns/:campaign_id/images", cc.authMiddleware.CheckToken("user", "admin"), cc.uploadCampaignImageHandler)
}

func NewCampaignsController(campaignUseCase usecase.CampaignsUseCase, rg *gin.RouterGroup, authMiddleware middleware.AuthMiddleware) *campaignController {
	return &campaignController{
		campaignUseCase: campaignUseCase,
		router:          rg,
		authMiddleware:  authMiddleware,
	}
}

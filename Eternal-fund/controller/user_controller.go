package controller

import (
	"eternal-fund/middleware"
	"eternal-fund/model"
	"eternal-fund/model/dto"
	commonresponse "eternal-fund/model/dto/common_response"
	"eternal-fund/usecase"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type userController struct {
	userUseCase    usecase.UserUseCase
	router         *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func (u *userController) listHandler(ctx *gin.Context) {
	page, er := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, er2 := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	// validation query params
	if er != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, er.Error())
	}
	if er2 != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, er2.Error())
	}

	listData, paging, err := u.userUseCase.FindAll(page, size)
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	var data []interface{}
	for _, b := range listData {
		data = append(data, b)
	}
	commonresponse.SendManyResponse(ctx, data, paging, "ok")
}

func (u *userController) getByIdHandler(ctx *gin.Context) {
    idStr := ctx.Param("user_id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        // Jika terjadi error saat parsing, kirim respons error ke client
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid user ID format",
        })
        return // Pastikan untuk menghentikan eksekusi dengan return
    }

    data, err := u.userUseCase.FindById(id)
    if err != nil {
        // Jika user tidak ditemukan, kirim respons error ke client
        ctx.JSON(http.StatusNotFound, gin.H{
            "error": "User not found",
        })
        return // Pastikan untuk menghentikan eksekusi dengan return
    }

    // Jika tidak ada error, kirim data user sebagai respons
    ctx.JSON(http.StatusOK, &dto.SingleResponse{
        Status: dto.Status{
            Code:    http.StatusOK,
            Message: "ok",
        },
        Data: data,
    })
}

func (u *userController) registerHandler(ctx *gin.Context) {
	var input model.RegisterUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, err := u.userUseCase.RegisterUser(input)
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	commonresponse.SendSingleResponse(ctx, user, "User registered successfully")
}

func (u *userController) updateUserHandler(ctx *gin.Context) {
    var input model.User
    if err := ctx.ShouldBindJSON(&input); err != nil {
        commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
        return
    }

    userIdStr := ctx.Param("user_id")
    userId, err := strconv.Atoi(userIdStr)
    if err != nil {
        commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID")
        return
    }

    input.ID = userId // Set the ID to the input struct
    updatedUser, err := u.userUseCase.UpdateUser(userId, input)
    if err != nil {
        commonresponse.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
        return
    }

    commonresponse.SendSingleResponse(ctx, updatedUser, "User updated successfully")
}


func (u *userController) saveAvatarHandler(ctx *gin.Context) {
    userIdStr := ctx.Param("id")
    userId, err := strconv.Atoi(userIdStr)
    if err != nil {
        commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID")
        return
    }

    file, err := ctx.FormFile("avatar")
    if err != nil {
        commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
        return
    }

    // Validasi ekstensi file
    ext := filepath.Ext(file.Filename)
    if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
        commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid file type")
        return
    }

    // Pembuatan nama file yang unik
    newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
    fileLocation := "images/avatars/" + newFileName

    // Pembuatan direktori jika belum ada
    if err := os.MkdirAll("images/avatars", os.ModePerm); err != nil {
        commonresponse.SendErrorResponse(ctx, http.StatusInternalServerError, "Could not create directory")
        return
    }

    // Menyimpan file yang diunggah
    if err := ctx.SaveUploadedFile(file, fileLocation); err != nil {
        commonresponse.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
        return
    }

    // Menyimpan lokasi file avatar dalam database
    user, err := u.userUseCase.SaveAvatar(userId, fileLocation)
    if err != nil {
        commonresponse.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
        return
    }

    commonresponse.SendSingleResponse(ctx, user, "Avatar saved successfully")
}

func (u *userController) isEmailAvailableHandler(ctx *gin.Context) {
	var input model.CheckEmailInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	isAvailable, err := u.userUseCase.IsEmailAvailable(input)
	if err != nil {
		commonresponse.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, &dto.SingleResponse{
		Status: dto.Status{
			Code:    http.StatusOK,
			Message: "ok",
		},
		Data: map[string]bool{"is_available": isAvailable},
	})
}

func (u *userController) Routing() {
	u.router.GET("/users", u.authMiddleware.CheckToken("user"), u.listHandler)
	u.router.GET("/users/:user_id", u.authMiddleware.CheckToken("user", "admin"), u.getByIdHandler)
	u.router.POST("/register", u.registerHandler)
	u.router.PUT("/users/:user_id", u.authMiddleware.CheckToken("user", "admin"), u.updateUserHandler)
	u.router.POST("/users/:id/avatar", u.authMiddleware.CheckToken("user"), u.saveAvatarHandler)
	u.router.POST("/users/check-email", u.isEmailAvailableHandler)

}

func NewUserController(userUc usecase.UserUseCase, rg *gin.RouterGroup, authMiddle middleware.AuthMiddleware) *userController {
	return &userController{
		userUseCase:    userUc,
		router:         rg,
		authMiddleware: authMiddle,
	}
}

package commonresponse

import (
	"eternal-fund/model/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendSingleResponse(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, &dto.SingleResponse{
		Status: dto.Status{
			Code:    http.StatusOK,
			Message: message,
		},
		Data: data,
	})
}

func SendManyResponse(c *gin.Context, data []interface{}, paging dto.Paging, message string) {
	c.JSON(http.StatusOK, &dto.ManyResponse{
		Status: dto.Status{
			Code:    http.StatusOK,
			Message: message,
		},
		Data:   data,
		Paging: paging,
	})
}

func SendErrorResponse(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, &dto.SingleResponse{
		Status: dto.Status{
			Code:    code,
			Message: message,
		},
	})
}

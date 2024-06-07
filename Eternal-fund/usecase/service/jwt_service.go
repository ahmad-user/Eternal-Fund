package service

import (
	"eternal-fund/config"
	"eternal-fund/model"
	"eternal-fund/model/dto"
	"eternal-fund/utils"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	CreateToken(user model.User) (dto.AuthResponDto, error)
	ValidateToken(token string) (jwt.MapClaims, error)
}

type jwtService struct {
	co config.TokenConfig
}

// CreateToken implements JwtService.
func (j *jwtService) CreateToken(user model.User) (dto.AuthResponDto, error) {
	claims := utils.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.co.IssuerName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.co.ExpiresTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role:     user.Role,
		UserId: strconv.Itoa(user.ID),
	}

	token := jwt.NewWithClaims(j.co.SigningMethod, claims)
	ss, err := token.SignedString(j.co.SignatureKey)
	if err != nil {
		return dto.AuthResponDto{}, fmt.Errorf("failed create access token")
	}
	return dto.AuthResponDto{Token: ss}, nil
}

// ValidateToken implements JwtService.
func (j *jwtService) ValidateToken(tokenHeader string) (jwt.MapClaims, error) {

	log.Println("---- kepanggil ---", tokenHeader)

	// token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
	// 	return j.co.SignatureKey, nil
	// })

	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		return j.co.SignatureKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to verify token when parsing")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to verify token when claims")
	}

	return claims, nil

}

func NewJwtService(c config.TokenConfig) JwtService {
	return &jwtService{co: c}
}

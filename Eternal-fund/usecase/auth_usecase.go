package usecase

import (
	"eternal-fund/model/dto"
	"eternal-fund/usecase/service"
)

type AuthUseCase interface {
	Login(payload dto.AuthReqDto) (dto.AuthResponDto, error)
}
type authUseCase struct {
	jwtService service.JwtService
	userUC   UserUseCase
}

// Login implements AuthUseCase.
func (a *authUseCase) Login(payload dto.AuthReqDto) (dto.AuthResponDto, error) {
	user, err := a.userUC.FindByEmail(payload.Email)
	if err != nil {
		return dto.AuthResponDto{}, err
	}
	token, err := a.jwtService.CreateToken(user)
	if err != nil {
		return dto.AuthResponDto{}, err
	}
	return token, nil
}
func NewAuthUseCase(jwtService service.JwtService, userUc UserUseCase) AuthUseCase {
	return &authUseCase{jwtService: jwtService, userUC: userUc}
}

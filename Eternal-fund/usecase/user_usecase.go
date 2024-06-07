package usecase

import (
	"database/sql"
	"errors"
	"eternal-fund/model"
	"eternal-fund/model/dto"
	"eternal-fund/repository"

	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	repo repository.UserRepo
}

func (u *userUseCase) RegisterUser(input model.RegisterUserInput) (model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}

	user := model.User{
		Name:         input.Name,
		Occupation:   input.Occupation,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		Role:         "user",
	}

	return u.repo.Save(user)
}

func (u *userUseCase) UpdateUser(userId int, input model.User) (model.User, error) {
    user, err := u.repo.FindById(userId)
    if err != nil {
        return user, err
    }
    user.Name = input.Name
    user.Occupation = input.Occupation
    user.Email = input.Email
    updatedUser, err := u.repo.Update(user)
    if err != nil {
        return updatedUser, err
    }
    return updatedUser, nil
}


func (u *userUseCase) SaveAvatar(userId int, fileLocation string) (model.User, error) {
	return u.repo.SaveAvatar(userId, fileLocation)
}

func (u *userUseCase) IsEmailAvailable(input model.CheckEmailInput) (bool, error) {
	email := input.Email
	_, err := u.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

// findAll implements AuthorUseCase.
func (u *userUseCase) FindAll(page int, size int) ([]model.User, dto.Paging, error) {
	return u.repo.FindAll(page, size)
}

// findById implements userUseCase.
func (u *userUseCase) FindById(id int) (model.User, error) {
	return u.repo.FindById(id)
}

func (u *userUseCase) FindByEmail(email string) (model.User, error) {
	return u.repo.FindByEmail(email)
}

type UserUseCase interface {
	RegisterUser(input model.RegisterUserInput) (model.User, error)
	UpdateUser(userId int, input model.User) (model.User, error)
	SaveAvatar(userId int, fileLocation string) (model.User, error)
	IsEmailAvailable(input model.CheckEmailInput) (bool, error)
	FindAll(page int, size int) ([]model.User, dto.Paging, error)
	FindById(id int) (model.User, error)
	FindByEmail(email string) (model.User, error)
}

func NewUserUseCase(repo repository.UserRepo) UserUseCase {
	return &userUseCase{repo: repo}
}

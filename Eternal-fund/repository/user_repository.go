package repository

import (
	"database/sql"
	"eternal-fund/model"
	"eternal-fund/model/dto"
	"log"
	"math"
	"time"
)

type userRepo struct {
	db *sql.DB
}

func (u *userRepo) Save(user model.User) (model.User, error) {
	query := "INSERT INTO users (name, occupation, email, password_hash, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id, created_at, updated_at"
	var id int
	var createdAt, updatedAt time.Time
	err := u.db.QueryRow(query, user.Name, user.Occupation, user.Email, user.PasswordHash, user.Role).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return user, err
	}
	user.ID = id
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt
	return user, nil
}

func (u *userRepo) Update(user model.User) (model.User, error) {
	query := "UPDATE users SET name = $1, occupation = $2, email = $3, updated_at = NOW() WHERE id = $4 RETURNING id, name, occupation, email, updated_at"
	err := u.db.QueryRow(query, user.Name, user.Occupation, user.Email, user.ID).Scan(&user.ID, &user.Name, &user.Occupation, &user.Email, &user.UpdatedAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u *userRepo) SaveAvatar(userId int, fileLocation string) (model.User, error) {
	query := "UPDATE users SET avatar_file_name = $1, updated_at = NOW() WHERE id = $2 RETURNING id, name, occupation, email, password_hash, avatar_file_name, role, created_at, updated_at"
	var user model.User
	err := u.db.QueryRow(query, fileLocation, userId).Scan(
		&user.ID, &user.Name, &user.Occupation, &user.Email, &user.PasswordHash, &user.AvatarFileName, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u *userRepo) FindAll(page int, size int) ([]model.User, dto.Paging, error) {
	var listData []model.User
	var rows *sql.Rows

	//rumus pagination
	offset := (page - 1) * size

	var err error
	rows, err = u.db.Query("SELECT * FROM users limit $1 offset $2", size, offset)
	if err != nil {
		return nil, dto.Paging{}, err
	}

	totalRows := 0
	err = u.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalRows)
	if err != nil {
		return nil, dto.Paging{}, err
	}

	for rows.Next() {
		var user model.User
		var avatarFileName sql.NullString

		err := rows.Scan(&user.ID, &user.Name, &user.Occupation, &user.Email, &user.PasswordHash, &user.AvatarFileName, &user.Role, &user.CreatedAt, &user.UpdatedAt)

		if err != nil {
			log.Println(err.Error())
		}

		if avatarFileName.Valid {
			user.AvatarFileName = &avatarFileName.String
		} else {
			user.AvatarFileName = nil
		}

		listData = append(listData, user)
	}

	paging := dto.Paging{
		Page:       page,
		Size:       size,
		TotalRows:  totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return listData, paging, nil

}

func (u *userRepo) FindById(id int) (model.User, error) {
	var user model.User
	var avatarFileName sql.NullString

	err := u.db.QueryRow("SELECT * FROM users WHERE id=$1", id).Scan(&user.ID, &user.Name, &user.Occupation, &user.Email, &user.PasswordHash, &user.AvatarFileName, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return model.User{}, err
	}

	if avatarFileName.Valid {
		user.AvatarFileName = &avatarFileName.String
	} else {
		user.AvatarFileName = nil
	}

	return user, nil

}

func (u *userRepo) FindByEmail(email string) (model.User, error) {
	var user model.User
	var avatarFileName sql.NullString

	err := u.db.QueryRow("SELECT * FROM users WHERE email=$1", email).
		Scan(&user.ID, &user.Name, &user.Occupation, &user.Email, &user.PasswordHash, &user.AvatarFileName, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return model.User{}, err
	}

	if avatarFileName.Valid {
		user.AvatarFileName = &avatarFileName.String
	} else {
		user.AvatarFileName = nil
	}

	return user, nil
}

type UserRepo interface {
	Save(user model.User) (model.User, error)
	Update(user model.User) (model.User, error)
	SaveAvatar(userId int, fileLocation string) (model.User, error)
	FindAll(page int, size int) ([]model.User, dto.Paging, error)
	FindById(id int) (model.User, error)
	FindByEmail(email string) (model.User, error)
}

func NewUserRepo(database *sql.DB) UserRepo {
	return &userRepo{db: database}
}

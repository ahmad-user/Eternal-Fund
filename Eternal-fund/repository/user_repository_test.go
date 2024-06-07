package repository

import (
	"database/sql"
	"eternal-fund/model"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UsersRepoTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSql sqlmock.Sqlmock
	repo    UserRepo
}

func (suite *UsersRepoTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New()
	suite.mockDB = db
	suite.mockSql = mock
	suite.repo = NewUserRepo(suite.mockDB)
}

func (suite *UsersRepoTestSuite) TestFindByEmail_Success() {
	email := "existing@example.com"
	expectedUser := model.User{
		ID:             0,
		Name:           "Test User",
		Occupation:     "Tester",
		Email:          email,
		PasswordHash:   "hashed_password",
		AvatarFileName: nil,
		Role:           "user",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Mock query database dan hasilnya
	suite.mockSql.ExpectQuery(`SELECT \* FROM users WHERE email=\$1`).
		WithArgs(email).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "occupation", "email", "password_hash", "avatar_file_name", "role", "created_at", "updated_at"}).
				AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Occupation, expectedUser.Email, expectedUser.PasswordHash, nil, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt),
		)
	user, err := suite.repo.FindByEmail(email)
	assert.NoError(suite.T(), err, "Diharapkan tidak ada error")
	assert.Equal(suite.T(), expectedUser.ID, user.ID, "ID pengguna diharapkan sesuai")
	assert.Equal(suite.T(), expectedUser.Name, user.Name, "Nama pengguna diharapkan sesuai")
	assert.Equal(suite.T(), expectedUser.Email, user.Email, "Email pengguna diharapkan sesuai")
	assert.Equal(suite.T(), expectedUser.PasswordHash, user.PasswordHash, "Hash kata sandi pengguna diharapkan sesuai")
	assert.Equal(suite.T(), expectedUser.Role, user.Role, "Peran pengguna diharapkan sesuai")
	assert.Equal(suite.T(), expectedUser.CreatedAt, user.CreatedAt, "Waktu pembuatan pengguna diharapkan sesuai")
	assert.Equal(suite.T(), expectedUser.UpdatedAt, user.UpdatedAt, "Waktu pembaruan pengguna diharapkan sesuai")
	assert.Nil(suite.T(), user.AvatarFileName, "Nama file avatar diharapkan nil")
	err = suite.mockSql.ExpectationsWereMet()
	assert.NoError(suite.T(), err, "Ada ekspektasi yang tidak terpenuhi")
}

func (suite *UsersRepoTestSuite) TestFindByID_Success() {
	userID := 1
	expectedUser := model.User{
		ID:             userID,
		Name:           "Test User",
		Occupation:     "Tester",
		Email:          "test@example.com",
		PasswordHash:   "hashed_password",
		AvatarFileName: nil,
		Role:           "user",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	suite.mockSql.ExpectQuery("^SELECT \\* FROM users WHERE id=\\$1").
		WithArgs(userID).
		WillReturnRows(
			suite.mockSql.NewRows([]string{"id", "name", "occupation", "email", "password_hash", "avatar_file_name", "role", "created_at", "updated_at"}).
				AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Occupation, expectedUser.Email, expectedUser.PasswordHash, nil, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt),
		)
	user, err := suite.repo.FindById(userID)
	assert.NoError(suite.T(), err, "Expected no error")
	assert.Equal(suite.T(), expectedUser.ID, user.ID, "Expected user ID to match")
	assert.Equal(suite.T(), expectedUser.Name, user.Name, "Expected user name to match")
	assert.Equal(suite.T(), expectedUser.Email, user.Email, "Expected user email to match")
	err = suite.mockSql.ExpectationsWereMet()
	assert.NoError(suite.T(), err, "There were unfulfilled expectations")
}

func (suite *UsersRepoTestSuite) TestSaveAvatar_Success() {
	userID := 1
	fileLocation := "/path/to/avatar.jpg"
	expectedUser := model.User{
		ID:             userID,
		Name:           "Test User",
		Occupation:     "Tester",
		Email:          "test@example.com",
		PasswordHash:   "hashed_password",
		AvatarFileName: &fileLocation, // Assuming the avatar file name is updated to this location
		Role:           "user",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	suite.mockSql.ExpectQuery("^UPDATE users SET avatar_file_name = \\$1, updated_at = NOW\\(\\) WHERE id = \\$2 RETURNING id, name, occupation, email, password_hash, avatar_file_name, role, created_at, updated_at").
		WithArgs(fileLocation, userID).
		WillReturnRows(
			suite.mockSql.NewRows([]string{"id", "name", "occupation", "email", "password_hash", "avatar_file_name", "role", "created_at", "updated_at"}).
				AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Occupation, expectedUser.Email, expectedUser.PasswordHash, expectedUser.AvatarFileName, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt),
		)
	user, err := suite.repo.SaveAvatar(userID, fileLocation)
	assert.NoError(suite.T(), err, "Expected no error")
	assert.Equal(suite.T(), expectedUser.ID, user.ID, "Expected user ID to match")
	assert.Equal(suite.T(), expectedUser.Name, user.Name, "Expected user name to match")
	assert.Equal(suite.T(), expectedUser.Email, user.Email, "Expected user email to match")
	assert.Equal(suite.T(), expectedUser.AvatarFileName, user.AvatarFileName, "Expected avatar file name to match")
	err = suite.mockSql.ExpectationsWereMet()
	assert.NoError(suite.T(), err, "There were unfulfilled expectations")
}

func (suite *UsersRepoTestSuite) TestSave_Success() {
	expectedUser := model.User{
		Name:         "Test User",
		Occupation:   "Tester",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		Role:         "user",
	}
	suite.mockSql.ExpectQuery("^INSERT INTO users \\(name, occupation, email, password_hash, role, created_at, updated_at\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, NOW\\(\\), NOW\\(\\)\\) RETURNING id, created_at, updated_at").
		WithArgs(expectedUser.Name, expectedUser.Occupation, expectedUser.Email, expectedUser.PasswordHash, expectedUser.Role).
		WillReturnRows(
			suite.mockSql.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(1, time.Now(), time.Now()),
		)
	user, err := suite.repo.Save(expectedUser)
	assert.NoError(suite.T(), err, "Expected no error")
	assert.NotNil(suite.T(), user.ID, "Expected user ID to be set")
	assert.Equal(suite.T(), expectedUser.Name, user.Name, "Expected user name to match")
	assert.Equal(suite.T(), expectedUser.Email, user.Email, "Expected user email to match")
	err = suite.mockSql.ExpectationsWereMet()
	assert.NoError(suite.T(), err, "There were unfulfilled expectations")
}

func (suite *UsersRepoTestSuite) TestUpdate_Success() {
	userID := 0
	updatedName := ""
	updatedOccupation := ""
	updatedEmail := ""
	suite.mockSql.ExpectQuery("^UPDATE users SET name = \\$1, occupation = \\$2, email = \\$3, updated_at = NOW\\(\\) WHERE id = \\$4 RETURNING updated_at").
		WithArgs(updatedName, updatedOccupation, updatedEmail, userID).
		WillReturnRows(
			suite.mockSql.NewRows([]string{"updated_at"}).
				AddRow(time.Now()),
		)
	user, err := suite.repo.Update(model.User{
		ID:           userID,
		Name:         updatedName,
		Occupation:   updatedOccupation,
		Email:        updatedEmail,
	})
	assert.NoError(suite.T(), err, "Expected no error")
	assert.NotNil(suite.T(), user.UpdatedAt, "Expected UpdatedAt to be set")
	assert.Equal(suite.T(), userID, user.ID, "Expected user ID to match")
	assert.Equal(suite.T(), updatedName, user.Name, "Expected user name to match")
	assert.Equal(suite.T(), updatedEmail, user.Email, "Expected user email to match")
	err = suite.mockSql.ExpectationsWereMet()
	assert.NoError(suite.T(), err, "There were unfulfilled expectations")
}

func (suite *UsersRepoTestSuite) TestFindAll_Success() {
	page := 1
	size := 10
	expectedUsers := []model.User{
		{
			ID:             1,
			Name:           "User 1",
			Occupation:     "Developer",
			Email:          "user1@example.com",
			PasswordHash:   "hashed_password_1",
			AvatarFileName: nil,
			Role:           "user",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             2,
			Name:           "User 2",
			Occupation:     "Tester",
			Email:          "user2@example.com",
			PasswordHash:   "hashed_password_2",
			AvatarFileName: nil,
			Role:           "user",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	userRows := sqlmock.NewRows([]string{"id", "name", "occupation", "email", "password_hash", "avatar_file_name", "role", "created_at", "updated_at"})
	for _, user := range expectedUsers {
		userRows.AddRow(user.ID, user.Name, user.Occupation, user.Email, user.PasswordHash, user.AvatarFileName, user.Role, user.CreatedAt, user.UpdatedAt)
	}
	suite.mockSql.ExpectQuery(`SELECT \* FROM users limit \$1 offset \$2`).
		WithArgs(size, (page-1)*size).
		WillReturnRows(userRows)
	suite.mockSql.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(expectedUsers)))
	users, paging, err := suite.repo.FindAll(page, size)
	assert.NoError(suite.T(), err, "Expected no error")
	assert.Equal(suite.T(), len(expectedUsers), len(users), "Expected number of users to match")
	for i, expectedUser := range expectedUsers {
		assert.Equal(suite.T(), expectedUser.ID, users[i].ID, "Expected user ID to match")
		assert.Equal(suite.T(), expectedUser.Name, users[i].Name, "Expected user name to match")
		assert.Equal(suite.T(), expectedUser.Email, users[i].Email, "Expected user email to match")
	}
	assert.Equal(suite.T(), page, paging.Page, "Expected page number to match")
	assert.Equal(suite.T(), size, paging.Size, "Expected page size to match")
	assert.Equal(suite.T(), len(expectedUsers), paging.TotalRows, "Expected total rows to match")
	assert.Equal(suite.T(), 1, paging.TotalPages, "Expected total pages to match") // Assuming page size of 10 for simplicity
}

func TestUserRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UsersRepoTestSuite))
}

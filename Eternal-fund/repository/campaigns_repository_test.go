package repository

import (
	"database/sql"
	"eternal-fund/model"
	"eternal-fund/model/dto"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CampaignsRepoTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSql sqlmock.Sqlmock
	repo    CampaignsRepo
}

func (suite *CampaignsRepoTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New()
	suite.mockDB = db
	suite.mockSql = mock
	suite.repo = NewCampaignsRepo(suite.mockDB)
}

var expectedCampaigns = []model.Campaigns{
	{
		ID:                1,
		User_id:           101,
		Name:              "Campaign 1",
		Short_description: "Short description 1",
		Description:       "Description 1",
		Perks:             "Perks 1",
		Backer_count:      10,
		Goal_amount:       1000,
		Current_amount:    500,
		Slug:              "campaign-1",
		Created_at:        time.Now(),
		Updated_at:        time.Now(),
	},
	{
		ID:                2,
		User_id:           102,
		Name:              "Campaign 2",
		Short_description: "Short description 2",
		Description:       "Description 2",
		Perks:             "Perks 2",
		Backer_count:      20,
		Goal_amount:       2000,
		Current_amount:    1500,
		Slug:              "campaign-2",
		Created_at:        time.Now(),
		Updated_at:        time.Now(),
	},
}

func (suite *CampaignsRepoTestSuite) TestGetAll_success() {
	page := 1
	size := 2
	offset := (page - 1) * size

	expectPaging := dto.Paging{
		Page:       page,
		Size:       size,
		TotalRows:  5,
		TotalPages: 3,
	}
	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "short_description", "description", "perks", "backer_count", "goal_amount", "current_amount", "slug", "created_at", "updated_at"}).
		AddRow(expectedCampaigns[0].ID, expectedCampaigns[0].User_id, expectedCampaigns[0].Name, expectedCampaigns[0].Short_description, expectedCampaigns[0].Description, expectedCampaigns[0].Perks, expectedCampaigns[0].Backer_count, expectedCampaigns[0].Goal_amount, expectedCampaigns[0].Current_amount, expectedCampaigns[0].Slug, expectedCampaigns[0].Created_at, expectedCampaigns[0].Updated_at).
		AddRow(expectedCampaigns[1].ID, expectedCampaigns[1].User_id, expectedCampaigns[1].Name, expectedCampaigns[1].Short_description, expectedCampaigns[1].Description, expectedCampaigns[1].Perks, expectedCampaigns[1].Backer_count, expectedCampaigns[1].Goal_amount, expectedCampaigns[1].Current_amount, expectedCampaigns[1].Slug, expectedCampaigns[1].Created_at, expectedCampaigns[1].Updated_at)

	suite.mockSql.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM campaigns limit $1 offset $2`)).
		WithArgs(size, offset).WillReturnRows(rows)

	totalRows := sqlmock.NewRows([]string{"COUNT"}).AddRow(5)
	suite.mockSql.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM campaigns")).WillReturnRows(totalRows)

	campaigns, paging, err := suite.repo.FindAllCampaigns(page, size)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectPaging, paging)
	assert.Equal(suite.T(), expectedCampaigns, campaigns)
}

func (suite *CampaignsRepoTestSuite) TestFindById_Success() {
	expectedCampaign := expectedCampaigns[0]

	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "short_description", "description", "perks", "backer_count", "goal_amount", "current_amount", "slug", "created_at", "updated_at"}).
		AddRow(expectedCampaign.ID, expectedCampaign.User_id, expectedCampaign.Name, expectedCampaign.Short_description,
			expectedCampaign.Description, expectedCampaign.Perks, expectedCampaign.Backer_count, expectedCampaign.Goal_amount,
			expectedCampaign.Current_amount, expectedCampaign.Slug, expectedCampaign.Created_at, expectedCampaign.Updated_at)
	expectedQuery := `SELECT \* FROM campaigns where id=\$1`

	suite.mockSql.ExpectQuery(expectedQuery).
		WithArgs(expectedCampaign.ID).
		WillReturnRows(rows)

	actualCampaign, err := suite.repo.FindByIdCampaigns(expectedCampaign.ID)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedCampaign.ID, actualCampaign.ID)
	assert.Equal(suite.T(), expectedCampaign.User_id, actualCampaign.User_id)
	assert.Equal(suite.T(), expectedCampaign.Name, actualCampaign.Name)
	assert.Equal(suite.T(), expectedCampaign.Short_description, actualCampaign.Short_description)
	assert.Equal(suite.T(), expectedCampaign.Description, actualCampaign.Description)
	assert.Equal(suite.T(), expectedCampaign.Perks, actualCampaign.Perks)
	assert.Equal(suite.T(), expectedCampaign.Backer_count, actualCampaign.Backer_count)
	assert.Equal(suite.T(), expectedCampaign.Goal_amount, actualCampaign.Goal_amount)
	assert.Equal(suite.T(), expectedCampaign.Current_amount, actualCampaign.Current_amount)
	assert.Equal(suite.T(), expectedCampaign.Slug, actualCampaign.Slug)

	err = suite.mockSql.ExpectationsWereMet()
	assert.NoError(suite.T(), err)
}

func (suite *CampaignsRepoTestSuite) TestDelete_Success() {
	id := 1

	suite.mockSql.ExpectExec(regexp.QuoteMeta("DELETE FROM campaigns WHERE id = $1")).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := suite.repo.DeleteCampaigns(id)
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), suite.mockSql.ExpectationsWereMet())
}

func (suite *CampaignsRepoTestSuite) TestFindByUserID_Success() {
	userID := 123

	expectedCampaigns := []model.Campaigns{
		{ID: 1, User_id: userID, Name: "Campaign 1" /* other fields */},
		{ID: 2, User_id: userID, Name: "Campaign 2" /* other fields */},
	}
	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "short_description", "description", "perks", "backer_count", "goal_amount", "current_amount", "slug", "created_at", "updated_at"})
	for _, campaign := range expectedCampaigns {
		rows.AddRow(campaign.ID, campaign.User_id, campaign.Name, campaign.Short_description,
			campaign.Description, campaign.Perks, campaign.Backer_count, campaign.Goal_amount,
			campaign.Current_amount, campaign.Slug, campaign.Created_at, campaign.Updated_at)
	}

	expectedQuery := `SELECT \* FROM campaigns WHERE user_id = \$1`
	suite.mockSql.ExpectQuery(expectedQuery).
		WithArgs(userID).
		WillReturnRows(rows)
	actualCampaigns, err := suite.repo.FindByUserID(userID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), len(expectedCampaigns), len(actualCampaigns), "Number of campaigns returned should match")

	for i := range expectedCampaigns {
		assert.Equal(suite.T(), expectedCampaigns[i].ID, actualCampaigns[i].ID)
		assert.Equal(suite.T(), expectedCampaigns[i].User_id, actualCampaigns[i].User_id)
		assert.Equal(suite.T(), expectedCampaigns[i].Name, actualCampaigns[i].Name)
	}
	err = suite.mockSql.ExpectationsWereMet()
	assert.NoError(suite.T(), err)
}

func (suite *CampaignsRepoTestSuite) TestMarkAllImagesAsNonPrimary_Success() {
	campaignID := 123
	expectedQuery := `UPDATE campaign_images SET is_primary = false WHERE campaign_id = \$1`

	suite.mockSql.ExpectExec(expectedQuery).
		WithArgs(campaignID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	success, err := suite.repo.MarkAllImagesAsNonPrimary(campaignID)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), success, "Expected success to be true")
	err = suite.mockSql.ExpectationsWereMet()
	assert.NoError(suite.T(), err)
}

func (suite *CampaignsRepoTestSuite) TestCreateCampaigns_Success() {
	// Prepare mock campaign data for insertion
	mockCampaign := model.Campaigns{
		User_id:           0,
		Name:              "",
		Short_description: "",
		Description:       "",
		Perks:             "",
		Backer_count:      0,
		Goal_amount:       0,
		Current_amount:    0,
		Slug:              "",
	}
	expectedQuery := `INSERT INTO campaigns`
	expectedCampaignID := 0

	suite.mockSql.ExpectPrepare(expectedQuery)
	suite.mockSql.ExpectQuery(expectedQuery).
		WithArgs(mockCampaign.User_id, mockCampaign.Name, mockCampaign.Short_description, mockCampaign.Description, mockCampaign.Perks,
			mockCampaign.Backer_count, mockCampaign.Goal_amount, mockCampaign.Current_amount, mockCampaign.Slug).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedCampaignID))
	createdCampaign, err := suite.repo.CreateCampaigns(mockCampaign)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedCampaignID, createdCampaign.ID)
	assert.Equal(suite.T(), mockCampaign.User_id, createdCampaign.User_id)
	assert.Equal(suite.T(), mockCampaign.Name, createdCampaign.Name)
	assert.Equal(suite.T(), mockCampaign.Short_description, createdCampaign.Short_description)
	assert.Equal(suite.T(), mockCampaign.Description, createdCampaign.Description)
	assert.Equal(suite.T(), mockCampaign.Perks, createdCampaign.Perks)
	assert.Equal(suite.T(), mockCampaign.Backer_count, createdCampaign.Backer_count)
	assert.Equal(suite.T(), mockCampaign.Goal_amount, createdCampaign.Goal_amount)
	assert.Equal(suite.T(), mockCampaign.Current_amount, createdCampaign.Current_amount)
	assert.Equal(suite.T(), mockCampaign.Slug, createdCampaign.Slug)

	// Verify expectations
	err = suite.mockSql.ExpectationsWereMet()
	assert.NoError(suite.T(), err)
}

// func (suite *CampaignsRepoTestSuite) TestUpdate_Success() {
// 	updatedCampaign := model.Campaigns{
// 		ID:                70,
// 		User_id:           2,
// 		Name:              "Kampanye Diperbarui",
// 		Short_description: "Deskripsi Singkat Diperbarui",
// 		Description:       "Deskripsi Diperbarui",
// 		Perks:             "Manfaat Diperbarui",
// 		Backer_count:      15,
// 		Goal_amount:       1500,
// 		Current_amount:    800,
// 		Slug:              "kampanye-diperbarui",
// 		Created_at:        time.Now(),
// 		Updated_at:        time.Now(),
// 	}
// 	// Memeriksa eksekusi persiapan query
// 	suite.mockSql.ExpectPrepare(regexp.QuoteMeta(`
// 		UPDATE campaigns
// 		SET user_id = $1, name = $2, short_description = $3, description = $4,
// 		perks = $5, backer_count = $6, goal_amount = $7, current_amount = $8, slug = $9, updated_at = NOW()
// 		WHERE id = $10
// 		RETURNING id, user_id, name, short_description, description, perks, backer_count, goal_amount, current_amount, slug, created_at, updated_at
// 	`))

// 	// Memeriksa eksekusi query update
// 	suite.mockSql.ExpectQuery(regexp.QuoteMeta("UPDATE campaigns")).WithArgs(
// 		updatedCampaign.User_id, updatedCampaign.Name, updatedCampaign.Short_description, updatedCampaign.Description,
// 		updatedCampaign.Perks, updatedCampaign.Backer_count, updatedCampaign.Goal_amount, updatedCampaign.Current_amount,
// 		updatedCampaign.Slug, updatedCampaign.ID,
// 	).WillReturnRows(sqlmock.NewRows([]string{
// 		"id", "user_id", "name", "short_description", "description", "perks", "backer_count",
// 		"goal_amount", "current_amount", "slug", "created_at", "updated_at",
// 	}).AddRow(
// 		updatedCampaign.ID, updatedCampaign.User_id, updatedCampaign.Name, updatedCampaign.Short_description,
// 		updatedCampaign.Description, updatedCampaign.Perks, updatedCampaign.Backer_count, updatedCampaign.Goal_amount,
// 		updatedCampaign.Current_amount, updatedCampaign.Slug, updatedCampaign.Created_at, updatedCampaign.Updated_at,
// 	))

// 	// Memanggil metode update kampanye
// 	actualUpdatedCampaign, err := suite.repo.UpdateCampaigns(updatedCampaign.ID)

// 	// Memeriksa apakah tidak ada error
// 	assert.NoError(suite.T(), err)

// 	// Membandingkan hasil yang diharapkan dengan hasil aktual
// 	assert.Equal(suite.T(), updatedCampaign, actualUpdatedCampaign)

// 	// Memeriksa bahwa semua ekspektasi query telah terpenuhi
// 	assert.Nil(suite.T(), suite.mockSql.ExpectationsWereMet())
// }

// func (suite *CampaignsRepoTestSuite) TestCreate_Success() {
// 	expectedCampaign := model.Campaigns{
// 		ID:                1,
// 		User_id:           101,
// 		Name:              "Campaign 1",
// 		Short_description: "Short description 1",
// 		Description:       "Description 1",
// 		Perks:             "Perks 1",
// 		Backer_count:      10,
// 		Goal_amount:       1000,
// 		Current_amount:    500,
// 		Slug:              "campaign-1",
// 		Created_at:        time.Now(),
// 		Updated_at:        time.Now(),
// 	}
// 	stmt := `INSERT INTO campaigns (user_id, name, short_description, description, perks,  backer_count, goal_amount,
// 		current_amount, slug, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8,$9, NOW(), NOW()) RETURNING id`
// 	suite.mockSql.ExpectPrepare(stmt).
// 		ExpectQuery().
// 		WithArgs(expectedCampaign.User_id, expectedCampaign.Name, expectedCampaign.Short_description,
// 			expectedCampaign.Description, expectedCampaign.Perks, expectedCampaign.Backer_count,
// 			expectedCampaign.Goal_amount, expectedCampaign.Current_amount, expectedCampaign.Slug).
// 		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
// 	createdCampaign, err := suite.repo.CreateCampaigns(expectedCampaign)
// 	assert.NoError(suite.T(), err)
// 	assert.NotNil(suite.T(), createdCampaign)
// 	assert.Equal(suite.T(), expectedCampaign.ID, createdCampaign.ID)
// 	assert.Equal(suite.T(), expectedCampaign.User_id, createdCampaign.User_id)
// 	assert.Equal(suite.T(), expectedCampaign.Name, createdCampaign.Name)
// 	assert.Equal(suite.T(), expectedCampaign.Short_description, createdCampaign.Short_description)
// 	assert.Equal(suite.T(), expectedCampaign.Description, createdCampaign.Description)
// 	assert.Equal(suite.T(), expectedCampaign.Perks, createdCampaign.Perks)
// 	assert.Equal(suite.T(), expectedCampaign.Backer_count, createdCampaign.Backer_count)
// 	assert.Equal(suite.T(), expectedCampaign.Goal_amount, createdCampaign.Goal_amount)
// 	assert.Equal(suite.T(), expectedCampaign.Current_amount, createdCampaign.Current_amount)
// 	assert.Equal(suite.T(), expectedCampaign.Slug, createdCampaign.Slug)
// 	assert.Nil(suite.T(), suite.mockSql.ExpectationsWereMet())
// }

func TestCampaignsRepoTestSuite(t *testing.T) {
	suite.Run(t, new(CampaignsRepoTestSuite))
}

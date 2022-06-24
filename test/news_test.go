package persistence

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ntm/internal/domain/news"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type NewsRepoTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	repo news.Repository
}

func (suite *NewsRepoTestSuite) SetupTest() {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	gdb, err1 := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{}) // open gorm db
	suite.Equal(nil, err1)

	repo := news.NewRepository(gdb)
	suite.mock = mock
	suite.repo = repo
}

func (suite *NewsRepoTestSuite) TestSaveNewsSuite() {
	id := uint(1)
	news := news.News{Title: "mutual fund is safe investment type", Writer: "budi", Status: "publish", Content: "mutual fund is safe investment typ", TopicID: 1}
	const sql = `INSERT INTO "news" (.+) RETURNING`

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(sql).WithArgs(news.Title, news.Writer, news.Content, news.Status, news.TopicID).
		WillReturnRows(sqlmock.NewRows([]string{"publish_at", "created_at", "updated_at", "deleted_at", "id"}).AddRow(nil, time.Now(), time.Now(), time.Now(), id))
	suite.mock.ExpectCommit()

	news, err := suite.repo.Upsert(context.Background(), news)
	suite.Empty(err)
	suite.Equal(id, news.ID)
	if err != nil {
		suite.T().Errorf("error '%v' was not expected, while inserting a row", err)
	}

	if err := suite.mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestNewsRepoTestSuite(t *testing.T) {
	suite.Run(t, new(NewsRepoTestSuite))
}

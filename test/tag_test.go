package persistence

// Basic imports
import (
	"context"
	"database/sql/driver"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ntm/internal/domain/tag"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TagRepoTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	repo tag.Repository
}

func (suite *TagRepoTestSuite) SetupTest() {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}

	gdb, err1 := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	suite.Equal(nil, err1)

	repo := tag.NewRepository(gdb)
	suite.mock = mock
	suite.repo = repo
}

func (suite *TagRepoTestSuite) TestSaveTagSuite() {
	id := uint(1)
	tag := tag.Tag{Tag: "fund", CreatedAt: time.Now()}
	const sql = `INSERT INTO "tags" ("tag","created_at") VALUES ($1,$2) RETURNING "created_at","updated_at","id"`

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(sql).WithArgs(tag.Tag, tag.CreatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(1, time.Now(), time.Now()))
	suite.mock.ExpectCommit()

	tag, err := suite.repo.Upsert(context.Background(), tag)
	suite.Empty(err)
	suite.Equal(id, tag.ID)
	if err != nil {
		suite.T().Errorf("error '%v' was not expected, while inserting a row", err)
	}

	if err := suite.mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *TagRepoTestSuite) TestGetTagSuite() {
	rows := sqlmock.
		NewRows([]string{"id", "tag", "updated_at", "created_at"}).
		AddRow(1, "fund", nil, time.Now())

	const sql = `SELECT * FROM "tags" WHERE tag = $1`
	const q = "fund"

	suite.mock.
		ExpectQuery(sql).
		WithArgs(q).
		WillReturnRows(rows)

	ctx := context.WithValue(context.Background(), tag.ContextKey("tags_filter"), tag.Filter{Tag: q})
	tag, err := suite.repo.GetAll(ctx)
	suite.NotEmpty(tag)
	suite.Empty(err)
	if err != nil {
		suite.T().Errorf("error '%v' was not expected, while getting a row", err)
	}

	if err := suite.mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *TagRepoTestSuite) TestGetTagByIdSuite() {
	rows := sqlmock.
		NewRows([]string{"id", "tag", "updated_at", "created_at"}).
		AddRow(1, "fund", nil, time.Now())

	const sql = `SELECT * FROM "tags" WHERE "tags"."id" = $1 ORDER BY "tags"."id" LIMIT 1`
	const id = 1

	suite.mock.
		ExpectQuery(sql).
		WithArgs(id).
		WillReturnRows(rows)

	tag, err := suite.repo.GetByID(context.Background(), id)
	suite.NotEmpty(tag)
	suite.Empty(err)
	if err != nil {
		suite.T().Errorf("error '%v' was not expected, while getting a row", err)
	}

	if err := suite.mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *TagRepoTestSuite) TestDeleteTagSuite() {
	const sql = `DELETE FROM "tags" WHERE "tags"."id" = $1`
	const id = 1

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(sql).WithArgs(id).WillReturnResult(driver.RowsAffected(1))
	suite.mock.ExpectCommit()

	err := suite.repo.DeleteByID(context.Background(), id)
	suite.Empty(err)
	if err != nil {
		suite.T().Errorf("error '%v' was not expected, while getting a row", err)
	}

	if err := suite.mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *TagRepoTestSuite) TestGetAllTagSuite() {
	rows := sqlmock.
		NewRows([]string{"id", "tag", "updated_at", "created_at"}).
		AddRow(1, "fund", nil, time.Now()).AddRow(1, "investment", nil, time.Now())

	const sql = `SELECT * FROM "tags"`

	suite.mock.
		ExpectQuery(sql).
		WillReturnRows(rows)

	ctx := context.WithValue(context.Background(), tag.ContextKey("tags_filter"), tag.Filter{})
	tag, err := suite.repo.GetAll(ctx)
	suite.Empty(err)
	suite.NotEmpty(tag)
	if err != nil {
		suite.T().Errorf("error '%v' was not expected, while getting a row", err)
	}

	if err := suite.mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTagRepoTestSuite(t *testing.T) {
	suite.Run(t, new(TagRepoTestSuite))
}

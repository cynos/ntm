package persistence

import (
	"context"
	"database/sql/driver"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ntm/internal/domain/topic"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TopicRepoTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	repo topic.Repository
}

func (suite *TopicRepoTestSuite) SetupTest() {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}

	gdb, err1 := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	suite.Equal(nil, err1)

	repo := topic.NewRepository(gdb)
	suite.mock = mock
	suite.repo = repo
}

func (suite *TopicRepoTestSuite) TestSaveTopicSuite() {
	id := uint(1)
	topic := topic.Topic{Topic: "investment", CreatedAt: time.Now()}
	const sql = `INSERT INTO "topics" ("topic","created_at") VALUES ($1,$2) RETURNING "created_at","updated_at","id"`

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(sql).WithArgs(topic.Topic, topic.CreatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(id, time.Now()))
	suite.mock.ExpectCommit()

	topic, err := suite.repo.Upsert(context.Background(), topic)
	suite.Empty(err)
	suite.Equal(id, topic.ID)
	if err != nil {
		suite.T().Errorf("error '%v' was not expected, while inserting a row", err)
	}

	if err := suite.mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *TopicRepoTestSuite) TestGetTopicSuite() {
	rows := sqlmock.
		NewRows([]string{"id", "topic", "updated_at", "created_at"}).
		AddRow(1, "investment", nil, time.Now())

	const sql = `SELECT * FROM "topics" WHERE topic = $1`
	const q = "investment"

	suite.mock.
		ExpectQuery(sql).
		WithArgs(q).
		WillReturnRows(rows)

	ctx := context.WithValue(context.Background(), topic.ContextKey("topics_filter"), topic.Filter{Topic: q})
	topic, err := suite.repo.GetAll(ctx)
	suite.NotEmpty(topic)
	suite.Empty(err)
	if err != nil {
		suite.T().Errorf("error '%v' was not expected, while getting a row", err)
	}

	if err := suite.mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *TopicRepoTestSuite) TestGetTopicByIdSuite() {
	rows := sqlmock.
		NewRows([]string{"id", "tag", "updated_at", "created_at"}).
		AddRow(1, "fund", nil, time.Now())

	const sql = `SELECT * FROM "topics" WHERE "topics"."id" = $1 ORDER BY "topics"."id" LIMIT 1`
	const id = 1

	suite.mock.
		ExpectQuery(sql).
		WithArgs(id).
		WillReturnRows(rows)

	topic, err := suite.repo.GetByID(context.Background(), id)
	suite.NotEmpty(topic)
	suite.Empty(err)
	if err != nil {
		suite.T().Errorf("error '%v' was not expected, while getting a row", err)
	}

	if err := suite.mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *TopicRepoTestSuite) TestDeleteTopicSuite() {
	const sql = `DELETE FROM "topics" WHERE "topics"."id" = $1`
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

func (suite *TopicRepoTestSuite) TestGetAllTopicSuite() {
	rows := sqlmock.
		NewRows([]string{"id", "topic", "updated_at", "created_at"}).
		AddRow(1, "fund", nil, time.Now()).AddRow(1, "investment fund", nil, time.Now())

	const sqlQuery = `SELECT * FROM "topics"`

	suite.mock.
		ExpectQuery(sqlQuery).
		WillReturnRows(rows)

	ctx := context.WithValue(context.Background(), topic.ContextKey("topics_filter"), topic.Filter{})
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

func TestTopicRepoTestSuite(t *testing.T) {
	suite.Run(t, new(TopicRepoTestSuite))
}

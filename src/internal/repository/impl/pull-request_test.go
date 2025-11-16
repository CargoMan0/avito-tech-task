package impl

import (
	"context"
	"github.com/CargoMan0/avito-tech-task/internal/domain"
	"github.com/CargoMan0/avito-tech-task/internal/repository/impl/repofixtures"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

var (
	nonExistentTestID = uuid.MustParse("789eb75c-e37d-47e1-8874-514fc3672087")
	testTime          = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	testPR            = &domain.PullRequest{
		ID:                uuid.MustParse("7fc69137-3b84-4620-8847-1d7cc6b713bd"),
		MergedAt:          &testTime,
		Name:              "test-pr-1",
		AuthorID:          uuid.MustParse("8a0a5a3f-ace7-4803-8486-9e3b54bfa30a"),
		Status:            domain.PullRequestStatusOpen,
		NeedMoreReviewers: false,
		Reviewers:         testReviewersArray,
	}

	testTeam = &domain.Team{
		Name:  "test-team",
		Users: testReviewersArray,
	}

	testReviewersArray = []domain.User{
		{
			ID:       uuid.MustParse("7fc69137-3b84-4620-8847-1d7cc6b710bd"),
			TeamName: "test-team",
			Name:     "test-pr-1",
			IsActive: true,
		},
		{
			ID:       uuid.MustParse("7fc69137-3b84-4620-8847-1d7cc6b711bd"),
			TeamName: "test-team",
			Name:     "test-pr-1",
			IsActive: true,
		},
	}
)

type pullRequestsRepositoryTestSuite struct {
	suite.Suite
	repo *PullRequestRepository

	teamFixture         *repofixtures.TeamFixture
	pullRequestsFixture *repofixtures.PullRequestFixture
	testDB              *repofixtures.TestDB
}

func TestPaymentRepository(t *testing.T) {
	suite.Run(t, new(pullRequestsRepositoryTestSuite))
}

// SetupSuite function calls only one time before testing flow.
func (p *pullRequestsRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	testDB, err := repofixtures.NewTestDB(ctx)
	if err != nil {
		p.T().Fatalf("failed to setup test database: %v", err)
	}

	p.testDB = testDB

	sqlDB := testDB.GetSQLDB()
	p.repo = NewPullRequestRepository(sqlDB)
	p.pullRequestsFixture = repofixtures.NewPullRequestFixture(sqlDB)
	p.teamFixture = repofixtures.NewTeamFixture(sqlDB)
}

// TearDownTest function calls after each separate test and clears data
// from the previous test. This function is needed to provide
// data isolation between tests.
func (p *pullRequestsRepositoryTestSuite) TearDownTest() {
	// Clearing test data from the previous test.
	p.pullRequestsFixture.MustClearTestData()
}

// TearDownSuite function calls only one time after testing flow.
func (p *pullRequestsRepositoryTestSuite) TearDownSuite() {
	err := p.testDB.ShutDown()
	if err != nil {
		p.T().Fatalf("failed to shut down test database: %v", err)
	}

	p.testDB = nil
	p.repo = nil
}

func (p *pullRequestsRepositoryTestSuite) TestPullRequestRepository_PullRequestExists() {
	ctx := context.Background()

	// Prepare team
	p.teamFixture.MustPrepareTestTeam(ctx, testTeam)

	// Prepare PR
	p.pullRequestsFixture.MustPrepareTestPullRequest(testPR)

	// Check if exists for existing payment, should result no error and true.
	exists, err := p.repo.PullRequestExists(ctx, testPR.ID)
	p.Require().NoError(err)
	p.Require().True(exists)

	// Check if exists for non-existent payment, should result no error and false.
	exists, err = p.repo.PullRequestExists(ctx, nonExistentTestID)
	p.Require().NoError(err)
	p.Require().False(exists)

	// If uuid.Nil is passed, should result no error and false.
	exists, err = p.repo.PullRequestExists(ctx, uuid.Nil)
	p.Require().NoError(err)
	p.Require().False(exists)
}

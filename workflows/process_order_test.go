package workflows

import (
	"idempotence-by-validation/activities"

	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"

)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) Test_Workflow() {
	env := s.NewTestWorkflowEnvironment()
	env.RegisterActivity(activities.GetToken)
	orderID := mock.AnythingOfType("string")
	s.Assertions.NotEmpty(orderID)
	//env.OnActivity(activities.CheckFraud, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
	env.OnActivity(activities.GetToken).
		Return(func(ctx context.Context) (string, error) {

			return activities.GetToken(ctx)
		})

	//assert we got a token
	//s.Assertions.NotEmpty()

	env.OnActivity(activities.GetReservation, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(func(ctx context.Context) (string, error) {

			return activities.GetReservation(ctx, "asdf", "TOKEN-12345")
		})

	env.ExecuteWorkflow(ProcessOrder)

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())

	var result string
	s.NoError(env.GetWorkflowResult(&result))
	s.Equal("Branch 0 done in 1ns.", result)
	env.AssertExpectations(s.T())
}

func Test_GetToken_Activity(t *testing.T) {
	testSuite := testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()
	
	env.RegisterActivity(activities.GetToken)


	token, err := env.ExecuteActivity(activities.GetToken)
	assert.NoError(t, err)

	assert.NotEmpty(t, token)
}

func Test_GetReservation_Activity(t *testing.T) {
	testSuite := testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()
	
	env.RegisterActivity(activities.GetReservation)


	reservation, err := env.ExecuteActivity(activities.GetReservation, "order-1234", "TOKEN-12345")
	assert.NoError(t, err)

	assert.NotEmpty(t, reservation)
}
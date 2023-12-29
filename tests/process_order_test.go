package workflows

import (
	"idempotence-by-validation/activities"
	//"context"
	"fmt"
	"math/rand"
	"testing"
	//"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"

	"idempotence-by-validation/ticket"
	"idempotence-by-validation/workflows"
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
	env.RegisterActivity(activities.CreateTicket)
	env.RegisterActivity(activities.GetReservation)
	env.RegisterActivity(activities.ValidateTicket)
	input := ticket.TicketOrder{
		OrderID: fmt.Sprintf("order-%d", rand.Intn(99999)),
	}
	s.Assertions.NotEmpty(input.OrderID)

	//env.OnActivity(activities.GetToken).
	//	Return(func(ctx context.Context) (string, error) {
	//
	//		return activities.GetToken(ctx)
	//	})

	env.ExecuteWorkflow(workflows.ProcessOrder, input)

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())

	var result string
	s.NoError(env.GetWorkflowResult(&result))
	s.Equal("Order Managed", result)
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

func Test_ValidateTicket_Activity(t *testing.T) {
	testSuite := testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()

	env.RegisterActivity(activities.ValidateTicket)

	reservation, err := env.ExecuteActivity(activities.ValidateTicket, "test-24594", "134740", "TOKEN-12345")
	assert.NoError(t, err)

	assert.NotEmpty(t, reservation)
}
package workflows

import (
	"idempotence-by-validation/activities"
	"idempotence-by-validation/ticket"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	u "idempotence-by-validation/utils"
)

// ProcessOrder is a function that handles the inventory workflow.
// It takes a workflow context as input and returns an error if any.
func ProcessOrder(ctx workflow.Context, order ticket.TicketOrder) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info(u.ColorGreen, order.OrderID, " - Starting ProcessOrder, Order Details: ", order, u.ColorReset)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	token := ""
	// get token - retryable like normal, it's failure-prone and idempotent
	err := workflow.ExecuteActivity(ctx, activities.GetToken).Get(ctx, &token)
	if err != nil {
		logger.Error(u.ColorRed, "GetToken activity failed.", "Error", err)
		return "", err
	}

	// get reservation - retryable like normal, it's failure-prone and idempotent
	reservation := ""
	err = workflow.ExecuteActivity(ctx, activities.GetReservation, order.OrderID, token).Get(ctx, &reservation)
	if err != nil {
		logger.Error(u.ColorRed, "GetReservation activity failed.", "Error", err, u.ColorReset)
		return "", err
	}

	// idempotency loop: retry the create ticket execution at most once, validate it, and if it doesn't succeed, try to create ticket again


	for ticketFound := false; !ticketFound;  {

		// RetryPolicy specifies how to automatically handle retries if an Activity fails.
		// we want to only run once
		// "CREATE-TICKET-ERROR", "CREATE-TICKET-TIMEOUT"
		activityretrypolicy := &temporal.RetryPolicy{
			MaximumAttempts:        1, // zero retries
			NonRetryableErrorTypes: []string{"CREATE-TICKET-ERROR", "CREATE-TICKET-TIMEOUT"},
		}

		activityoptions := workflow.ActivityOptions{
			StartToCloseTimeout: time.Minute,         // Timeout options specify when to automatically timeout Activity functions.
			RetryPolicy:         activityretrypolicy, // Temporal retries failed Activities by default.
		}

		originalActivityOptions := workflow.GetActivityOptions(ctx)
		// Apply the options.
		ctx = workflow.WithActivityOptions(ctx, activityoptions)

		ticketCreateErr := workflow.ExecuteActivity(ctx, activities.CreateTicket, order.OrderID, reservation, token).Get(ctx, &order.Ticket)
		if ticketCreateErr != nil {
			logger.Error(u.ColorYellow, "CreateTicket activity failed... Or did it?", "Error", ticketCreateErr, u.ColorReset)

			
			delay := 15 // wait a bit for things to settle
			logger.Debug(u.ColorYellow, "ProcessOrder: Sleeping between activity calls - to wait for the ticket system to settle", u.ColorReset)
			logger.Info(u.ColorYellow, "ProcessOrder: workflow.Sleep duration", delay, "seconds", u.ColorReset)
			workflow.Sleep(ctx, time.Duration(delay)*time.Second)
		}		

		logger.Debug(u.ColorYellow, "ProcessOrder: Restoring original Activity Options", u.ColorReset)
		ctx = workflow.WithActivityOptions(ctx, originalActivityOptions)		

		logger.Info(u.ColorBlue, "ProcessOrder: Validating Ticket Exists for", order.OrderID, u.ColorReset)
		err = workflow.ExecuteActivity(ctx, activities.ValidateTicket, order.OrderID, reservation, token).Get(ctx, &ticketFound)
		if err != nil {
			logger.Error(u.ColorRed, "ValidateTicket activity failed.", "Error", err, u.ColorReset)
			return "", err
		}

		if ticketFound {
			logger.Info(u.ColorGreen, "ProcessOrder: Ticket Found for Order: ", order.OrderID, u.ColorReset)
			break
		}
		logger.Info(u.ColorYellow, "ProcessOrder: Ticket Not Found for Order: ", order.OrderID, ", trying again.", u.ColorReset)

	}

	logger.Info(u.ColorGreen, "ProcessOrder completed. Order Details: ", order, " confirmed.", u.ColorReset)

	return "Order Managed", nil

}

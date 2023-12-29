package activities

import (
	"context"
	"idempotence-by-validation/ticket"

	"go.temporal.io/sdk/activity"
)

/* ValidateTicket
 *   This activity attempts to find a ticket to see if it exists
 *
 * Takes a context.Context, an order ID, a reservation, and a token as parameters
 * Returns an reservation and an error if something went bad.
 */
func ValidateTicket(ctx context.Context, orderID string, reservation string, token string) (bool, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("ValidateTicket Activity started")

	ticketExists, err := ticket.ValidateTicket(orderID, reservation, token)
	if err != nil {
		return false, err
	}

	logger.Info("ValidateTicket Activity completed successfully")
	return ticketExists, nil
}

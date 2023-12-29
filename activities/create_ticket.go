package activities

import (
	"context"
	"idempotence-by-validation/ticket"

	u "idempotence-by-validation/utils"

	"go.temporal.io/sdk/activity"
)

/* CreateTicket
 *   This activity retrieves a reservation record from the reservation & ticket system
 *
 * Takes a context.Context, an order ID, a reservation, and a token as parameters
 * Returns an reservation and an error if something went bad.
 */
func CreateTicket(ctx context.Context, orderID string, reservation string, token string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info(u.ColorYellow, "CreateTicket Activity started: ", orderID, u.ColorReset)

	ticket, err := ticket.CreateTicket(orderID, reservation, token)
	if err != nil {
		logger.Info(u.ColorRed, "CreateTicket Activity errored:", err, u.ColorReset)
		return "", err
	}

	logger.Info(u.ColorYellow, "CreateTicket Activity completed successfully with ticket", ticket, u.ColorReset)
	return ticket, nil
}

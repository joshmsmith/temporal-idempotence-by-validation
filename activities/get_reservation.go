package activities

import (
	"context"
	"idempotence-by-validation/ticket"

	"go.temporal.io/sdk/activity"
)

/* GetReservation
 *   This activity retrieves a reservation record from the reservation & ticket system
 *
 * Takes a context.Context, a order ID, and a token as parameters
 * Returns an reservation and an error if something went bad.
 */
func GetReservation(ctx context.Context, orderID string, token string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("GetReservation Activity started")

	

	reservation, err := ticket.GetReservation(orderID, token)
	if err != nil {
		return "", err
	}

	logger.Info("GetReservation Activity completed successfully")
	return reservation, nil
}

package activities

import (
	"context"
	"idempotence-by-validation/ticket"


	"go.temporal.io/sdk/activity"
)

/* GetToken
 *   This activity retrieves a token for calling the ticket reservation system
 *
 * Takes a context.Context and an inventory.Order as parameters.
 * Returns an indicator if there is a duplicate oder and an error.
 */
func GetToken(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("GetToken Activity started")

	token, err := ticket.GetToken()
	if err != nil {
		return "", err
	}

	logger.Info("GetToken Activity completed successfully")
	

	return token, nil
}

package utils

import (
	"log"

	"go.temporal.io/sdk/workflow"
)

/* UpsertSearchAttribute in Temporal Workflow */
func UpsertSearchAttribute(ctx workflow.Context, attribute string, value string) (err error) {

	attributes := map[string]interface{}{
		attribute: value,
	}
	upserterr := workflow.UpsertSearchAttributes(ctx, attributes)
	if upserterr != nil {
		log.Println("Start: Failed to Upsert Search Attributes", upserterr)
	}
	return upserterr
}

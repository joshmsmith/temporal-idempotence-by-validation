package main

import (
	"log"
	"idempotence-by-validation/activities"
	u "idempotence-by-validation/utils"
	"idempotence-by-validation/workflows"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

var TicketOrderManagementTransferTaskQueueName = os.Getenv("TICKET_ORDER_MANAGEMENT_TASK_QUEUE")

// main is the entry point of the program.
// No parameters.
// No return values.
func main() {
	log.Printf("%sGo worker starting.%s", u.ColorGreen, u.ColorReset)

	// Load the Temporal Cloud from env
	clientOptions, err := u.LoadClientOptions(u.SDKMetrics)
	if err != nil {
		log.Fatalf("Failed to load Temporal Cloud environment: %v", err)
	}

	log.Printf("%sGo worker connecting to server.%s", u.ColorGreen, u.ColorReset)
	temporalClient, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer temporalClient.Close()

	temporalWorker := worker.New(temporalClient, TicketOrderManagementTransferTaskQueueName, worker.Options{})

	RegisterWFOptions := workflow.RegisterOptions{
		Name: "CreateTIcket",
	}
	temporalWorker.RegisterWorkflowWithOptions(workflows.ProcessOrder, RegisterWFOptions)

	// activities 
	temporalWorker.RegisterActivity(activities.CreateTicket)
	temporalWorker.RegisterActivity(activities.GetReservation)
	temporalWorker.RegisterActivity(activities.ValidateTicket)
	temporalWorker.RegisterActivity(activities.GetToken)


	// Start listening to the task queue.
	err = temporalWorker.Run(worker.InterruptCh())
	if err != nil {
		log.Fatal("Unable to start worker", err)
	}
}

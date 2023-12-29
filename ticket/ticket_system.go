package ticket

import (
	"encoding/json"

	"errors"
	"fmt"
	"idempotence-by-validation/utils"
	"io"
	"math/rand"
	"os"
	//"github.com/google/uuid"
)

type ticket struct {
	OrderID     string `json:"orderID"`
	TicketID    string `json:"ticketID"`
	PaymentInfo string `json:"paymentInfo"`
}

// ReadJSON reads a JSON file and returns the unmarshalled data as an Inventory struct.
//
// It takes a filename string as a parameter.
// It returns an Inventory struct and an error.
func ReadJSON(filename string) (ticket, error) {
	file, err := os.Open(filename)
	if err != nil {
		return ticket{}, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return ticket{}, err
	}

	var data ticket
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return ticket{}, err
	}

	return data, nil
}

// GetToken returns a token to use to access the reservation & ticketing system
//
// It returns an token and an error.
func GetToken() (string, error) {
	token := "12348972134-1234213"
	// simulate a random token error
	if utils.IsError() {
		return "", errors.New("RANDOM ERROR RETRIEVING TOKEN: TOKEN SERVER FAILED")
	}

	return token, nil
}

// GetReservation returns a reservation to use to access & update the reservation & ticketing system
//
// It returns an token and an error.
func GetReservation(orderID string, token string) (string, error) {
	//	reservation := "1123581321-ASDFQWERTY"
	reservation := fmt.Sprintf("%d-%s", rand.Intn(999999), "ASDFQWERTY")
	
	// simulate a random reservation error
	if utils.IsError() {
		return "", errors.New("RANDOM ERROR RETRIEVING RESERVATION: RESERVATION SYSTEM FAILED")
	}

	return reservation, nil
}

// CreateTicket creates a ticket and returns the ticket ID from the ticketing system
//
// It returns an token and an error.
func CreateTicket(orderID string, reservation string, token string) (string, error) {

	// get PCI payment info for orderID -- retryable
	creditCardInfo, err := GetPCIInfo(reservation, token)
	if err != nil {
		return "", err
	}

	// create request to create ticket with payment info
	ticketOrder := ticket{
		OrderID:     orderID,
		TicketID:    fmt.Sprintf("TICKET-%d", rand.Intn(99999)),
		PaymentInfo: creditCardInfo,
	}

	// simulate a random error before creating a ticket
	if utils.IsErrorPrettyLikely() {
		return "", errors.New("CREATE-TICKET-ERROR")
	}

	filename := fmt.Sprintf("%s%s-reserve-%s.json", os.Getenv("DATABASEPATH"), orderID, reservation)
	// call to issue tickets with PCI info -- non retryable
	err = UpdateJSON(filename, ticketOrder)

	if err != nil {
		return "", err
	}

	// simulate a random error  after creating a ticket
	if utils.IsErrorPrettyLikely() {
		return "", errors.New("CREATE-TICKET-TIMEOUT")
	}

	return ticketOrder.TicketID, nil
}

// ValidateTicket creates a ticket and returns true if the ticket is found in  the ticketing system
//
// It returns a boolean indicator if the ticket exists or not and an error.
func ValidateTicket(orderID string, reservation string, token string) (string, error) {

	filename := fmt.Sprintf("%s%s-reserve-%s.json", os.Getenv("DATABASEPATH"), orderID, reservation)

	data, err := ReadJSON(filename)
	if err != nil {
		fmt.Println("Error reading JSON:", err)
		return "", nil
	}

	// simulate a random error
	if utils.IsError() {
		return "", errors.New("RANDOM ERROR VALIDATING TICKET EXISTS")
	}

	if data.OrderID == orderID {
		return data.TicketID, nil
	}
	return "", nil

}

// GetPCIInfo creates a ticket and returns true if the ticket is found in  the ticketing system
//
// It returns an token and an error.
func GetPCIInfo(reservation string, token string) (string, error) {

	//creditcardInfo := "VISA-5552-1223-2345-7890" // fake that we got it
	if utils.IsError() {
		return "", errors.New("PCI-RETRIEVAL-ERROR")
	}
	//  "VISA-5552-1223-2345-7890"
	creditPCIInfo := fmt.Sprintf("VISA-%d-%d-%d-%d", rand.Intn(9999), rand.Intn(9999), rand.Intn(9999), rand.Intn(9999))

	return creditPCIInfo, nil
}

// SearchOrder searches for an order in the database by its order ID.
//
// Parameters:
// - orderID: a string representing the order ID to search for.
//
// Returns:
// - bool: true if the ticket is found, false otherwise.
func SearchOrder(orderID string) bool {

	data, err := ReadJSON(os.Getenv("DATABASE"))
	if err != nil {
		fmt.Println("Error reading JSON:", err)
		return false
	}

	if data.OrderID == orderID {
		return true
	}
	return false
}

// GetInStock retrieves the quantity of a product that is currently in stock.
//
// It takes a string parameter, productID, which represents the unique identifier of the product.
// The function returns an integer value representing the quantity of the product that is currently in stock.
/* func GetInStock(productID string) (int, error) {

	data, err := ReadJSON(os.Getenv("DATABASE"))
	if err != nil {
		fmt.Println("Error reading JSON:", err)
		return 0, err
	}
	if data.ProductID == productID {
		return data.InStock, err
	}
	return 0, nil
}
*/

// UpdateStock updates the stock of a product for a given order in the database.
//
// Parameters:
// - orderID: the ID of the order.
// - productID: the ID of the product.
// - inStock: the new stock value for the product.
//
// Return type: error.
/*func UpdateStock(orderID string, productID string, inStock int) error {

	data, err := ReadJSON(os.Getenv("DATABASE"))
	if err != nil {
		fmt.Println("Error reading JSON:", err)
		return err
	}
	data.OrderID = orderID
	data.InStock = inStock
	data.ProductID = productID
	UpdateJSON(os.Getenv("DATABASE"), data)
	return nil
}
*/

// UpdateJSON updates the given JSON file with the provided data.
//
// The function returns an error if any error occurs during the file opening,
// encoding, or closing process. Otherwise, it returns nil.
func UpdateJSON(filename string, data ticket) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ") // Pretty print the JSON
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return nil
}

/*
func SupplierOrder(quantity int) error {

	data, err := ReadJSON(os.Getenv("DATABASE"))
	if err != nil {
		fmt.Println("Error reading JSON:", err)
		return err
	}
	GetInStock(data.ProductID)
	// how much should we order?
	toOrder := quantity - data.InStock

	// order that much
	fmt.Println("Ordering from supplier: ", toOrder)
	data.InStock = data.InStock + toOrder
	UpdateJSON(os.Getenv("DATABASE"), data)
	return nil
}
*/

package client_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/qualitorque/terraform-provider-torque/client"
)

// Hello returns a greeting for the named person.
func Hello(name string) (string, error) {
	// If no name was given, return an error with a message.
	if name == "" {
		return name, errors.New("empty name")
	}
	// Create a message using a random format.
	message := "wow"
	//message := fmt.Sprint(randomFormat())
	return message, nil
}

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHelloName(t *testing.T) {
	url := "https://portal.qtorque.io/"
	space := "samples"
	token := "my-awesome-value"
	client, _ := client.NewClient(&url, &space, &token)
	result, _ := client.GetUserDetails("my@email.com")
	fmt.Println(result)
}

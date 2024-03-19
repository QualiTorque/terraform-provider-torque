package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) CreateAccount(AccountEmail string, AccountFirstName string, AccountLastName string, AccountPassword string, AccountCompany string) error {
	fmt.Println(c.HostURL + "api/accounts")

	account := Account{
		Email:     AccountEmail,
		FirstName: AccountFirstName,
		LastName:  AccountLastName,
		Password:  AccountPassword,
		Company:   AccountCompany,
	}

	payload, err := json.Marshal(account)
	if err != nil {
		log.Fatalf("impossible to marshall space: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/accounts/", c.HostURL), bytes.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	_, err = c.doRequest(req, &c.Token)
	if err != nil {
		return err
	}

	return nil
}

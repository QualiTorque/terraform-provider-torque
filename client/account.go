package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) CreateAccount(ParentAccount string, AccountName string, AccountPassword string, AccountCompany string) error {
	fmt.Println(c.HostURL + "api/accounts/" + ParentAccount + "/subaccounts")

	account := Account{
		ParentAccount: ParentAccount,
		AccountName:   AccountName,
		Password:      AccountPassword,
		Company:       AccountCompany,
	}

	payload, err := json.Marshal(account)
	if err != nil {
		log.Fatalf("impossible to marshall Account: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/accounts/%s/subaccounts", c.HostURL, ParentAccount), bytes.NewReader(payload))
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

func (c *Client) RemoveAccount(name string) error {
	fmt.Println(c.HostURL + "api/accounts/" + name)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/accounts/%s", c.HostURL, name), nil)
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

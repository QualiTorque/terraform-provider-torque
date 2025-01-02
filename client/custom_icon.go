package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
)

func (c *Client) UploadCustomIcon(space_name string, file_path string) error {
	file, err := os.Open(file_path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreatePart(
		textproto.MIMEHeader{
			"Content-Disposition": []string{fmt.Sprintf(`form-data; name="file"; filename="%s"`, filepath.Base(file_path))},
			"Content-Type":        []string{"image/svg+xml"},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy file data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/spaces/%s/blueprint_icons", c.HostURL, space_name), body)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	_, err = c.doRequest(req, &c.Token)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetCustomIcons(space_name string, file_path string) ([]TorqueSpaceCustomIcon, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/spaces/%s/blueprint_icons", c.HostURL, space_name), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}
	icons := []TorqueSpaceCustomIcon{}
	err = json.Unmarshal(body, &icons)
	if err != nil {
		return nil, err
	}
	return icons, nil
}

func (c *Client) GetCustomIcon(space_name string, file_path string) (*TorqueSpaceCustomIcon, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/spaces/%s/blueprint_icons", c.HostURL, space_name), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}
	icons := []TorqueSpaceCustomIcon{}
	err = json.Unmarshal(body, &icons)
	if err != nil {
		return nil, err
	}
	icon := TorqueSpaceCustomIcon{}
	fileName := filepath.Base(file_path)
	for _, icon_item := range icons {
		if fileName == icon_item.FileName {
			icon = icon_item
			return &icon, nil
		}
	}

	return nil, fmt.Errorf("icon %s not found", fileName)
}

func (c *Client) DeleteCustomIcon(space_name string, key string) error {
	type deleteCustomIconRequest struct {
		Key string `json:"key"`
	}
	request := deleteCustomIconRequest{
		Key: key,
	}
	payload, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("impossible to marshall custom icon request: %s", err)
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%sapi/spaces/%s/blueprint_icons", c.HostURL, space_name), bytes.NewReader(payload))
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

package goaviatrix

import (
	"fmt"
	"bytes"
	"encoding/json"
	"crypto/tls"
	"strings"
	"net/http"
	"errors"
	"github.com/ajg/form"

)

type LoginResp struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
	CID     string `json:"CID"`
}

type ApiResp struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
}


type Client struct {
	HTTPClient *http.Client
	Username string
	Password string
	CID string
	baseUrl string
}

func (c *Client) Login() error {
	path := c.baseUrl + fmt.Sprintf("?action=login&username=%s&password=%s", c.Username, c.Password)
	resp, err := c.Get(path, nil) 
	if err != nil {
		return err
	}
	var data LoginResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if(!data.Return){
		return errors.New(data.Reason)
	}
	c.CID = data.CID
	return nil
}

func NewClient(username string, password string, controllerIP string, HttpClient *http.Client) (*Client, error) {
	client := &Client{Username: username, Password: password, HTTPClient: HttpClient}
	return client.init(controllerIP)
}

func (c *Client) init(controllerIP string) (*Client, error) {
	c.baseUrl = "https://"+controllerIP+"/v1/api"

	if c.HTTPClient == nil {
	    tr := &http.Transport{
    	    TLSClientConfig: &tls.Config{
    	    	InsecureSkipVerify: true,
    	    },
		}
		c.HTTPClient = &http.Client{Transport: tr}
	}
	if err := c.Login(); err!=nil {
		return nil, err
	}

	return c, nil
}

// Get issues an HTTP GET request.
func (c *Client) Get(path string, i interface{}) (*http.Response, error) {
	return c.Request("GET", path, i)
}

// Post issues an HTTP POST request with the given interface form-encoded.
func (c *Client) Post(path string, i interface{}) (*http.Response, error) {
	return c.Request("POST", path, i)
}


// Put issues an HTTP PUT request with the given interface form-encoded.
func (c *Client) Put(path string, i interface{}) (*http.Response, error) {
	return c.Request("PUT", path, i)
}

// Delete issues an HTTP DELETE request.
func (c *Client) Delete(path string, i interface{}) (*http.Response, error) {
	return c.Request("GET", path, i)
}


// Request makes an HTTP request with the given interface being encoded as
// form data.
func (c *Client) Request(verb, path string, i interface{}) (*http.Response, error) {
	buf := new(bytes.Buffer)
	if(i!=nil) {
		if err := form.NewEncoder(buf).KeepZeros(true).DelimitWith('|').Encode(i); err != nil {
			return nil, err
		}
	}
	body := buf.String()
	req, err := http.NewRequest(verb, path, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c.HTTPClient.Do(req)
}





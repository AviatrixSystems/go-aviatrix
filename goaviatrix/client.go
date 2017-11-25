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
	"github.com/google/go-querystring/query"
	"log"
	"io/ioutil"
	"time"
)

type LoginResp struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	Reason  string `json:"reason"`
	CID     string `json:"CID"`
}

type ApiResp struct {
	Return  bool   `json:"return"`
	Reason  string `json:"reason"`
}

type ApiRequest struct {
	CID string `form:"CID,omitempty" json:"CID" url:"CID"`
	Action string `form:"action,omitempty" json:"action" url:"action"`
}

type Client struct {
	HTTPClient *http.Client
	Username string
	Password string
	CID string
	ControllerIP string
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
	log.Printf("[TRACE] CID is '%s'.", data.CID)
	c.CID = data.CID
	return nil
}

func NewClient(username string, password string, controllerIP string, HttpClient *http.Client) (*Client, error) {
	client := &Client{Username: username, Password: password, HTTPClient: HttpClient, ControllerIP: controllerIP}
	return client.init(controllerIP)
}

func (c *Client) init(controllerIP string) (*Client, error) {
	if len(controllerIP) == 0 {
		return nil, fmt.Errorf("Aviatrix: Client: Controller IP is not set")
	}
	
	c.baseUrl = "https://" + controllerIP + "/v1/api"

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

//
func (c *Client) Do(verb string, req interface{}) (*http.Response, []byte, error) {
	var err error
	var resp *http.Response
	var url string
	var body []byte
	respdata := new(ApiResp)

	// do request
	var loop int = 0
	for {
		url = c.baseUrl
		loop = loop + 1
		if verb == "GET" {
			// prepare query string
			v, _ := query.Values(req)
			url = url + "?" + v.Encode()
			resp, err = c.Request(verb, url, nil)
		} else {
			resp, err = c.Request(verb, url, req)
		}

		// check response for error
		if err != nil {
			return resp, nil, err
		}
		log.Printf("[TRACE] %s %s: %d", verb, url, resp.StatusCode)
		// decode the json response and look for errors to retry
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			body, _ = ioutil.ReadAll(resp.Body)
			if err = json.Unmarshal(body, respdata); err != nil {
				return resp, body, err
			}
			if (!respdata.Return) {
				if respdata.Reason == "CID is invalid or expired." && loop < 2 {
					log.Printf("[TRACE] re-login (expired CID)")
					time.Sleep(500 * time.Millisecond)
					if err = c.Login(); err != nil {
						return resp, body, err
					}
					// loop around again using new CID
				} else {
					return resp, body, errors.New(respdata.Reason)
				}
			} else {
				return resp, body, nil
			}
		} else {
			return resp, body, errors.New("Status code")
		}
	}

	return resp, body, err
}

// Request makes an HTTP request with the given interface being encoded as
// form data.
func (c *Client) Request(verb string, path string, i interface{}) (*http.Response, error) {
	log.Printf("[TRACE] %s %s", verb, path)
	var req *http.Request
	var err error
	if (i != nil) {
		buf := new(bytes.Buffer)
		if err = form.NewEncoder(buf).KeepZeros(true).DelimitWith('|').Encode(i); err != nil {
			return nil, err
		}
		body := buf.String()
		log.Printf("[TRACE] %s %s Body: %s", verb, path, body)
		reader := strings.NewReader(body)
		req, err = http.NewRequest(verb, path, reader)
		if err == nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}

	} else {
		req, err = http.NewRequest(verb, path, nil)
	}

	if err != nil {
		return nil, err
	}
	return c.HTTPClient.Do(req)
}


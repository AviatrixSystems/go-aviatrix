package goaviatrix

import (
	"fmt"
	"bytes"
	"encoding/json"
	"crypto/tls"
	//"io"
	"io/ioutil"
	"strings"
	"net/http"
	"net/url"
	"github.com/ajg/form"
	//"github.com/hashicorp/go-cleanhttp"
	//"github.com/mitchellh/mapstructure"

)

type APIResp struct {
	Return  bool   `json:"return"`
	Results string `json:"results"`
	CID     string `json:"CID"`
	Reason  string `json:"reason"`
}

type Client struct {
	HTTPClient *http.Client
	Username string
	Password string
	Cid string
	url *url.URL
}

func (c *Client) Login() error {
	fmt.Println(3)
	path := fmt.Sprintf("action=login&username=%s&password=%s", c.Username, c.Password)
	fmt.Println(path)
	resp, err := c.Get(path, nil)
	if err != nil {
		fmt.Println(err)
		fmt.Println(4)
		return err
	}
	fmt.Println(resp)
    //c.Cid = 1
	//b:= *APIResp{}
	var data APIResp
	//err = json.NewDecoder(resp.Body).Decode(b)
	//body, err := ioutil.ReadAll(resp.Body)
	//err = json.Unmarshal(body, &data)
	body, err := ioutil.ReadAll(resp.Body) // read the body of the request                                                  
	if err != nil {
		fmt.Println(5)
		return err
	}
	if err := json.Unmarshal(body, &data); err != nil { // unmarshall body contents as a type Candidate                                          
		fmt.Println(5)
		return err
	}
	c.Cid = data.CID
	fmt.Println("CID======", c.Cid)
	fmt.Println("reason======", data.Reason)
	//return b, nil

    return nil
}

func NewClient(username string, password string, baseUrl string) (*Client, error) {
	fmt.Println(1)
	client := &Client{Username: username, Password: password}
	return client.init(baseUrl)
}

func (c *Client) init(baseUrl string) (*Client, error) {
	fmt.Println(2)
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	c.url = u
	//fmt.Println("url====",u)
	if c.HTTPClient == nil {
	    tr := &http.Transport{
    	    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		c.HTTPClient = &http.Client{Transport: tr}
    	
		//c.HTTPClient = cleanhttp.DefaultClient()
		//tls_config := &tls.Config{InsecureSkipVerify: true}
		//transport := cleanhttp.DefaultTransport()
		//transport.TLSClientConfig = tls_config
		//c.HTTPClient.Transport = transport
	}
	c.Login()
	return c, nil
}
// Get issues an HTTP GET request.
func (c *Client) Get(p string, ro *RequestOptions) (*http.Response, error) {
	fmt.Println("path: ",p)
	//return c.Request("GET", p, ro)
	return c.HTTPClient.Get("https://"+c.url.Host+c.url.Path+"?"+p)
}

// Head issues an HTTP HEAD request.
func (c *Client) Head(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request("HEAD", p, ro)
}

// Post issues an HTTP POST request.
func (c *Client) Post(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request("POST", p, ro)
}

// PostForm issues an HTTP POST request with the given interface form-encoded.
func (c *Client) PostForm(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestForm("POST", p, i, ro)
}

// Put issues an HTTP PUT request.
func (c *Client) Put(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request("PUT", p, ro)
}

// PutForm issues an HTTP PUT request with the given interface form-encoded.
func (c *Client) PutForm(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestForm("PUT", p, i, ro)
}

// Delete issues an HTTP DELETE request.
func (c *Client) Delete(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request("DELETE", p, ro)
}


func (c *Client) Request(verb, p string, ro *RequestOptions) (*http.Response, error) {
	req, err := c.RawRequest(verb, p, ro)
	if err != nil {
		return nil, err
	}

	resp, err := checkResp(c.HTTPClient.Do(req))
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RequestForm makes an HTTP request with the given interface being encoded as
// form data.
func (c *Client) RequestForm(verb, p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	if ro == nil {
		ro = new(RequestOptions)
	}

	if ro.Headers == nil {
		ro.Headers = make(map[string]string)
	}
	ro.Headers["Content-Type"] = "application/x-www-form-urlencoded"

	buf := new(bytes.Buffer)
	if err := form.NewEncoder(buf).KeepZeros(true).DelimitWith('|').Encode(i); err != nil {
		return nil, err
	}
	body := buf.String()

	ro.Body = strings.NewReader(body)
	ro.BodyLength = int64(len(body))

	return c.Request(verb, p, ro)
}

func checkResp(resp *http.Response, err error) (*http.Response, error) {
	// If the err is already there, there was an error higher up the chain, so
	// just return that.
	if err != nil {
		return resp, err
	}

	switch resp.StatusCode {
	case 200, 201, 202, 204, 205, 206:
		return resp, nil
	default:
		return resp, nil
	}
}

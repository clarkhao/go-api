package request

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	re "restapi/utils/error"
)

type Items map[string]interface{}

// String returns the every properties of Test.
func (item Items) String() string {
	var strs string = ""
	for key, value := range item {
		strs += fmt.Sprintf("\n\t%s: %v", key, value)
	}
	strs += "\n"
	return strs
}

type Request interface {
	GetRequest(url string) ([]Items, error)
	PostRequest(url string, data map[string]string) ([]Items, error)
}

type ClientBase struct {
	baseUrl string
	client  *http.Client
	req     *http.Request
}

var Client *ClientBase = &ClientBase{
	baseUrl: "",
	client: &http.Client{
		Timeout: time.Duration(time.Duration(1) * time.Second),
	},
	req: nil,
}

func (client *ClientBase) NewClient(baseUrl string) *ClientBase {
	client.baseUrl = baseUrl
	return client
}

func (client *ClientBase) SetRequest(method string, headers map[string]string, data map[string]string, path string) *ClientBase {
	//query parameters or body
	formData := url.Values{}
	for k, v := range data {
		formData.Set(k, v)
	}
	queryString := formData.Encode()
	// for url path
	var str strings.Builder
	str.WriteString(client.baseUrl)
	str.WriteString(path)
	switch method {
	case "GET":
		str.WriteString("?")
		str.WriteString(queryString)
		urlstr := str.String()
		client.req, _ = http.NewRequest("GET", urlstr, nil)
	case "POST":
		urlstr := str.String()
		body := strings.NewReader(queryString)
		client.req, _ = http.NewRequest("POST", urlstr, body)
	default:
		str.WriteString("?")
		str.WriteString(queryString)
		urlstr := str.String()
		client.req, _ = http.NewRequest("GET", urlstr, nil)
	}
	for k, v := range headers {
		client.req.Header.Add(k, v)
	}
	return client
}

func (client ClientBase) String() string {
	return fmt.Sprintf("Method is %s and Access URL is %v", client.req.Method, client.req.URL)
}

func (client *ClientBase) GetRequest() ([]Items, error) {
	var result []Items
	res, err := http.DefaultClient.Do(client.req)
	if err != nil || res.StatusCode >= 400 {
		return nil, re.RequestError{
			Code:  res.StatusCode,
			Msg:   fmt.Sprintf("when fetching data from %s in request/GetRequest", client.req.URL),
			Cause: err,
		}
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	err = json.Unmarshal([]byte(string(body)), &result)
	if err != nil {
		return nil, re.RequestError{
			Code:  500,
			Msg:   fmt.Sprintln("failed to extract data from response body in request/GetRequest"),
			Cause: err,
		}
	}
	return result, nil
}

func (client *ClientBase) PostRequest() (Items, error) {
	var result Items
	res, err := http.DefaultClient.Do(client.req)
	if err != nil || res.StatusCode >= 400 {
		return nil, re.RequestError{
			Code:  res.StatusCode,
			Msg:   fmt.Sprintf("when posting data to %s in request/PostRequest", client.req.URL),
			Cause: err,
		}
	}
	defer res.Body.Close()
	resBody, _ := io.ReadAll(res.Body)
	log.Println(string(resBody))
	err = json.Unmarshal([]byte(string(resBody)), &result)
	if err != nil {
		return nil, re.RequestError{
			Code:  500,
			Msg:   fmt.Sprintln("failed to extract data from response body in request/PostRequest"),
			Cause: err,
		}
	}
	return result, nil
}

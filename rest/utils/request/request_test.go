package request_test

import (
	"encoding/json"
	"restapi/utils/logger"
	"restapi/utils/request"
	"testing"
	"time"

	re "restapi/utils/error"
)

type ClientUser struct {
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	ID        string    `json:"id"`
}

func ExtractClientUser(items []request.Items) ([]ClientUser, error) {
	var result = make([]ClientUser, len(items))
	jsonData, err := json.Marshal(items)
	if err != nil {
		return result, re.RequestError{
			Code:  500,
			Msg:   "from ExtractClientUser func when marshaling",
			Cause: err,
		}
	}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return result, re.RequestError{
			Code:  500,
			Msg:   "from ExtractClientUser func when unmarshaling",
			Cause: err,
		}
	}
	return result, nil
}

func TestGetRequest(t *testing.T) {
	cases := []struct {
		name         string
		base         string
		path         string
		method       string
		headers      map[string]string
		data         map[string]string
		expectedCode int
	}{
		{
			name:         "GETList",
			base:         "https://641b10fb9b82ded29d494d1c.mockapi.io/api/posts",
			path:         "/user",
			method:       "GET",
			headers:      map[string]string{},
			data:         map[string]string{},
			expectedCode: 200,
		},
		{
			name:         "GETListWithBadPath",
			base:         "https://641b10fb9b82ded29d494d1c.mockapi.io/api/posts",
			path:         "/users",
			method:       "GET",
			headers:      map[string]string{},
			data:         map[string]string{},
			expectedCode: 404,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			items, err := request.Client.NewClient(tc.base).SetRequest(tc.method, tc.headers, tc.data, tc.path).GetRequest()
			if err != nil {
				e, ok := err.(re.RequestError)
				if ok {
					logger.Log.Error(e)
				}
				if tc.expectedCode != e.Code {
					t.Errorf("failed to get user list")
				}
			} else {
				_, err := ExtractClientUser(items)
				if err != nil {
					t.Errorf("failed to extract user data")
				}
			}
		})
	}
}

func TestPostRequest(t *testing.T) {
	cases := []struct {
		name         string
		base         string
		path         string
		method       string
		headers      map[string]string
		data         map[string]string
		expectedCode int
	}{
		{
			name:         "POSTList",
			base:         "https://641b10fb9b82ded29d494d1c.mockapi.io/api/posts",
			path:         "/user",
			method:       "POST",
			headers:      map[string]string{},
			expectedCode: 201,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			items, err := request.Client.NewClient(tc.base).SetRequest(tc.method, tc.headers, tc.data, tc.path).PostRequest()
			if err != nil {
				e, ok := err.(re.RequestError)
				if ok {
					logger.Log.Error(e)
				}
				if tc.expectedCode != e.Code {
					t.Errorf("failed to get user list")
				}
			} else {
				logger.Log.Info(items)
			}
		})
	}
}

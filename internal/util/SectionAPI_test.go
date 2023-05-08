package util_test

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/dpc-sdp/bay-section-ip-controller/internal/util"
	"github.com/stretchr/testify/assert"
)

type MockRoundTripper struct {
	Response *http.Response
	Err      error
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}

func TestGetApplications(t *testing.T) {
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body: ioutil.NopCloser(strings.NewReader(`[
			{
				"id": 1,
				"href": "/api/v1/account/1/application/1",
				"application_name": "www.test.com"
			}]`)),
	}
	mockTransport := &MockRoundTripper{
		Response: response,
		Err:      nil,
	}

	client := &http.Client{Transport: mockTransport}

	section := util.SectionAPI{Client: client}
	section.GetApplications()

	assert.Equal(t, len(section.Applications), 1)
	assert.Equal(t, section.Applications[0].ApplicationName, "www.test.com")
	assert.Equal(t, section.Applications[0].Id, 1)
}

func TestGetEnvironments(t *testing.T) {
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body: ioutil.NopCloser(strings.NewReader(`[
			{
				"id": 1,
				"environment_name": "Production",
				"href": "/api/v1/account/1/application/1/environment/Production",
				"is_hosted": true,
				"rum_token": "xxxx",
				"domains": []
			},
			{
				"id": 2,
				"environment_name": "Development",
				"href": "/api/v1/account/1/application/1/environment/Development",
				"is_hosted": false,
				"domains": []
			},
			{
				"id": 3,
				"environment_name": "Develop",
				"href": "/api/v1/account/3/application/3/environment/Develop",
				"is_hosted": true,
				"rum_token": "xxxx",
				"domains": []
			},
			{
				"id": 4,
				"environment_name": "Master",
				"href": "/api/v1/account/1/application/1/environment/Master",
				"is_hosted": true,
				"rum_token": "xxxx",
				"domains": []
			}
		]`)),
	}
	mockTransport := &MockRoundTripper{
		Response: response,
		Err:      nil,
	}
	client := &http.Client{Transport: mockTransport}

	section := util.SectionAPI{Client: client}
	app := util.SectionApp{
		Id:              1,
		ApplicationName: "www.test.com",
	}

	section.GetApplicationEnvironments(&app)

	assert.Equal(t, len(app.Environments), 4)

	_, ok := app.Environments["Production"]
	assert.True(t, ok)
}

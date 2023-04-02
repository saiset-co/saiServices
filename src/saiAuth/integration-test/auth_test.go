package integration_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webmakom-com/saiAuth/auth"
	"github.com/webmakom-com/saiAuth/config"
	"github.com/webmakom-com/saiAuth/models"
	"github.com/webmakom-com/saiAuth/utils/saiStorageUtil"
	"go.uber.org/zap"
)

// config values for saiAuth
const (
	host          = "localhost:8800"
	token         = "12345"
	responseFalse = "false"
	responseTrue  = "true"
	// register
	registerPath              = "http://" + host + "/register"
	registerKey               = "user"
	registerPassword          = "123456"
	registerResponseSucessStr = "{\"Status\":\"Ok\"}"
	registerResponseFalse     = "false"

	// login
	loginPath = "http://" + host + "/login"

	// access
	accessPath   = "http://" + host + "/access"
	accessCol    = "tokens"
	accessMethod = "get"
)

var (
	cfg     config.Configuration
	logger  *zap.Logger
	err     error
	storage saiStorageUtil.Database
	manager auth.Manager

	accessToken string // for saving access token from logn endpoint and use this variable in access endpoint
)

type registerRequest struct {
	Key      string `json:"key"`
	Password string `json:"password"`
}

type accessRequest struct {
	Collection string `json:"collection"`
	Method     string `json:"method"`
}

// test register endpoint
func TestRegister(t *testing.T) {
	r := registerRequest{
		Key:      registerKey,
		Password: registerPassword,
	}
	b, err := json.Marshal(&r)
	if err != nil {
		t.Errorf("TestRegister - marshal test body : %s", err.Error())
		t.FailNow()
	}

	req, err := http.NewRequest("GET", registerPath, bytes.NewBuffer(b))
	if err != nil {
		t.Errorf("TestRegister - http.NewRequest : %s", err.Error())
		t.FailNow()
	}

	req.Header.Add("Token", token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("TestRegister - http.Do : %s", err.Error())
		t.FailNow()
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("TestRegister - ioutil.ReadAll : %s", err.Error())
		t.FailNow()
	}

	if string(body) == "false" {
		t.Errorf("TestRegister - error")
		t.FailNow()
	}

	assert := assert.New(t)

	s, err := strconv.Unquote(string(body))
	if err != nil {
		t.Errorf("TestRegister - unqoute : %s", err.Error())
		t.FailNow()
	}

	assert.Equal(registerResponseSucessStr, s, "response status")
	assert.Equal(200, resp.StatusCode, "status code")

}

// test register error if user already exists
func TestRegisterAlreadyExists(t *testing.T) {

	type registerRequest struct {
		Key      string `json:"key"`
		Password string `json:"password"`
	}
	r := registerRequest{
		Key:      registerKey,
		Password: registerPassword,
	}
	b, err := json.Marshal(&r)
	if err != nil {
		t.Errorf("TestRegisterAlreadyExists - marshal test body : %s", err.Error())
		t.FailNow()
	}

	req, err := http.NewRequest("GET", registerPath, bytes.NewBuffer(b))
	if err != nil {
		t.Errorf("TestRegisterAlreadyExists - http.NewRequest : %s", err.Error())
		t.FailNow()
	}

	req.Header.Add("Token", token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("TestRegisterAlreadyExists - http.Do : %s", err.Error())
		t.FailNow()
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("TestRegisterAlreadyExists - ioutil.ReadAll : %s", err.Error())
		t.FailNow()
	}

	assert := assert.New(t)

	assert.Equal(responseFalse, string(body), "response status")
	assert.Equal(200, resp.StatusCode, "status code")

}

// test login endpoint
func TestLogin(t *testing.T) {
	r := registerRequest{
		Key:      registerKey,
		Password: registerPassword,
	}
	b, err := json.Marshal(&r)
	if err != nil {
		t.Errorf("TestLogin - marshal test body : %s", err.Error())
		t.FailNow()
	}

	req, err := http.NewRequest("GET", loginPath, bytes.NewBuffer(b))
	if err != nil {
		t.Errorf("TestLogin - http.NewRequest : %s", err.Error())
		t.FailNow()
	}
	req.Header.Add("Token", token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("TestLogin - http.Do : %s", err.Error())
		t.FailNow()
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("TestLogin - ioutil.ReadAll : %s", err.Error())
		t.FailNow()
	}

	if string(body) == "false" {
		t.Errorf("TestLogin - error")
		t.FailNow()
	}

	loginResp := models.LoginResponse{}

	err = json.Unmarshal(body, &loginResp)
	if err != nil {
		t.Errorf("TestLogin - unmarshal : %s", err.Error())
		t.FailNow()
	}

	accessToken = loginResp.AccessToken.Name

	assert := assert.New(t)

	assert.NotEmpty(loginResp.AccessToken.Name)
	assert.NotEmpty(loginResp.AccessToken.Expiration)
	assert.NotEmpty(loginResp.RefreshToken.Name)
	assert.NotEmpty(loginResp.RefreshToken.Expiration)

	assert.Equal(200, resp.StatusCode, "status code")

}

// test wrong/empty login in login endpoint
func TestWrongLogin(t *testing.T) {
	r := registerRequest{
		Key: registerKey, // empty password
	}
	b, err := json.Marshal(&r)
	if err != nil {
		t.Errorf("TestWrongLogin - marshal test body : %s", err.Error())
		t.FailNow()
	}

	req, err := http.NewRequest("GET", loginPath, bytes.NewBuffer(b))
	if err != nil {
		t.Errorf("TestWrongLogin - http.NewRequest : %s", err.Error())
		t.FailNow()
	}
	req.Header.Add("Token", token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("TestWrongLogin - http.Do : %s", err.Error())
		t.FailNow()
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("TestWrongLogin - ioutil.ReadAll : %s", err.Error())
		t.FailNow()
	}

	assert := assert.New(t)
	assert.Equal(responseFalse, string(body), "response status")
	assert.Equal(200, resp.StatusCode, "status code")

}

// test access endpoint
func TestAccess(t *testing.T) {
	r := accessRequest{
		Collection: accessCol,
		Method:     accessMethod,
	}
	b, err := json.Marshal(&r)
	if err != nil {
		t.Errorf("TestAccess - marshal test body : %s", err.Error())
		t.FailNow()
	}

	req, err := http.NewRequest("GET", accessPath, bytes.NewBuffer(b))
	if err != nil {
		t.Errorf("TestAccess - http.NewRequest : %s", err.Error())
		t.FailNow()
	}
	req.Header.Add("Token", accessToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("TestAccess - http.Do : %s", err.Error())
		t.FailNow()
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("TestAccess - ioutil.ReadAll : %s", err.Error())
		t.FailNow()
	}

	if string(body) == "false" {
		t.Errorf("TestAccess - error")
		t.FailNow()
	}

	assert := assert.New(t)

	assert.Equal(responseTrue, string(body), "response status")
	assert.Equal(200, resp.StatusCode, "status code")

}

// test access endpoint
func TestAccessEmptyToken(t *testing.T) {
	r := accessRequest{
		Collection: accessCol,
		Method:     accessMethod,
	}
	b, err := json.Marshal(&r)
	if err != nil {
		t.Errorf("TestAccess - marshal test body : %s", err.Error())
		t.FailNow()
	}

	req, err := http.NewRequest("GET", accessPath, bytes.NewBuffer(b))
	if err != nil {
		t.Errorf("TestAccess - http.NewRequest : %s", err.Error())
		t.FailNow()
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("TestAccess - http.Do : %s", err.Error())
		t.FailNow()
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("TestAccess - ioutil.ReadAll : %s", err.Error())
		t.FailNow()
	}

	assert := assert.New(t)

	assert.Equal(responseFalse, string(body), "response status")
	assert.Equal(200, resp.StatusCode, "status code")

}

// func (s Server) Auth(h HandlerRequest) interface{} {
// 	return s.AuthManager.Auth(h.getInterface(), h.Token)
// }

// func (s Server) Password(h HandlerRequest) interface{} {
// 	return s.AuthManager.Password(h.getInterface(), h.Token)
// }

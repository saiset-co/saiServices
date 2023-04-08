package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

var (
	registerKey = 0
	loginKey    = 0
)

type Targeter struct {
	name string
	T    vegeta.Targeter
}

func main() {
	config := Load()

	rate := vegeta.Rate{Freq: 10, Per: time.Second}
	duration := 1 * time.Minute

	registerTargeter := RegisterTargeter(&config)
	testAndReport("register", registerTargeter, rate, duration)
	loginTargeter := LoginTargeter(&config)
	testAndReport("login", loginTargeter, rate, duration)

	accessToken, err := getAccessToken(&config)
	if err != nil {
		log.Fatal(err)
	}
	accessTargeter := AccessTargeter(accessToken, &config)
	testAndReport("access", accessTargeter, rate, duration)
}

func testAndReport(name string, targeter vegeta.Targeter, rate vegeta.ConstantPacer, duration time.Duration) {
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Whatever name") {
		metrics.Add(res)
	}
	metrics.Close()

	reporter := vegeta.NewTextReporter(&metrics)
	fmt.Printf("\n")
	fmt.Printf("Name : %s\nSettings :\n Rate : %+v\n Duration : %+v\n", name, rate, duration)
	reporter(os.Stdout)
}

// register endpoint targeter
func RegisterTargeter(c *Configuration) vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		tgt.Method = "GET"

		tgt.URL = fmt.Sprintf("http://%s:%s/register", c.HttpServer.Host, c.HttpServer.Port)

		type req struct {
			Key      string `json:"key"`
			Password string `json:"password"`
		}

		registerKey++
		r := req{
			Key:      strconv.Itoa(registerKey),
			Password: "123456",
		}

		b, _ := json.Marshal(r)

		tgt.Body = b

		header := http.Header{}
		header.Add("Token", "12345")
		header.Add("Content-Type", "application/json")
		tgt.Header = header

		return nil
	}
}

// login endpoint targeter
func LoginTargeter(c *Configuration) vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		tgt.Method = "GET"

		tgt.URL = fmt.Sprintf("http://%s:%s/login", c.HttpServer.Host, c.HttpServer.Port)

		type req struct {
			Key      string `json:"key"`
			Password string `json:"password"`
		}
		loginKey++

		r := req{
			Key:      strconv.Itoa(loginKey),
			Password: "123456",
		}

		b, _ := json.Marshal(r)

		tgt.Body = b

		header := http.Header{}
		header.Add("Token", "12345")
		header.Add("Content-Type", "application/json")
		tgt.Header = header

		return nil
	}
}

// access endpoint targeter
func AccessTargeter(accessToken string, c *Configuration) vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		tgt.Method = "GET"

		tgt.URL = fmt.Sprintf("http://%s:%s/access", c.HttpServer.Host, c.HttpServer.Port)

		type req struct {
			Collection string `json:"collection"`
			Method     string `json:"method"`
		}

		r := req{
			Collection: "tokens",
			Method:     "get",
		}

		b, _ := json.Marshal(r)

		tgt.Body = b

		header := http.Header{}
		header.Add("Token", c.Token)
		header.Add("Content-Type", "application/json")
		tgt.Header = header

		return nil
	}
}

// get access tokens for tests
func getAccessToken(c *Configuration) (string, error) {
	type req struct {
		Key      string `json:"key"`
		Password string `json:"password"`
	}

	r := req{
		Key:      "1",
		Password: "123456",
	}

	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}

	resp, err := http.Post("http://localhost:8800/login", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	loginResp := LoginResponse{}

	err = json.Unmarshal(b, &loginResp)
	if err != nil {
		return "", err
	}
	return loginResp.AccessToken.Name, nil
}

package main

import (
	"fmt"

	"github.com/tkanos/gonfig"
)

type Permission struct {
	Exists   bool
	Methods  map[string]bool
	Required map[string]any
}

type Configuration struct {
	HttpServer struct {
		Host string
		Port string
	}
	HttpsServer struct {
		Host string
		Port string
	}
	Address struct {
		Url string
	}
	SocketServer struct {
		Host string
		Port string
	}
	Salt    string
	Token   string
	Storage struct {
		Token string
		Url   string
		Auth  struct {
			Email    string
			Password string
		}
	}
	Operations []string
	StartBlock int
	WebSocket  struct {
		Token string
		Url   string
	}
	Contract struct {
		Address string
		ABI     string
	}
	Geth  []string
	Sleep int
	Roles map[string]struct {
		Exists      bool
		Permissions map[string]Permission
	}
	AccessTokenExp  int64
	RefreshTokenExp int64
	DefaultRole     string
	EnableProfiling bool
	ProfilingPort   int64
}

func Load() Configuration {
	var config Configuration
	err := gonfig.GetConf("./../saiAuth/config.json", &config)

	if err != nil {
		fmt.Println("Configuration problem:", err)
		panic(err)
	}

	return config
}

type LoginResponse struct {
	*AccessToken  `json:"at"`
	*RefreshToken `json:"rt"`
	User          map[string]interface{} `json:"user,omitempty"`
}

// Access token representation for unmarshal
type AccessToken struct {
	ID          string                  `json:"_id,omitempty"`
	Type        string                  `json:"type,omitempty"`
	Name        string                  `json:"name"`
	Expiration  int64                   `json:"expiration"`
	InternalID  string                  `json:"internal_id,omitempty"`
	User        map[string]interface{}  `json:"user,omitempty"`
	Permissions []map[string]Permission `json:"permissions,omitempty"`
}

// User representation inside access token
type User struct {
	ID         string `json:"_id,omitempty"`
	InternalID string `json:"internal_id,omitempty"`
}

// Refresh token representation
type RefreshToken struct {
	ID          string                  `json:"_id,omitempty"`
	Type        string                  `json:"type,omitempty"`
	Name        string                  `json:"name"`
	Expiration  int64                   `json:"expiration"`
	InternalID  string                  `json:"internal_id,omitempty"`
	AccessToken *AccessToken            `json:"access_token,omitempty"`
	Permissions []map[string]Permission `json:"permissions,omitempty"`
}

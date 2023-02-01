package internal

import (
	"github.com/saiset-co/saiService"
)

func (is InternalService) Handlers() saiService.Handler {
	return saiService.Handler{
		"move": saiService.HandlerElement{
			Name:        "move",
			Description: "Move data to another storage",
			Function: func(data interface{}) (interface{}, error) {
				return is.Move()
			},
		},
	}
}

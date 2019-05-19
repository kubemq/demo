package main

import (
	"github.com/labstack/echo"
)

type CreateAccount struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (s *Server) register(c echo.Context) error {
	r := NewResponse(c)
	acc := &CreateAccount{}
	err := c.Bind(acc)
	if err != nil {
		return r.SetError(err).Send()
	}

	_, err = s.kube.SendCommand(c.Request().Context(), s.config.UsersChannel, "register", acc)
	if err != nil {
		return r.SetError(err).Send()
	}
	return r.Send()

}

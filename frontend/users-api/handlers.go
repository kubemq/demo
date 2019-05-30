package main

import (
	"errors"
	"github.com/labstack/echo"
)

func (s *Server) register(c echo.Context) error {
	r := NewResponse(c, s.kube, "register", "query")
	acc := &NewUser{}
	err := c.Bind(acc)
	if err != nil {
		return r.SetError(err).Send()
	}
	r.SetRequestBody(acc)
	resp, err := s.kube.SendQuery(c.Request().Context(), s.cfg.UsersChannel, "register", acc)
	if err != nil {
		return r.SetError(err).Send()
	}
	if !resp.Executed {
		return r.SetError(errors.New(resp.Error)).Send()
	}
	user, err := getUser(resp.Body)
	if err != nil {
		return r.SetError(err).Send()
	}

	r.SetResponseBody(user)
	return r.Send()
}

func (s *Server) login(c echo.Context) error {
	r := NewResponse(c, s.kube, "login", "query")
	lr := &LoginRequest{}
	err := c.Bind(lr)
	if err != nil {
		return r.SetError(err).Send()
	}
	r.SetRequestBody(lr)
	resp, err := s.kube.SendQuery(c.Request().Context(), s.cfg.UsersChannel, "login", lr)
	if err != nil {
		return r.SetError(err).Send()
	}
	if !resp.Executed {
		return r.SetError(errors.New(resp.Error)).Send()
	}
	loginResp, err := getLoginResponse(resp.Body)
	if err != nil {
		return r.SetError(err).Send()
	}

	r.SetResponseBody(loginResp)
	return r.Send()
}

func (s *Server) verify(c echo.Context) error {
	r := NewResponse(c, s.kube, "verify", "command")
	vr := &VerifyRegistration{}
	err := c.Bind(vr)
	if err != nil {
		return r.SetError(err).Send()
	}
	r.SetRequestBody(vr)
	resp, err := s.kube.SendCommand(c.Request().Context(), s.cfg.UsersChannel, "verify_registration", vr)
	if err != nil {
		return r.SetError(err).Send()
	}
	if !resp.Executed {
		return r.SetError(errors.New(resp.Error)).Send()
	}
	return r.Send()
}

func (s *Server) logout(c echo.Context) error {
	r := NewResponse(c, s.kube, "logout", "command")
	lo := &LogoutRequest{}
	err := c.Bind(lo)
	if err != nil {
		return r.SetError(err).Send()
	}
	r.SetRequestBody(lo)
	resp, err := s.kube.SendCommand(c.Request().Context(), s.cfg.UsersChannel, "logout", lo)
	if err != nil {
		return r.SetError(err).Send()
	}
	if !resp.Executed {
		return r.SetError(errors.New(resp.Error)).Send()
	}
	return r.Send()
}

//
//func (s *Server) passwordReset(c echo.Context) error {
//	r := NewResponse(c)
//	pr := &PasswordResetRequest{}
//	err := c.Bind(pr)
//	if err != nil {
//		return r.SetError(err).Send()
//	}
//	r.SetRequestBody(pr)
//	resp, err := s.kube.SendCommand(c.Request().Context(), s.cfg.UsersChannel, "password_reset_request", pr)
//	if err != nil {
//		return r.SetError(err).Send()
//	}
//	if !resp.Executed {
//		return r.SetError(errors.New(resp.Error)).Send()
//	}
//	return r.Send()
//}

//
//func (s *Server) passwordChange(c echo.Context) error {
//	r := NewResponse(c)
//	pc := &PasswordChangeRequest{}
//	err := c.Bind(pc)
//	if err != nil {
//		return r.SetError(err).Send()
//	}
//	r.SetRequestBody(pc)
//	resp, err := s.kube.SendCommand(c.Request().Context(), s.cfg.UsersChannel, "password_change_request", pc)
//	if err != nil {
//		return r.SetError(err).Send()
//	}
//	if !resp.Executed {
//		return r.SetError(errors.New(resp.Error)).Send()
//	}
//	return r.Send()
//}
//
//func (s *Server) lock(c echo.Context) error {
//	r := NewResponse(c)
//	lc := &LockRequest{}
//	err := c.Bind(lc)
//	if err != nil {
//		return r.SetError(err).Send()
//	}
//	r.SetRequestBody(lc)
//	resp, err := s.kube.SendCommand(c.Request().Context(), s.cfg.UsersChannel, "lock", lc)
//	if err != nil {
//		return r.SetError(err).Send()
//	}
//	if !resp.Executed {
//		return r.SetError(errors.New(resp.Error)).Send()
//	}
//	return r.Send()
//}
//
//func (s *Server) unlock(c echo.Context) error {
//	r := NewResponse(c)
//	lc := &UnlockRequest{}
//	err := c.Bind(lc)
//	if err != nil {
//		return r.SetError(err).Send()
//	}
//	r.SetRequestBody(lc)
//	resp, err := s.kube.SendCommand(c.Request().Context(), s.cfg.UsersChannel, "unlock", lc)
//	if err != nil {
//		return r.SetError(err).Send()
//	}
//	if !resp.Executed {
//		return r.SetError(errors.New(resp.Error)).Send()
//	}
//	return r.Send()
//}

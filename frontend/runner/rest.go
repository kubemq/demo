package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty"
)

type Response struct {
	IsError bool            `json:"is_error"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func Register(url string, newUser *NewUser) (*User, error) {

	resp := &Response{}
	uri := fmt.Sprintf("%s/register", url)
	r, err := resty.R().SetBody(newUser).SetResult(resp).SetError(resp).Post(uri)
	if err != nil {
		return nil, err
	}
	if !r.IsSuccess() || resp.IsError {
		return nil, errors.New(resp.Message)
	}
	user := &User{}
	err = json.Unmarshal(resp.Data, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func SendVerify(url, name, token string) error {

	resp := &Response{}
	uri := fmt.Sprintf("%s/verify", url)
	vr := &VerifyRegistration{
		Name:  name,
		Token: token,
	}
	r, err := resty.R().SetBody(vr).SetResult(resp).SetError(resp).Post(uri)
	if err != nil {
		return err
	}
	if !r.IsSuccess() || resp.IsError {
		return errors.New(resp.Message)
	}

	return nil
}

func SendLogin(url, name, password string) (*LoginResponse, error) {

	resp := &Response{}
	uri := fmt.Sprintf("%s/login", url)
	lo := &LoginRequest{
		Name:     name,
		Password: password,
	}
	r, err := resty.R().SetBody(lo).SetResult(resp).SetError(resp).Post(uri)
	if err != nil {
		return nil, err
	}
	if !r.IsSuccess() || resp.IsError {
		return nil, errors.New(resp.Message)
	}
	lr := &LoginResponse{}

	err = json.Unmarshal(resp.Data, lr)
	if err != nil {
		return nil, err
	}
	return lr, nil
}

func SendLogout(url, userId string) error {

	resp := &Response{}
	uri := fmt.Sprintf("%s/logout", url)
	lo := &LogoutRequest{
		UserId: userId,
	}
	r, err := resty.R().SetBody(lo).SetResult(resp).SetError(resp).Post(uri)
	if err != nil {
		return err
	}
	if !r.IsSuccess() || resp.IsError {
		return errors.New(resp.Message)
	}

	return nil
}

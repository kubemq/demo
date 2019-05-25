package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

var newUserCounter = 0

func createNewUser() *NewUser {
	newUserCounter++

	return &NewUser{
		Name:     fmt.Sprintf("newUser%d@some-email.com", newUserCounter),
		Password: fmt.Sprintf("password_%d", newUserCounter),
		Email:    fmt.Sprintf("newUser%d@some-email.com", newUserCounter),
	}

}
func script(url string) {
	newUser := createNewUser()
	user, err := Register(url, newUser)
	if err != nil {
		log.Println(fmt.Sprintf("Register New User:\n%sError: %s\n", PrettyJson(newUser), err.Error()))
		return
	}
	log.Println(fmt.Sprintf("Register New User:\n%sResponse: %s\n", PrettyJson(newUser), PrettyJson(user)))

	err = SendVerify(url, user.Name, user.Token)
	if err != nil {
		log.Println(fmt.Sprintf("Verify New User:\nName: %s\nToken: %s\nError: %s\n", user.Name, user.Token, err.Error()))
		return
	}
	log.Println(fmt.Sprintf("Verify New User:\nName: %s\nToken: %s\nResponse: OK\n", user.Name, user.Token))

	lr, err := SendLogin(url, newUser.Name, newUser.Password)
	if err != nil {
		log.Println(fmt.Sprintf("Login User:\nName: %s\nPassword: %s\nError: %s\n", newUser.Name, newUser.Password, err.Error()))
		return
	}
	log.Println(fmt.Sprintf("Login User:\nName: %s\nPassword: %s\nResponse: %s\n", user.Name, newUser.Password, PrettyJson(lr)))

	err = SendLogout(url, lr.UserId)
	if err != nil {
		log.Println(fmt.Sprintf("Logout User:\nUser Id: %s\nErrorr: %s\n", lr.UserId, err.Error()))
		return
	}
	log.Println(fmt.Sprintf("Logout User:\nUser Id: %s\nResponse: OK\n", lr.UserId))

}

func run(ctx context.Context, url string, interval int) {

	for {
		select {
		case <-time.After(time.Duration(interval) * time.Second):
			script(url)
		case <-ctx.Done():

		}
	}
}

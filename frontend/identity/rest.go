package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Server struct {
	
	echoWebServer *echo.Echo
	kube *KubeMQ
	config *Config
	
}

func NewServer(kube *KubeMQ, config *Config) (*Server, error) {
	s := &Server{
		echoWebServer: nil,
		kube:kube,
		config:config,
	}
	
	e := echo.New()
	s.echoWebServer = e

	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	e.HTTPErrorHandler = s.customHTTPErrorHandler
	e.POST("/register", func(c echo.Context) error {
		return register(c)
	})
	//e.GET("/", func(c echo.Context) error {
	//	return c.String(http.StatusOK, "hi")
	//})
	//
	//// route to websocket of subscribe to commands
	//e.GET("/subscribe/events", func(c echo.Context) error {
	//	return nil
	//})
	//
	//// route to websocket of subscribe to request
	//e.GET("/subscribe/requests", func(c echo.Context) error {
	//
	//	return s.handlerSubscribeToRequests(c, appConfig)
	//})
	//
	//// route to websocket of sending stream of messages
	//e.GET("/send/stream", func(c echo.Context) error {
	//
	//	return s.handlerSendMessageStream(c)
	//
	//})
	//
	//// route to post message
	//e.POST("/send/event", func(c echo.Context) error {
	//	return s.handlerSendMessage(c)
	//})
	//
	//// route to post request
	//e.POST("/send/request", func(c echo.Context) error {
	//	return s.handlerSendRequest(c)
	//})
	//
	//// route to post response
	//e.POST("/send/response", func(c echo.Context) error {
	//	return s.handlerSendResponse(c)
	//})

	go func () {
		_ = s.echoWebServer.Start(":" + port)
	}()

	return s, nil
}



type ErrorMessage struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

func (s *Server) customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		c.JSON(code, &ErrorMessage{
			ErrorCode: he.Code,
			Message:   he.Message.(string),
		})

	}
}

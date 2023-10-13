package main

import (
	"github.com/gin-gonic/gin"
)

type ApiServer struct {
	client string
	store  Store
}

// NewApiServer returns a new ApiServer instance
func NewApiServer() *ApiServer {
	env := GetEnv()
	return &ApiServer{
		client: env.Listen,
		store:  NewDbInstance(),
	}
}

func (api *ApiServer) Home(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{})
}

func (api *ApiServer) Start() {
	api.store.preStart()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.GET("/", api.Home)

	err := r.Run(api.client)
	if err != nil {
		panic(err)
	}
}

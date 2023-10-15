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

func (api *ApiServer) handleGetPastes(c *gin.Context) {
    pastes, err := api.store.GetAllPastes()
    if err != nil {
        c.String(200, "Error getting pastes")
    }

    c.HTML(200, "paste.html", gin.H{
        "pastes": pastes,
    })
}

func (api *ApiServer) handleGetPastesByUserName(c *gin.Context) {
    pastes, err := api.store.GetPastesByUserName(c.Param("username"))
    if err != nil {
        c.String(200, "Error getting pastes")
    }

    c.HTML(200, "paste.html", gin.H{
        "pastes": pastes,
    })
}

func (api *ApiServer) handlePasteCreation(c *gin.Context) {
    // This actually has to be from the session or cookie
    username := "noob"

    req := CreatePasteRequest{
        Content: c.PostForm("content"),
        Lang: c.PostForm("lang"),
    }

    err := api.store.CreatePaste(username, req)
    if err != nil {
        c.String(200, "Error creating paste")
    }

    c.String(200, "Paste created!")
}

func (api *ApiServer) Start() {
	api.store.preStart()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.GET("/", api.Home)
    r.GET("/pastes", api.handleGetPastes)
    r.GET("/pastes/:username", api.handleGetPastesByUserName)
    r.POST("/create_paste", api.handlePasteCreation)

	err := r.Run(api.client)
	if err != nil {
		panic(err)
	}
}

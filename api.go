package main

import (
	"fmt"
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
		Lang:    c.PostForm("lang"),
	}

	err := api.store.CreatePaste(username, req)
	fmt.Println(err)
	if err != nil {
		c.String(200, "Error creating paste")
	}

	c.String(200, "Paste created!")
}

func (api *ApiServer) handleUserCreation(c *gin.Context) {
	req := CreateUserRequest{
		Username: c.PostForm("username"),
		Email:    c.PostForm("email"),
		Password: c.PostForm("password"),
	}

	err := api.store.CreateUser(req)
	if err != nil {
		c.String(200, "Error creating user")
	} else {
		c.String(200, "User created!")
	}
}

func (api *ApiServer) handleUserLogin(c *gin.Context) {
	userName := c.PostForm("username")
	password := c.PostForm("password")

	actualPassword, err := api.store.GetUserPassword(userName)
	if err != nil {
		c.String(200, "Error getting user password")
	}

	if !MatchPasswords(password, actualPassword) {
		c.String(200, "Passwords don't match")
	}

	// Create session
	sessionId, err := api.store.CreateSession(CreateSessionRequest{Username: userName})
	if err != nil {
		c.String(200, "Error creating session")
	}

	c.SetCookie("session_id", sessionId, 3600, "/", "localhost", false, true)
	c.String(200, "Logged in!")
}

func (api *ApiServer) isAuth(sessionId string) (string, error) {
	return api.store.GetSession(sessionId)
}

func (api *ApiServer) handleAuth(c *gin.Context) {
	sessionId, err := c.Cookie("session_id")
	if err != nil {
		sessionId = c.GetHeader("session_id")
		if sessionId == "" {
			c.AbortWithStatus(401)
		}
	}

	_, err = api.store.GetSession(sessionId)
	if err != nil {
		c.AbortWithStatus(401)
	}

	c.Next()
}

func (api *ApiServer) handleLoginPage(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{})
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
	r.POST("/create_user", api.handleUserCreation)
	r.POST("/login_user", api.handleUserLogin)
	r.GET("/auth", api.handleAuth)
	r.GET("/login", api.handleLoginPage)

	err := r.Run(api.client)
	if err != nil {
		panic(err)
	}
}

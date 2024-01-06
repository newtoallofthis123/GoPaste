package main

type Env struct {
	Db       string
	User     string
	Password string
	Host     string
	Port     string
	Listen   string
}

type Paste struct {
	PasteId   string
	Username  string
	Content   string
	Lang      string
	CreatedAt string
}

type CreatePasteRequest struct {
	Content string
	Lang    string
}

type CreateUserRequest struct {
	Username string
	Email    string
	Password string
}

type CreateSessionRequest struct {
	Username string
}

type User struct {
	Username  string
	Email     string
	FullName  string
	Bio       string
	CreatedAt string
}

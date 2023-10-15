package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// DbInstance A Struct to hold the Database Instance
// Any interaction with the database through the application
// will be done through this struct only
type DbInstance struct {
	db *sql.DB
}

type Store interface {
	preStart()
    GetAllPastes() ([]Paste, error)
    GetPastesByUserName(string) ([]Paste, error)
    CreatePaste(string, CreatePasteRequest) error
}

// NewDbInstance : Constructs the connection string from
// the environment variables and returns a pointer to the
// DbInstance struct with a connection to the database.
func NewDbInstance() *DbInstance {
	env := GetEnv()
	var connStr = fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", env.User, env.Password, env.Db, env.Host, env.Port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return &DbInstance{
		db: db,
	}
}

// preStart : The Database Initialization function
func (pq *DbInstance) preStart() {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			username VARCHAR(50) PRIMARY KEY NOT NULL,
		    full_name VARCHAR(225) NOT NULL,
			bio TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS sessions (
		    session_id VARCHAR(50) PRIMARY KEY,
		    username VARCHAR(50) REFERENCES users(username) NOT NULL,
		    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    valid BOOLEAN DEFAULT TRUE
		);

		CREATE TABLE IF NOT EXISTS pastes (
			paste_id VARCHAR(50) PRIMARY KEY,
			username VARCHAR(50) REFERENCES users(username) NOT NULL,
			content TEXT NOT NULL,
			lang VARCHAR(50) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := pq.db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func (pq *DbInstance) GetAllPastes() ([]Paste, error){
    query := `
    SELECT * from pastes;
    `

    rows, err := pq.db.Query(query)
    if err != nil {
        return nil, err
    }

    var pastes []Paste
    for rows.Next() {
        var paste Paste
        rows.Scan(&paste.PasteId, &paste.Username, &paste.Content, &paste.Lang, &paste.CreatedAt)
        pastes = append(pastes, paste)
    }

    return pastes, nil
}

func (pq *DbInstance) GetPastesByUserName(username string) ([]Paste, error){
    query := `
    SELECT * from pastes WHERE username = $1;
    `

    rows, err := pq.db.Query(query, username)
    if err != nil {
        return nil, err
    }

    var pastes []Paste
    for rows.Next() {
        var paste Paste
        rows.Scan(&paste.PasteId, &paste.Username, &paste.Content, &paste.Lang, &paste.CreatedAt)
        pastes = append(pastes, paste)
    }

    return pastes, nil
}

func (pq *DbInstance) CreatePaste(username string, req CreatePasteRequest) error {
    query := `
    INSERT INTO pastes (paste_id, username, content, lang) VALUES ($1, $2, $3, $4);
    `

    _, err := pq.db.Exec(query, RanHash(10), username, req.Content, req.Lang)
    if err != nil {
        return err
    }

    return nil
}

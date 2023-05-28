package apiserver

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gopherschool/http-rest-api/internal/app/store/sqlstore"
	"github.com/gorilla/sessions"
)

// Start ...
func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	srv := newServer(store, sessionStore)

	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	var success bool
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE login=? AND password=?", "asd", "asd").Scan(&success)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("log.Logger: %v\n", success)
	return db, nil
}

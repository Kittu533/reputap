package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/reputap")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("UPDATE users SET password_hash = ? WHERE email = ?", string(hash), "admin@reputap.id")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Password updated for admin@reputap.id")
}

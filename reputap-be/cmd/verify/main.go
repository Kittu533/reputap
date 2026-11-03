package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/reputap")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// List tables
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("📋 Tables in 'reputap':")
	fmt.Println("─────────────────────────")
	i := 1
	for rows.Next() {
		var table string
		rows.Scan(&table)
		fmt.Printf("  %2d. %s\n", i, table)
		i++
	}

	// Check seed data
	fmt.Println("\n👤 Seeded Users:")
	fmt.Println("─────────────────────────")
	userRows, _ := db.Query("SELECT u.name, u.email, r.name as role FROM users u JOIN roles r ON u.role_id = r.id")
	defer userRows.Close()
	for userRows.Next() {
		var name, email, role string
		userRows.Scan(&name, &email, &role)
		fmt.Printf("  %s (%s) — role: %s\n", name, email, role)
	}

	// Check roles
	fmt.Println("\n🔑 Roles:")
	fmt.Println("─────────────────────────")
	roleRows, _ := db.Query("SELECT name, description FROM roles")
	defer roleRows.Close()
	for roleRows.Next() {
		var name, desc string
		roleRows.Scan(&name, &desc)
		fmt.Printf("  • %s — %s\n", name, desc)
	}

	// Check permissions count
	var permCount int
	db.QueryRow("SELECT COUNT(*) FROM permissions").Scan(&permCount)
	fmt.Printf("\n🛡️  Permissions: %d total\n", permCount)

	var rpCount int
	db.QueryRow("SELECT COUNT(*) FROM role_permissions").Scan(&rpCount)
	fmt.Printf("🔗 Role-Permission mappings: %d total\n", rpCount)
}

package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite" // The blank import is used to register the driver
)

///////// VIBE CODE EXAMPLE ///////////

// Item represents a simple data structure.
type Item struct {
	ID   int
	Name string
}

func main() {
	// For simplicity, we'll remove the old database file on each run.
	// In a real app, you would not do this.

	// 1. Open a connection to the SQLite database.
	// The database/sql package is the standard Go interface for SQL databases.
	// The Open function doesn't create the database file immediately or establish a connection.
	// It just prepares a database handle for later use.
	db, err := sql.Open("sqlite", "./items.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close() // Ensure the database connection is closed when the function exits.

	log.Println("Database opened successfully.")

	// 2. Create a table.
	// We use db.Exec() for statements that don't return rows.
	createTableSQL := `CREATE TABLE IF NOT EXISTS items (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	log.Println("Table 'items' created successfully.")

	// 3. Insert some data.
	// Use prepared statements with '?' placeholders to prevent SQL injection.
	item1 := "Bread"
	item2 := "Milk"
	insertSQL := `INSERT INTO items(name) VALUES (?)`
	_, err = db.Exec(insertSQL, item1)
	if err != nil {
		log.Fatalf("Failed to insert item 1: %v", err)
	}
	_, err = db.Exec(insertSQL, item2)
	if err != nil {
		log.Fatalf("Failed to insert item 2: %v", err)
	}

	log.Println("Inserted 2 items successfully.")

	// 4. Query the data.
	// We use db.Query() for statements that are expected to return multiple rows.
	rows, err := db.Query("SELECT id, name FROM items")
	if err != nil {
		log.Fatalf("Failed to query items: %v", err)
	}
	defer rows.Close() // Important: always close the rows iterator.

	log.Println("Reading items from the database:")
	var items []Item
	for rows.Next() {
		var item Item
		// Scan copies the columns from the current row into the variables.
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}
		items = append(items, item)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating rows: %v", err)
	}

	// Print the results.
	for _, item := range items {
		log.Printf(" - ID: %d, Name: %s\n", item.ID, item.Name)
	}
}

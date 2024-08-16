package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"
	
	_ "github.com/go-sql-driver/mysql"
)



func main() {
	fmt.Println("Welcome enter your code hereh !")
	dsn := "nizar:dfskj_878$*=@tcp(127.0.0.1:3306)/ecommerce"

    // Open a connection to the database
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Error opening database: %v", err)
    }
    defer db.Close()

    // Ping the database to check if it's reachable
    if err := db.Ping(); err != nil {
        log.Fatalf("Error pinging database: %v", err)
    }
    fmt.Println("Successfully connected to the database!")


	query := `
		SELECT 
			ced.CustomerID,
			SUM(ced.Quantity * cp.Price) AS TotalCA
		FROM 
			CustomerEventData ced
		JOIN 
			ContentPrice cp ON ced.ContentID = cp.ContentID
		WHERE 
			ced.EventTypeID = 6 
			AND ced.EventDate >= '2020-04-01'
		GROUP BY 
			ced.CustomerID
		ORDER BY 
			TotalCA DESC;
	`
	//Execute the query
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}
	defer rows.Close()

// Map to store CustomerID and their TotalCA
customerMap := make(map[int]float64)

// Iterate through the result set and populate the map
for rows.Next() {
	var customerID int
	var totalCA float64
	if err := rows.Scan(&customerID, &totalCA); err != nil {
		log.Fatalf("Error scanning row: %v", err)
	}
	customerMap[customerID] = totalCA
}

// Check for errors after row iteration
if err := rows.Err(); err != nil {
	log.Fatalf("Error iterating rows: %v", err)
}

// Print 10 random entries from the map
printRandomEntries(customerMap, 10)



}










/////////////////////////////////////////////

















func printRandomEntries(customerMap map[int]float64, count int) {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Convert map keys to a slice for random selection
	keys := make([]int, 0, len(customerMap))
	for k := range customerMap {
		keys = append(keys, k)
	}

	// Shuffle the keys slice to randomize the order
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	// Print the specified number of entries
	for i := 0; i < count && i < len(keys); i++ {
		customerID := keys[i]
		totalCA := customerMap[customerID]
		fmt.Printf("CustomerID: %d, TotalCA: %.2f\n", customerID, totalCA)
	}
}
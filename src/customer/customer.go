package customer

import (
    "database/sql"
    "fmt"
   
    "math/rand"
)

// GetCustomerCA retrieves customer data from the database and returns a map of CustomerID to TotalCA.
func GetCustomerCA(db *sql.DB) (map[int]float64, error) {
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
    
    // Execute the query
    rows, err := db.Query(query) 
    if err != nil {
        return nil, fmt.Errorf("error executing query: %w", err)
    }
    defer rows.Close()

    // Map to store CustomerID and their TotalCA
    customerMap := make(map[int]float64)

    // Iterate through the result set and populate the map
    for rows.Next() {
        var customerID int
        var totalCA float64
        if err := rows.Scan(&customerID, &totalCA); err != nil {
            return nil, fmt.Errorf("error scanning row: %w", err)
        }
        customerMap[customerID] = totalCA
    }



    // Check for errors after row iteration
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating rows: %w", err)
    }

    return customerMap, nil
}




// PrintRandomEntries prints a specified number of random entries from the customer map.
func PrintRandomEntries(customerMap map[int]float64, count int) {
    keys := make([]int, 0, len(customerMap))
    for k := range customerMap {
        keys = append(keys, k)
    }

    rand.Shuffle(len(keys), func(i, j int) {
        keys[i], keys[j] = keys[j], keys[i]
    })

    for i := 0; i < count && i < len(keys); i++ {
        customerID := keys[i]
        fmt.Printf("CustomerID: %d, TotalCA: %.2f\n", customerID, customerMap[customerID])
    }
}
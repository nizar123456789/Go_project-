package export 

import("fmt"
"database/sql"
"My_project/src/quantiles"
"log"
)


type CustomerExport struct {
    CustomerID int64
    Revenue    float64
	Email string 
}

// This function is responsable for getting the email of the customer based on the customer ID

func GetCustomerEmail(db *sql.DB, customerID int64) (string, error) {
    var email string
    query := `
        SELECT ChannelValue
        FROM customerData
        WHERE CustomerID = ?;
    `

    err := db.QueryRow(query, customerID).Scan(&email)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", fmt.Errorf("no email found for CustomerID %d", customerID)
        }
        return "", fmt.Errorf("error retrieving email: %w", err)
    }

    return email, nil
}



// 

func SaveOrUpdateCustomerData(db *sql.DB, customers []quantile.Customer) error {
    tableName := "test_export_date"

    // Create the table if it does not exist
    createTableQuery := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            CustomerID BIGINT PRIMARY KEY,
            CA FLOAT,
            Email VARCHAR(255),
            Date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    `, tableName)

    _, err := db.Exec(createTableQuery)
    if err != nil {
        log.Printf("Error creating table: %v\n", err)
        return fmt.Errorf("error creating table: %w", err)
    }

    // Insert or update customer data
    for _, customer := range customers {
        // Attempt to update the existing record
        updateQuery := fmt.Sprintf(`
            UPDATE %s
            SET CA = ?, Date = CURRENT_TIMESTAMP
            WHERE CustomerID = ?;
        `, tableName)
        
        result, err := db.Exec(updateQuery, customer.Revenue, customer.CustomerID)
        if err != nil {
            log.Printf("Error updating customer data for CustomerID %d: %v\n", customer.CustomerID, err)
            return fmt.Errorf("error updating customer data: %w", err)
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil {
            log.Printf("Error getting rows affected for CustomerID %d: %v\n", customer.CustomerID, err)
            return fmt.Errorf("error getting rows affected: %w", err)
        }

        // If no rows were affected, it means the customer doesn't exist, so insert a new record
        if rowsAffected == 0 {
            Email, err := GetCustomerEmail(db, customer.CustomerID)
            if err != nil {
                log.Printf("Error retrieving email for CustomerID %d: %v\n", customer.CustomerID, err)
                return fmt.Errorf("error retrieving email: %w", err)
            }

            // Before inserting, check if the customer exists (safeguard)
            var exists bool
            checkQuery := fmt.Sprintf(`
                SELECT EXISTS(
                    SELECT 1 FROM %s WHERE CustomerID = ?
                );
            `, tableName)

            err = db.QueryRow(checkQuery, customer.CustomerID).Scan(&exists)
            if err != nil {
                log.Printf("Error checking existence of CustomerID %d: %v\n", customer.CustomerID, err)
                return fmt.Errorf("error checking existence of CustomerID: %w", err)
            }

            if !exists {
                insertQuery := fmt.Sprintf(`
                    INSERT INTO %s (CustomerID, CA, Email, Date)
                    VALUES (?, ?, ?, CURRENT_TIMESTAMP);
                `, tableName)
                
                _, err = db.Exec(insertQuery, customer.CustomerID, customer.Revenue, Email)
                if err != nil {
                    log.Printf("Error inserting customer data for CustomerID %d: %v\n", customer.CustomerID, err)
                    return fmt.Errorf("error inserting customer data: %w", err)
                }
            } else {
                log.Printf("CustomerID %d already exists, skipping insert.\n", customer.CustomerID)
            }
        }
    }

    return nil
}

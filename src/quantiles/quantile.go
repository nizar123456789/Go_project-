package quantile 

import("fmt"
 "database/sql")

// defining the structure of Customer and QuantileInfo

type Customer struct {
    CustomerID int64
    Revenue    float64
}

type QuantileInfo struct {
    NumClients int
    MaxRevenue float64
    MinRevenue float64
}

// this function get all the customers  from the database ordered by their revenues in descending order  and return a slice 
func FetchCustomerData(db *sql.DB) ([]Customer, error) {
    rows, err := db.Query(`
        SELECT CustomerID, SUM(Price * Quantity) AS Revenue
        FROM CustomerEventData
        JOIN ContentPrice ON CustomerEventData.ContentID = ContentPrice.ContentID
        GROUP BY CustomerID
        ORDER BY Revenue DESC
    `)

    // check if the query is executed successfully if there is an error then the function returns nil
    if err != nil {
        return nil, err
    }
    // ensure that the rows set result is closed after  processing  
    defer rows.Close()
    
    var customers []Customer
    // loop throw each row and append to the customers slices a customer object in each iteration
    for rows.Next() {
        var customer Customer
        if err := rows.Scan(&customer.CustomerID, &customer.Revenue); err != nil {
            return nil, err
        }
        customers = append(customers, customer)
    }
    return customers, rows.Err()
}



// takes a slice of customers as parameter and then return a map where each key is the quantile and its value is QuantileInfo object
func CalculateQuantiles(customers []Customer) map[float64]QuantileInfo {
    quantiles := map[float64]QuantileInfo{}
    // set the quantile step to (2.5%)
    quantile := 0.025
    // calculate the number of quantiles
    numQuantiles := int(1 / quantile)
    numCustomers := len(customers)
    // if there are no customers return an empty map
    if numCustomers == 0 {
        return quantiles
    }
    // iterate over each quantile 

    for i := 0; i < numQuantiles; i++ {
        startIndex := i * (numCustomers / numQuantiles)
        endIndex := (i + 1) * (numCustomers / numQuantiles)
        if i == numQuantiles-1 {
            endIndex = numCustomers
        }
        //Slices the customer list to get only the customers in the current quantile.

        quantileCustomers := customers[startIndex:endIndex]
        if len(quantileCustomers) == 0 {
            continue
        }

        var maxRevenue, minRevenue float64
        maxRevenue = quantileCustomers[0].Revenue
        minRevenue = quantileCustomers[0].Revenue
        //Iterates over each customer in the current quantile.
        for _, cust := range quantileCustomers {
            if cust.Revenue > maxRevenue {
                maxRevenue = cust.Revenue
            }
            if cust.Revenue < minRevenue {
                minRevenue = cust.Revenue
            }
        }

        quantiles[quantile*100] = QuantileInfo{
            NumClients: len(quantileCustomers),
            MaxRevenue: maxRevenue,
            MinRevenue: minRevenue,
        }

        quantile += 0.025
    }

    return quantiles
}
// Print the quantile info 
func PrintQuantileInfo(quantiles map[float64]QuantileInfo) {
    for q, info := range quantiles {
        fmt.Printf("Quantile %v%%:\n", q)
        fmt.Printf("  Number of clients: %d\n", info.NumClients)
        fmt.Printf("  Max revenue: %.2f\n", info.MaxRevenue)
        fmt.Printf("  Min revenue: %.2f\n", info.MinRevenue)
    }
}


func TopNCustomers(customers []Customer, N int) []Customer {
    if N > len(customers) {
        N = len(customers)
    }
    return customers[:N]
}

func PrintTopNCustomers(customers []Customer) {
    fmt.Printf("Top %d Customers:\n", len(customers))
    for i, cust := range customers {
        fmt.Printf("%d. CustomerID: %d, Revenue: %.2f\n", i+1, cust.CustomerID, cust.Revenue)
    }
}
package main

import (
	"fmt"
    "database/sql"
    "My_project/src/logger"
    "log"
    _ "github.com/go-sql-driver/mysql"
	"My_project/src/quantiles"
	"My_project/src/customer"
    "My_project/src/export"
)

// Define data structures here


func main() {
	//Connect to the database called ecommerce 
	dsn := "Condidat2020:dfskj_878$*=@tcp(127.0.0.1:3306)/ecommerce"

    // Open a connection to the database
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        logger.LogError("This is an error message")
    }
    defer db.Close()
    logger.Initialize()
    
    // Ping the database to check if it's reachable
    if err := db.Ping(); err != nil {
        logger.LogError("This is an error message")
    }
    fmt.Println("Successfully connected to the database!")
    logger.LogInfo("Successfully connected to the database")

    // Call the GetCustomerCA function to get the map of CustomerID to TotalCA
    customerMap, err := customer.GetCustomerCA(db)
    if err != nil {
        logger.LogError("This is an error message")
    }
    logger.LogInfo("Get the map of CustomerID to TotalCA")
    // Call the PrintRandomEntries function to print 10 random entries from the customerMap
    customer.PrintRandomEntries(customerMap, 10)
    logger.LogInfo("Print randomly 10 entries from the customerMap")


    // Fetch and process customer data
    customers, err := quantile.FetchCustomerData(db)
    if err != nil {
        log.Fatal(err)
    }
    logger.LogInfo("CustomerData fetched successfully")
    // Calculate quantiles and print results
	
    quantiles := quantile.CalculateQuantiles(customers)
    quantile.PrintQuantileInfo(quantiles)
    N := 5
    topCustomers := quantile.TopNCustomers(customers, N)
    logger.LogInfo("Get the top Customer done  ")
    // Step 5: Print top N customers
    quantile.PrintTopNCustomers(topCustomers)
    logger.LogInfo("print the  TOP N customer")
export.SaveOrUpdateCustomerData(db,customers)

logger.LogInfo("Filling the table test_export_Data is completed ")

}











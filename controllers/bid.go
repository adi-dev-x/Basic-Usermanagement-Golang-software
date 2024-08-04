package controllers

import (
	//"fmt"
	//"errors"
	"fmt"
	"net/http"

	//"os/user"
	"projectgo/database"
	"projectgo/model"

	"github.com/gin-gonic/gin"
	// "gorm.io/gorm"
)

// Get total bids
func GetBid(c *gin.Context) {
	var toolProduct []model.Bid
	database.DB.Find(&toolProduct)
	//fmt.Println(menus)

	//Prpare menu data for response ,including only desire fields
	tools := make([]gin.H, len(toolProduct))
	for i, toolitem := range toolProduct {
		tools[i] = gin.H{
			"menuID": toolitem.ID,

			"Date":           toolitem.Date,
			"Status":         toolitem.Status,
			"Payment_Status": toolitem.PaymentStatus,
			"ProductId":      toolitem.ProductId,
			"User":           toolitem.Users,
			"Units":          toolitem.Units,
			"Extra_rate":     toolitem.ExtraRate,
			"Vendor":         toolitem.Vendor,
		}
	}
	c.JSON(200, gin.H{
		"status":   "Success",
		"message":  "Tools details fetched successfully",
		"menulist": tools,
	})
}

// Get users bids
func GetUserBid(c *gin.Context) {
	userID := c.Param("id")
	fmt.Println("this is the id    ", userID)

	// Prepare the SQL query
	query := `SELECT id, name, time, date, status, payment_status, product_id, users, units, extra_rate, vendor FROM bids WHERE users = $1`

	// Get the underlying *sql.DB from the GORM connection
	db, err := database.DB.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Execute the SQL query
	fmt.Println("this is the query and id", query, "   ** ", userID)
	rows, err := db.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	// Define a slice to hold the results
	var results []map[string]interface{}

	// Iterate over the rows and fetch column names dynamically
	columns, err := rows.Columns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch rows
	for rows.Next() {
		// Create a slice to hold field values for the current row
		values := make([]interface{}, len(columns))
		// Create a slice to hold pointers to each field value
		valuePointers := make([]interface{}, len(columns))

		// Populate valuePointers with pointers to each field value
		for i := range values {
			valuePointers[i] = &values[i]
		}

		// Scan the current row into the value pointers
		if err := rows.Scan(valuePointers...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Construct a map for the current row
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			rowMap[col] = values[i]
		}

		// Append the row map to the results slice
		results = append(results, rowMap)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the retrieved bids
	c.JSON(http.StatusOK, gin.H{"bids": results})
}

// Get vendors bids
func GetVendorBid(c *gin.Context) {
	userID := c.Param("id")
	fmt.Println("this is the id    ", userID)

	// Prepare the SQL query
	query := `SELECT id, name, time, date, status, payment_status, product_id, users, units, extra_rate, vendor FROM bids WHERE vendor = $1`

	// Get the underlying *sql.DB from the GORM connection
	db, err := database.DB.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Execute the SQL query
	fmt.Println("this is the query and id", query, "   ** ", userID)
	rows, err := db.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	// Define a slice to hold the results
	var results []map[string]interface{}

	// Iterate over the rows and fetch column names dynamically
	columns, err := rows.Columns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch rows
	for rows.Next() {
		// Create a slice to hold field values for the current row
		values := make([]interface{}, len(columns))
		// Create a slice to hold pointers to each field value
		valuePointers := make([]interface{}, len(columns))

		// Populate valuePointers with pointers to each field value
		for i := range values {
			valuePointers[i] = &values[i]
		}

		// Scan the current row into the value pointers
		if err := rows.Scan(valuePointers...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Construct a map for the current row
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			rowMap[col] = values[i]
		}

		// Append the row map to the results slice
		results = append(results, rowMap)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the retrieved bids
	c.JSON(http.StatusOK, gin.H{"bids": results})
}

// Adding bid
func AddBid(c *gin.Context) {
	var table model.Bid
	if err := c.BindJSON(&table); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("this is the data  ", table)
	if table.Date == "" ||
		table.ProductId == 0 || table.Users == 0 || table.Units == 0 ||
		table.ExtraRate == 0 || table.Vendor == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields must have values"})
		return
	}
	if err := database.DB.Create(&table).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{
		"message": "Table created successfully",
		"table": gin.H{
			"toolsID":        table.ID,
			"name":           table.Name,
			"Date":           table.Date,
			"Status":         table.Status,
			"Payment_Status": table.PaymentStatus,
			"ProductId":      table.ProductId,
			"User":           table.Users,
			"Units":          table.Units,
			"Extra_rate":     table.ExtraRate,
			"Vendor":         table.Vendor,
		},
	})
}

// Update the tools
func UpdateBidStatus(c *gin.Context) {
	var bid model.Bid
	bidID := c.Param("id")

	if err := database.DB.First(&bid, bidID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "No bid with this ID"})
		return
	}

	if err := c.BindJSON(&bid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	if err := database.DB.Save(&bid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":  "Success",
		"message": "bid detailes have been updated sucessfully sucessfully",
		"data":    bid,
	})

}

// Delete the tools for admin with authentication
func DeleteBid(c *gin.Context) {
	id := c.Param("id")
	var tool model.Bid

	if err := database.DB.First(&tool, id).Error; err != nil {
		c.JSON(400, gin.H{
			"status":  "Failed",
			"message": "Menu id Not Found",
			"data":    err.Error(),
		})
		return
	}
	database.DB.Delete(&tool)
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "Tools Product Removed Successfully",
		"data":    id,
	})

}

package controllers

import (
	//"fmt"
	"fmt"
	"net/http"
	"projectgo/database"
	"projectgo/model"

	"github.com/gin-gonic/gin"
)

// Get tools services
func GetProduct(c *gin.Context) {
	var toolProduct []model.ProductModel
	database.DB.Find(&toolProduct)
	//fmt.Println(menus)

	//Prpare menu data for response ,including only desire fields
	tools := make([]gin.H, len(toolProduct))
	for i, toolitem := range toolProduct {
		tools[i] = gin.H{
			"menuID":   toolitem.ID,
			"name":     toolitem.Name,
			"category": toolitem.Category,
			"price":    toolitem.Amount,
			"units":    toolitem.Units,
			"tax":      toolitem.Tax,
		}
	}
	c.JSON(200, gin.H{
		"status":   "Success",
		"message":  "Tools details fetched successfully",
		"menulist": tools,
	})
}
func GetVendProduct(c *gin.Context) {
	userID := c.Param("id")
	fmt.Println("this is the id    ", userID)

	// Prepare the SQL query
	query := `SELECT  name, category, status, tax, amount, units,vendor FROM product_models WHERE vendor = $1`

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
	c.JSON(http.StatusOK, gin.H{"products": results})
}

// Create tools
func AddProduct(c *gin.Context) {
	var table model.ProductModel
	if err := c.BindJSON(&table); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if table.Name == "" || table.Category == "" || table.Amount == 0 || table.Units == 0 || table.Tax == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All required fields must be filled"})
		return
	}
	if err := database.DB.Create(&table).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{
		"message": "Table created successfully",
		"table": gin.H{
			"toolsID":         table.ID,
			"category":        table.Category,
			"name of Product": table.Name,
			"amount":          table.Amount,
			"units":           table.Units,
			"tax":             table.Tax,
		},
	})
}

// Update the tools
func UpdateProduct(c *gin.Context) {
	var tools model.ProductModel
	if err := c.BindJSON(&tools); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	toolsID := c.Param("id")
	fmt.Println(toolsID)
	var existingTools model.ProductModel

	if err := database.DB.First(&existingTools, toolsID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Menu not Found"})
		return
	}

	//update the fiels of existing menulist
	existingTools.ID = tools.ID
	existingTools.Category = tools.Category
	existingTools.Name = tools.Name
	// existingTools.amount = tools.amount
	existingTools.Status = tools.Status
	fmt.Println("TOOL ID", existingTools)
	//save the updated tools item to the database

	if err := database.DB.Save(&existingTools).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to update menu"})
		return
	}

	response := gin.H{
		"toolsID":  tools.ID,
		"Category": tools.Category,
		"amount":   tools.Amount,
		"Status":   tools.Status,
		"units":    tools.Units,
	}
	c.JSON(200, gin.H{
		"status":  "Success",
		"message": "Menu Details Updated successfully",
		"data":    response,
	})
}

// Delete the tools for admin with authentication
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var tool model.ProductModel

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

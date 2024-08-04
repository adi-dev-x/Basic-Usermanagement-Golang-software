package main

import (
	"fmt"

	"projectgo/authentication"
	"projectgo/controllers"
	"projectgo/database"
	"projectgo/email"

	"github.com/gin-gonic/gin"
)

func main() {
	email.SendEmailWithOTP("adithyanunni258@gmail.com")
	err := database.InitRedis()
	if err != nil {
		fmt.Printf("Error initializing Redis: %s\n", err.Error())
	} else {
		fmt.Println("Redis connection successful!")
	}
	r := gin.Default()
	database.DBconnect()
	// //admin

	r.POST("/adminSignup", controllers.PostAdmin)
	r.POST("/adminLogin", controllers.AdminLogin)
	r.POST("/adminotpverify", controllers.AdminOtpVerify)
	///user
	r.POST("/userLogin", controllers.UserLogin)
	r.POST("/userSignup", controllers.UserSignup)
	r.POST("/otpverify", controllers.UserOtpVerify)
	//vendor
	r.POST("/vendorLogin", controllers.VenLogin)
	r.POST("/vendorSignup", controllers.VenSignup)
	r.POST("/otpverifyvendor", controllers.VenOtpVerify)
	//creditor
	r.POST("/creditorLogin", controllers.CreLogin)
	r.POST("/creditorSignup", controllers.CreSignup)
	r.POST("/otpverifycreditor", controllers.CreOtpVerify)

	admin := r.Group("/admin")
	admin.Use(authentication.VenAuthMiddleware())
	{
		admin.GET("/productList", controllers.GetProduct)
		admin.POST("/addproduct", controllers.AddProduct)
		admin.PUT("/updateproduct/:id", controllers.UpdateProduct)
		admin.DELETE("/deleteproduct/:id", controllers.DeleteProduct)
	}
	user := r.Group("/user")
	user.Use(authentication.UserAuthMiddleware())
	{
		user.GET("/ProductList/:id", controllers.GetUserBid)
		user.GET("/ProductList", controllers.GetProduct)
		user.POST("/addbid", controllers.AddBid)
	}
	vendor := r.Group("/vendor")
	vendor.Use(authentication.VenAuthMiddleware())
	{
		vendor.GET("/BidList/:id", controllers.GetVendorBid)
		vendor.GET("/productList/:id", controllers.GetVendProduct)
		vendor.POST("/addproduct", controllers.AddProduct)
		vendor.PUT("/updateproduct/:id", controllers.UpdateProduct)
		vendor.POST("/UpdateBidStatus/:id", controllers.UpdateBidStatus)
		vendor.DELETE("/deleteproduct/:id", controllers.DeleteProduct)
	}
	creditor := r.Group("/creditor")
	creditor.Use(authentication.VenAuthMiddleware())
	{
		creditor.GET("/toolsList", controllers.GetProduct)
	}
	r.Run(":3000")
}

package main

import (
	"net/http"
	"restapisch/controllers/Admin"
	"restapisch/controllers/LoginRegister"
	"restapisch/controllers/userController"
	"restapisch/models"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	models.ConnectDatabase()
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	r.POST("/sch/login/", LoginRegister.Login)
	r.POST("/sch/register/", LoginRegister.Register)
	r.GET("/sch/getalldata/:id", userController.GetData)
	r.GET("/admin/getdataquick2", Admin.GetData2)
	r.GET("/admin/getdata/", Admin.GetAllData)
	r.POST("/admin/update/", Admin.UpdateData)
	r.GET("/admin/hapus/:id", Admin.DeleteUser)
	r.GET("/sch/Duser/:id", userController.Detail_pendaftars)
	r.POST("/sch/updatedata/", userController.UpdateData)
	r.POST("/sch/Cekemail/", LoginRegister.Cekemail)
	r.POST("/sch/updpassword/", LoginRegister.UpdatePassword)
	r.GET("/admin/quick", Admin.CountAsal)

	r.Run()
}

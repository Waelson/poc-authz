package main

import (
	"database/sql"
	"net/http"

	"github.com/Waelson/applications-api/internal/db"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	r := gin.Default()

	r.GET("/applications", getApplications)
	r.GET("/applications/:id", getApplicationByID)
	r.POST("/applications", createApplication)
	r.PUT("/applications/:id", updateApplication)
	r.DELETE("/applications/:id", deleteApplication)

	r.Run() // listen and serve on 0.0.0.0:8080
}

// Implemente as funções de manipulação aqui...
func getApplications(c *gin.Context) {
	rows, err := db.DB.Query("SELECT * FROM tb_application")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	applications := make([]db.Application, 0)
	for rows.Next() {
		var app db.Application
		if err := rows.Scan(&app.ID, &app.UserID, &app.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		applications = append(applications, app)
	}
	c.JSON(http.StatusOK, applications)
}

func getApplicationByID(c *gin.Context) {
	var app db.Application
	id := c.Param("id")
	query := "SELECT * FROM tb_application WHERE application_id = ?"
	if err := db.DB.QueryRow(query, id).Scan(&app.ID, &app.UserID, &app.Name); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, app)
}

func createApplication(c *gin.Context) {
	var newApp db.Application
	if err := c.BindJSON(&newApp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO tb_application (application_id, user_id, name) VALUES (?,?,?)"
	_, err := db.DB.Exec(query, newApp.ID, newApp.UserID, newApp.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newApp)
}
func updateApplication(c *gin.Context) {
	id := c.Param("id")
	var updateApp db.Application
	if err := c.BindJSON(&updateApp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "UPDATE tb_application SET user_id = ?, name = ? WHERE application_id = ?"
	_, err := db.DB.Exec(query, updateApp.UserID, updateApp.Name, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	updateApp.ID = id
	c.JSON(http.StatusOK, updateApp)
}

func deleteApplication(c *gin.Context) {
	id := c.Param("id")
	query := "DELETE FROM tb_application WHERE application_id = ?"
	_, err := db.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Application deleted"})
}

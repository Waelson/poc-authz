package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Waelson/operators-api/internal/authz"
	"github.com/Waelson/operators-api/internal/db"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	r := gin.Default()

	r.POST("/operators", createOperator)
	r.GET("/operators", getOperators)
	r.GET("/operators/:id", getOperatorByID)
	r.PUT("/operators/:id", updateOperator)
	r.DELETE("/operators/:id", deleteOperator)

	r.Run(":9090") // listen and serve on 0.0.0.0:8080
}

func createOperator(c *gin.Context) {
	var newOperator db.Operator
	if err := c.BindJSON(&newOperator); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	relationshipJSON, err := json.Marshal(newOperator.Relationship)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO tb_operator (root_id, operator_id, role_id, relationship) VALUES (?, ?, ?, ?)"
	_, err = db.DB.Exec(query, newOperator.RootID, newOperator.OperatorID, newOperator.RoleID, relationshipJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, r := range newOperator.Relationship {
		authzReq := authz.AuthzRequest{
			Resource: authz.Resource{
				Namespace: "operator/application",
				ID:        r.Resource,
			},
			Relation: r.RelationType,
			Subject: authz.Subject{
				Namespace: "operator/user",
				ID:        fmt.Sprint(newOperator.OperatorID),
			},
		}

		fmt.Printf("Imprimindo authzReq: %+v\n", authzReq)

		err = authz.WriteRelationship(authzReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			break
		}
	}

	c.JSON(http.StatusCreated, newOperator)
}

func getOperators(c *gin.Context) {
	rows, err := db.DB.Query("SELECT * FROM tb_operator")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	operators := make([]db.Operator, 0)
	for rows.Next() {
		var op db.Operator
		var relationshipJSON string
		if err := rows.Scan(&op.RootID, &op.OperatorID, &op.RoleID, &relationshipJSON); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		json.Unmarshal([]byte(relationshipJSON), &op.Relationship)
		operators = append(operators, op)
	}

	c.JSON(http.StatusOK, operators)
}

func getOperatorByID(c *gin.Context) {
	id := c.Param("id")
	var op db.Operator
	var relationshipJSON string

	query := "SELECT * FROM tb_operator WHERE operator_id = ?"
	row := db.DB.QueryRow(query, id)
	if err := row.Scan(&op.RootID, &op.OperatorID, &op.RoleID, &relationshipJSON); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Operator not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	json.Unmarshal([]byte(relationshipJSON), &op.Relationship)

	c.JSON(http.StatusOK, op)
}

func updateOperator(c *gin.Context) {
	id := c.Param("id")
	var updateOp db.Operator
	if err := c.BindJSON(&updateOp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	relationshipJSON, _ := json.Marshal(updateOp.Relationship)
	query := "UPDATE tb_operator SET root_id = ?, role_id = ?, relationship = ? WHERE operator_id = ?"
	_, err := db.DB.Exec(query, updateOp.RootID, updateOp.RoleID, relationshipJSON, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Operator updated"})
}

func deleteOperator(c *gin.Context) {
	id := c.Param("id")
	query := "DELETE FROM tb_operator WHERE operator_id = ?"
	_, err := db.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Operator deleted"})
}

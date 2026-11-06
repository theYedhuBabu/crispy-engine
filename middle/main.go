package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type auditEntry struct {
	Initiator string    `json:"initiator"`
	Value     string    `json:"value"`
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time	`json:"created_at"`
	Success	  bool		`json:"success"`
}

func main() {

	router := gin.Default()

	router.GET("/middleware/route", middlewareHandler)

	router.Run(":8081")
}

func middlewareHandler(c *gin.Context) {

	value, err := c.Query("value")
	initiator, err := c.Query("initiator")
	id, err := c.Query("id")

	

	fmt.Println(initiator)
	fmt.Println(id)

	result, err := getValueFromDB(value, id)

	payload := auditEntry{
		Initiator: initiator,
		Value: value,
		Id: id,
		CreatedAt: time.Now(),
		Success: err==nil,
	}

	go writeAudit(payload)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}


	c.JSON(http.StatusOK, result)
}



func getValueFromDB(value string, id uuid.UUID) (gin.H, error) {

}
	

func writeAudit(payload auditEntry) {}
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {

	router := gin.Default()

	router.GET("/api/start", startHandler)

	router.Run(":8080")
}

func startHandler(c *gin.Context) {

	value, err := c.Query("value")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Parameter value",
			"value": value,
		})
	}

	result, err := processValue(value)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}



func processValue(value string) (gin.H, error) {
	//send request to middle service
	id := uuid.New()

	params := url.Values{}
	params.Set("value", value)
	params.Set("initiator", "entrypoint")
	params.Set("id", id)
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(middleURL() + "/middleware/route" + params.Encode())

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Middleware returned %d", resp.StatusCode)
	}

	var result gin.H

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result

}

func middleURL() string {
	if u := os.Getenv("MIDDLE_URL"); u != "" {
		return u
	}
	return "http://localhost:8081"
}

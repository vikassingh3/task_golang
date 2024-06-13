package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/cetec")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()

	r.GET("/person/:person_id/info", getPersonInfo)
	r.POST("/person/create", createPerson)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

type PersonInfo struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	City        string `json:"city"`
	State       string `json:"state"`
	Street1     string `json:"street1"`
	Street2     string `json:"street2"`
	ZipCode     string `json:"zip_code"`
}

func getPersonInfo(c *gin.Context) {
	personID := c.Param("person_id")

	var info PersonInfo
	err := db.QueryRow(`SELECT p.name, ph.number, a.city, a.state, a.street1, a.street2, a.zip_code
	FROM person p
	INNER JOIN phone ph ON p.id = ph.person_id
	INNER JOIN address_join aj ON p.id = aj.person_id
	INNER JOIN address a ON aj.address_id = a.id
	WHERE p.id = ?`, personID).Scan(&info.Name, &info.PhoneNumber, &info.City, &info.State, &info.Street1, &info.Street2, &info.ZipCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get person info"})
		return
	}

	c.JSON(http.StatusOK, info)
}

func createPerson(c *gin.Context) {
	var person struct {
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
		City        string `json:"city"`
		State       string `json:"state"`
		Street1     string `json:"street1"`
		Street2     string `json:"street2"`
		ZipCode     string `json:"zip_code"`
	}

	if err := c.BindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	_, err := db.Exec("INSERT INTO person(name) VALUES(?)", person.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create person"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Person created successfully"})
}

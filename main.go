package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Contact struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

var db *sql.DB

func initDB() {
	connStr := "user=postgres dbname=contacts sslmode=disable password=partner"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS contacts (
			id SERIAL PRIMARY KEY,
			name TEXT,
			phone TEXT,
			email TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func getContacts(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, phone, email FROM contacts")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching contacts"})
		return
	}
	defer rows.Close()

	var contacts []Contact
	for rows.Next() {
		var contact Contact
		err := rows.Scan(&contact.ID, &contact.Name, &contact.Phone, &contact.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning contacts"})
			return
		}
		contacts = append(contacts, contact)
	}
	c.JSON(http.StatusOK, contacts)
}

func getContactByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var contact Contact
	err = db.QueryRow("SELECT id, name, phone, email FROM contacts WHERE id = $1", id).Scan(&contact.ID, &contact.Name, &contact.Phone, &contact.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	c.JSON(http.StatusOK, contact)
}

func createContact(c *gin.Context) {
	var contact Contact
	if err := c.BindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	_, err := db.Exec("INSERT INTO contacts (name, phone, email) VALUES ($1, $2, $3)", contact.Name, contact.Phone, contact.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating contact"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Contact created successfully"})
}

func updateContact(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var contact Contact
	if err := c.BindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	_, err = db.Exec("UPDATE contacts SET name = $1, phone = $2, email = $3 WHERE id = $4", contact.Name, contact.Phone, contact.Email, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating contact"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact updated successfully"})
}

func deleteContact(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	_, err = db.Exec("DELETE FROM contacts WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting contact"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted successfully"})
}

func main() {
	initDB()
	r := gin.Default()

	// API Endpoints
	r.GET("/contacts", getContacts)
	r.GET("/contacts/:id", getContactByID)
	r.POST("/contacts", createContact)
	r.PUT("/contacts/:id", updateContact)
	r.DELETE("/contacts/:id", deleteContact)

	// Serve HTML and static files
	r.Static("/static", "./static")
	r.LoadHTMLFiles("templates/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.Run(":8080") // Start the server on port 8080
}

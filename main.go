package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type User_JSON struct { // JSON
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Introduction string `json:"introduction"`
}

type User struct { // DB
	Id           int
	Name         string
	Email        string
	Password     string
	Introduction string
	CreatedAt    time.Time
}

func main() {
	r := gin.Default()

	u := r.Group("/user")

	u.POST("/create", CreateUser)

	r.Run()
}

func CreateUser(c *gin.Context) {
	db, err := sql.Open("mysql", "root:password@(localhost:3306)/local?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var user User_JSON
	c.ShouldBindJSON(&user)

	insert, err := db.Prepare("INSERT INTO user(name, email, password, introduction) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	insert.Exec(user.Name, user.Email, user.Password, user.Introduction)
}

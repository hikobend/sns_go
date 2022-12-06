package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
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

type UserName struct { // 名前のみ
	Id   int
	Name string
}

func main() {
	r := gin.Default()

	u := r.Group("/user")

	u.POST("/create", CreateUser)
	u.GET("/gets", GetsUser)
	u.GET("/get/:id", GetByIdUser)
	u.PATCH("/update/:id", UpdateUserName)
	u.DELETE("/delete/:id", DeleteUser)

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

func GetsUser(c *gin.Context) {
	db, err := sql.Open("mysql", "root:password@(localhost:3306)/local?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select id, name from user")
	if err != nil {
		log.Fatal(err)
	}
	var resultUser []UserName

	for rows.Next() {
		user := UserName{}
		if err := rows.Scan(&user.Id, &user.Name); err != nil {
			log.Fatal(err)
		}
		resultUser = append(resultUser, user)
	}

	c.JSON(http.StatusOK, resultUser)
}

func GetByIdUser(c *gin.Context) {
	db, err := sql.Open("mysql", "root:password@(localhost:3306)/local?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatalln(err)
	}

	var getuser UserName

	if err = db.QueryRow("SELECT id, name FROM user where id = ?", id).Scan(&getuser.Id, &getuser.Name); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, getuser)
}

func UpdateUserName(c *gin.Context) {
	db, err := sql.Open("mysql", "root:password@(localhost:3306)/local?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var json User_JSON
	c.ShouldBindJSON(&json)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatalln(err)
	}

	update, err := db.Prepare("UPDATE user SET name = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	update.Exec(json.Name, id)
}

func DeleteUser(c *gin.Context) {
	db, err := sql.Open("mysql", "root:password@(localhost:3306)/local?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatalln(err)
	}

	delete, err := db.Prepare("DELETE FROM user WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	delete.Exec(id)
}

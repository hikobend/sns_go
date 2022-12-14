package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hikobend/sns_go/file"
	"golang.org/x/crypto/bcrypt"
)

type User_JSON struct {
	Name         string `json:"name" db:"name" validate:"required"`
	Email        string `json:"email" db:"email" validate:"required,email"`
	Password     string `json:"password" db:"password" validate:"required,min=8,max=50"`
	Introduction string `json:"introduction" db:"introduction"`
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
	m := r.Group("/message")

	u.POST("/create", CreateUser)
	u.GET("/gets", GetsUser)
	u.GET("/get/:id", GetByIdUser)
	u.PATCH("/update/:id", UpdateUserName)
	u.DELETE("/delete/:id", DeleteUser)

	m.POST("/create", file.CreateMessage)
	m.GET("/gets", file.GetsAllMessage)
	m.GET("/get/:id", file.GetMessage)
	m.GET("/user/:id", file.GetByUserIdAllMessage)

	r.Run()
}

// 暗号(Hash)化
func PasswordEncrypt(p string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(hash), err
}

func CreateUser(c *gin.Context) {
	db, err := sql.Open("mysql", "root:password@(localhost:3306)/local?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var user User_JSON
	validate := validator.New()
	c.ShouldBindJSON(&user)

	err = validate.Struct(&user) //バリデーションを実行し、NGの場合、ここでエラーが返る。
	if err != nil {
		log.Fatal(err)
	}

	insert, err := db.Prepare("INSERT INTO user(name, email, password, introduction) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	password, _ := PasswordEncrypt(user.Password)
	insert.Exec(user.Name, user.Email, password, user.Introduction)
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

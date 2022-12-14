package file

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Message_JSON struct { // JSON
	UserId  int    `json:"user_id"`
	Content string `json:"content"`
}

type Message struct { // DB
	Id        int
	UserId    int
	Content   string
	CreatedAt time.Time
}

type MessageContent struct {
	Id      int
	Content string
}

func CreateMessage(c *gin.Context) {
	db, err := sql.Open("mysql", "root:password@(localhost:3306)/local?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var message Message_JSON
	c.ShouldBindJSON(&message)

	insert, err := db.Prepare("INSERT INTO message(user_id, content) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	insert.Exec(message.UserId, message.Content)
}

func GetsAllMessage(c *gin.Context) {
	db, err := sql.Open("mysql", "root:password@(localhost:3306)/local?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select id, content from message")
	if err != nil {
		log.Fatal(err)
	}
	var resultMessage []MessageContent

	for rows.Next() {
		message := MessageContent{}
		if err := rows.Scan(&message.Id, &message.Content); err != nil {
			log.Fatal(err)
		}
		resultMessage = append(resultMessage, message)
	}

	c.JSON(http.StatusOK, resultMessage)
}

func GetMessage(c *gin.Context) {
	db, err := sql.Open("mysql", "root:password@(localhost:3306)/local?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	user_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatalln(err)
	}

	var getmessage MessageContent

	if err = db.QueryRow("SELECT id, content FROM message where id = ?", user_id).Scan(&getmessage.Id, &getmessage.Content); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, getmessage)
}

func GetByUserIdAllMessage(c *gin.Context) {
	db, err := sql.Open("mysql", "root:password@(localhost:3306)/local?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatalln(err)
	}

	rows, err := db.Query("select id, content from message where user_id = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	var resultMessage []MessageContent

	for rows.Next() {
		message := MessageContent{}
		if err := rows.Scan(&message.Id, &message.Content); err != nil {
			log.Fatal(err)
		}
		resultMessage = append(resultMessage, message)
	}

	c.JSON(http.StatusOK, resultMessage)
}

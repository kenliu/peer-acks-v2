package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	router := gin.Default()

	createStaticRoutes(router)

	// initialize the DB
	//TODO fix this to point to correct user and db
	//datasource := "postgresql://maxroach@localhost:26257/bank?ssl=true&sslmode=require&sslrootcert=certs/ca.crt&sslkey=certs/client.maxroach.key&sslcert=certs/client.maxroach.crt""
	datasource := "postgresql://maxroach@localhost:26257/peeracks?sslmode=disable"

	db, err := sql.Open("postgres", datasource)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	defer db.Close()

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", fetchAcks(db))
	})

	router.GET("/acks", func(c *gin.Context) {
		c.JSON(http.StatusOK, fetchAcks(db))
	})

	router.POST("/acks", func(c *gin.Context) {
		message := c.PostForm("message")
		log.Println(message)

		//TODO handle case where msg is empty
		//TODO trim spaces
		db.Exec("INSERT INTO acks (message, updated_at) values ($1, current_timestamp)", message)

		c.HTML(http.StatusOK, "ack_submitted.tmpl", fetchAcks(db))
	})

	router.DELETE("/acks", func(c *gin.Context) {
		//TODO
	})

	// my acks
	router.GET("/myacks", func(c *gin.Context) {
		c.HTML(http.StatusOK, "myacks.tmpl", fetchAcks(db))
	})

	// report
	router.GET("/report")

	router.Run(":8888") // listen and serve on 0.0.0.0:8080
}

func createStaticRoutes(router *gin.Engine) {
	// create static routes
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")
	router.StaticFile("/radiator", "./templates/radiator.html")
	router.Static("/resources", "./resources")

	router.LoadHTMLGlob("templates/*")
}

func deleteAck(c *gin.Context) {
	//TODO delete ack from DB
}

func fetchAcks(db *sql.DB) gin.H {
	var messages []string

	rows, err := db.Query("select message from acks order by updated_at desc") //TODO add date range of previous week
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var message string
	log.Println("fetching acks:")
	for rows.Next() {
		err := rows.Scan(&message)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(message)
		messages = append(messages, message)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return gin.H{
		"acks": messages,
	}
}

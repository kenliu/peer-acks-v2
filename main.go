package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type AcksForm struct {
	ack string `form:"ack"`
}

func main() {
	router := gin.Default()

	createStaticRoutes(router)

	router.GET("/", index)

	// acks
	router.GET("/acks", listAcks)
	router.POST("/ack", createAck)
	router.DELETE("/ack", deleteAck)

	// my acks
	router.GET("/myacks", showMyAcks)

	// report
	router.GET("/report")

	router.Run() // listen and serve on 0.0.0.0:8080
}

func createStaticRoutes(router *gin.Engine) {
	// create static routes
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")
	router.StaticFile("/radiator.html", "./templates/radiator.html")
	router.Static("/resources", "./resources")

	router.LoadHTMLGlob("templates/*")
}

////////////
// acks
func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", fetchAcks())
}

func listAcks(c *gin.Context) {
	c.JSON(http.StatusOK, fetchAcks())
}

func createAck(c *gin.Context) {
	fmt.Println("posting ack")
	var form AcksForm

	// This will infer what binder to use depending on the content-type header.
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(form.ack)

	//TODO write ack to DB
	c.HTML(http.StatusOK, "ack_submitted.tmpl", nil)
}

func deleteAck(c *gin.Context) {
	//TODO delete ack from DB
}

////////////
// myacks
func showMyAcks(c *gin.Context) {
	c.HTML(http.StatusOK, "myacks.tmpl", fetchAcks())
}

func initDb() {
	db, err := sql.Open("postgres",
		"postgresql://maxroach@localhost:26257/bank?ssl=true&sslmode=require&sslrootcert=certs/ca.crt&sslkey=certs/client.maxroach.key&sslcert=certs/client.maxroach.crt")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	defer db.Close()

	// Create the "acks" table.
	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS acks (id INT PRIMARY KEY, sender_email STRING, message STRING, created_at TIMESTAMP)"); err != nil {
		log.Fatal(err)
	}
}

func fetchAcks() gin.H {
	return gin.H{
		"acks": []string{"shout out to a", "shout out to b", "shout out to c"}, //hard coded for now
	}
}

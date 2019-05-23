package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	createStaticRoutes(router)

	router.GET("/", showMyAcks) // TODO change this to a redirect

	// acks
	router.GET("/acks", listAcks)
	router.DELETE("/ack", deleteAck)

	// my acks
	router.GET("/myacks", showMyAcks)

	// report
	router.GET("/report")
	// TODO finish this

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
func listAcks(c *gin.Context) {
	c.JSON(http.StatusOK, fetchAcks())
}

func deleteAck(c *gin.Context) {

}

////////////
// myacks
func showMyAcks(c *gin.Context) {
	c.HTML(http.StatusOK, "myacks.tmpl", fetchAcks())
}

func fetchAcks() gin.H {
	return gin.H{
		"acks": []string{"a", "b", "c"}, //hard coded for now
	}
}
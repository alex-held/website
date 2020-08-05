package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var db map[string]string

func init() {
	db = map[string]string{
		"alex":  "motorcycle",
		"peter": "boat",
		"hans":  "car",
	}
}


func main() {
	port := os.Getenv("PORT")
	mode := os.Getenv("GIN_MODE")
	if env := mode; env == "release" {
		gin.DisableConsoleColor()
	}

	router := NewRouter()
	router.ConfigureRoutes()


	err := router.Run(port) // 0.0.0.0:8080
	if err != nil {
		_, _ = fmt.Fprintf(gin.DefaultWriter, "error occured while starting the server. %#v", err)
	}
}

type Router gin.Engine

func (router *Router) Run(portNumbers ...string) error {
	var portStrings []string
	for _, port := range portNumbers {
		portStrings = append(portStrings, ":" + port)
	}
	engine := (*gin.Engine)(router)
	return engine.Run(portStrings...)
}

func (router *Router) ConfigureRoutes() Router  {
	router.GET("/index", handleIndex)
	router.GET("/home", handleHome)
	router.GET("/hello/:message", handleEcho)

	user := router.Group("/user")
	{
		user.GET("/:name", handleUser) // user.GET("/", handleUser)
	}
	//goland:noinspection GoVetCopyLock
	return *router
}

//goland:noinspection GoVetCopyLock
func NewRouter() Router {
	engine := gin.Default()
	engine.Use(gin.Logger())
	engine.Static("/assets", "/app/assets")
	engine.StaticFile("favicon.ico", "/app/assets/favicon.ico")
	engine.LoadHTMLGlob("templates/*.html.tmpl")
	router := Router(*engine)
	return router
}

func handleHome(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html.tmpl", gin.H{
		"Title":   "home",
		"Message": "hello friends",
	})
}

func handleIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html.tmpl", gin.H{
		"Title": "index",
	})
}

func handleEcho(c *gin.Context) {
	c.String(200, c.Param("message"))
}

func handleUser(c *gin.Context) {
	user := c.Params.ByName("name")

	if vehicle, ok := db[user]; ok {
		c.JSON(http.StatusOK, gin.H{"user": user,
			"vehicle": vehicle})
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user,
			"vehicle": "no vehicle"})
	}
}

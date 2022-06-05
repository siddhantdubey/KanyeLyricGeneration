package main

//set up a basic api with gin
import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
        "github.com/gin-contrib/cors"
)

func sentence(c *gin.Context) {
	//input is a file path to "../data/lyrics/alllyrics.txt"
	input := "../data/lyrics/alllyrics.txt"
	m, err := NewChainFromFile(input, 2)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		c.String(http.StatusOK, m.GenSentence(10, false) + "\n")
	}
}

func main() {
	//set up the router
	r := gin.Default()
        config := cors.DefaultConfig()
        config.AllowAllOrigins = true
        r.Use(cors.New(config))
    	r.GET("/sentence", sentence)
	//start the server
	r.Run(":8080")
}

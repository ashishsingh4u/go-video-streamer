package main

import (
	"fmt"
	"go-video-streamer/config"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.ForceConsoleColor()

	router := gin.Default()

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.CustomRecovery(func(ctx *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			ctx.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}))

	var conf config.Configuration
	if err := config.GetConfig(&conf); err != nil {
		panic(fmt.Sprintf("Couldn't read the configuration file. Error: %s", err.Error()))
	}

	machineIP := fmt.Sprintf("%s:%s", conf.SERVER_IP, conf.PORT)
	log.Printf("Server will be starting on %s\n", machineIP)

	router.Static("/browse", "./static/browse")
	router.Run(machineIP)
}

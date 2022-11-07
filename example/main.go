package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ilaziness/gweb"
)

func main() {
	g := gweb.NewDefault()
	g.GET("/", func(g *gweb.Context) {
		fmt.Println("index")
		g.String(200, "index")
	})
	g.GET("/hello", func(g *gweb.Context) {
		fmt.Println("hello")
		g.JSON(200, map[string]any{"a": 1, "b": "hello"})
	})

	v1 := g.Group("/v1")
	{
		v1.GET("/", func(c *gweb.Context) {
			c.HTML(http.StatusOK, "<h1>Hello gWeb</h1>")
		})
		v1.GET("/hello", func(c *gweb.Context) {
			c.HTML(http.StatusOK, "<h1>Hello gWeb v1</h1>")
		})
	}

	v2 := g.Group("/v2")
	{
		v2.GET("/", func(c *gweb.Context) {
			c.HTML(http.StatusOK, "<h1>Hello gWeb, v2</h1>")
		})
		v2.GET("/hello", func(c *gweb.Context) {
			c.HTML(http.StatusOK, "<h1>Hello gWeb, v2, v2</h1>")
		})
		// 应用中间件
		v2.Use(Logger())
	}

	g.Use(Logger2())

	log.Fatalln(g.Run(":8080"))
}

func Logger() gweb.HandleFunc {
	return func(c *gweb.Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func Logger2() gweb.HandleFunc {
	return func(c *gweb.Context) {
		log.Println("middleware logger2")
	}
}

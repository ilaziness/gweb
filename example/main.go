package main

import (
	"fmt"
	"github.com/ilaziness/gweb"
	"log"
)

func main() {
	g := gweb.NewDefault()
	g.GET("/", func(g *gweb.Context) {
		fmt.Println("index")
		g.String(200, "123")
	})
	g.GET("/hello", func(g *gweb.Context) {
		fmt.Println("hello")
		g.JSON(200, map[string]any{"a": 1, "b": "hello"})
	})
	log.Fatalln(g.Run(":8080"))
}

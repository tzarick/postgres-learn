// server
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func indexHandler(c *fiber.Ctx) error {
	return c.SendString("index")
}
func postHandler(c *fiber.Ctx) error {
	return c.SendString("post")
}
func putHandler(c *fiber.Ctx) error {
	return c.SendString("put")
}
func deleteHandler(c *fiber.Ctx) error {
	return c.SendString("delete")
}

func main() {
	fmt.Println("Hello")
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	app := fiber.New()

	app.Get("/", indexHandler)
	app.Post("/", postHandler)
	app.Put("/update", putHandler)
	app.Delete("/delete", deleteHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}

	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}

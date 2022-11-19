// server
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type movie struct {
	title  string
	length int
}

type response struct {
	Item string
}

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	var res string
	var movieTitles []string

	// display all short film titles

	rows, err := db.Query("SELECT title FROM user_film_dump WHERE length < 50")
	if err != nil {
		log.Fatalln(err)
		c.JSON("An error occured")
	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(&res)
		movieTitles = append(movieTitles, res)
	}

	fmt.Printf("movieTitles: %v\n", len(movieTitles))

	return c.Render("index", fiber.Map{"Movies": movieTitles})
}
func postHandler(c *fiber.Ctx, db *sql.DB) error {
	response := response{}
	newMovie := movie{}

	// add user input to DB with default length value

	if err := c.BodyParser(&response); err != nil {
		log.Printf("An error occured: %v", err)
		return c.SendString(err.Error())
	}

	newMovie.title = response.Item
	newMovie.length = 10
	if newMovie.title != "" {
		_, err := db.Exec("INSERT into user_film_dump (title, length) VALUES ($1, $2)", newMovie.title, newMovie.length)
		if err != nil {
			log.Fatalf("An error occured while executing query: %v", err)
		}
	}

	return c.Redirect("/")
}
func putHandler(c *fiber.Ctx, db *sql.DB) error {
	oldTitle := c.Query("oldTitle")
	newTitle := c.Query("newTitle")

	fmt.Println(oldTitle, newTitle)

	db.Exec("UPDATE user_film_dump SET title = $1 WHERE title = $2", newTitle, oldTitle)

	return c.Redirect("/")
}
func deleteHandler(c *fiber.Ctx, db *sql.DB) error {
	titleToDelete := c.Query("title")

	db.Exec("DELETE from user_film_dump WHERE title = $1", titleToDelete)

	return c.Redirect("/")
}

// connect to our postgres DB. Params are specified in .env file. We are using the lib/pq postgres driver package
func connectToDB() *sql.DB {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PWD")
	dbLocation := fmt.Sprintf("%s:%v", os.Getenv("DB_IP"), os.Getenv("DB_PORT"))
	dbName := "dvdrental"

	connectionStr := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", user, password, dbLocation, dbName)

	// connect to db
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		log.Fatal(err)
	}

	// ensure proper connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func main() {
	fmt.Println("Hello")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}

	db := connectToDB()

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})

	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c, db)
	})
	app.Post("/", func(c *fiber.Ctx) error {
		return postHandler(c, db)
	})
	app.Put("/update", func(c *fiber.Ctx) error {
		return putHandler(c, db)
	})
	app.Delete("/delete", func(c *fiber.Ctx) error {
		return deleteHandler(c, db)
	})

	app.Static("/", "./public")

	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}

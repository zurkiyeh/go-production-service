// cmd dir should have the entry point into the application. That's the standard
// If you want to add a CLI functionality for example, you can add a sub dir just for that

package main

import (
	"context"
	"fmt"

	"github.com/zurkiyeh/go-production-service/internal/comment"
	"github.com/zurkiyeh/go-production-service/internal/db"
)

// Reposnible for instantiation and start up for our a function
// This is best to handle the errors happening in the different layers of our application. For example an error occured while connecting to the db? -> you can handle by forward to a monitoring tool or something. This way your main app doesnt panic
func Run() error {
	fmt.Println("Application starting...")

	db, err := db.NewDatabase()
	if err != nil {
		fmt.Println("Couldnt connect to db")
		return err
	}
	if err := db.MigrateDB(); err != nil {
		fmt.Println("Failed to migrate db")
		return err
	}
	fmt.Println("Database connection was successful")

	cmtService := comment.NewService(db)

	cmtService.PostComment(context.Background(), comment.Comment{
		ID:     "f9a9f7bc-ed87-4aec-b5c6-ec37f6222147",
		Author: "Me",
		Slug:   "",
		Body:   "You're a donkey",
	})

	// fmt.Println(cmtService.GetComment(
	// 	context.Background(),
	// 	"f9a9f7bc-ed87-4aec-b5c6-ec37f6222147",
	// ))

	return nil
}

func main() {
	if err := Run(); err != nil {
		fmt.Printf("ERROR: %v", err)
	}
	fmt.Println("Hello to REST API project ")
}

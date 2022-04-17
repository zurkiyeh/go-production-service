// cmd dir should have the entry point into the application. That's the standard
// If you want to add a CLI functionality for example, you can add a sub dir just for that

package main

import "fmt"

// Reposnible for instantiation and start up for our a function
// This is best to handle the errors happening in the different layers of our application. For example an error occured while connecting to the db? -> you can handle by forward to a monitoring tool or something. This way your main app doesnt panic
func Run() error {
	fmt.Println("Application starting...")
	return nil
}

func main() {
	if err := Run(); err != nil {
		fmt.Println("ERROR: %v", err)
	}
	fmt.Println("Hello to REST API project ")
}

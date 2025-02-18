package main

import "fmt"

type User struct {
	ID int
}

func main() {
	testChan := make(chan *User)
	go func() {
		testChan <- nil
	}()
	fmt.Println(<-testChan)

}

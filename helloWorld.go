package main

import "fmt"

type HelloWorldService interface {
	HelloWorld() (string, error)
}

type helloWorldService struct{}

func (helloWorldService) HelloWorld() (string, error) {
	str := "I am learning go"
	fmt.Print(str)
	return str, nil
}
package main

import (
	"errors"
	"gophermart/internal/config"
	"gophermart/internal/repository"
	"gophermart/internal/routes"
	"gophermart/internal/utils"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	c, err := config.New()
	if err != nil {
		return errors.New("failed to initialize config: " + err.Error())
	}
	u, err := utils.New(c)
	if err != nil {
		return errors.New("failed to initialize utils: " + err.Error())
	}
	repo, err := repository.New(u)

	r := routes.Init(u)

	return r.Run(c.Address)
}

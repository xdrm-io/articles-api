package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/xdrm-io/aicra"
	"github.com/xdrm-io/aicra/validator"
	"github.com/xdrm-io/articles-api/model"
	"github.com/xdrm-io/articles-api/services"
	"github.com/xdrm-io/articles-api/storage"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("fatal: %s", err)
	}
}

func run() error {
	b := &aicra.Builder{}
	b.Input(validator.AnyType{})
	b.Input(validator.BoolType{})
	b.Input(validator.UintType{})
	b.Input(validator.IntType{})
	b.Input(validator.StringType{})

	b.Output("int", reflect.TypeOf(int(0)))
	b.Output("uint", reflect.TypeOf(uint(0)))
	b.Output("string", reflect.TypeOf(""))
	b.Output("user", reflect.TypeOf(model.User{}))
	b.Output("[]user", reflect.TypeOf([]model.User{}))
	b.Output("article", reflect.TypeOf(model.Article{}))
	b.Output("[]article", reflect.TypeOf([]model.Article{}))

	config, err := os.OpenFile("api/definition.json", os.O_RDONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot open config file: %w", err)
	}

	err = b.Setup(config)
	config.Close()
	if err != nil {
		return fmt.Errorf("cannot setup builder: %w", err)
	}

	b.WithContext(services.AuthMiddleware)

	db := &storage.DB{}
	if err := db.Open(); err != nil {
		return fmt.Errorf("database: %w", err)
	}
	defer db.Close()

	handlers := services.NewHandler(db)
	err = handlers.Wire(b)
	if err != nil {
		return fmt.Errorf("wire: %w", err)
	}

	server, err := b.Build()
	if err != nil {
		return fmt.Errorf("cannot build server: %w", err)
	}

	// 5. listen and serve
	log.Printf("[server] up and running at 0.0.0.0:4242")
	return http.ListenAndServe("0.0.0.0:4242", server)
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/xdrm-io/aicra"
	"github.com/xdrm-io/aicra/validator"
	"github.com/xdrm-io/articles-api/services"
	"github.com/xdrm-io/articles-api/storage"
	"github.com/xdrm-io/articles-api/types"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("fatal: %s", err)
	}
}

func run() error {
	b := &aicra.Builder{}
	b.Validate(validator.AnyType{})
	b.Validate(validator.BoolType{})
	b.Validate(validator.UintType{})
	b.Validate(validator.IntType{})
	b.Validate(validator.StringType{})

	b.Validate(types.UserType{})
	b.Validate(types.UsersType{})
	b.Validate(types.ArticleType{})
	b.Validate(types.ArticlesType{})

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

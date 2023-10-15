package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alexandremarti/golang/rest-kafka/internal/infra/akafka"
	"github.com/alexandremarti/golang/rest-kafka/internal/infra/repository"
	"github.com/alexandremarti/golang/rest-kafka/internal/infra/web"
	"github.com/alexandremarti/golang/rest-kafka/internal/usecase"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// db, err := sql.Open("mysql", "root:roo@tcp(host.docker.internal:3306)/products")
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/products")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repository := repository.NewProductRepositoryMysql(db)
	createProductUseCase := usecase.NewCreateProductUseCase(repository)
	listProductsUseCase := usecase.NewListProductsUseCase(repository)

	productHandlers := web.NewProductHandlers(createProductUseCase, listProductsUseCase)

	r := chi.NewRouter()
	r.Post("/products", productHandlers.CreateProductHandler)
	r.Get("/products", productHandlers.ListProductsHandler)

	go http.ListenAndServe(":8000", r)

	msgChan := make(chan *kafka.Message)
	// go akafka.Consume([]string{"products"}, "host.docker.internal:9094", msgChan)
	go akafka.Consume([]string{"products"}, "localhost:9093", msgChan)

	for msg := range msgChan {
		dto := usecase.CreateProductInputDto{}
		err := json.Unmarshal(msg.Value, &dto)
		if err != nil {
			// logar o erro
			fmt.Println(err)
		}
		_, err = createProductUseCase.Execute(dto)
		if err != nil {
			fmt.Println(err)
		}
	}

}

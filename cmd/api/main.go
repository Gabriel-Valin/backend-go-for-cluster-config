package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	ctx := context.Background()

	// conexões vêm de variáveis de ambiente (injetadas via Secret/ConfigMap)
	db, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer db.Close()

	// cria a tabela na subida (simplificação para o exemplo)
	db.Exec(ctx, `CREATE TABLE IF NOT EXISTS pedidos (id serial primary key, item text)`)

	rabbit, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("rabbitmq: %v", err)
	}
	defer rabbit.Close()
	ch, _ := rabbit.Channel()
	defer ch.Close()
	ch.QueueDeclare("pedidos", true, false, false, false, nil)

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/pedidos", func(w http.ResponseWriter, r *http.Request) {
		item := r.URL.Query().Get("item")
		if item == "" {
			item = "item-sem-nome"
		}

		// 1. grava no Postgres
		if _, err := db.Exec(ctx, "INSERT INTO pedidos (item) VALUES ($1)", item); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// 2. publica no RabbitMQ
		ch.PublishWithContext(ctx, "", "pedidos", false, false,
			amqp.Publishing{ContentType: "text/plain", Body: []byte(item)})

		json.NewEncoder(w).Encode(map[string]string{"status": "criado", "item": item})
	})

	log.Println("api ouvindo na :8080")
	http.ListenAndServe(":8080", nil)
}

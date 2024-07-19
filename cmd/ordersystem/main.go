package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gilbertom/desafio-clean-architecture/configs"
	"github.com/gilbertom/desafio-clean-architecture/internal/event/handler"
	"github.com/gilbertom/desafio-clean-architecture/internal/infra/graph"
	"github.com/gilbertom/desafio-clean-architecture/internal/infra/grpc/pb"
	"github.com/gilbertom/desafio-clean-architecture/internal/infra/grpc/service"
	"github.com/gilbertom/desafio-clean-architecture/internal/infra/web/webserver"
	"github.com/gilbertom/desafio-clean-architecture/pkg/events"
	"github.com/go-chi/chi/middleware"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel(configs.MQHost, configs.MQPort, configs.MQUser, configs.MQPassword)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	queryOrderUseCase := NewQueryOrderUseCase(db)

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webserver.Router.Use(middleware.Logger)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webQueryOrderHandler := NewWebQueryOrderHandler(db)
	webserver.AddHandler("/order", "POST", webOrderHandler.Create)
	webserver.AddHandler("/order", "GET", webQueryOrderHandler.Query)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *queryOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		QueryOrderUseCase:  *queryOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}


func getRabbitMQChannel(host, port, user, pass string) *amqp.Channel {
    addr := fmt.Sprintf("amqp://%s:%s@%s:%s", user, pass, host, port)
    log.Println("Connecting to RabbitMQ at", addr)
    var conn *amqp.Connection
    var err error

    for i := 0; i < 5; i++ {
        conn, err = amqp.Dial(addr)
        if err == nil {
            break
        }
        log.Printf("Failed to connect to RabbitMQ, retrying in 5 seconds... (%d/5)\n", i+1)
        time.Sleep(10 * time.Second)
    }

    if err != nil {
        panic(fmt.Sprintf("Failed to connect to RabbitMQ after 5 attempts: %v", err))
    }

    ch, err := conn.Channel()
    if err != nil {
        panic(fmt.Sprintf("Failed to open a channel: %v", err))
    }

		err = setupRabbitMQ(conn)
		if err != nil {
			panic(fmt.Sprintf("Failed to setup RabbitMQ: %v", err))
		}

		return ch
}

func setupRabbitMQ(conn *amqp.Connection) error {
    channel, err := conn.Channel()
    if err != nil {
        return fmt.Errorf("failed to open a channel: %v", err)
    }
    defer channel.Close()

    // Declare the exchange
    err = channel.ExchangeDeclare(
        "amq.direct", // name
        "direct",     // type
        true,         // durable
        false,        // auto-deleted
        false,        // internal
        false,        // no-wait
        nil,          // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare exchange: %v", err)
    }

    // Declare the queue
    _, err = channel.QueueDeclare(
        "order", // name
        true,    // durable
        false,   // delete when unused
        false,   // exclusive
        false,   // no-wait
        nil,     // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare queue: %v", err)
    }

    // Bind the queue to the exchange
    err = channel.QueueBind(
        "order",      // queue name
        "order",      // routing key
        "amq.direct", // exchange
        false,
        nil,
    )
    if err != nil {
        return fmt.Errorf("failed to bind queue: %v", err)
    }

    return nil
}
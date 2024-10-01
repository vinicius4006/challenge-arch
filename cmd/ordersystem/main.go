package main

import (
	"challenge-arch/configs"
	"challenge-arch/internal/event"
	"challenge-arch/internal/event/handler"
	"challenge-arch/internal/infra/database"
	"challenge-arch/internal/infra/graph"
	"challenge-arch/internal/infra/grpc/pb"
	"challenge-arch/internal/infra/grpc/service"
	"challenge-arch/internal/infra/web"
	"challenge-arch/internal/infra/web/webserver"
	"challenge-arch/internal/usecase"
	"challenge-arch/pkg/events"
	"database/sql"
	"fmt"
	"net"
	"net/http"

	handlerGraphql "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName)
	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})
	orderRepository := database.NewOrderRepository(db)
	orderCreated := event.NewOrderCreated()

	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository, orderCreated, eventDispatcher)
	listOrderUseCase := usecase.NewListOrdersUseCase(orderRepository, orderCreated, eventDispatcher)

	webserver := webserver.NewWebServer(configs.WebServerPort)

	webOrderHandler := web.NewWebOrderHandler(eventDispatcher, orderRepository, orderCreated)
	webserver.AddHandler("/order", webOrderHandler.Create)
	webserver.AddHandler("/list", webOrderHandler.List)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := handlerGraphql.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}

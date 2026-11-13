package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	tgrpc "dmitrii.com/car-service/grpc"
	httpw "dmitrii.com/car-service/http"
	"dmitrii.com/car-service/implementation"
	"dmitrii.com/car-service/pb"
	"dmitrii.com/car-service/repository"
	"dmitrii.com/car-service/transport"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	httpAddr := flag.String("http.addr", ":8080", "HTTP listen address")
	flag.Parse()
	// Настройка логгера
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", "car",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}
	_ = level.Info(logger).Log("msg", "service started")
	defer func() {
		_ = level.Info(logger).Log("msg", "service ended")
	}()
	// Создание сервиса
	repo := repository.NewMemoryRepository()
	svc := implementation.NewService(repo, logger)
	// Создание HTTP-обработчика
	endpoints := transport.MakeEndpoints(svc)
	handler := httpw.NewService(endpoints, logger)
	grpcServer := tgrpc.NewGRPCServer(endpoints)
	// Запуск сервера
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	go func() {
		_ = level.Info(logger).Log("transport", "HTTP", "addr",
			*httpAddr)
		server := &http.Server{
			Addr:    *httpAddr,
			Handler: handler,
		}
		errs <- server.ListenAndServe()
	}()

	go func() {
		listener, err := net.Listen("tcp", ":9090")
		if err != nil {
			errs <- err
			return
		}
		baseServer := grpc.NewServer()
		pb.RegisterCarServiceServer(baseServer, grpcServer)
		reflection.Register(baseServer) // Для удобной работы grpcurl

		fmt.Println("gRPC сервер запущен на :9090")
		errs <- baseServer.Serve(listener)
	}()
	_ = level.Error(logger).Log("exit", <-errs)
}

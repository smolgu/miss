package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/smolgu/miss/models"
	"google.golang.org/grpc"
)

func main() {
	go runGRPCServer()
	time.Sleep(1 * time.Second)
	runProxy()
}

func runGRPCServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 4455))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	models.RegisterLoveServer(grpcServer, &server{})
	grpcServer.Serve(lis)
}

func runProxy() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	endpoint := "localhost:4455"
	err := models.RegisterLoveHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		log.Fatalf("err: %s", err)
		return
	}

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}

// Server API for Love service

type server struct {
}

// Sends a greeting
func (server) VkAuth(context.Context, *models.VkAuthRequest) (*models.VkAuthReply, error) {

	return nil, nil
}

// User return user info by their id
func (server) User(context.Context, *models.UserRequest) (*models.UserReply, error) {
	return nil, nil
}

// User return list of users
func (server) RandomUsers(context.Context, *models.RandomRequest) (*models.UsersReply, error) {
	return nil, nil
}

// Vote vote for user
func (server) Vote(context.Context, *models.VoteRequest) (*models.VoteReply, error) {
	return nil, nil
}

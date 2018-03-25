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
	"github.com/smolgu/miss/pkg/vk"
	"google.golang.org/grpc"
)

func main() {

	err := models.NewContext()
	if err != nil {
		log.Fatalf("err models.NewContext: %s", err)
	}

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
func (server) VkAuth(_ context.Context, req *models.VkAuthRequest) (*models.VkAuthReply, error) {
	vkAccessToken := req.GetVkToken()
	vkID, err := vk.CheckToken(vkAccessToken)
	if err != nil {
		return nil, err
	}
	user, err := models.Users.GetByVkID(vkID)
	if err != nil {
		return nil, err
	}
	sessionID, err := models.Sessions.New(user.GetId())
	if err != nil {
		return nil, err
	}
	return &models.VkAuthReply{
		Token: sessionID,
	}, nil
}

// User return user info by their id
func (server) GetUser(ctx context.Context, req *models.UserRequest) (*models.User, error) {
	userID := req.GetUserId()
	return models.Users.Get(userID)
}

// User return list of users
func (server) RandomUsers(_ context.Context, req *models.RandomRequest) (*models.UsersReply, error) {
	userID, err := models.Sessions.Check(req.Token)
	if err != nil {
		return nil, err
	}
	users, err := models.Users.Random(userID)
	if err != nil {
		return nil, err
	}
	return &models.UsersReply{
		Users: users,
	}, nil
}

// Vote vote for user
func (server) Vote(_ context.Context, req *models.VoteRequest) (*models.VoteReply, error) {
	userID, err := models.Sessions.Check(req.Token)
	if err != nil {
		return nil, err
	}
	matched, err := models.Votes.Vote(req.GetUserId(), req.VoteType, userID)
	if err != nil {
		return nil, err
	}
	return &models.VoteReply{
		Matched: matched,
	}, nil
}

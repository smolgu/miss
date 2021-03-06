package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/smolgu/miss/cmd"
	"github.com/smolgu/miss/models"
	"github.com/smolgu/miss/pkg/errors"
	"github.com/smolgu/miss/pkg/setting"
	"github.com/smolgu/miss/pkg/vk"

	"github.com/go-playground/validator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func main() {
	app := &cli.App{
		Action: run,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name: "dev",
			},
			cli.StringFlag{
				Name:  "config",
				Value: "conf/app.yaml",
			},
		},
		Commands: []cli.Command{
			cmd.Photos,
			cmd.Profile,
		},
	}

	app.Run(os.Args)
}

func run(ctx *cli.Context) {

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	setting.Dev = ctx.Bool("dev")

	err := setting.NewContext(ctx.String("config"))
	if err != nil {
		log.Fatalf("err setting.NewContext: %s", err)
	}

	err = models.NewContext()
	if err != nil {
		log.Fatalf("err models.NewContext: %s", err)
	}

	go runGRPCServer()
	time.Sleep(1 * time.Second)
	runProxy()
}

func runGRPCServer() {
	log.Printf("listn grpc at :4455")
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

	log.Printf("start listen proxy at :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}

// Server API for Love service

type server struct {
}

func (srv server) vkAuth(ctx context.Context, req *models.VkAuthRequest, user *models.User) (*models.VkAuthReply, error) {
	sessionID, err := models.Sessions.New(user.GetId())
	if err != nil {
		return nil, err
	}
	return &models.VkAuthReply{
		Token: sessionID,
		User:  user,
	}, nil
}

// Sends a greeting
func (srv server) VkAuth(ctx context.Context, req *models.VkAuthRequest) (*models.VkAuthReply, error) {
	vkAccessToken := req.GetVkToken()
	vkID, err := vk.CheckToken(vkAccessToken)
	if err != nil {
		return nil, err
	}
	user, err := models.Users.GetByVkID(vkID)
	if err != nil {
		if err == errors.ErrNotFound {
			return srv.VkRegistre(ctx, req)
		}
		return nil, err
	}
	return srv.vkAuth(ctx, req, user)
}

func (srv server) VkRegistre(ctx context.Context, req *models.VkAuthRequest) (*models.VkAuthReply, error) {
	vkUser, err := vk.GetUser(req.GetVkToken(), 0)
	if err != nil {
		return nil, err
	}
	user, err := models.Users.CreateByVKUser(vkUser)
	if err != nil {
		return nil, err
	}
	return srv.vkAuth(ctx, req, user)
}

// GetUser return user info by their id
func (server) GetUser(ctx context.Context, req *models.UserRequest) (*models.User, error) {
	userID := req.GetUserId()
	return models.Users.Get(userID)
}

// RandomUsers return list of users
func (server) RandomUsers(_ context.Context, req *models.RandomRequest) (*models.UsersReply, error) {
	validate := validator.New()
	errs := validate.Struct(req)
	if errs != nil {
		return nil, errors.FormatValidatorError(errs)
	}
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

func (server) GetDialogs(_ context.Context, req *models.DialogsRequest) (*models.DialogsReply, error) {
	userID, err := models.Sessions.Check(req.Token)
	if err != nil {
		log.Printf("session check %v", errors.Sprint(err))
		return nil, err
	}
	dialogs, err := models.Dialogs.Dialogs(userID)
	if err != nil {
		return nil, err
	}
	return &models.DialogsReply{
		Dialogs: dialogs,
	}, nil
}

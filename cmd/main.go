package main

import (
	"blog-api/internal/config"
	"blog-api/internal/helper"
	"blog-api/internal/proto"
	"blog-api/internal/repository"
	"blog-api/internal/service"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	//init config
	cfg := config.InitConfig()
	defer cfg.Free()

	//init repository
	blogRepo := repository.NewBlog(cfg.DB, cfg.Env)

	//init helper
	s3Helper := helper.NewS3(cfg.AwsSession, cfg.Env)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", cfg.Env.Port))
	if err != nil {
		log.Fatal(err)
		return
	}
	s := grpc.NewServer()

	//init service
	fmt.Printf("start on port: %v\n", cfg.Env.Port)
	proto.RegisterBlogServiceServer(s, service.NewBlog(s3Helper, blogRepo))
	err = s.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}

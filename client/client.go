package main

import (
	"blog-api/internal/proto"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()
	c := proto.NewBlogServiceClient(cc)
	//create(c)
	download(c)
}

func download(c proto.BlogServiceClient) {
	resp, err := c.GetContentByTitle(context.TODO(), &proto.ContentByTitleReq{
		Title: "My golang story",
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	if resp == nil {
		log.Fatal(errors.New("response is nil"))
		return
	}
	for _, val := range resp.Content {
		fmt.Printf("Message: %v", val.Message)
		fmt.Printf("filename: %v\n", val.Filename)
		if len(val.Filename) > 0 && val.Image != nil {
			f, err := os.Create(val.Filename)
			if err != nil {
				log.Fatal(err)
				return
			}
			_, err = f.Write(val.Image)
			if err != nil {
				_ = f.Close()
				log.Fatal(err)
				return
			}
			_ = f.Close()
		}
	}

}

func create(c proto.BlogServiceClient) {
	data, err := ioutil.ReadFile("./golang-pic.png")
	if err != nil {
		log.Fatal(err)
	}
	var contents []*proto.Content
	contents = append(contents, &proto.Content{
		Message:  "<h1> let check my picture </h1>",
		Filename: "golang-pic.png",
		Image:    data,
	})

	resp, err := c.Create(context.TODO(), &proto.CreateReq{
		Title:   "My golang story",
		Content: contents,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}

package service

import (
	"blog-api/internal/helper"
	"blog-api/internal/proto"
	"blog-api/internal/repository"
	"context"
	"errors"
	"fmt"
)

const (
	contentIsRequired = "contents is required"
)

type blog struct {
	s3Helper helper.AwsS3
	blogRepo repository.Blog
}

func NewBlog(s3Helper helper.AwsS3, blogRepo repository.Blog) *blog {
	return &blog{
		s3Helper: s3Helper,
		blogRepo: blogRepo,
	}
}

func (s *blog) Create(_ context.Context, req *proto.CreateReq) (*proto.CreateResp, error) {
	entities, images, err := validateCreateReq(req)
	if err != nil {
		fmt.Printf("validateCreateReq: %v", err)
		return nil, err
	}

	//upload image
	err = s.uploadImage(images)
	if err != nil {
		fmt.Printf("s.uploadImgae: %v", err)
		return nil, err
	}

	//start create content
	err = s.blogRepo.Create(entities)
	if err != nil {
		fmt.Printf("blogRepo.Create: %v", err)
		return nil, err
	}

	return &proto.CreateResp{
		Result: &proto.CommonResult{
			Status: "OK",
		},
	}, nil
}

func validateCreateReq(req *proto.CreateReq) (*repository.BlogEntity, map[string][]byte, error) {
	if req == nil {
		return nil, nil, errors.New(contentIsRequired)
	}
	if len(req.Title) == 0 {
		return nil, nil, errors.New("title is required")
	}

	var r = make(map[string][]byte)
	var c []repository.Content
	for i, val := range req.Content {
		if len(val.Message) == 0 {
			return nil, nil, fmt.Errorf("message at row: %v is reqruied", i+1)
		}
		content := repository.Content{Message: val.Message}
		if len(val.Image) > 0 {
			r[val.Filename] = val.Image
			content.Filename = val.Filename
		}
		c = append(c, content)
	}

	return &repository.BlogEntity{
		Title:    req.Title,
		Contents: c,
	}, r, nil
}

func (s *blog) uploadImage(images map[string][]byte) error {
	if len(images) == 0 {
		//no image found
		return nil
	}
	//upload file to s3
	for filename, b := range images {
		err := s.s3Helper.Upload(filename, b)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *blog) GetContentByTitle(_ context.Context, req *proto.ContentByTitleReq) (*proto.ContentByTitleResp, error) {
	if len(req.Title) == 0 {
		return nil, errors.New("title is required")
	}
	entity, err := s.blogRepo.FindByTitle(req.Title)
	if err != nil {
		return nil, err
	}
	//case not found
	if entity == nil {
		return nil, fmt.Errorf("titile %v is not found", req.Title)
	}

	//download file from AwsS3
	var contents []*proto.Content
	for _, val := range entity.Contents {
		var tmp proto.Content
		if len(val.Filename) > 0 {
			b, err := s.s3Helper.Download(val.Filename)
			if err != nil {
				return nil, err
			}
			tmp.Filename = val.Filename
			tmp.Image = b
		}
		tmp.Message = val.Message
		contents = append(contents, &tmp)
	}

	return &proto.ContentByTitleResp{
		Title:   req.Title,
		Content: contents,
	}, nil
}

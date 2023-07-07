package service

import (
	"errors"
	"fmt"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"github.com/zhayt/clean-arch-tmp-forum/internal/repository"
	"github.com/zhayt/clean-arch-tmp-forum/logger"
	"log"
	"strconv"
	"strings"
)

type PostService struct {
	repo repository.Post
	l    *logger.Logger
}

func NewPostService(repo repository.Post, l *logger.Logger) *PostService {
	return &PostService{
		repo: repo,
		l:    l,
	}
}

var InvalidData = errors.New("invalid date")

func (s *PostService) CreatePost(post model.Post) error {
	post.Title = strings.TrimSpace(post.Title)
	post.Description = strings.TrimSpace(post.Description)

	if post.Title == "" || len(post.Title) >= 40 {
		s.l.Error.Printf("Title:", post.Title)
		return fmt.Errorf("title error: %w", InvalidData)
	}

	if post.Description == "" {
		s.l.Error.Printf("Content:", post.Description)
		return fmt.Errorf("content error: %w", InvalidData)
	}

	if len(post.Category) == 0 {
		s.l.Error.Printf("Category:", post.Title)
		return fmt.Errorf("empty category error: %w", InvalidData)
	}

	return s.repo.CreatePost(post)
}

func (s *PostService) ShowAllPosts() ([]model.Post, error) {
	posts, err := s.repo.ShowAllPosts()
	if err != nil {
		return []model.Post{}, err
	}

	return posts, nil
}

func (s *PostService) GetPostByID(id string) (*model.Post, error) {
	_, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	post, err := s.repo.GetPostByID(id)
	if err != nil {
		fmt.Printf("service: %s", err)
		return nil, err
	}
	return post, nil
}

func (s *PostService) ChangePost(newpost, oldPost model.Post, user model.User) error {
	if user.Username != oldPost.Author {
		return fmt.Errorf("Uncorrect change post author ")
	}

	newpost.Title = strings.TrimSpace(newpost.Title)
	newpost.Description = strings.TrimSpace(newpost.Description)
	if newpost.Title == "" && len(newpost.Title) >= 40 {
		return fmt.Errorf("Title invalid: %w", InvalidData)
	}

	if newpost.Description == "" {
		return fmt.Errorf("Description invalid: %w", InvalidData)
	}

	if err := s.repo.ChangePost(newpost, oldPost.Id); err != nil {
		return fmt.Errorf("CHANGE :%w", err)
	}

	return nil
}

func (s *PostService) ShowMyPosts(userId int) ([]model.Post, error) {
	posts, err := s.repo.ShowMyPosts(userId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return posts, nil
}

func (s *PostService) ShowMyCommentPosts(userId int) ([]model.Post, error) {
	posts, err := s.repo.ShowMyCommentPosts(userId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return posts, nil
}

func (s *PostService) ShowMyLikedPosts(userId int) ([]model.Post, error) {
	posts, err := s.repo.ShowMyLikedPosts(userId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return posts, nil
}

func (s *PostService) GetPostsByCategory(categoryID int) ([]model.Post, error) {
	if categoryID < 1 && categoryID > 5 {
		return []model.Post{}, InvalidData
	}

	posts, err := s.repo.GetPostsByCategory(categoryID)
	if err != nil {
		return posts, err
	}

	return posts, nil
}

func (s *PostService) DeletePost(userID int, postID int) error {
	if err := s.repo.DeletePost(userID, postID); err != nil {
		return err
	}

	return nil
}

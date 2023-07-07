package service

import (
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"github.com/zhayt/clean-arch-tmp-forum/internal/repository"
	"github.com/zhayt/clean-arch-tmp-forum/logger"
)

type Authorization interface {
	CreateUser(model.User) error
	GenerateToken(login, password string) (model.User, error)
	GetUserByToken(string) (model.User, error)
	DeleteToken(string) error
}

type Post interface {
	CreatePost(model.Post) error
	ShowAllPosts() ([]model.Post, error)
	GetPostByID(id string) (*model.Post, error)
	GetPostsByCategory(categoryID int) ([]model.Post, error)
	ChangePost(newPost, oldPost model.Post, user model.User) error
	ShowMyPosts(userId int) ([]model.Post, error)
	ShowMyCommentPosts(userId int) ([]model.Post, error)
	ShowMyLikedPosts(userId int) ([]model.Post, error)
	DeletePost(userID int, postID int) error
}
type Comment interface {
	CreateComment(model.Comment) error
	GetCommentByPostID(int) (*[]model.Comment, error)
}

type Like interface {
	SetPostLike(model.Like) error
	SetCommentLike(model.Like) error
}

type Dislike interface {
	SetPostDislike(model.Dislike) error
	SetCommentDislike(model.Dislike) error
}

type Service struct {
	Authorization
	Post
	Comment
	Like
	Dislike
}

func NewService(repos *repository.Repository, l *logger.Logger) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, l),
		Post:          NewPostService(repos.Post, l),
		Comment:       NewCommentService(repos.Comment, l),
		Like:          NewLikeService(repos.Like, l),
		Dislike:       NewDislikeService(repos.Dislike, l),
	}
}

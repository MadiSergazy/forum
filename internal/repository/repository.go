package repository

import (
	"github.com/zhayt/clean-arch-tmp-forum/config"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"github.com/zhayt/clean-arch-tmp-forum/internal/repository/sqlite"
	"github.com/zhayt/clean-arch-tmp-forum/logger"
	"time"
)

type Authorization interface {
	CreateUser(model.User) error
	GetUser(string) (model.User, error)
	GetUserByToken(string) (model.User, error)
	SaveTokens(string, time.Time, string) error
	DeleteToken(string) error
}

type Post interface {
	CreatePost(model.Post) error
	ShowAllPosts() ([]model.Post, error)
	GetPostByID(string) (*model.Post, error)
	GetPostsByCategory(int) ([]model.Post, error)
	ChangePost(model.Post, int) error
	ShowMyPosts(int) ([]model.Post, error)
	ShowMyCommentPosts(int) ([]model.Post, error)
	ShowMyLikedPosts(int) ([]model.Post, error)
	DeletePost(int, int) error
}

type Comment interface {
	CreateComment(model.Comment) error
	GetCommentByPostID(int) (*[]model.Comment, error)
}

type Dislike interface {
	SetPostDislike(model.Dislike) error
	SetCommentDislike(model.Dislike) error
	LikeDislike
}

type Like interface {
	SetPostLike(model.Like) error
	SetCommentLike(model.Like) error
	LikeDislike
}

type LikeDislike interface {
	CheckPostLike(int, int) error
	DeletePostLike(int, int) error
	CheckPostDislike(int, int) error
	DeletePostDislike(int, int) error
	CheckCommentLike(int, int) error
	DeleteCommentLike(int, int) error
	CheckCommentDislike(int, int) error
	DeleteCommentDislike(int, int) error
}

type Repository struct {
	Authorization
	Post
	Comment
	Like
	Dislike
}

func NewRepository(cfg *config.Config, l *logger.Logger) (*Repository, error) {
	db, err := sqlite.Dial(cfg)
	if err != nil {
		return nil, err
	}

	if err = sqlite.InnitDB(db); err != nil {
		return nil, err
	}

	return &Repository{
		Authorization: sqlite.NewAuthRepo(db),
		Post:          sqlite.NewPostRepo(db, l),
		Comment:       sqlite.NewCommentRepo(db),
		Like:          sqlite.NewLikeRepo(db),
		Dislike:       sqlite.NewDislikeRepo(db),
	}, nil
}

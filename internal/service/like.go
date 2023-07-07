package service

import (
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"github.com/zhayt/clean-arch-tmp-forum/internal/repository"
	"github.com/zhayt/clean-arch-tmp-forum/logger"
)

type LikeService struct {
	repo repository.Like
	l    *logger.Logger
}

func NewLikeService(repo repository.Like, l *logger.Logger) *LikeService {
	return &LikeService{
		repo: repo,
		l:    l,
	}
}

func (l *LikeService) SetPostLike(like model.Like) error {
	if l.repo.CheckPostDislike(like.PostID, like.UserID) == nil {
		l.repo.DeletePostDislike(like.PostID, like.UserID)
	}
	if l.repo.CheckPostLike(like.PostID, like.UserID) != nil {
		return l.repo.SetPostLike(like)
	} else {
		return l.repo.DeletePostLike(like.PostID, like.UserID)
	}
}

func (l *LikeService) SetCommentLike(like model.Like) error {
	if l.repo.CheckCommentDislike(like.CommentId, like.UserID) == nil {
		l.repo.DeleteCommentDislike(like.CommentId, like.UserID)
	}
	if l.repo.CheckCommentLike(like.CommentId, like.UserID) != nil {
		return l.repo.SetCommentLike(like)
	} else {
		return l.repo.DeleteCommentLike(like.CommentId, like.UserID)
	}
}

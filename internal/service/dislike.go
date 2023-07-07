package service

import (
	models "github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"github.com/zhayt/clean-arch-tmp-forum/internal/repository"
	"github.com/zhayt/clean-arch-tmp-forum/logger"
)

type DislikeService struct {
	repo repository.Dislike
	l    *logger.Logger
}

func NewDislikeService(repo repository.Dislike, l *logger.Logger) *DislikeService {
	return &DislikeService{
		repo: repo,
		l:    l,
	}
}

func (d *DislikeService) SetPostDislike(dislike models.Dislike) error {
	if d.repo.CheckPostLike(dislike.PostID, dislike.UserID) == nil {
		d.repo.DeletePostLike(dislike.PostID, dislike.UserID)
	}
	if d.repo.CheckPostDislike(dislike.PostID, dislike.UserID) != nil {
		return d.repo.SetPostDislike(dislike)
	} else {
		return d.repo.DeletePostDislike(dislike.PostID, dislike.UserID)
	}
}

func (d *DislikeService) SetCommentDislike(dislike models.Dislike) error {
	if d.repo.CheckCommentLike(dislike.CommentId, dislike.UserID) == nil {
		d.repo.DeleteCommentLike(dislike.CommentId, dislike.UserID)
	}
	if d.repo.CheckCommentDislike(dislike.CommentId, dislike.UserID) != nil {
		return d.repo.SetCommentDislike(dislike)
	} else {
		return d.repo.DeleteCommentDislike(dislike.CommentId, dislike.UserID)
	}
}

package service

import (
	"fmt"
	models "github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"github.com/zhayt/clean-arch-tmp-forum/internal/repository"
	"github.com/zhayt/clean-arch-tmp-forum/logger"
	"strings"
)

type CommentService struct {
	repo repository.Comment
	l    *logger.Logger
}

func NewCommentService(repo repository.Comment, l *logger.Logger) *CommentService {
	return &CommentService{
		repo: repo,
		l:    l,
	}
}

func (c *CommentService) CreateComment(commnet models.Comment) error {
	commnet.Text = strings.TrimSpace(commnet.Text)
	if commnet.Text == "" {
		return fmt.Errorf("emty comment text: %w", InvalidData)
	}

	return c.repo.CreateComment(commnet)
}

func (c *CommentService) GetCommentByPostID(id int) (*[]models.Comment, error) {
	return c.repo.GetCommentByPostID(id)
}

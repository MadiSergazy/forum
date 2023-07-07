package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
)

type CommentRepo struct {
	db *sql.DB
}

func NewCommentRepo(db *sql.DB) *CommentRepo {
	return &CommentRepo{
		db: db,
	}
}

func (c *CommentRepo) CreateComment(comment model.Comment) error {
	query := `INSERT INTO comment (postId, userId, author, text, date) VALUES ($2, $3, $4, $5, $6)`

	if _, err := c.db.Exec(query, comment.PostId, comment.UserId, comment.Author, comment.Text, comment.Date); err != nil {
		return fmt.Errorf("couldn't create comment: %w", err)
	}

	return nil
}

func (c *CommentRepo) GetCommentByPostID(postid int) (*[]model.Comment, error) {
	query := `select * from comment where postId=$1`

	rows, err := c.db.Query(query, postid)
	if err != nil {
		return nil, fmt.Errorf("couldn't get comment by post id: %w", err)
	}

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		if err := rows.Scan(&comment.Id, &comment.PostId, &comment.UserId, &comment.Author, &comment.Text, &comment.Like, &comment.Dislike, &comment.Date); err != nil {
			return nil, fmt.Errorf("couldn't get comment by post id, scan error: %w", err)
		}

		comments = append(comments, comment)
	}

	return &comments, nil
}

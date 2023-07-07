package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
)

type LikeRepo struct {
	db *sql.DB
}

func NewLikeRepo(db *sql.DB) *LikeRepo {
	return &LikeRepo{
		db: db,
	}
}

func (r *LikeRepo) SetPostLike(like model.Like) error {
	query := `INSERT INTO like(postId, userId, active) VALUES(?,?,?) `

	if _, err := r.db.Exec(query, like.PostID, like.UserID, like.Active); err != nil {
		return fmt.Errorf("couldn't set post like: %w", err)
	}

	query = `update post set like=(select count(*) from like where postId=?) where id=?`
	if _, err = r.db.Exec(query, like.PostID, like.PostID); err != nil {
		return fmt.Errorf("couldn't increase post like: %w", err)
	}

	return nil
}

func (r *LikeRepo) CheckPostLike(PostID, UserID int) error {
	query := `select id from like where userId = ? and postId = ?`

	row := r.db.QueryRow(query, UserID, PostID)

	var likeID int

	if err = row.Scan(&likeID); err != nil {
		return fmt.Errorf("couldn't check post like: %w", err)
	}

	return nil
}

func (r *LikeRepo) DeletePostLike(PostID, UserID int) error {
	query := `delete from like where userId = ? and postId = ?`

	if _, err = r.db.Exec(query, UserID, PostID); err != nil {
		return fmt.Errorf("couldn't delete post like: %w", err)
	}

	query = `update post set like=(select count(*) from like where postId=?) where id=?`

	if _, err = r.db.Exec(query, PostID, PostID); err != nil {
		return fmt.Errorf("couldn't decrease post like: %w", err)
	}

	return nil
}

func (r *LikeRepo) CheckPostDislike(PostID, UserID int) error {
	query := `select id from dislike where userId = ? and postId = ?`

	row := r.db.QueryRow(query, UserID, PostID)

	var dislikeID int

	if err = row.Scan(&dislikeID); err != nil {
		return fmt.Errorf("couldn't check post like: %w", err)
	}

	return nil
}

func (r *LikeRepo) DeletePostDislike(PostID, UserID int) error {
	query := `delete from dislike where userId = ? and postId = ?`

	if _, err = r.db.Exec(query, UserID, PostID); err != nil {
		return fmt.Errorf("couldn't delete post dislike: %w", err)
	}

	query = `update post set dislike=(select count(*) from dislike where postId=?) where id=?`

	if _, err = r.db.Exec(query, PostID, PostID); err != nil {
		return fmt.Errorf("couldn't decrease post like: %w", err)
	}

	return nil
}

// comment
func (r *LikeRepo) SetCommentLike(like model.Like) error {
	query := `INSERT INTO like(commentId,userId,active) VALUES(?, ?, ?) `

	if _, err := r.db.Exec(query, like.CommentId, like.UserID, like.Active); err != nil {
		return fmt.Errorf("couldn't set comment like: %w", err)
	}

	query = `update comment set like=(select count(*) from like where commentId=?) where id=?`

	if _, err = r.db.Exec(query, like.CommentId, like.CommentId); err != nil {
		return fmt.Errorf("couldn't increase comment like: %w", err)
	}

	return nil
}

func (r *LikeRepo) CheckCommentLike(CommentID, UserID int) error {
	query := `select id from like where userId = ? and commentId = ?`

	row := r.db.QueryRow(query, UserID, CommentID)

	var likeID int

	if err = row.Scan(&likeID); err != nil {
		return fmt.Errorf("couldn't check comment like: %w", err)
	}

	return nil
}

func (r *LikeRepo) DeleteCommentLike(CommentID, UserID int) error {
	query := `delete from like where userId = ? and commentId = ?`

	if _, err = r.db.Exec(query, UserID, CommentID); err != nil {
		return fmt.Errorf("couldn't delete comment like: %w", err)
	}

	query = `update comment set like=(select count(*) from like where commentId=?) where id=?`

	if _, err = r.db.Exec(query, CommentID, CommentID); err != nil {
		return fmt.Errorf("couldn't decrease comment like: %w", err)
	}

	return nil
}

func (r *LikeRepo) CheckCommentDislike(CommentID, UserID int) error {
	query := `select id from dislike where userId = ? and commentId = ?`

	row := r.db.QueryRow(query, UserID, CommentID)

	var likeID int

	if err = row.Scan(&likeID); err != nil {
		return fmt.Errorf("couldn't check comment dislike: %w", err)
	}

	return nil
}

func (r *LikeRepo) DeleteCommentDislike(CommentID, UserID int) error {
	query := `delete from dislike where userId = ? and commentId = ?`

	if _, err = r.db.Exec(query, UserID, CommentID); err != nil {
		return fmt.Errorf("couldn't delete comment dislike: %w", err)
	}

	query = `update comment set dislike=(select count(*) from dislike where commentId=?) where id=?`
	if _, err = r.db.Exec(query, CommentID, CommentID); err != nil {
		return fmt.Errorf("couldn't decrease comment dislike: %w", err)
	}

	return nil
}

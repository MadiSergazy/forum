package sqlite

import (
	"database/sql"
	"fmt"
	models "github.com/zhayt/clean-arch-tmp-forum/internal/model"
)

type DislikeRepo struct {
	db *sql.DB
}

func NewDislikeRepo(db *sql.DB) *DislikeRepo {
	return &DislikeRepo{
		db: db,
	}
}

func (r *DislikeRepo) SetPostDislike(dislike models.Dislike) error {
	query := `INSERT INTO dislike(postId, userId, active) VALUES(?,?,?) `

	if _, err = r.db.Exec(query, dislike.PostID, dislike.UserID, dislike.Active); err != nil {
		return fmt.Errorf("couldn't set post dislike: %w", err)
	}

	query = `update post set dislike=(select count(*) from dislike where postId=?) where id=?`
	if _, err := r.db.Exec(query, dislike.PostID, dislike.PostID); err != nil {
		return fmt.Errorf("couldn't increase post dislike: %w", err)
	}

	return nil
}

func (r *DislikeRepo) CheckPostDislike(PostID, UserID int) error {
	query := `select id from dislike where userId = ? and postId = ?`

	row := r.db.QueryRow(query, UserID, PostID)

	var dislikeID int

	if err = row.Scan(&dislikeID); err != nil {
		return fmt.Errorf("couldn't check post dislike: %w", err)
	}

	return nil
}

func (r *DislikeRepo) DeletePostDislike(PostID, UserID int) error {
	query := `delete from dislike where userId = ? and postId = ?`

	if _, err = r.db.Exec(query, UserID, PostID); err != nil {
		return fmt.Errorf("couldn't delete post dislike: %w", err)
	}

	query = `update post set dislike=(select count(*) from dislike where postId=?) where id=?`

	if _, err = r.db.Exec(query, PostID, PostID); err != nil {
		return fmt.Errorf("couldn't decrease post dislike: %w", err)
	}

	return nil
}

func (r *DislikeRepo) CheckPostLike(PostID, UserID int) error {
	query := `select id from like where userId = ? and postId = ?`

	row := r.db.QueryRow(query, UserID, PostID)

	var likeID int

	if err = row.Scan(&likeID); err != nil {
		return fmt.Errorf("couldn't check post like: %w", err)
	}

	return nil
}

func (r *DislikeRepo) DeletePostLike(PostID, UserID int) error {
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

func (r *DislikeRepo) SetCommentDislike(dislike models.Dislike) error {
	query := `INSERT INTO dislike(commentId,userId,active) VALUES(?, ?, ?) `

	if _, err = r.db.Exec(query, dislike.CommentId, dislike.UserID, dislike.Active); err != nil {
		return fmt.Errorf("couldn't set comment like: %w", err)
	}

	query = `update comment set dislike=(select count(*) from dislike where commentId=?) where id=?`

	if _, err = r.db.Exec(query, dislike.CommentId, dislike.CommentId); err != nil {
		return fmt.Errorf("couldn't increase comment like: %w", err)
	}

	return nil
}

func (r *DislikeRepo) CheckCommentDislike(CommentID, UserID int) error {
	query := `select id from dislike where userId = ? and commentId = ?`

	row := r.db.QueryRow(query, UserID, CommentID)

	var likeID int

	if err = row.Scan(&likeID); err != nil {
		return fmt.Errorf("couldn't check comment dislike: %w", err)
	}

	return nil
}

func (r *DislikeRepo) DeleteCommentDislike(CommentID, UserID int) error {
	query := `delete from dislike where userId = ? and commentId = ?`

	if _, err := r.db.Exec(query, UserID, CommentID); err != nil {
		return fmt.Errorf("couldn't delete post like: %w", err)
	}

	query = `update comment set dislike=(select count(*) from dislike where commentId=?) where id=?`
	if _, err = r.db.Exec(query, CommentID, CommentID); err != nil {
		return fmt.Errorf("couldn't decrease comment like: %w", err)
	}

	return nil
}

func (r *DislikeRepo) CheckCommentLike(CommentID, UserID int) error {
	query := `select id from like where userId = ? and commentId = ?`

	row := r.db.QueryRow(query, UserID, CommentID)

	var likeID int

	if err = row.Scan(&likeID); err != nil {
		return fmt.Errorf("couldn't check comment like: %w", err)
	}

	return nil
}

func (r *DislikeRepo) DeleteCommentLike(CommentID, UserID int) error {
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

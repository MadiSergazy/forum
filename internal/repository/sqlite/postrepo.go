package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"github.com/zhayt/clean-arch-tmp-forum/logger"
	"strconv"
	"strings"
)

var (
	rows *sql.Rows
	err  error
)

type PostRepo struct {
	db *sql.DB
	l  *logger.Logger
}

func NewPostRepo(db *sql.DB, l *logger.Logger) *PostRepo {
	return &PostRepo{
		db: db,
		l:  l,
	}
}

func (r *PostRepo) CreatePost(post model.Post) error {
	r.l.Info.Printf("Create post repo", post.Category)

	query := `INSERT INTO post (user_id, author, title, description, date) VALUES ($1, $2, $3, $4, $5)`

	res, err := r.db.Exec(query, post.AuthorId, post.Author, post.Title, post.Description, post.Date)
	if err != nil {
		return fmt.Errorf("couldn't create post: %w", err)
	}

	postID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("couldn't create post, get creted post id error: %w", err)
	}

	for _, category := range post.Category {
		categoryID, _ := strconv.Atoi(category)
		if _, err = r.db.Exec(`INSERT INTO post_category(post_id, category_id) VALUES (?, ?)`, postID, categoryID); err != nil {
			r.l.Error.Printf("couldn't create category for post: postID", postID, categoryID, err.Error())
		}
		r.l.Info.Println("Category for post added categoryID:", categoryID, "for post", postID)
	}

	return nil
}

func (r *PostRepo) ShowAllPosts() ([]model.Post, error) {
	rows, err := r.db.Query(`SELECT 
        post.id, post.title, post.description, post.like, post.dislike, post.user_id, post.author, 
        post.date, group_concat(category.title) AS category 
    FROM 
        post 
        INNER JOIN post_category ON post.id = post_category.post_id 
        INNER JOIN category ON category.id = post_category.category_id 
    GROUP BY 
        post.id`)
	if err != nil {
		return []model.Post{}, fmt.Errorf("couldn't get all post: %w", err)
	}

	var posts []model.Post
	for rows.Next() {
		post := new(model.Post)

		var category string

		if err = rows.Scan(&post.Id, &post.Title, &post.Description, &post.Like, &post.Dislike, &post.AuthorId, &post.Author, &post.Date, &category); err != nil {
			return nil, fmt.Errorf("couldn't get all post, scan error: %w", err)
		}

		r.l.Info.Println("Category:", category, "postID:", post.Id)
		post.Category = strings.Split(category, ",")
		r.l.Info.Println("Post", post)
		posts = append(posts, *post)
	}

	return posts, nil
}

func (r *PostRepo) GetPostByID(id string) (*model.Post, error) {
	rows, err := r.db.Query(`SELECT 
        post.id, post.title, post.description, post.like, post.dislike, post.user_id, post.author, 
        post.date, group_concat(category.title) AS category 
    FROM 
        post 
        INNER JOIN post_category ON post.id = post_category.post_id 
        INNER JOIN category ON category.id = post_category.category_id 
    GROUP BY 
        post.id
        HAVING post.id = ?`, id)
	if err != nil {
		return nil, fmt.Errorf("couldn't get post by id: %w", err)
	}

	var post model.Post

	var category string

	for rows.Next() {
		if err = rows.Scan(&post.Id, &post.Title, &post.Description, &post.Like, &post.Dislike, &post.AuthorId, &post.Author, &post.Date, &category); err != nil {
			return nil, fmt.Errorf("couldn't get all post, scan error: %w", err)
		}
	}

	post.Category = strings.Split(category, ",")

	return &post, nil
}

func (r *PostRepo) ChangePost(post model.Post, oldPostId int) error {
	query := `UPDATE post SET title=$1,description=$2 where id=$3;`

	if _, err := r.db.Exec(query, post.Title, post.Description, oldPostId); err != nil {
		return fmt.Errorf("couldn't change post: %w", err)
	}

	return nil
}

func (r *PostRepo) ShowMyPosts(userId int) ([]model.Post, error) {
	rows, err := r.db.Query(`SELECT 
        post.id, post.title, post.description, post.like, post.dislike, post.user_id, post.author, 
        post.date, group_concat(category.title) AS category 
    FROM 
        post 
        INNER JOIN post_category ON post.id = post_category.post_id 
        INNER JOIN category ON category.id = post_category.category_id 
    GROUP BY 
        post.id
        HAVING post.user_id = ?`, userId)
	if err != nil {
		return []model.Post{}, fmt.Errorf("couldn't show user post: %w", err)
	}

	var posts []model.Post

	for rows.Next() {
		post := new(model.Post)

		var category string

		if err = rows.Scan(&post.Id, &post.Title, &post.Description, &post.Like, &post.Dislike, &post.AuthorId, &post.Author, &post.Date, &category); err != nil {
			return nil, fmt.Errorf("couldn't get all post, scan error: %w", err)
		}

		post.Category = strings.Split(category, ",")

		posts = append(posts, *post)
	}

	return posts, nil
}

func (r *PostRepo) ShowMyCommentPosts(userId int) ([]model.Post, error) {
	rows, err := r.db.Query(`SELECT 
        post.id, post.title, post.description, post.like, post.dislike, post.user_id, post.author, 
        post.date, group_concat(category.title) AS category 
    FROM 
        post 
        INNER JOIN post_category ON post.id = post_category.post_id 
        INNER JOIN category ON category.id = post_category.category_id 
    GROUP BY 
        post.id
        HAVING post.id in (select postId from comment where userId=?)`, userId)
	if err != nil {
		return []model.Post{}, fmt.Errorf("couldn't get user commented post: %w", err)
	}

	var posts []model.Post

	for rows.Next() {
		post := new(model.Post)

		var category string

		if err = rows.Scan(&post.Id, &post.Title, &post.Description, &post.Like, &post.Dislike, &post.AuthorId, &post.Author, &post.Date, &category); err != nil {
			return nil, fmt.Errorf("couldn't get all post, scan error: %w", err)
		}

		post.Category = strings.Split(category, ",")

		posts = append(posts, *post)
	}

	return posts, nil
}

func (r *PostRepo) ShowMyLikedPosts(userId int) ([]model.Post, error) {
	rows, err := r.db.Query(`SELECT 
        post.id, post.title, post.description, post.like, post.dislike, post.user_id, post.author, 
        post.date, group_concat(category.title) AS category 
    FROM 
        post 
        INNER JOIN post_category ON post.id = post_category.post_id 
        INNER JOIN category ON category.id = post_category.category_id 
    GROUP BY 
        post.id
        HAVING post.id in (select postId from like where userId=? and postId not null)`, userId)
	if err != nil {
		return []model.Post{}, fmt.Errorf("couldn't get user liked post: %w", err)
	}

	var posts []model.Post

	for rows.Next() {
		post := new(model.Post)

		var category string

		if err = rows.Scan(&post.Id, &post.Title, &post.Description, &post.Like, &post.Dislike, &post.AuthorId, &post.Author, &post.Date, &category); err != nil {
			return nil, fmt.Errorf("couldn't get all post, scan error: %w", err)
		}

		post.Category = strings.Split(category, ",")

		posts = append(posts, *post)
	}

	return posts, nil
}

func (r *PostRepo) GetPostsByCategory(categoryID int) ([]model.Post, error) {
	var posts []model.Post

	r.l.Info.Printf("CategoryID-", categoryID)

	rows, err := r.db.Query(`SELECT 
        post.id, post.title, post.description, post.like, post.dislike, post.user_id, post.author, 
        post.date, group_concat(category.title) AS category 
    FROM 
        post 
        INNER JOIN post_category ON post.id = post_category.post_id 
        INNER JOIN category ON category.id = post_category.category_id 
    WHERE category_id = ?
    GROUP BY 
        post.id`, categoryID)

	if err != nil {
		return []model.Post{}, fmt.Errorf("couldn't get post by category: %w", err)
	}

	for rows.Next() {
		var post model.Post

		var category string

		if err = rows.Scan(&post.Id, &post.Title, &post.Description, &post.Like, &post.Dislike, &post.AuthorId, &post.Author, &post.Date, &category); err != nil {
			return nil, fmt.Errorf("couldn't get all post, scan error: %w", err)
		}

		post.Category = strings.Split(category, ",")

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepo) DeletePost(userID int, postID int) error {
	query := `delete from post where  id = ? AND user_id = ?`

	if _, err = r.db.Exec(query, postID, userID); err != nil {
		return fmt.Errorf("couldn't delete post: %w", err)
	}

	return nil
}

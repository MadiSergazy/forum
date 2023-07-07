package handler

import (
	"errors"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"github.com/zhayt/clean-arch-tmp-forum/internal/service"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(model.CtxUserKey).(model.User)

	if user.Username == "" {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}

	switch r.Method {
	case http.MethodGet:
		temp, err := template.ParseFiles("ui/createpost.html")
		if err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err = temp.Execute(w, temp); err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		title := r.FormValue("title")
		category := r.Form["category"]
		content := r.FormValue("content")

		post := model.Post{
			Title:       title,
			Description: content,
			AuthorId:    user.Id,
			Author:      user.Username,
			Date:        time.Now().Format("January 2, 2006 15:04:05"),
			Category:    category,
		}

		err := h.service.Post.CreatePost(post)
		if err != nil {
			h.l.Error.Printf("Create post error: %s", err.Error())
			if errors.Is(err, service.InvalidData) {
				errorHandler(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		h.l.Info.Printf("Post created")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		errorHandler(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}

}

type DispPost struct {
	Username string
	Post     model.Post
	Comments []model.Comment
	Auth     bool
}

// Не понятно зачем нужен свич кейс
func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	PostID := r.URL.Query().Get("id")
	user := r.Context().Value(model.CtxUserKey).(model.User)

	_, err := strconv.Atoi(PostID)
	if err != nil {
		errorHandler(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	post, err := h.service.GetPostByID(PostID)

	if err != nil || post.Title == "" {
		errorHandler(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	dispPost := DispPost{
		Username: user.Username,
		Post:     *post,
		Comments: []model.Comment{},
	}

	dispPost.Auth = dispPost.Username == dispPost.Post.Author
	switch r.Method {
	case http.MethodGet:
		comments, err := h.service.Comment.GetCommentByPostID(post.Id)
		if err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		dispPost.Comments = *comments
	case http.MethodPost:
		if user.Username == "" {
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		if r.FormValue("postLike") != "" {
			postId, _ := strconv.Atoi(r.FormValue("postLike"))
			like := model.Like{
				UserID: user.Id,
				PostID: postId,
			}

			if err := h.service.Like.SetPostLike(like); err != nil {
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/post/?id="+PostID, http.StatusSeeOther)
		} else if r.FormValue("postDislike") != "" {
			postId, _ := strconv.Atoi(r.FormValue("postDislike"))

			dislike := model.Dislike{
				UserID: user.Id,
				PostID: postId,
			}

			if err = h.service.Dislike.SetPostDislike(dislike); err != nil {
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/post/?id="+PostID, http.StatusSeeOther)
			return
		}

		if r.FormValue("commentLike") != "" {
			commentId, _ := strconv.Atoi(r.FormValue("commentLike"))
			like := model.Like{
				UserID:    user.Id,
				CommentId: commentId,
			}

			if err := h.service.Like.SetCommentLike(like); err != nil {
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/post/?id="+PostID, http.StatusSeeOther)

		} else if r.FormValue("commentDislike") != "" {
			commentId, _ := strconv.Atoi(r.FormValue("commentDislike"))

			dislike := model.Dislike{
				UserID:    user.Id,
				CommentId: commentId,
			}

			if err = h.service.Dislike.SetCommentDislike(dislike); err != nil {
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/post/?id="+PostID, http.StatusSeeOther)
			return
		}

		if user.Token == "" {
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		commentText := r.FormValue("comment")

		comment := model.Comment{
			UserId: user.Id,
			PostId: post.Id,
			Text:   commentText,
			Author: user.Username,
			Date:   time.Now().Format("01-02-2006 15:04:05"),
		}

		if err = h.service.Comment.CreateComment(comment); err != nil {
			if errors.Is(err, service.InvalidData) {
				http.Redirect(w, r, "/post/?id="+PostID, http.StatusSeeOther)
				return
			}

			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/post/?id="+PostID, http.StatusSeeOther)

	default:
		errorHandler(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	temp, err := template.ParseFiles("ui/postpage.html")
	if err != nil {
		errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err = temp.Execute(w, dispPost); err != nil {
		errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ChangePost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(model.CtxUserKey).(model.User)
	if user.Username == "" {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}

	postId := r.URL.Query().Get("id")
	oldPost, err := h.service.GetPostByID(postId)
	if err != nil {
		errorHandler(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		temp, err := template.ParseFiles("ui/changepost.html")
		if err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err = temp.Execute(w, oldPost); err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		category := r.Form["category"]

		title := r.FormValue("title")

		description := r.FormValue("content")

		newPost := model.Post{
			Title:       title,
			Description: description,
			Category:    category,
		}

		if err := h.service.Post.ChangePost(newPost, *oldPost, user); err != nil {
			if errors.Is(err, service.InvalidData) {
				errorHandler(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/post/?id="+postId, http.StatusSeeOther)
	default:
		errorHandler(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorHandler(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	user := r.Context().Value(model.CtxUserKey).(model.User)
	if user.Username == "" {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}

	postID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = h.service.Post.DeletePost(user.Id, postID)
	if err != nil {
		h.l.Error.Println("Delete post error: userID:", user.Id, "postID:", postID, "error:", err.Error())
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

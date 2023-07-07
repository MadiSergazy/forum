package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"github.com/zhayt/clean-arch-tmp-forum/internal/service"
	"github.com/zhayt/clean-arch-tmp-forum/logger"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Handler struct {
	service *service.Service
	l       *logger.Logger
}

func NewHandler(s *service.Service, l *logger.Logger) *Handler {
	return &Handler{
		service: s,
		l:       l,
	}
}

type Display struct {
	Username string
	Posts    []model.Post
	Category []string
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	user := r.Context().Value(model.CtxUserKey).(model.User)

	switch r.Method {
	case http.MethodGet:
		temp, err := template.ParseFiles("ui/homepage.html")
		if err != nil {
			h.l.Error.Println("Parse file error:")
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var posts []model.Post

		posts, err = h.service.Post.ShowAllPosts()
		if err != nil {
			h.l.Error.Printf("Show all posts error: %s", err.Error())
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		h.l.Info.Println("All post founded, count %v", len(posts))

		display := Display{
			Username: user.Username,
			Posts:    posts,
		}

		if err = temp.Execute(w, display); err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		if r.FormValue("category") != "" {
			// fmt.Println(category)
			categoryID, err := strconv.Atoi(r.FormValue("category"))
			if err != nil {
				h.l.Error.Printf("convert category id error", err.Error())
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			posts, err := h.service.Post.GetPostsByCategory(categoryID)
			if err != nil {
				h.l.Error.Printf("Get posts by category error", err.Error())

				if errors.Is(err, service.InvalidData) {
					errorHandler(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}

				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			h.l.Info.Printf("Posts by category founded id", categoryID, "amount", len(posts))

			categories := []string{"IT", "Sport", "Education", "News", "Health"}
			display := Display{
				Username: user.Username,
				Posts:    posts,
				Category: []string{categories[categoryID-1]},
			}

			temp, err := template.ParseFiles("ui/homepage.html")
			if err != nil {
				h.l.Error.Printf("Parse file error", err.Error())
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			h.l.Info.Printf("To display struct", display)
			if err = temp.Execute(w, display); err != nil {
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			return
		}

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

			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		} else if r.FormValue("postDislike") != "" {
			postId, _ := strconv.Atoi(r.FormValue("postDislike"))
			dislike := model.Dislike{
				UserID: user.Id,
				PostID: postId,
			}
			if err := h.service.Dislike.SetPostDislike(dislike); err != nil {
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		} else {
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		}

	default:
		errorHandler(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) MyPosts(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(model.CtxUserKey).(model.User)

	if user.Username == "" {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
	}

	switch r.Method {
	case http.MethodGet:

		var posts []model.Post
		posts, err := h.service.Post.ShowMyPosts(user.Id)
		if err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		display := Display{
			Username: user.Username,
			Posts:    posts,
		}
		temp, err := template.ParseFiles("ui/myposts.html")
		if err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err = temp.Execute(w, display); err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
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
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		} else if r.FormValue("postDislike") != "" {
			postId, _ := strconv.Atoi(r.FormValue("postDislike"))
			dislike := model.Dislike{
				UserID: user.Id,
				PostID: postId,
			}
			if err := h.service.Dislike.SetPostDislike(dislike); err != nil {
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		}
	default:
		errorHandler(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) MyCommentPosts(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(model.CtxUserKey).(model.User)

	if user.Username == "" {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
	}
	switch r.Method {
	case http.MethodGet:
		var posts []model.Post

		posts, err := h.service.Post.ShowMyCommentPosts(user.Id)
		if err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		display := Display{
			Username: user.Username,
			Posts:    posts,
		}

		temp, err := template.ParseFiles("ui/mycommentposts.html")
		if err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err = temp.Execute(w, display); err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
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

			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		} else if r.FormValue("postDislike") != "" {
			postId, _ := strconv.Atoi(r.FormValue("postDislike"))

			dislike := model.Dislike{
				UserID: user.Id,
				PostID: postId,
			}

			if err := h.service.Dislike.SetPostDislike(dislike); err != nil {
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		}
	default:
		errorHandler(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) MyLikedPosts(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(model.CtxUserKey).(model.User)

	if user.Username == "" {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
	}

	switch r.Method {
	case http.MethodGet:
		var posts []model.Post

		posts, err := h.service.Post.ShowMyLikedPosts(user.Id)
		if err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		display := Display{
			Username: user.Username,
			Posts:    posts,
		}

		temp, err := template.ParseFiles("ui/mylikedposts.html")
		if err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err = temp.Execute(w, display); err != nil {
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
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

			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		} else if r.FormValue("postDislike") != "" {
			postId, _ := strconv.Atoi(r.FormValue("postDislike"))

			dislike := model.Dislike{
				UserID: user.Id,
				PostID: postId,
			}

			if err := h.service.Dislike.SetPostDislike(dislike); err != nil {
				errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		}
	default:
		errorHandler(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

// Class struct represents a single class
type Class struct {
	ID         int
	Student_id int
	Class_name string
	Days       string
	Time_start string
	Time_end   string
}

// GetClasses retrieves a list of all classes from the database
func GetClasses() ([]Class, error) {
	fmt.Println("111111")
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		fmt.Println("sql.Open err")
		return nil, err
	}
	defer db.Close()
	fmt.Println("2222222")
	rows, err := db.Query("SELECT * FROM schedule")
	if err != nil {
		fmt.Println("Query  err")
		return nil, err
	}
	defer rows.Close()

	classes := []Class{}
	fmt.Println("333333333333")
	for rows.Next() {
		class := Class{}
		err := rows.Scan(&class.ID, &class.Student_id, &class.Class_name, &class.Days, &class.Time_start, &class.Time_end)
		if err != nil {
			fmt.Println("rows.Scan err:", err)
			return nil, err
		}
		classes = append(classes, class)
	}
	fmt.Println("44444444444444")
	if err = rows.Err(); err != nil {
		fmt.Println("rows.Err() err")
		return nil, err
	}
	fmt.Println("5555555555555")
	return classes, nil
}

func (h *Handler) ScheduleHandler(w http.ResponseWriter, r *http.Request) {
	classes, err := GetClasses()
	if err != nil {
		fmt.Println("Get class Err")
		errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("ui/schedule.html")
	if err != nil {
		fmt.Println("ParseFiles Err")
		errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, classes)
	if err != nil {
		fmt.Println("Execute Err")
		errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// Структура для представления информации о файле в базе данных
type File struct {
	ID   int
	Name string
	Path string
}

// Функция для рендеринга HTML шаблона
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	fmt.Println("renderTemplate")
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		fmt.Println("t, err := template.ParseFiles(tmpl)", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("after template.ParseFiles")
	err = t.Execute(w, data)
	if err != nil {
		fmt.Println("err = t.Execute(w, data)")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("after Execute.")
}

func (h *Handler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Инициализация базы данных SQLite
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создание таблицы для хранения информации о файлах
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS files (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			path TEXT
		);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	if r.Method == http.MethodGet {
		renderTemplate(w, "ui/upload.html", nil)
	}

	if r.Method == http.MethodPost {
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Генерация уникального имени файла
		fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + filepath.Ext(handler.Filename)
		fmt.Println("before creating Directory uploads")
		filePath := filepath.Join("uploads", fileName)
		fmt.Println("After creating Directory uploads")
		// Создание файла на сервере
		out, err := os.Create(filePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer out.Close()

		// Копирование содержимого файла
		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Сохранение информации о файле в базе данных
		insertQuery := "INSERT INTO files (name, path) VALUES (?, ?)"
		_, err = db.Exec(insertQuery, handler.Filename, filePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/files", http.StatusSeeOther)
	}
	//w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully")

}

func (h *Handler) GetFilesHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Запрос всех файлов из базы данных
	selectQuery := "SELECT id, name, path FROM files"
	rows, err := db.Query(selectQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		var file File
		err := rows.Scan(&file.ID, &file.Name, &file.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		files = append(files, file)
	}

	renderTemplate(w, "ui/files.html", files)
}

func (h *Handler) DownloadFilesHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fileIDStr := r.URL.Query().Get("id")
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Запрос информации о файле по его ID
	selectQuery := "SELECT id, name, path FROM files WHERE id = ?"
	row := db.QueryRow(selectQuery, fileID)

	var file File
	err = row.Scan(&file.ID, &file.Name, &file.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Открытие файла для чтения
	fileContent, err := os.Open(file.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer fileContent.Close()

	// Установка заголовка Content-Disposition для указания имени файла при скачивании
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Name))

	// Копирование содержимого файла в ответ
	_, err = io.Copy(w, fileContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

type Material struct {
	Name     string `json:"name"`
	Subtopic string `json:"subtopic"`
	Content  string `json:"content"`
}

func (h *Handler) AddMaterialHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		renderTemplate(w, "ui/material.html", nil)
	} else if r.Method == http.MethodPost {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get form values
		materialName := r.PostForm.Get("material-name")
		subtopic := r.PostForm.Get("subtopic")
		textContent := r.PostForm.Get("text-content")

		fmt.Println("materialName: ", materialName)
		fmt.Println("subtopic: ", subtopic)
		fmt.Println("textContent: ", textContent)

		// Insert material into the database
		db, err := sql.Open("sqlite3", "forum.db")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		stmt, err := db.Prepare("INSERT INTO materials (name, subtopic, content) VALUES (?, ?, ?)")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(materialName, subtopic, textContent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to the homepage
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

func (h *Handler) GetMaterialsHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve all materials from the database
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT name, subtopic FROM materials")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	materials := []Material{}
	for rows.Next() {
		var material Material
		err := rows.Scan(&material.Name, &material.Subtopic)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		materials = append(materials, material)
	}

	// Render the materials page
	fmt.Println("GETmaterials: ", materials)
	renderTemplate(w, "ui/materials.html", materials)
}

func (h *Handler) GetMaterialHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve material name from query parameters
	name := r.URL.Query().Get("name")
	fmt.Println("GETmaterial getNAME: ", name)
	// Retrieve material from the database
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var material Material
	err = db.QueryRow("SELECT name, subtopic, content FROM materials WHERE name = ?", name).Scan(&material.Name, &material.Subtopic, &material.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the material page
	fmt.Println("GETmaterial: ", material)
	renderTemplate(w, "ui/getmaterial.html", material)
}

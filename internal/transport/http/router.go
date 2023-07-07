package http

import "net/http"

func (s *Server) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.mid.LogRequest(s.mid.AuthorizationRequired(s.handler.Home)))
	mux.HandleFunc("/schedule", s.mid.LogRequest(s.handler.ScheduleHandler))
	mux.HandleFunc("/upload", s.mid.LogRequest(s.handler.UploadHandler))

	mux.HandleFunc("/files", s.mid.LogRequest(s.handler.GetFilesHandler))
	mux.HandleFunc("/download", s.mid.LogRequest(s.handler.DownloadFilesHandler))

	mux.HandleFunc("/add-material", s.mid.LogRequest(s.handler.AddMaterialHandler))
	mux.HandleFunc("/materials", s.mid.LogRequest(s.handler.GetMaterialsHandler))
	mux.HandleFunc("/getmaterial", s.mid.LogRequest(s.handler.GetMaterialHandler))

	mux.HandleFunc("/sign-up", s.mid.LogRequest(s.handler.SignUp))
	mux.HandleFunc("/sign-in", s.mid.LogRequest(s.mid.WithSessionBlocked(s.handler.SignIn)))
	mux.HandleFunc("/createpost", s.mid.LogRequest(s.mid.AuthorizationRequired(s.handler.CreatePost)))
	mux.HandleFunc("/logout", s.mid.LogRequest(s.mid.AuthorizationRequired(s.handler.Logout)))
	mux.HandleFunc("/post/", s.mid.LogRequest(s.mid.AuthorizationRequired(s.handler.Post)))
	mux.HandleFunc("/like-post", s.mid.LogRequest(s.mid.AuthorizationRequired(s.handler.LikePost)))
	mux.HandleFunc("/post/change/", s.mid.LogRequest(s.mid.AuthorizationRequired(s.handler.ChangePost)))
	mux.HandleFunc("/post/delete/", s.mid.LogRequest(s.mid.AuthorizationRequired(s.handler.DeletePost)))
	mux.HandleFunc("/myposts", s.mid.LogRequest(s.mid.AuthorizationRequired(s.handler.MyPosts)))
	mux.HandleFunc("/my-comment-posts", s.mid.LogRequest(s.mid.AuthorizationRequired(s.handler.MyCommentPosts)))
	mux.HandleFunc("/my-liked-posts", s.mid.LogRequest(s.mid.AuthorizationRequired(s.handler.MyLikedPosts)))
	mux.Handle("/ui/css/", http.StripPrefix("/ui/css/", http.FileServer(http.Dir("./ui/css/"))))

	return mux
}

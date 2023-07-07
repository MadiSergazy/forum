package model

type Like struct {
	Id        int
	PostID    int
	UserID    int
	CommentId int
	Active    bool
}

type Dislike struct {
	Id        int
	PostID    int
	UserID    int
	CommentId int
	Active    bool
}

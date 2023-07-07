package model

type Comment struct {
	Id      int
	Author  string
	Date    string
	PostId  int
	Text    string
	Like    int
	Dislike int
	UserId  int
}

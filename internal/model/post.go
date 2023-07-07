package model

type Post struct {
	Id          int
	Title       string
	Description string
	Like        int
	Dislike     int
	AuthorId    int
	Author      string
	Date        string
	Category    []string
	Path        string
}

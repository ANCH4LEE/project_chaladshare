package models

import "time"

// post
type Post struct {
	PostID       int       `json:"post_id"`
	AuthorUserID int       `json:"author_user_id"`
	Title        string    `json:"post_title"`
	Description  string    `json:"post_description"`
	Visibility   string    `json:"post_visibility"`
	DocumentID   *int      `json:"post_document_id"`
	SummaryID    *int      `json:"post_summary_id"`
	CreatedAt    time.Time `json:"post_created_at"`
	UpdatedAt    time.Time `json:"post_updated_at"`
}

// tag
type Tag struct {
	PostID int `json:"post_id"`
	TagID  int `json:"tag_id"`
}

// liked
type Like struct {
	UserID    int       `json:"like_user_id"`
	PostID    int       `json:"like_post_id"`
	CreatedAt time.Time `json:"like_created_at"`
}

// save post
type SavePost struct {
	UserID    int       `json:"save_user_id"`
	PostID    int       `json:"save_post_id"`
	CreatedAt time.Time `json:"save_created_at"`
}

// post stat
type PostStats struct {
	PostID         int       `json:"post_id"`
	LikeCount      int       `json:"like_count"`
	SaveCount      int       `json:"save_count"`
	LaveActivityAt time.Time `json:"last_activity_at"`
}

// for response
type PostResponse struct {
	Post
	AuthorID   int        `json:"author_id"`
	AuthorName string     `json:"author_name"`
	Tags       []string   `json:"tags"`
	LikeCount  int        `json:"like_count"`
	SaveCount  int        `json:"save_count"`
	IsLiked    bool       `json:"is_liked"`
	IsSaved    bool       `json:"is_savesd"`
	Stats      *PostStats `json:"stats"`
}

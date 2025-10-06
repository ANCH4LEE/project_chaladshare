package repository

import (
	"database/sql"
	"fmt"
)

type LikeRepository interface {
	LikePost(userID, postID int) error
	UnlikePost(userID, postID int) error
	IsPostLiked(userID, postID int) (bool, error)
	UpdateLikeCount(postID int) error
}

type likeRepository struct {
	db *sql.DB
}

func NewLikeRepository(db *sql.DB) LikeRepository {
	return &likeRepository{db: db}
}

// กด like
func (r *likeRepository) LikePost(userID, postID int) error {
	query := `
		INSERT INTO likes (like_user_id, like_post_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING;
	`
	if _, err := r.db.Exec(query, userID, postID); err != nil {
		return fmt.Errorf("failed to like post: %v", err)
	}
	return r.UpdateLikeCount(postID)
}

// ยกเลิก like
func (r *likeRepository) UnlikePost(userID, postID int) error {
	query := `DELETE FROM likes WHERE like_user_id=$1 AND like_post_id=$2`
	if _, err := r.db.Exec(query, userID, postID); err != nil {
		return fmt.Errorf("failed to unlike post: %v", err)
	}
	return r.UpdateLikeCount(postID)
}

// โพสต์ถูกกดไลก์หรือยัง
func (r *likeRepository) IsPostLiked(userID, postID int) (bool, error) {
	query := `SELECT 1 FROM likes WHERE like_user_id=$1 AND like_post_id=$2`
	var exists int
	err := r.db.QueryRow(query, userID, postID).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// อัปเดตจำนวนไลก์ใน post_stat
func (r *likeRepository) UpdateLikeCount(postID int) error {
	query := `
		UPDATE post_stats
		SET post_like_count = (
			SELECT COUNT(*) FROM likes WHERE like_post_id=$1
		),
		post_last_activity_at = NOW()
		WHERE post_stats_post_id=$1;
	`
	_, err := r.db.Exec(query, postID)
	return err
}

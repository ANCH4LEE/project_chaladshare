package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"chaladshare_backend/internal/posts/models"
)

type PostRepository interface {
	CreatePost(post *models.Post) (int, error)            //สร้างโพสต์ใหม่
	AddTags(postID int, tags []string) error              //เพิ่ม tags
	GetAllPosts() ([]models.PostResponse, error)          //ดึง all post
	GetPostByID(postID int) (*models.PostResponse, error) //ดึง each post
	UpdatePost(post *models.Post) error                   //update foe edit title, description
	DeletePost(postID int) error                          //delete post
	InitPostStats(postID int) error                       //มีโพสต์ใหม่ให้สร้างบันทึก like = 0, save = 0
}

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{db: db}
}

// 1.สร้างโพสต์ใหม่
func (r *postRepository) CreatePost(post *models.Post) (int, error) {
	query := `
		INSERT INTO posts (
			post_author_user_id,
			post_title,
			post_description,
			post_visibility,
			post_document_id,
			post_summary_id
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING post_id;
	`

	var postID int
	err := r.db.QueryRow(query,
		post.AuthorUserID,
		post.Title,
		post.Description,
		post.Visibility,
		post.DocumentID,
		post.SummaryID,
	).Scan(&postID)

	if err != nil {
		return 0, fmt.Errorf("failed to create post: %v", err)
	}
	return postID, nil
}

// 2.เพิ่ม tag
func (r *postRepository) AddTags(postID int, tags []string) error {
	if len(tags) == 0 {
		return nil
	}

	for _, tag := range tags {
		var tagID int
		//ตรวจว่ามีแท็กนี้ยัง
		checkQuery := `SELECT tag_id FROM tags WHERE tag_name = $1`
		err := r.db.QueryRow(checkQuery, tag).Scan(&tagID)

		if err == sql.ErrNoRows {
			//ไม่มี = สร้างใหม่
			insertTag := `INSERT INTO tags (tag_name) VALUES ($1) RETURNING tag_id`
			if err := r.db.QueryRow(insertTag, tag).Scan(&tagID); err != nil {
				return fmt.Errorf("failed to insert tag: %v", err)
			}
		} else if err != nil {
			return err
		}

		// เชื่อมแท็กเข้ากับโพสต์
		linkQuery := `INSERT INTO post_tags (post_tag_post_id, post_tag_tag_id)
		              VALUES ($1, $2) ON CONFLICT DO NOTHING`
		if _, err := r.db.Exec(linkQuery, postID, tagID); err != nil {
			return fmt.Errorf("failed to link tag: %v", err)
		}
	}
	return nil
}

// 3.สร้างข้อมูลในการเก็บ stat
func (r *postRepository) InitPostStats(postID int) error {
	query := `INSERT INTO post_stats (post_stats_post_id, post_like_count, post_save_count)
	          VALUES ($1, 0, 0)
	          ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(query, postID)
	return err
}

// 4.ดึง all post
func (r *postRepository) GetAllPosts() ([]models.PostResponse, error) {
	query := `
	SELECT 
		p.post_id, p.post_author_user_id,
		u.username AS author_name,
		p.post_title,
		p.post_description,
		p.post_visibility,
		p.post_document_id,
		p.post_summary_id,
		p.post_created_at,
		p.post_updated_at,
		ps.post_like_count,
		ps.post_save_count,
		COALESCE(array_agg(t.tag_name) FILTER (WHERE t.tag_name IS NOT NULL), '{}') AS tags
	FROM posts p
	JOIN users u ON u.user_id = p.post_author_user_id
	LEFT JOIN post_stats ps ON ps.post_stats_post_id = p.post_id
	LEFT JOIN post_tags pt ON pt.post_tag_post_id = p.post_id
	LEFT JOIN tags t ON t.tag_id = pt.post_tag_tag_id
	GROUP BY p.post_id, u.username, ps.post_like_count, ps.post_save_count
	ORDER BY p.post_created_at DESC;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.PostResponse
	for rows.Next() {
		var p models.PostResponse
		var tagsArr []byte

		err := rows.Scan(
			&p.PostID,
			&p.AuthorID,
			&p.AuthorName,
			&p.Title,
			&p.Description,
			&p.Visibility,
			&p.DocumentID,
			&p.SummaryID,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.LikeCount,
			&p.SaveCount,
			&tagsArr,
		)
		if err != nil {
			return nil, err
		}

		p.Tags = strings.Split(strings.Trim(string(tagsArr), "{}"), ",")
		posts = append(posts, p)
	}

	return posts, nil
}

// 5.each post
func (r *postRepository) GetPostByID(postID int) (*models.PostResponse, error) {
	query := `
	SELECT 
		p.post_id,
		p.post_author_user_id,
		u.username AS author_name,
		p.post_title,
		p.post_description,
		p.post_visibility,
		p.post_document_id,
		p.post_summary_id,
		p.post_created_at,
		p.post_updated_at,
		ps.post_like_count,
		ps.post_save_count,
		COALESCE(array_agg(t.tag_name) FILTER (WHERE t.tag_name IS NOT NULL), '{}') AS tags
	FROM posts p
	JOIN users u ON u.user_id = p.post_author_user_id
	LEFT JOIN post_stats ps ON ps.post_stats_post_id = p.post_id
	LEFT JOIN post_tags pt ON pt.post_tag_post_id = p.post_id
	LEFT JOIN tags t ON t.tag_id = pt.post_tag_tag_id
	WHERE p.post_id = $1
	GROUP BY p.post_id, u.username, ps.post_like_count, ps.post_save_count;
	`

	row := r.db.QueryRow(query, postID)
	var p models.PostResponse
	var tagsArr []byte

	err := row.Scan(
		&p.PostID,
		&p.AuthorID,
		&p.AuthorName,
		&p.Title,
		&p.Description,
		&p.Visibility,
		&p.DocumentID,
		&p.SummaryID,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.LikeCount,
		&p.SaveCount,
		&tagsArr,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	p.Tags = strings.Split(strings.Trim(string(tagsArr), "{}"), ",")
	return &p, nil
}

// 6.อัปเดตโพสต์
func (r *postRepository) UpdatePost(post *models.Post) error {
	query := `
	UPDATE posts
	SET post_title = $1,
	    post_description = $2,
	    post_visibility = $3,
	    post_updated_at = now()
	WHERE post_id = $4;
	`
	_, err := r.db.Exec(query, post.Title, post.Description, post.Visibility, post.PostID)
	return err
}

// 7.ลบโพสต์
func (r *postRepository) DeletePost(postID int) error {
	query := `DELETE FROM posts WHERE post_id = $1`
	_, err := r.db.Exec(query, postID)
	return err
}

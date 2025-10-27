package repository

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"chaladshare_backend/internal/posts/models"
)

type PostRepository interface {
	CreatePost(post *models.Post, tags []string) (int, error)
	UpdatePost(post *models.Post, tags []string) error
	DeletePost(postID int) error

	GetAllPosts() ([]models.PostResponse, error)
	GetPostByID(postID int) (*models.PostResponse, error)
	GetPostOwnerID(postID int) (int, error)
}

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) CreatePost(post *models.Post, tags []string) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if post.DocumentID == nil {
		return 0, fmt.Errorf("document_id is required")
	}
	var docArg interface{} = *post.DocumentID

	var sumArg interface{} = nil
	if post.SummaryID != nil { // ถ้าเป็น *int
		sumArg = *post.SummaryID
	}

	query := `INSERT INTO posts (post_author_user_id, post_title, post_description,
			  post_visibility, post_document_id, post_summary_id) 
			  SELECT $1, $2, $3, $4, $5, $6
			  FROM documents d
			  WHERE d.document_id = $5 AND d.document_user_id = $1
			  RETURNING post_id;`

	var postID int
	if err := tx.QueryRow(
		query,
		post.AuthorUserID, post.Title, post.Description,
		post.Visibility, docArg, sumArg,
	).Scan(&postID); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("invalid document_id or not owned by user")
		}
		return 0, fmt.Errorf("create post: %w", err)
	}

	if len(tags) > 0 {
		upsertTag := `INSERT INTO tags (tag_name) VALUES ($1) ON CONFLICT (tag_name) DO UPDATE
					  SET tag_name = EXCLUDED.tag_name RETURNING tag_id;`

		link := `INSERT INTO post_tags (post_tag_post_id, post_tag_tag_id)
				 VALUES ($1, $2) ON CONFLICT DO NOTHING;`

		for _, t := range tags {
			var tagID int
			if err := tx.QueryRow(upsertTag, t).Scan(&tagID); err != nil {
				return 0, fmt.Errorf("upsert tag %q: %w", t, err)
			}
			if _, err := tx.Exec(link, postID, tagID); err != nil {
				return 0, fmt.Errorf("link tag %q: %w", t, err)
			}
		}
	}

	initStats := `INSERT INTO post_stats (post_stats_post_id, post_like_count, post_save_count)
				  VALUES ($1, 0, 0) ON CONFLICT DO NOTHING;`
	if _, err := tx.Exec(initStats, postID); err != nil {
		return 0, fmt.Errorf("init stats: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}
	return postID, nil
}

func (r *postRepository) UpdatePost(post *models.Post, tags []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(`UPDATE posts SET post_title = $1,
        				 post_description = $2, post_visibility = $3, post_updated_at = now()
    					 WHERE post_id = $4;`,
		post.Title, post.Description, post.Visibility, post.PostID)
	if err != nil {
		return fmt.Errorf("update post: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}

	if tags != nil {
		if _, err := tx.Exec(`DELETE FROM post_tags WHERE post_tag_post_id = $1`, post.PostID); err != nil {
			return fmt.Errorf("clear old tags: %w", err)
		}

		if len(tags) > 0 {
			upsertTag := `INSERT INTO tags (tag_name)
						  VALUES ($1) ON CONFLICT (tag_name) DO UPDATE
						  SET tag_name = EXCLUDED.tag_name
						  RETURNING tag_id;`

			link := `INSERT INTO post_tags (post_tag_post_id, post_tag_tag_id)
					 VALUES ($1, $2) ON CONFLICT DO NOTHING;`

			for _, t := range tags {
				var tagID int
				if err := tx.QueryRow(upsertTag, t).Scan(&tagID); err != nil {
					return fmt.Errorf("upsert tag %q: %w", t, err)
				}
				if _, err := tx.Exec(link, post.PostID, tagID); err != nil {
					return fmt.Errorf("link tag %q: %w", t, err)
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

func (r *postRepository) DeletePost(postID int) error {
	query := `DELETE FROM posts WHERE post_id = $1`
	res, err := r.db.Exec(query, postID)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *postRepository) GetAllPosts() ([]models.PostResponse, error) {
	query := `SELECT p.post_id, p.post_author_user_id, u.username AS author_name,
		p.post_title, p.post_description, p.post_visibility,
		p.post_document_id, p.post_summary_id, p.post_created_at, p.post_updated_at,
		COALESCE(ps.post_like_count, 0) AS post_like_count,
		COALESCE(ps.post_save_count, 0) AS post_save_count,
		d.document_url AS document_file_url,
		ARRAY_REMOVE(ARRAY_AGG(DISTINCT t.tag_name), NULL) AS tags
	FROM posts p
	JOIN users u ON u.user_id = p.post_author_user_id
	LEFT JOIN post_stats ps ON ps.post_stats_post_id = p.post_id
	LEFT JOIN post_tags pt ON pt.post_tag_post_id = p.post_id
	LEFT JOIN tags t ON t.tag_id = pt.post_tag_tag_id
	LEFT JOIN documents d ON d.document_id = p.post_document_id
	GROUP BY p.post_id, u.username, ps.post_like_count, ps.post_save_count, d.document_url
	ORDER BY p.post_created_at DESC;`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.PostResponse
	for rows.Next() {
		var (
			p       models.PostResponse
			tags    pq.StringArray
			fileURL sql.NullString
			docID   sql.NullInt64
			sumID   sql.NullInt64
		)

		if err := rows.Scan(
			&p.PostID,
			&p.AuthorID,
			&p.AuthorName,
			&p.Title,
			&p.Description,
			&p.Visibility,
			&docID,
			&sumID,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.LikeCount,
			&p.SaveCount,
			&fileURL,
			&tags,
		); err != nil {
			return nil, err
		}

		if docID.Valid {
			v := int(docID.Int64)
			p.DocumentID = &v
		}
		if sumID.Valid {
			v := int(sumID.Int64)
			p.SummaryID = &v
		}
		if fileURL.Valid {
			p.FileURL = &fileURL.String
		}

		p.Tags = []string(tags)
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) GetPostByID(postID int) (*models.PostResponse, error) {
	query := `SELECT p.post_id, p.post_author_user_id, u.username AS author_name,
		p.post_title, p.post_description, p.post_visibility, p.post_document_id,
		p.post_summary_id, p.post_created_at, p.post_updated_at,
		COALESCE(ps.post_like_count, 0)  AS post_like_count,
		COALESCE(ps.post_save_count, 0)  AS post_save_count,
		d.document_url AS document_file_url,
		ARRAY_REMOVE(ARRAY_AGG(DISTINCT t.tag_name), NULL) AS tags
	FROM posts p
	JOIN users u ON u.user_id = p.post_author_user_id
	LEFT JOIN post_stats ps ON ps.post_stats_post_id = p.post_id
	LEFT JOIN post_tags pt ON pt.post_tag_post_id = p.post_id
	LEFT JOIN tags t ON t.tag_id = pt.post_tag_tag_id
	LEFT JOIN documents d ON d.document_id = p.post_document_id
	WHERE p.post_id = $1
	GROUP BY p.post_id, u.username, ps.post_like_count, ps.post_save_count, d.document_url;`

	row := r.db.QueryRow(query, postID)
	var (
		p       models.PostResponse
		tags    pq.StringArray
		fileURL sql.NullString
		docID   sql.NullInt64
		sumID   sql.NullInt64
	)

	if err := row.Scan(
		&p.PostID,
		&p.AuthorID,
		&p.AuthorName,
		&p.Title,
		&p.Description,
		&p.Visibility,
		&docID,
		&sumID,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.LikeCount,
		&p.SaveCount,
		&fileURL,
		&tags,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	if docID.Valid {
		v := int(docID.Int64)
		p.DocumentID = &v
	}
	if sumID.Valid {
		v := int(sumID.Int64)
		p.SummaryID = &v
	}
	if fileURL.Valid {
		p.FileURL = &fileURL.String
	}

	p.Tags = []string(tags)
	return &p, nil
}

func (r *postRepository) GetPostOwnerID(postID int) (int, error) {
	const query = `SELECT post_author_user_id FROM posts WHERE post_id = $1`
	var owner int
	if err := r.db.QueryRow(query, postID).Scan(&owner); err != nil {
		return 0, err
	}
	return owner, nil
}

package repository

import (
	"database/sql"
	"fmt"

	postmodels "chaladshare_backend/internal/posts/models"

	"github.com/lib/pq"
)

type RecommendReadRepo interface {
	ListUserRecommendations(userID, limit, offset int) ([]postmodels.PostResponse, error)
}

type recommendReadRepo struct {
	db *sql.DB
}

func NewRecommendReadRepo(db *sql.DB) RecommendReadRepo {
	return &recommendReadRepo{db: db}
}

func (r *recommendReadRepo) ListUserRecommendations(userID, limit, offset int) ([]postmodels.PostResponse, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid userID")
	}
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	q := `
WITH rec AS (
  SELECT rec_post_id, score, created_at
  FROM recommendations
  WHERE rec_user_id = $1
  ORDER BY score DESC, created_at DESC
  LIMIT $2 OFFSET $3
)
SELECT
  p.post_id, p.post_author_user_id, u.username AS author_name,
  p.post_title, p.post_description, p.post_visibility,
  p.post_document_id, p.post_created_at, p.post_updated_at,
  COALESCE(ps.post_like_count, 0) AS post_like_count,
  COALESCE(ps.post_save_count, 0) AS post_save_count,
  d.document_url  AS document_file_url,
  d.document_name AS document_name,
  p.post_cover_url, up.avatar_url,
  ARRAY_REMOVE(ARRAY_AGG(DISTINCT t.tag_name), NULL) AS tags,

  EXISTS (
    SELECT 1 FROM likes l
    WHERE l.like_user_id = $1 AND l.like_post_id = p.post_id
  ) AS is_liked,
  EXISTS (
    SELECT 1 FROM saved_posts sp
    WHERE sp.save_user_id = $1 AND sp.save_post_id = p.post_id
  ) AS is_saved

FROM rec
JOIN posts p ON p.post_id = rec.rec_post_id
JOIN users u ON u.user_id = p.post_author_user_id
LEFT JOIN post_stats ps ON ps.post_stats_post_id = p.post_id
LEFT JOIN post_tags pt ON pt.post_tag_post_id = p.post_id
LEFT JOIN tags t ON t.tag_id = pt.post_tag_tag_id
LEFT JOIN documents d ON d.document_id = p.post_document_id
LEFT JOIN user_profiles up ON up.profile_user_id = u.user_id

-- กันเคส visibility friends / public ให้เหมือน feed
WHERE
  (
    p.post_author_user_id = $1
    OR p.post_visibility = 'public'
    OR (
      p.post_visibility = 'friends'
      AND EXISTS (
        SELECT 1
        FROM friendships f
        WHERE
          f.user_id  = LEAST(p.post_author_user_id, $1)
          AND f.friend_id = GREATEST(p.post_author_user_id, $1)
      )
    )
  )

GROUP BY
  p.post_id, u.username,
  ps.post_like_count, ps.post_save_count,
  d.document_url, d.document_name,
  p.post_cover_url, up.avatar_url,
  rec.score, rec.created_at

ORDER BY rec.score DESC, rec.created_at DESC;
`

	rows, err := r.db.Query(q, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]postmodels.PostResponse, 0, limit)

	for rows.Next() {
		var (
			p         postmodels.PostResponse
			tags      pq.StringArray
			fileURL   sql.NullString
			docName   sql.NullString
			coverURL  sql.NullString
			avatarURL sql.NullString
			docID     sql.NullInt64
			isLiked   bool
			isSaved   bool
		)

		if err := rows.Scan(
			&p.PostID, &p.AuthorID, &p.AuthorName,
			&p.Title, &p.Description, &p.Visibility,
			&docID, &p.CreatedAt, &p.UpdatedAt,
			&p.LikeCount, &p.SaveCount,
			&fileURL, &docName, &coverURL, &avatarURL, &tags,
			&isLiked, &isSaved,
		); err != nil {
			return nil, err
		}

		if docID.Valid {
			v := int(docID.Int64)
			p.DocumentID = &v
		}
		if fileURL.Valid {
			p.FileURL = &fileURL.String
		}
		if docName.Valid {
			p.DocumentName = &docName.String
		}
		if coverURL.Valid {
			p.CoverURL = &coverURL.String
		}
		if avatarURL.Valid {
			p.AvatarURL = &avatarURL.String
		}

		p.Tags = []string(tags)
		p.IsLiked = isLiked
		p.IsSaved = isSaved

		out = append(out, p)
	}

	return out, rows.Err()
}

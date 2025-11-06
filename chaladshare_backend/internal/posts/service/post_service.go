package service

import (
	"fmt"
	"strings"

	"chaladshare_backend/internal/posts/models"
	"chaladshare_backend/internal/posts/repository"
)

type PostService interface {
	CreatePost(post *models.Post, tags []string) (int, error)
	UpdatePost(post *models.Post, tags []string) error
	DeletePost(postID int) error

	GetAllPosts() ([]models.PostResponse, error)
	GetPostByID(postID int) (*models.PostResponse, error)
	CountByUserID(userID int) (int, error)

	IsOwner(postID int, userID int) (bool, error)
	ViewPost(viewerID, postID int) (bool, string, error)
	Friends(viewerID, authorID int) (bool, error)
}

type postService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) PostService {
	return &postService{postRepo: postRepo}
}

// สร้างโพสต์ใหม่
func (s *postService) CreatePost(post *models.Post, tags []string) (int, error) {
	if post.AuthorUserID <= 0 {
		return 0, fmt.Errorf("invalid author")
	}
	if strings.TrimSpace(post.Title) == "" {
		return 0, fmt.Errorf("post_title is required")
	}
	if post.Visibility != models.VisibilityPublic {
		return 0, fmt.Errorf("unsupported visibility")
	}

	normTags := normalizeTags(tags)
	postID, err := s.postRepo.CreatePost(post, normTags)
	if err != nil {
		return 0, fmt.Errorf("failed to create post: %w", err)
	}
	return postID, nil
}

func (s *postService) UpdatePost(post *models.Post, tags []string) error {
	if post.PostID <= 0 {
		return fmt.Errorf("invalid post_id")
	}
	if strings.TrimSpace(post.Title) == "" {
		return fmt.Errorf("post_title is required")
	}
	if post.Visibility != models.VisibilityPublic {
		return fmt.Errorf("unsupported visibility")
	}
	var normTags []string
	if tags != nil {
		normTags = normalizeTags(tags)
	}
	if err := s.postRepo.UpdatePost(post, normTags); err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	return nil
}

func (s *postService) DeletePost(postID int) error {
	return s.postRepo.DeletePost(postID)
}

func normalizeTags(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))

	const (
		maxTags = 10
		maxLen  = 30
	)
	for _, t := range in {
		tag := strings.TrimSpace(t)
		tag = strings.TrimPrefix(tag, "#")
		tag = strings.ToLower(tag)

		if tag == "" || len(tag) > maxLen {
			continue
		}
		valid := true
		for _, r := range tag {
			if !(r >= 'a' && r <= 'z') && !(r >= '0' && r <= '9') && r != '_' && r != '-' {
				valid = false
				break
			}
		}
		if !valid {
			continue
		}

		if _, dup := seen[tag]; dup {
			continue
		}
		seen[tag] = struct{}{}
		out = append(out, tag)
		if len(out) >= maxTags {
			break
		}
	}
	return out
}

func (s *postService) GetAllPosts() ([]models.PostResponse, error) {
	return s.postRepo.GetAllPosts()
}

// each post by ID
func (s *postService) GetPostByID(postID int) (*models.PostResponse, error) {
	return s.postRepo.GetPostByID(postID)
}

func (s *postService) CountByUserID(userID int) (int, error) {
	return s.postRepo.CountByUserID(userID)
}

func (s *postService) IsOwner(postID int, userID int) (bool, error) {
	ownerID, err := s.postRepo.GetPostOwnerID(postID)
	if err != nil {
		return false, fmt.Errorf("cannot get post owner: %w", err)
	}
	return ownerID == userID, nil
}

func (s *postService) ViewPost(viewerID, postID int) (bool, string, error) {
	post, err := s.GetPostByID(postID)
	if err != nil {
		return false, "error", fmt.Errorf("get post: %w", err)
	}
	if post == nil {
		return false, "not_found", nil
	}

	authorID := post.AuthorID
	vis := strings.ToLower(strings.TrimSpace(post.Visibility))
	if vis == "" {
		vis = models.VisibilityPublic
	}
	if viewerID == authorID {
		return true, "owner", nil
	}

	switch vis {
	case models.VisibilityPublic:
		return true, "public", nil
	case models.VisibilityFriends:
		ok, err := s.Friends(viewerID, authorID)
		if err != nil {
			return false, "error", err
		}
		if ok {
			return true, "friends", nil
		}
		return false, "friends_only", nil
	default:
		return false, "denied", nil
	}
}

func (s *postService) Friends(viewerID, authorID int) (bool, error) {
	return false, nil
}

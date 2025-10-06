package service

import (
	"fmt"

	"chaladshare_backend/internal/posts/models"
	"chaladshare_backend/internal/posts/repository"
)

type PostService interface {
	CreatePost(post *models.Post, tags []string) (int, error)
	GetAllPosts() ([]models.PostResponse, error)
	GetPostByID(postID int) (*models.PostResponse, error)
	UpdatePost(post *models.Post) error
	DeletePost(postID int) error
}

type postService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) PostService {
	return &postService{postRepo: postRepo}
}

// สร้างโพสต์ใหม่
func (s *postService) CreatePost(post *models.Post, tags []string) (int, error) {
	postID, err := s.postRepo.CreatePost(post)
	if err != nil {
		return 0, fmt.Errorf("failed to create post: %v", err)
	}

	if err := s.postRepo.AddTags(postID, tags); err != nil {
		return 0, fmt.Errorf("failed to add tags: %v", err)
	}

	if err := s.postRepo.InitPostStats(postID); err != nil {
		return 0, fmt.Errorf("failed to init stats: %v", err)
	}

	return postID, nil
}

// ดึง all post
func (s *postService) GetAllPosts() ([]models.PostResponse, error) {
	return s.postRepo.GetAllPosts()
}

// each post by ID
func (s *postService) GetPostByID(postID int) (*models.PostResponse, error) {
	return s.postRepo.GetPostByID(postID)
}

// update for edit post
func (s *postService) UpdatePost(post *models.Post) error {
	return s.postRepo.UpdatePost(post)
}

// delete post
func (s *postService) DeletePost(postID int) error {
	return s.postRepo.DeletePost(postID)
}

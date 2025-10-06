package service

import "chaladshare_backend/internal/posts/repository"

type LikeService interface {
	LikePost(userID, postID int) error
	UnlikePost(userID, postID int) error
	IsPostLiked(userID, postID int) (bool, error)
}

type likeService struct {
	likeRepo repository.LikeRepository
}

func NewLikeService(likeRepo repository.LikeRepository) LikeService {
	return &likeService{likeRepo: likeRepo}
}

// กด like
func (s *likeService) LikePost(userID, postID int) error {
	return s.likeRepo.LikePost(userID, postID)
}

// ยกเลิก like
func (s *likeService) UnlikePost(userID, postID int) error {
	return s.likeRepo.UnlikePost(userID, postID)
}

// ตรวจสอบ
func (s *likeService) IsPostLiked(userID, postID int) (bool, error) {
	return s.likeRepo.IsPostLiked(userID, postID)
}

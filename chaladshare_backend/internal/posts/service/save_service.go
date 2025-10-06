package service

import "chaladshare_backend/internal/posts/repository"

type SaveService interface {
	SavePost(userID, postID int) error
	UnsavePost(userID, postID int) error
	IsPostSaved(userID, postID int) (bool, error)
}

type saveService struct {
	saveRepo repository.SaveRepository
}

func NewSaveService(saveRepo repository.SaveRepository) SaveService {
	return &saveService{saveRepo: saveRepo}
}

// บันทึกโพสต์
func (s *saveService) SavePost(userID, postID int) error {
	return s.saveRepo.SavePost(userID, postID)
}

//ยกเลิกการบันทึก
func (s *saveService) UnsavePost(userID, postID int) error {
	return s.saveRepo.UnsavePost(userID, postID)
}

//ตรวจสอบ
func (s *saveService) IsPostSaved(userID, postID int) (bool, error) {
	return s.saveRepo.IsPostSaved(userID, postID)
}

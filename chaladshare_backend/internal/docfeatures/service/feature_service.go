package service

import (
	"fmt"

	"chaladshare_backend/internal/docfeatures/models"
	"chaladshare_backend/internal/docfeatures/repository"
)

type FeatureService interface {
	CreateQueued(documentID int) error
	MarkProcessing(documentID int) error
	SaveResult(input models.SaveResult) error
	MarkFailed(documentID int, msg string) error
	GetByDocumentID(documentID int) (*models.DocumentFeature, error)
	ProcessDocument(documentID int, pdfPath string)
}

type featureService struct {
	repo repository.DocFeaturesRepo
}

func NewFeatureService(repo repository.DocFeaturesRepo) FeatureService {
	return &featureService{repo: repo}
}

func (s *featureService) CreateQueued(documentID int) error {
	if documentID <= 0 {
		return fmt.Errorf("invalid documentID")
	}
	return s.repo.CreateQueued(documentID)
}

func (s *featureService) MarkProcessing(documentID int) error {
	if documentID <= 0 {
		return fmt.Errorf("invalid documentID")
	}
	return s.repo.MarkProcessing(documentID)
}

func (s *featureService) SaveResult(input models.SaveResult) error {
	if input.DocumentID <= 0 {
		return fmt.Errorf("invalid documentID")
	}
	return s.repo.SaveResult(input)
}

func (s *featureService) MarkFailed(documentID int, msg string) error {
	if documentID <= 0 {
		return fmt.Errorf("invalid documentID")
	}
	if msg == "" {
		msg = "unknown error"
	}
	return s.repo.MarkFailed(documentID, msg)
}

func (s *featureService) GetByDocumentID(documentID int) (*models.DocumentFeature, error) {
	if documentID <= 0 {
		return nil, fmt.Errorf("invalid documentID")
	}
	return s.repo.GetByDocumentID(documentID)
}

func (s *featureService) ProcessDocument(documentID int, pdfPath string) {
	if err := s.repo.MarkProcessing(documentID); err != nil {
		_ = s.repo.MarkFailed(documentID, err.Error())
		return
	}

	resp, err := callColabExtract(documentID, pdfPath)
	if err != nil {
		_ = s.repo.MarkFailed(documentID, err.Error())
		return
	}

	ct := resp.ContentText

	if err := s.repo.SaveResult(models.SaveResult{
		DocumentID:       documentID,
		StyleLabel:       resp.StyleLabel,
		StyleVectorV16:   resp.StyleVectorV16,
		ContentText:      &ct,
		ContentEmbedding: resp.Embedding,
		ClusterID:        resp.ClusterID,
	}); err != nil {
		_ = s.repo.MarkFailed(documentID, err.Error())
		return
	}
}

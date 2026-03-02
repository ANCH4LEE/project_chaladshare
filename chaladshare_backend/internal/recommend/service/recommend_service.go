package service

import (
	"fmt"
	"log"

	"chaladshare_backend/internal/connect"
	recmodels "chaladshare_backend/internal/recommend/models"
	recrepo "chaladshare_backend/internal/recommend/repository"
)

type RecommendService interface {
	RecomputeFromLikes(userID int) error
	OnLikeHook(userID int)
}

type svc struct {
	repo     recrepo.RecommendRepo
	aiClient *connect.Client
}

func NewRecommendService(repo recrepo.RecommendRepo, aiClient *connect.Client) RecommendService {
	return &svc{repo: repo, aiClient: aiClient}
}

func (s *svc) OnLikeHook(userID int) {
	go func() {
		if err := s.RecomputeFromLikes(userID); err != nil {
			log.Printf("[RECOMMEND] recompute error user=%d: %v", userID, err)
		}
	}()
}

func (s *svc) RecomputeFromLikes(userID int) error {
	if s.aiClient == nil {
		return fmt.Errorf("ai client is nil")
	}
	if userID <= 0 {
		return fmt.Errorf("invalid userID")
	}

	seeds, pairs, err := s.repo.ListSeedsFromLikes(userID, 5)
	if err != nil {
		return err
	}
	if len(seeds) == 0 || len(pairs) == 0 {
		return nil
	}

	cands, err := s.repo.ListCandidatesBySeedPairs(userID, pairs, 800)
	if err != nil {
		return err
	}
	if len(cands) == 0 {
		return nil
	}

	req := recmodels.ColabRecommendFromLikedReq{
		Seeds:            seeds,
		Candidates:       cands,
		TopK:             10,
		BoostSameCluster: 0.05,
		MaxPerCluster:    4,
	}

	resp, err := s.aiClient.RecommendFromLiked(req)
	if err != nil {
		return err
	}

	if resp == nil || len(resp.Recommendations) == 0 {
		return nil
	}
	return s.repo.ReplaceUserRecommendations(userID, resp.Recommendations)
}

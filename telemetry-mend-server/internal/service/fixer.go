package service

import (
	"context"
	"telemetry-mend-server/internal/ai"
	"telemetry-mend-server/internal/models"

	"github.com/uptrace/bun"
)

type FixerService struct {
	db              *bun.DB
	analyzerService *AnalyzerService
	aiProvider      ai.Provider
}

func NewFixerService(db *bun.DB, analyzer *AnalyzerService, ai ai.Provider) *FixerService {
	return &FixerService{
		db:              db,
		analyzerService: analyzer,
		aiProvider:      ai,
	}
}

func (s *FixerService) GenerateAndStoreFix(ctx context.Context, clusterID int64) (*models.SuggestedFix, error) {
	// 1. Get Code Context
	contexts, err := s.analyzerService.GetContextForCluster(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	// 2. Get Cluster Info
	cluster := new(models.ErrorCluster)
	if err := s.db.NewSelect().Model(cluster).Where("id = ?", clusterID).Scan(ctx); err != nil {
		return nil, err
	}

	// 3. Prepare AI Prompt Context
	var contextBuilder string
	for _, c := range contexts {
		contextBuilder += "File: " + c.FilePath + "\n"
		contextBuilder += "Code:\n" + c.Content + "\n\n"
	}

	// 4. Generate Fix
	resp, err := s.aiProvider.GenerateFix(ctx, cluster.LogTemplate, contextBuilder)
	if err != nil {
		return nil, err
	}

	// 5. Store Fix
	fix := &models.SuggestedFix{
		ClusterID:    clusterID,
		Explanation:  resp.Explanation,
		FilePath:     resp.FilePath,
		OriginalCode: resp.OriginalCode,
		FixedCode:    resp.FixedCode,
		GitPatch:     resp.GitPatch,
		Status:       "pending",
	}

	_, err = s.db.NewInsert().Model(fix).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return fix, nil
}

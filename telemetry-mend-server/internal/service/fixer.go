package service

import (
	"context"
	"fmt"
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
	// 1. Get Cluster & Application Info
	cluster := new(models.ErrorCluster)
	if err := s.db.NewSelect().
		Model(cluster).
		Relation("Application").
		Where("ec.id = ?", clusterID).
		Scan(ctx); err != nil {
		return nil, err
	}

	if cluster.Application == nil || cluster.Application.RepoURL == "" {
		return nil, fmt.Errorf("fix generation skipped: no repository URL configured for this application")
	}

	// 2. Get LLM Settings
	settings := new(models.Settings)
	err := s.db.NewSelect().Model(settings).Limit(1).Scan(ctx)
	
	var provider ai.Provider
	if err != nil || settings.LLMProvider == "mock" {
		provider = &ai.MockProvider{}
	} else {
		provider = ai.NewOpenAIProvider(settings.LLMBaseURL, settings.LLMAPIKey, settings.LLMModel)
	}

	// 3. Get Code Context
	contexts, err := s.analyzerService.GetContextForCluster(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	if len(contexts) == 0 {
		return nil, fmt.Errorf("fix generation skipped: could not find relevant code context in repository")
	}

	// 4. Prepare AI Prompt Context
	var contextBuilder string
	for _, c := range contexts {
		contextBuilder += "File: " + c.FilePath + "\n"
		contextBuilder += "Code:\n" + c.Content + "\n\n"
	}

	// 5. Generate Fix
	resp, err := provider.GenerateFix(ctx, cluster.LogTemplate, contextBuilder)
	if err != nil {
		return nil, err
	}

	// 6. Store Fix
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

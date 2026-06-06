package service

import (
	"context"
	"fmt"
	"strings"
	"telemetry-mend-server/internal/models"
	"telemetry-mend-server/internal/scm"

	"github.com/uptrace/bun"
)

type AnalyzerService struct {
	db        *bun.DB
	gitClient *scm.GitClient
}

func NewAnalyzerService(db *bun.DB, gitClient *scm.GitClient) *AnalyzerService {
	return &AnalyzerService{
		db:        db,
		gitClient: gitClient,
	}
}

type CodeContext struct {
	FilePath string
	Content  string
	Line     int
}

func (s *AnalyzerService) GetContextForCluster(ctx context.Context, clusterID int64) ([]CodeContext, error) {
	cluster := new(models.ErrorCluster)
	err := s.db.NewSelect().
		Model(cluster).
		Relation("Application").
		Where("ec.id = ?", clusterID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	// Get latest log entry for this cluster to get commit hash
	lastLog := new(models.LogEntry)
	err = s.db.NewSelect().
		Model(lastLog).
		Where("cluster_id = ?", clusterID).
		Order("created_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	locations := scm.ParseStackTrace(lastLog.LogBody)
	if len(locations) == 0 {
		return nil, fmt.Errorf("no code locations found in stack trace")
	}

	repoPath, err := s.gitClient.CloneOrFetch(ctx, cluster.Application.RepoURL, cluster.AppID)
	if err != nil {
		return nil, err
	}

	var contexts []CodeContext
	// Just get the first few locations to avoid overwhelming context
	for i, loc := range locations {
		if i > 2 {
			break
		}
		content, err := s.gitClient.GetFileContent(repoPath, lastLog.CommitHash, loc.FilePath)
		if err != nil {
			continue
		}

		// Extract a window around the line
		lines := strings.Split(content, "\n")
		start := loc.Line - 20
		if start < 0 {
			start = 0
		}
		end := loc.Line + 20
		if end > len(lines) {
			end = len(lines)
		}

		windowContent := strings.Join(lines[start:end], "\n")

		contexts = append(contexts, CodeContext{
			FilePath: loc.FilePath,
			Content:  windowContent,
			Line:     loc.Line,
		})
	}

	return contexts, nil
}

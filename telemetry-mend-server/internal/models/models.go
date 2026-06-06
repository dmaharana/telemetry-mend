package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Application struct {
	bun.BaseModel `bun:"table:applications,alias:a"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	Name        string    `bun:"name,notnull" json:"name"`
	Language    string    `bun:"language" json:"language"`
	RepoURL     string    `bun:"repo_url,notnull" json:"repo_url"`
	DefaultBranch string  `bun:"default_branch,notnull,default:'main'" json:"default_branch"`
	CreatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
}

type Deployment struct {
	bun.BaseModel `bun:"table:deployments,alias:d"`

	ID            int64     `bun:"id,pk,autoincrement" json:"id"`
	AppID         int64     `bun:"app_id,notnull" json:"app_id"`
	Environment   string    `bun:"environment,notnull" json:"environment"`
	Branch        string    `bun:"branch,notnull" json:"branch"`
	CommitHash    string    `bun:"commit_hash,notnull" json:"commit_hash"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`

	Application   *Application `bun:"rel:belongs-to,join:app_id=id" json:"application"`
}

type ErrorCluster struct {
	bun.BaseModel `bun:"table:error_clusters,alias:ec"`

	ID            int64     `bun:"id,pk,autoincrement" json:"id"`
	AppID         int64     `bun:"app_id,notnull" json:"app_id"`
	Fingerprint   string    `bun:"fingerprint,notnull,unique" json:"fingerprint"`
	LogTemplate   string    `bun:"log_template,notnull" json:"log_template"`
	LastSeen      time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"last_seen"`
	Count         int       `bun:"count,notnull,default:1" json:"count"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`

	Application   *Application `bun:"rel:belongs-to,join:app_id=id" json:"application"`
}

type LogEntry struct {
	bun.BaseModel `bun:"table:log_entries,alias:le"`

	ID            int64     `bun:"id,pk,autoincrement" json:"id"`
	AppID         int64     `bun:"app_id,notnull" json:"app_id"`
	ClusterID     int64     `bun:"cluster_id" json:"cluster_id"`
	Environment   string    `bun:"environment,notnull" json:"environment"`
	CommitHash    string    `bun:"commit_hash,notnull" json:"commit_hash"`
	LogBody       string    `bun:"log_body,notnull" json:"log_body"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`

	Cluster       *ErrorCluster `bun:"rel:belongs-to,join:cluster_id=id" json:"cluster"`
}

type SuggestedFix struct {
	bun.BaseModel `bun:"table:suggested_fixes,alias:sf"`

	ID            int64     `bun:"id,pk,autoincrement" json:"id"`
	ClusterID     int64     `bun:"cluster_id,notnull" json:"cluster_id"`
	Explanation   string    `bun:"explanation" json:"explanation"`
	FilePath      string    `bun:"file_path" json:"file_path"`
	OriginalCode  string    `bun:"original_code" json:"original_code"`
	FixedCode     string    `bun:"fixed_code" json:"fixed_code"`
	GitPatch      string    `bun:"git_patch" json:"git_patch"`
	Status        string    `bun:"status,notnull,default:'pending'" json:"status"` // pending, applied, rejected
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`

	Cluster       *ErrorCluster `bun:"rel:belongs-to,join:cluster_id=id" json:"cluster"`
}

func CreateTables(ctx context.Context, db *bun.DB) error {
	models := []interface{}{
		(*Application)(nil),
		(*Deployment)(nil),
		(*ErrorCluster)(nil),
		(*LogEntry)(nil),
		(*SuggestedFix)(nil),
	}

	for _, model := range models {
		_, err := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

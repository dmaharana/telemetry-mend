package ai

import (
	"context"
)

type FixResponse struct {
	Explanation  string `json:"explanation"`
	FilePath     string `json:"file_path"`
	OriginalCode string `json:"original_code"`
	FixedCode    string `json:"fixed_code"`
	GitPatch     string `json:"git_patch"`
}

type Provider interface {
	GenerateFix(ctx context.Context, errorLog string, codeContext string) (*FixResponse, error)
}

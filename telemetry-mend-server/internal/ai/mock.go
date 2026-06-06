package ai

import (
	"context"
)

type MockProvider struct{}

func (m *MockProvider) GenerateFix(ctx context.Context, errorLog string, codeContext string) (*FixResponse, error) {
	return &FixResponse{
		Explanation:  "This is a mock fix generated for the error: " + errorLog,
		FilePath:     "main.go",
		OriginalCode: "// original code here",
		FixedCode:    "// fixed code here",
		GitPatch:     "--- a/main.go\n+++ b/main.go\n@@ -1,1 +1,1 @@\n-// original code here\n+// fixed code here",
	}, nil
}

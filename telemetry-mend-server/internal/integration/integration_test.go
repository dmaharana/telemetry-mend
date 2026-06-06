package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"telemetry-mend-server/internal/db"
	"telemetry-mend-server/internal/handlers"
	"telemetry-mend-server/internal/models"
	"telemetry-mend-server/internal/scm"
	"telemetry-mend-server/internal/ai"
	"telemetry-mend-server/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestFullFlow(t *testing.T) {
	ctx := context.Background()
	
	// 1. Setup
	database, err := db.InitDB(ctx, ":memory:")
	assert.NoError(t, err)
	defer database.Close()

	gitClient, _ := scm.NewGitClient("./test-repos")
	analyzer := service.NewAnalyzerService(database, gitClient)
	mockAI := &ai.MockProvider{}
	fixer := service.NewFixerService(database, analyzer, mockAI)

	appHandler := handlers.NewAppHandler(database)
	logHandler := handlers.NewLogHandler(database)
	clusterHandler := handlers.NewClusterHandler(database, fixer)

	r := chi.NewRouter()
	r.Post("/apps", appHandler.Create)
	r.Post("/logs/ingest", logHandler.Ingest)
	r.Get("/clusters/{id}", clusterHandler.Get)
	r.Post("/clusters/{id}/fix", clusterHandler.GenerateFix)

	// 2. Create App
	appReq := models.Application{
		Name:    "test-app",
		RepoURL: "https://github.com/test/repo",
	}
	appBody, _ := json.Marshal(appReq)
	req, _ := http.NewRequest("POST", "/apps", bytes.NewBuffer(appBody))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var createdApp models.Application
	json.Unmarshal(rr.Body.Bytes(), &createdApp)
	assert.NotZero(t, createdApp.ID)

	// 3. Ingest Log
	logReq := handlers.IngestLogRequest{
		AppID:       createdApp.ID,
		Environment: "prod",
		CommitHash:  "abcdef123456",
		LogBody:     "panic: runtime error: index out of range\n\ngoroutine 1 [running]:\nmain.main()\n\t/home/user/main.go:42 +0x24",
	}
	logBody, _ := json.Marshal(logReq)
	req, _ = http.NewRequest("POST", "/logs/ingest", bytes.NewBuffer(logBody))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusAccepted, rr.Code)

	var ingestResp map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &ingestResp)
	clusterID := int64(ingestResp["cluster_id"].(float64))

	// 4. Verify Cluster
	req, _ = http.NewRequest("GET", "/clusters/"+strconv.FormatInt(clusterID, 10), nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// 5. Generate Fix (Mocked)
	// Expecting 500 because git clone will fail for fake repo URL.
	req, _ = http.NewRequest("POST", "/clusters/"+strconv.FormatInt(clusterID, 10)+"/fix", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

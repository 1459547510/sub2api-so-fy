package handler

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func soraCompatMultipartRequest(t *testing.T, fields map[string]string, pngField string) *http.Request {
	t.Helper()
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	for key, value := range fields {
		require.NoError(t, writer.WriteField(key, value))
	}
	if pngField != "" {
		part, err := writer.CreateFormFile(pngField, "reference.png")
		require.NoError(t, err)
		_, err = part.Write(append([]byte("\x89PNG\r\n\x1a\n"), make([]byte, 24)...))
		require.NoError(t, err)
	}
	require.NoError(t, writer.Close())
	request := httptest.NewRequest(http.MethodPost, "/v1/videos", buffer)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	return request
}

func TestSoraVideoCompatCreateTranslatesMultipartForm(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerVideoJobRepo{}
	h := &OpenAIGatewayHandler{videoJobService: newHandlerVideoJobService(repo)}
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = soraCompatMultipartRequest(t, map[string]string{
		"model": "seedance-2.0", "prompt": "waves", "seconds": "8", "size": "1920x1080",
	}, "")
	setHandlerVideoAuth(c)

	h.SoraVideoCompatCreate(c)

	require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
	require.NotNil(t, repo.job)
	require.Equal(t, "1080p", repo.job.Resolution)
	require.Equal(t, "16:9", repo.job.AspectRatio)
	require.Equal(t, 8, repo.job.DurationSeconds)

	var body map[string]any
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
	require.Equal(t, repo.job.JobID, body["id"])
	require.Equal(t, "video", body["object"])
	require.Equal(t, "seedance-2.0", body["model"])
	require.Equal(t, "queued", body["status"])
	require.Equal(t, "8", body["seconds"])
	require.Equal(t, "1920x1080", body["size"])
	require.EqualValues(t, 0, body["progress"])
}

func TestSoraVideoCompatCreateMetadataOverridesAndImageURLs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerVideoJobRepo{}
	h := &OpenAIGatewayHandler{videoJobService: newHandlerVideoJobService(repo)}
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = soraCompatMultipartRequest(t, map[string]string{
		"model": "seedance-2.0", "prompt": "waves", "size": "1920x1080",
		"metadata": `{"resolution":"720p","aspect_ratio":"9:16","audio":true,"image_urls":["https://example.com/ref.png"],"ignored_key":"x"}`,
	}, "")
	setHandlerVideoAuth(c)

	h.SoraVideoCompatCreate(c)

	require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
	require.NotNil(t, repo.job)
	require.Equal(t, "720p", repo.job.Resolution)
	require.Equal(t, "9:16", repo.job.AspectRatio)
	require.True(t, repo.job.Audio)
	require.Equal(t, "url", repo.job.ImageSource)
}

func TestSoraVideoCompatCreateStoresInputReference(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerVideoJobRepo{}
	store := service.NewVideoInputStore(t.TempDir(), 8080)
	h := &OpenAIGatewayHandler{videoJobService: newHandlerVideoJobService(repo), videoInputHandler: NewVideoInputHandler(store)}
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = soraCompatMultipartRequest(t, map[string]string{
		"model": "seedance-2.0", "prompt": "waves", "seconds": "8", "size": "1280x720",
	}, "input_reference")
	setHandlerVideoAuth(c)

	h.SoraVideoCompatCreate(c)

	require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
	require.NotNil(t, repo.job)
	require.NotEmpty(t, repo.job.LocalInputName, "input_reference must be tracked as a local input token")
	require.Equal(t, "local", repo.job.ImageSource)
}

func TestSoraVideoCompatCreateRejectsBadSizeAndSeconds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &OpenAIGatewayHandler{videoJobService: newHandlerVideoJobService(&handlerVideoJobRepo{})}

	for name, fields := range map[string]map[string]string{
		"bad size":    {"model": "seedance-2.0", "prompt": "waves", "size": "huge"},
		"bad seconds": {"model": "seedance-2.0", "prompt": "waves", "seconds": "4.5"},
	} {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request = soraCompatMultipartRequest(t, fields, "")
		setHandlerVideoAuth(c)

		h.SoraVideoCompatCreate(c)

		require.Equal(t, http.StatusBadRequest, recorder.Code, name)
		require.Contains(t, recorder.Body.String(), "invalid_request_error", name)
	}
}

func TestParseSoraVideoSize(t *testing.T) {
	for _, tc := range []struct {
		size       string
		resolution string
		aspect     string
	}{
		{"1920x1080", "1080p", "16:9"},
		{"720x1280", "720p", "9:16"},
		{"720*720", "720p", "1:1"},
		{"2560x1080", "1080p", "21:9"},
		{"640x480", "480p", "4:3"},
		{"720p", "720p", ""},
		{"4k", "2160p", ""},
	} {
		resolution, aspect, err := parseSoraVideoSize(tc.size)
		require.NoError(t, err, tc.size)
		require.Equal(t, tc.resolution, resolution, tc.size)
		require.Equal(t, tc.aspect, aspect, tc.size)
	}

	_, _, err := parseSoraVideoSize("huge")
	require.Error(t, err)
	_, _, err = parseSoraVideoSize("0x720")
	require.Error(t, err)
}

func TestSoraVideoStatusMapsEveryNativeStatus(t *testing.T) {
	require.Equal(t, "queued", soraVideoStatus(service.VideoJobPending))
	require.Equal(t, "in_progress", soraVideoStatus(service.VideoJobRunning))
	require.Equal(t, "in_progress", soraVideoStatus(service.VideoJobSettling))
	require.Equal(t, "completed", soraVideoStatus(service.VideoJobCompleted))
	require.Equal(t, "failed", soraVideoStatus(service.VideoJobFailed))
	require.Equal(t, "failed", soraVideoStatus(service.VideoJobCanceled))
	require.Equal(t, "queued", soraVideoStatus("unexpected"))
}

func TestLeoVideoJobSoraReturnsVideoObjectAndHidesProviderNames(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerVideoJobRepo{job: &service.VideoJob{
		JobID: "vidjob_failed", APIKeyID: 2, Status: service.VideoJobFailed,
		RequestedModel: "seedance-2.0", Resolution: "720p", AspectRatio: "16:9", DurationSeconds: 8,
		ErrorMessage: "LeoStudio failed after Leonardo.ai rejected the Krea request upstream",
	}}
	h := &OpenAIGatewayHandler{videoJobService: newHandlerVideoJobService(repo)}
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/videos/vidjob_failed", nil)
	c.Params = gin.Params{{Key: "request_id", Value: "vidjob_failed"}}
	setHandlerVideoAuth(c)

	h.LeoVideoJobSora(c)

	require.Equal(t, http.StatusOK, recorder.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
	require.Equal(t, "vidjob_failed", body["id"])
	require.Equal(t, "video", body["object"])
	require.Equal(t, "failed", body["status"])
	require.Equal(t, "1280x720", body["size"])
	require.NotNil(t, body["error"])

	lower := strings.ToLower(recorder.Body.String())
	for _, forbidden := range []string{"leonardo", "leostudio", "krea", "trioma"} {
		require.NotContains(t, lower, forbidden)
	}
}

func TestLeoVideoJobSoraCompletedAndScoped(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerVideoJobRepo{job: &service.VideoJob{
		JobID: "vidjob_done", APIKeyID: 2, Status: service.VideoJobCompleted,
		RequestedModel: "seedance-2.0", Resolution: "1080p", AspectRatio: "9:16", DurationSeconds: 5,
	}}
	h := &OpenAIGatewayHandler{videoJobService: newHandlerVideoJobService(repo)}

	for _, apiKeyID := range []int64{999, 2} {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request = httptest.NewRequest(http.MethodGet, "/v1/videos/vidjob_done", nil)
		c.Params = gin.Params{{Key: "request_id", Value: "vidjob_done"}}
		setHandlerVideoAuthWithKey(c, apiKeyID)
		h.LeoVideoJobSora(c)
		if apiKeyID == 999 {
			require.Equal(t, http.StatusNotFound, recorder.Code)
			continue
		}
		require.Equal(t, http.StatusOK, recorder.Code)
		var body map[string]any
		require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
		require.Equal(t, "completed", body["status"])
		require.EqualValues(t, 100, body["progress"])
		require.Equal(t, "1080x1920", body["size"])
	}
}

func TestSoraVideoCompatCreateKeepsJSONContractUntouched(t *testing.T) {
	// JSON clients on POST /v1/videos keep the native 202 async contract; the
	// dialect switch happens on Content-Type in the route dispatcher, so the
	// native handler must still accept JSON with Prefer: respond-async.
	gin.SetMode(gin.TestMode)
	repo := &handlerVideoJobRepo{}
	h := &OpenAIGatewayHandler{videoJobService: newHandlerVideoJobService(repo)}
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/videos", strings.NewReader(`{"model":"seedance","prompt":"waves"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("Prefer", "respond-async")
	setHandlerVideoAuth(c)

	h.LeoVideoGeneration(c)

	require.Equal(t, http.StatusAccepted, recorder.Code)
	require.Contains(t, recorder.Body.String(), "job_id")
	require.Contains(t, recorder.Body.String(), "status_url")
}

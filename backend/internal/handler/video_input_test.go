package handler

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestVideoInputHandlerUploadAndLoopbackRead(t *testing.T) {
	store := service.NewVideoInputStore(t.TempDir(), 8080)
	handler := NewVideoInputHandler(store)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", "frame.png")
	require.NoError(t, err)
	_, err = part.Write(videoPNGBytes())
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/v1/videos/uploads", func(c *gin.Context) {
		c.Set(string(middleware2.ContextKeyAPIKey), &service.APIKey{ID: 2})
		handler.Upload(c)
	})
	router.GET("/internal/video-inputs/:token", handler.GetInternal)
	req := httptest.NewRequest(http.MethodPost, "/v1/videos/uploads", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)
	require.Contains(t, resp.Body.String(), "upload_id")

	var upload struct {
		UploadID string `json:"upload_id"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &upload))
	token := upload.UploadID
	getReq := httptest.NewRequest(http.MethodGet, "/internal/video-inputs/"+token, nil)
	getReq.RemoteAddr = "127.0.0.1:51234"
	getResp := httptest.NewRecorder()
	router.ServeHTTP(getResp, getReq)
	require.Equal(t, http.StatusOK, getResp.Code)
	require.Equal(t, videoPNGBytes(), getResp.Body.Bytes())

	getReq = httptest.NewRequest(http.MethodGet, "/internal/video-inputs/"+token, nil)
	getReq.RemoteAddr = "10.0.0.9:51234"
	getResp = httptest.NewRecorder()
	router.ServeHTTP(getResp, getReq)
	require.Equal(t, http.StatusNotFound, getResp.Code)
}

func TestVideoInputHandlerRequiresAPIKeyForUpload(t *testing.T) {
	handler := NewVideoInputHandler(service.NewVideoInputStore(t.TempDir(), 8080))
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/videos/uploads", strings.NewReader(""))
	handler.Upload(c)
	require.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestVideoInputHandlerPreparesLeoMultipartImageRequest(t *testing.T) {
	store := service.NewVideoInputStore(t.TempDir(), 8080)
	require.NoError(t, store.SetPublicBaseURL("https://api.example.com"))
	handler := NewVideoInputHandler(store)
	parsed := &service.OpenAIImagesRequest{
		Endpoint:  "/v1/images/edits",
		Multipart: true,
		Model:     "GPT Image-2",
		Prompt:    "keep the subject",
		N:         1,
		Size:      "1024x1024",
		Uploads: []service.OpenAIImagesUpload{{
			FileName:    "reference.png",
			ContentType: "image/png",
			Data:        videoPNGBytes(),
		}},
	}

	body, tokens, err := handler.PrepareLeoMultipartImageRequest(parsed)
	require.NoError(t, err)
	require.Len(t, tokens, 1)
	require.False(t, parsed.Multipart)
	require.Equal(t, "application/json", parsed.ContentType)
	require.Equal(t, "GPT Image-2", gjson.GetBytes(body, "model").String())
	require.Equal(t, "keep the subject", gjson.GetBytes(body, "prompt").String())
	require.Equal(t, "https://api.example.com/media/image-inputs/"+tokens[0], gjson.GetBytes(body, "image_urls.0").String())
	_, err = store.Open(tokens[0])
	require.NoError(t, err)
}

func TestVideoInputHandlerPublicImageRead(t *testing.T) {
	store := service.NewVideoInputStore(t.TempDir(), 8080)
	input, err := store.Save(bytes.NewReader(videoPNGBytes()))
	require.NoError(t, err)
	handler := NewVideoInputHandler(store)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/media/image-inputs/:token", handler.GetPublicImage)
	req := httptest.NewRequest(http.MethodGet, "/media/image-inputs/"+input.Token, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)
	require.Equal(t, "image/png", resp.Header().Get("Content-Type"))
	require.Equal(t, "no-store", resp.Header().Get("Cache-Control"))
	require.Equal(t, "nosniff", resp.Header().Get("X-Content-Type-Options"))
	require.Equal(t, videoPNGBytes(), resp.Body.Bytes())
}

func TestVideoInputHandlerUploadsVideoAndAudioFields(t *testing.T) {
	store := service.NewVideoInputStore(t.TempDir(), 8080)
	handler := NewVideoInputHandler(store)

	for _, test := range []struct {
		field       string
		filename    string
		content     []byte
		mediaType   string
		responseKey string
	}{
		{field: "video", filename: "reference.mp4", content: handlerMP4Bytes(), mediaType: "video", responseKey: "video_url"},
		{field: "audio", filename: "reference.mp3", content: handlerMP3Bytes(), mediaType: "audio", responseKey: "audio_url"},
	} {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile(test.field, test.filename)
		require.NoError(t, err)
		_, err = part.Write(test.content)
		require.NoError(t, err)
		require.NoError(t, writer.Close())

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.POST("/v1/videos/uploads", func(c *gin.Context) {
			c.Set(string(middleware2.ContextKeyAPIKey), &service.APIKey{ID: 2})
			handler.Upload(c)
		})
		req := httptest.NewRequest(http.MethodPost, "/v1/videos/uploads", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		var response map[string]any
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &response))
		require.Equal(t, test.mediaType, response["media_type"])
		require.NotEmpty(t, response["media_url"])
		require.NotEmpty(t, response[test.responseKey])
	}
}

func handlerMP4Bytes() []byte {
	data := make([]byte, 16)
	binary.BigEndian.PutUint32(data[0:4], uint32(len(data)))
	copy(data[4:8], "ftyp")
	copy(data[8:12], "isom")
	return data
}

func handlerMP3Bytes() []byte {
	const frameLength = 417
	const frameCount = 77
	data := make([]byte, frameLength*frameCount)
	for offset := 0; offset < len(data); offset += frameLength {
		data[offset] = 0xff
		data[offset+1] = 0xfb
		data[offset+2] = 0x90
		data[offset+3] = 0x64
	}
	return data
}

func videoPNGBytes() []byte {
	return []byte("\x89PNG\r\n\x1a\nvideo")
}

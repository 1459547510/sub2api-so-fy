package service

import (
	"bytes"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestRewriteOpenAIImagesN_RewritesJSONAndMultipart(t *testing.T) {
	jsonBody, contentType, err := RewriteOpenAIImagesN([]byte(`{"model":"gpt-image-2","n":4}`), "application/json", 1)
	require.NoError(t, err)
	require.Equal(t, int64(1), gjson.GetBytes(jsonBody, "n").Int())
	require.Equal(t, "application/json", contentType)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	require.NoError(t, writer.WriteField("model", "gpt-image-2"))
	require.NoError(t, writer.WriteField("n", "4"))
	part, err := writer.CreateFormFile("image", "ref.png")
	require.NoError(t, err)
	_, err = part.Write([]byte("png"))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	rewritten, rewrittenType, err := RewriteOpenAIImagesN(body.Bytes(), writer.FormDataContentType(), 1)
	require.NoError(t, err)
	parsed := &OpenAIImagesRequest{Endpoint: openAIImagesEditsEndpoint, ContentType: rewrittenType, N: 1, Multipart: true, Body: rewritten}
	require.NoError(t, parseOpenAIImagesMultipartRequest(rewritten, rewrittenType, parsed))
	require.Equal(t, 1, parsed.N)
}

func TestOpenAIImagesNValidationRejectsUnsafeFanout(t *testing.T) {
	for _, body := range [][]byte{
		[]byte(`{"n":11}`),
		[]byte(`{"n":1.5}`),
		[]byte(`{"n":0}`),
	} {
		err := parseOpenAIImagesJSONRequest(body, &OpenAIImagesRequest{})
		require.Error(t, err)
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	require.NoError(t, writer.WriteField("n", "11"))
	require.NoError(t, writer.Close())
	err := parseOpenAIImagesMultipartRequest(body.Bytes(), writer.FormDataContentType(), &OpenAIImagesRequest{})
	require.ErrorContains(t, err, "between 1 and 10")
}

func TestGrokMediaImageCountValidationRejectsUnsafeFanout(t *testing.T) {
	for _, body := range [][]byte{
		[]byte(`{"n":11}`),
		[]byte(`{"n":1.5}`),
		[]byte(`{"n":0}`),
		[]byte(`{"n":"2"}`),
	} {
		info := ParseGrokMediaRequest("application/json", body)
		require.Error(t, info.ValidateImageCount())
	}

	info := ParseGrokMediaRequest("application/json", []byte(`{"n":10}`))
	require.NoError(t, info.ValidateImageCount())
	require.Equal(t, 10, info.N)
}

func TestMergeOpenAIImageResponses_AppendsDataAndSumsUsage(t *testing.T) {
	merged, err := MergeOpenAIImageResponses([][]byte{
		[]byte(`{"created":1,"model":"gpt-image-2","data":[{"b64_json":"one"}],"usage":{"output_tokens":2,"details":{"images":1}}}`),
		[]byte(`{"created":2,"data":[{"b64_json":"two"}],"usage":{"output_tokens":3,"details":{"images":1}}}`),
	})
	require.NoError(t, err)
	require.Len(t, gjson.GetBytes(merged, "data").Array(), 2)
	require.Equal(t, int64(5), gjson.GetBytes(merged, "usage.output_tokens").Int())
	require.Equal(t, int64(2), gjson.GetBytes(merged, "usage.details.images").Int())
	require.Equal(t, "gpt-image-2", gjson.GetBytes(merged, "model").String())
}

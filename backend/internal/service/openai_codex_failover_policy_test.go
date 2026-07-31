package service

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenAIStreamFailedPlanGatedModelDoesNotFailover(t *testing.T) {
	payload := []byte(`{"type":"response.failed","response":{"status":"failed","error":{"code":"model_not_supported","message":"The 'gpt-5.2' model is not supported when using Codex with a ChatGPT account."}}}`)

	require.False(t, openAIStreamFailedEventShouldFailover(payload, "The 'gpt-5.2' model is not supported when using Codex with a ChatGPT account."))
}

func TestOpenAIStreamFailedTransientProcessingStillFailsOver(t *testing.T) {
	payload := []byte(`{"type":"response.failed","error":{"message":"An error occurred while processing your request. You can retry your request, or contact us through our help center at help.openai.com. Please include the request ID req_123."}}`)

	require.True(t, openAIStreamFailedEventShouldFailover(payload, "An error occurred while processing your request. You can retry your request, or contact us through our help center at help.openai.com. Please include the request ID req_123."))
}

func TestOpenAICodexPlanGatedHTTP400IsNotAClassifiedFailover(t *testing.T) {
	svc := &OpenAIGatewayService{}
	body := []byte(`{"detail":"The 'gpt-5.2' model is not supported when using Codex with a ChatGPT account."}`)

	require.False(t, svc.shouldFailoverOpenAIUpstreamResponse(http.StatusBadRequest, "", body))
	require.True(t, isOpenAICodexPlanGatedModelError(http.StatusBadRequest, body))
}

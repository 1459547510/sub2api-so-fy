package service

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateLeoVideoRequestAppliesSeedanceDefaults(t *testing.T) {
	info, err := ValidateLeoVideoRequest([]byte(`{"model":"seedance-2.0-mini","prompt":"waves"}`))

	require.NoError(t, err)
	require.Equal(t, "720p", info.Resolution)
	require.Equal(t, 8, info.DurationSeconds)
	require.Equal(t, "16:9", info.AspectRatio)
}

func TestValidateLeoVideoRequestRejectsUnsupportedModelParameters(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "fast 1080p",
			body: `{"model":"seedance-2.0-fast","prompt":"waves","resolution":"1080p"}`,
			want: "resolution is not supported",
		},
		{
			name: "standard 720p 9:21",
			body: `{"model":"seedance-2.0","prompt":"waves","resolution":"720p","aspect_ratio":"9:21"}`,
			want: "aspect_ratio is not supported",
		},
		{
			name: "mini non 16:9 aspect",
			body: `{"model":"seedance-2.0-mini","prompt":"waves","resolution":"720p","aspect_ratio":"9:16"}`,
			want: "aspect_ratio is not supported",
		},
		{
			name: "duration below range",
			body: `{"model":"seedance-2.0","prompt":"waves","duration":3}`,
			want: "duration must be a whole number from 4 through 15",
		},
		{
			name: "fractional duration",
			body: `{"model":"seedance-2.0","prompt":"waves","duration":4.5}`,
			want: "duration must be a whole number from 4 through 15",
		},
		{
			name: "standard 1080p duration above range",
			body: `{"model":"seedance-2.0","prompt":"waves","resolution":"1080p","duration":13}`,
			want: "duration must be a whole number from 4 through 12",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateLeoVideoRequest([]byte(tt.body))
			require.ErrorContains(t, err, tt.want)
		})
	}
}

func TestValidateLeoVideoRequestAcceptsResolutionSpecificAspect(t *testing.T) {
	info, err := ValidateLeoVideoRequest([]byte(`{"model":"seedance-2.0","prompt":"waves","resolution":"1080p","aspect_ratio":"9:21","duration":12}`))

	require.NoError(t, err)
	require.Equal(t, "1080p", info.Resolution)
	require.Equal(t, "9:21", info.AspectRatio)
	require.Equal(t, 12, info.DurationSeconds)
}

func TestValidateLeoVideoRequestKeepsFifteenSecondsForOtherSupportedModes(t *testing.T) {
	for _, body := range []string{
		`{"model":"seedance-2.0","prompt":"waves","resolution":"720p","duration":15}`,
		`{"model":"seedance-2.0-fast","prompt":"waves","resolution":"720p","duration":15}`,
	} {
		info, err := ValidateLeoVideoRequest([]byte(body))

		require.NoError(t, err)
		require.Equal(t, 15, info.DurationSeconds)
	}
}

func TestValidateLeoVideoRequestRejectsPromptAndGuidanceLimits(t *testing.T) {
	longPrompt := `{"model":"seedance-2.0","prompt":"` + strings.Repeat("界", 5001) + `"}`
	_, err := ValidateLeoVideoRequest([]byte(longPrompt))
	require.ErrorContains(t, err, "prompt supports at most 5000 characters")

	tooManyImages := []byte(`{
		"model":"seedance-2.0","prompt":"waves",
		"image_urls":["https://example.com/1.png","https://example.com/2.png"],
		"guidances":{"image_reference":[{},{},{}]}
	}`)
	_, err = ValidateLeoVideoRequest(tooManyImages)
	require.ErrorContains(t, err, "guidances.image_reference supports at most 4")

	duplicateStart := []byte(`{
		"model":"seedance-2.0","prompt":"waves",
		"image_url":"https://example.com/start.png",
		"start_frame_url":"https://example.com/other.png"
	}`)
	_, err = ValidateLeoVideoRequest(duplicateStart)
	require.ErrorContains(t, err, "start frame must be supplied only once")
}

func TestValidateLeoVideoRequestAcceptsMediaAndAudioReferenceURLs(t *testing.T) {
	body := []byte(`{
		"model":"seedance-2.0","prompt":"waves",
		"guidances":{
			"image_reference":[{"image":{"url":"https://example.com/reference.png","type":"UPLOADED"}}],
			"video_reference_base":[{"video":{"url":"https://example.com/reference.mp4","type":"UPLOADED"}}],
			"audio_reference":[{"audio":{"url":"https://example.com/reference.mp3","type":"UPLOADED"}}]
		}
	}`)

	_, err := ValidateLeoVideoRequest(body)
	require.NoError(t, err)
}

func TestValidateLeoVideoRequestValidatesMediaAndAudioReferences(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "video generated URL",
			body: `{"model":"seedance-2.0","prompt":"waves","guidances":{"video_reference_base":[{"video":{"url":"https://example.com/reference.mp4","type":"GENERATED"}}]}}`,
			want: "video.url requires type UPLOADED",
		},
		{
			name: "audio duration with URL",
			body: `{"model":"seedance-2.0","prompt":"waves","guidances":{"image_reference":[{"image":{"url":"https://example.com/reference.png"}}],"audio_reference":[{"audio":{"url":"https://example.com/reference.mp3","duration":3}}]}}`,
			want: "audio.duration must be omitted when audio.url is used",
		},
		{
			name: "audio without visual reference",
			body: `{"model":"seedance-2.0","prompt":"waves","guidances":{"audio_reference":[{"audio":{"id":"33333333-3333-3333-3333-333333333333"}}]}}`,
			want: "guidances.audio_reference requires an image_reference or video_reference_base",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateLeoVideoRequest([]byte(tt.body))
			require.ErrorContains(t, err, tt.want)
		})
	}
}

func TestValidateLeoVideoRequestUsesMappedModelSpec(t *testing.T) {
	body := []byte(`{"model":"public-seedance","prompt":"waves","resolution":"1080p"}`)

	_, err := ValidateLeoVideoRequest(body)
	require.NoError(t, err)
	_, err = ValidateLeoVideoRequestForModel(body, "seedance-2.0-fast")
	require.ErrorContains(t, err, "resolution is not supported")

	body = []byte(`{"model":"public-seedance","prompt":"waves","resolution":"1080p","duration":13}`)
	_, err = ValidateLeoVideoRequestForModel(body, "seedance-2.0")
	require.ErrorContains(t, err, "duration must be a whole number from 4 through 12")
}

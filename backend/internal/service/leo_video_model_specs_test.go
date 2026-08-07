package service

import (
	"fmt"
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

func TestValidateLeoVideoRequestAppliesLTXDefaults(t *testing.T) {
	for _, model := range []string{"ltx-2.3-pro", "ltx-2.3-fast"} {
		info, err := ValidateLeoVideoRequest([]byte(`{"model":"` + model + `","prompt":"waves"}`))

		require.NoError(t, err)
		require.Equal(t, "1080p", info.Resolution)
		require.Equal(t, 6, info.DurationSeconds)
		require.Equal(t, "16:9", info.AspectRatio)
	}
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
			name: "mini unsupported aspect",
			body: `{"model":"seedance-2.0-mini","prompt":"waves","resolution":"720p","aspect_ratio":"4:3"}`,
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
			name: "zero duration",
			body: `{"model":"seedance-2.0","prompt":"waves","duration":0}`,
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

func TestValidateLeoVideoRequestSupportsLatestModels(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "mini portrait at 480p",
			body: `{"model":"seedance-2.0-mini","prompt":"waves","resolution":"480p","aspect_ratio":"9:16","duration":4}`,
		},
		{
			name: "happy horse with prompt enhancement",
			body: `{"model":"happy-horse-1.1","prompt":"waves","resolution":"1080p","duration":3,"prompt_enhance":"AUTO","guidances":{"image_reference":[{"image":{"url":"https://example.com/reference.png"}}]}}`,
		},
		{
			name: "grok square resolution",
			body: `{"model":"grok-imagine-1.5","prompt":"waves","resolution":"544p","aspect_ratio":"1:1","duration":15,"start_frame_url":"https://example.com/start.png"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateLeoVideoRequest([]byte(tt.body))
			require.NoError(t, err)
		})
	}
}

func TestValidateLeoVideoRequestSupportsLTXModels(t *testing.T) {
	tests := []string{
		`{"model":"ltx-2.3-pro","prompt":"waves","resolution":"2160p","duration":10,"aspect_ratio":"16:9","audio":true,"prompt_enhance":"ON","start_frame_url":"https://example.com/start.png","end_frame_url":"https://example.com/end.png"}`,
		`{"model":"ltx-2.3-fast","prompt":"waves","resolution":"1440p","duration":20,"aspect_ratio":"16:9","prompt_enhance":"AUTO"}`,
	}

	for _, body := range tests {
		info, err := ValidateLeoVideoRequest([]byte(body))
		require.NoError(t, err)
		require.Contains(t, []string{"1440p", "2160p"}, info.Resolution)
	}
}

func TestValidateLeoVideoRequestRejectsUnsupportedLTXParameters(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{name: "pro duration", body: `{"model":"ltx-2.3-pro","prompt":"waves","duration":12}`, want: "duration must be one of 6, 8, 10"},
		{name: "fast odd duration", body: `{"model":"ltx-2.3-fast","prompt":"waves","duration":7}`, want: "duration must be one of 6, 8, 10, 12, 14, 16, 18, 20"},
		{name: "aspect ratio", body: `{"model":"ltx-2.3-pro","prompt":"waves","aspect_ratio":"9:16"}`, want: "aspect_ratio is not supported"},
		{name: "image reference", body: `{"model":"ltx-2.3-pro","prompt":"waves","guidances":{"image_reference":[{"image":{"url":"https://example.com/reference.png"}}]}}`, want: "guidances.image_reference supports at most 0"},
		{name: "video reference", body: `{"model":"ltx-2.3-fast","prompt":"waves","guidances":{"video_reference_base":[{"video":{"url":"https://example.com/reference.mp4"}}]}}`, want: "guidances.video_reference_base supports at most 0"},
		{name: "audio reference", body: `{"model":"ltx-2.3-fast","prompt":"waves","guidances":{"audio_reference":[{"audio":{"url":"https://example.com/reference.mp3"}}]}}`, want: "guidances.audio_reference supports at most 0"},
		{name: "seed", body: `{"model":"ltx-2.3-pro","prompt":"waves","seed":1}`, want: "seed and mode are not supported"},
		{name: "mode", body: `{"model":"ltx-2.3-fast","prompt":"waves","mode":"fast"}`, want: "seed and mode are not supported"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateLeoVideoRequest([]byte(tt.body))
			require.ErrorContains(t, err, tt.want)
		})
	}
}

func TestValidateLeoVideoRequestRejectsLatestModelGuidanceLimits(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "happy horse end frame",
			body: `{"model":"happy-horse-1.1","prompt":"waves","end_frame_url":"https://example.com/end.png"}`,
			want: "guidances.end_frame supports at most 0",
		},
		{
			name: "happy horse audio",
			body: `{"model":"happy-horse-1.1","prompt":"waves","guidances":{"audio_reference":[{"audio":{"url":"https://example.com/reference.mp3"}}],"image_reference":[{"image":{"url":"https://example.com/reference.png"}}]}}`,
			want: "guidances.audio_reference supports at most 0",
		},
		{
			name: "grok requires start frame",
			body: `{"model":"grok-imagine-1.5","prompt":"waves"}`,
			want: "start frame is required",
		},
		{
			name: "grok reference image",
			body: `{"model":"grok-imagine-1.5","prompt":"waves","start_frame_url":"https://example.com/start.png","guidances":{"image_reference":[{"image":{"url":"https://example.com/reference.png"}}]}}`,
			want: "reference images cannot be combined",
		},
		{
			name: "happy horse prompt enhancement with start frame",
			body: `{"model":"happy-horse-1.1","prompt":"waves","prompt_enhance":"ON","start_frame_url":"https://example.com/start.png"}`,
			want: "prompt_enhance ON is not supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateLeoVideoRequest([]byte(tt.body))
			require.ErrorContains(t, err, tt.want)
		})
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

	mixedFrameAndReferences := []byte(`{
		"model":"seedance-2.0","prompt":"waves",
		"start_frame_url":"https://example.com/start.png",
		"image_urls":["https://example.com/reference.png"]
	}`)
	_, err = ValidateLeoVideoRequest(mixedFrameAndReferences)
	require.ErrorContains(t, err, "reference images cannot be combined")
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

	body = []byte(`{"model":"public-ltx","prompt":"waves"}`)
	info, err := ValidateLeoVideoRequestForModel(body, leoLTX23FastUpstreamModelID)
	require.NoError(t, err)
	require.Equal(t, "1080p", info.Resolution)
	require.Equal(t, 6, info.DurationSeconds)
}

func TestValidateLeoVideoRequestSupportsLeoStudioNewModels(t *testing.T) {
	tests := []struct {
		model      string
		resolution string
		duration   int
		aspect     string
	}{
		{"hailuo-03", "1440p", 5, "16:9"},
		{"gemini-omni-flash", "720p", 3, "9:16"},
		{"kling-2.1", "1080p", 5, "16:9"},
		{"kling-2.5", "720p", 10, "1:1"},
		{"kling-2.5-turbo-standard", "720p", 5, "9:16"},
		{"kling-2.6", "auto", 5, "auto"},
		{"kling-video-o-1", "1080p", 3, "16:9"},
		{"kling-3.0", "2160p", 15, "9:16"},
		{"kling-3.0-turbo", "auto", 3, "auto"},
		{"kling-video-o-3", "1080p", 15, "1:1"},
		{"veo-3.1-generate-001", "2160p", 8, "9:16"},
		{"veo-3.1-fast-generate-001", "1080p", 6, "16:9"},
		{"veo-3.1-lite", "720p", 4, "16:9"},
	}
	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			startFrame := ""
			if tt.model == "kling-2.1" || tt.model == "kling-2.5-turbo-standard" {
				startFrame = `,"start_frame_url":"https://example.com/start.png"`
			}
			body := fmt.Sprintf(`{"model":%q,"prompt":"waves","resolution":%q,"duration":%d,"aspect_ratio":%q%s}`, tt.model, tt.resolution, tt.duration, tt.aspect, startFrame)
			_, err := ValidateLeoVideoRequest([]byte(body))
			require.NoError(t, err)
		})
	}
}

func TestValidateLeoVideoRequestSupportsSeedanceMiniAudio(t *testing.T) {
	_, err := ValidateLeoVideoRequest([]byte(`{"model":"seedance-2.0-mini","prompt":"waves","audio":true,"guidances":{"image_reference":[{"image":{"url":"https://example.com/reference.png","type":"UPLOADED"}}],"audio_reference":[{"audio":{"url":"https://example.com/reference.mp3","type":"UPLOADED"}}]}}`))
	require.NoError(t, err)
}

func TestValidateLeoVideoRequestEnforcesNewGuidanceConstraints(t *testing.T) {
	_, err := ValidateLeoVideoRequest([]byte(`{"model":"hailuo-03","prompt":"waves","end_frame_url":"https://example.com/end.png"}`))
	require.ErrorContains(t, err, "end frame requires a start frame")

	_, err = ValidateLeoVideoRequest([]byte(`{"model":"hailuo-03","prompt":"waves","start_frame_url":"https://example.com/start.png","end_frame_url":"https://example.com/end.png","image_urls":["https://example.com/reference.png"]}`))
	require.ErrorContains(t, err, "reference images cannot be combined")

	_, err = ValidateLeoVideoRequest([]byte(`{"model":"hailuo-03","prompt":"waves","image_urls":["https://example.com/reference.png"],"guidances":{"audio_reference":[{"audio":{"url":"https://example.com/a.mp3","type":"UPLOADED"}},{"audio":{"url":"https://example.com/b.mp3","type":"UPLOADED"}},{"audio":{"url":"https://example.com/c.mp3","type":"UPLOADED"}}]}}`))
	require.NoError(t, err)

	_, err = ValidateLeoVideoRequest([]byte(`{"model":"hailuo-03","prompt":"waves","image_urls":["https://example.com/reference.png"],"guidances":{"audio_reference":[{"audio":{"id":"33333333-3333-3333-3333-333333333333","duration":10}},{"audio":{"id":"44444444-4444-4444-4444-444444444444","duration":6}}]}}`))
	require.ErrorContains(t, err, "at most 15 seconds total")

	_, err = ValidateLeoVideoRequest([]byte(`{"model":"kling-video-o-3","prompt":"waves","guidances":{"video_reference_base":[{"video":{"url":"https://example.com/reference.mp4","type":"UPLOADED"}}]}}`))
	require.ErrorContains(t, err, "video.url requires type UPLOADED")

	_, err = ValidateLeoVideoRequest([]byte(`{"model":"kling-video-o-3","prompt":"waves","duration":11,"guidances":{"video_reference_base":[{"video":{"id":"33333333-3333-3333-3333-333333333333","type":"GENERATED"}}]}}`))
	require.ErrorContains(t, err, "at most 10 seconds")

	_, err = ValidateLeoVideoRequest([]byte(`{"model":"gemini-omni-flash","prompt":"waves","audio":true}`))
	require.ErrorContains(t, err, "audio is not supported")
}

package service

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/tidwall/gjson"
)

const leoVideoMaxPromptLength = 5000

var leoVideoAssetIDPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

type leoVideoModelSpec struct {
	resolutions                      []string
	defaultResolution                string
	aspects                          map[string][]string
	allowedDurations                 []int
	maxDurationByResolution          map[string]int
	minDuration                      int
	maxDuration                      int
	defaultDuration                  int
	maxPromptLength                  int
	maxStartFrames                   int
	maxEndFrames                     int
	maxImageRefs                     int
	maxVideoRefs                     int
	maxAudioRefs                     int
	maxImageRefsWithVideo            int
	maxDurationWithVideo             int
	maxAudioRefSeconds               float64
	frameImageTypes                  []string
	videoReferenceTypes              []string
	requiresStartFrame               bool
	endFrameRequiresStart            bool
	framesExcludeOtherRef            bool
	audioRefRequiresMedia            bool
	supportsAudio                    bool
	supportsPromptEnhance            bool
	rejectsPromptEnhanceOnStartFrame bool
	rejectsSeedAndMode               bool
}

var defaultLeoVideoModelSpec = leoVideoModelSpec{
	defaultResolution:     "720p",
	minDuration:           4,
	maxDuration:           15,
	defaultDuration:       8,
	maxPromptLength:       leoVideoMaxPromptLength,
	maxStartFrames:        1,
	maxEndFrames:          1,
	maxImageRefs:          4,
	maxVideoRefs:          3,
	maxAudioRefs:          1,
	audioRefRequiresMedia: true,
	supportsAudio:         true,
	framesExcludeOtherRef: true,
}

var leoSeedance25ModelSpec = leoVideoModelSpec{
	resolutions:       []string{"480p", "720p"},
	defaultResolution: "720p",
	aspects: map[string][]string{
		"480p": {"16:9", "9:16", "1:1", "4:3", "3:4", "21:9"},
		"720p": {"16:9", "9:16", "1:1", "4:3", "3:4", "21:9"},
	},
	minDuration:           4,
	maxDuration:           30,
	defaultDuration:       8,
	maxPromptLength:       5000,
	maxStartFrames:        1,
	maxEndFrames:          1,
	maxImageRefs:          30,
	maxVideoRefs:          10,
	maxAudioRefs:          10,
	maxAudioRefSeconds:    30.2,
	audioRefRequiresMedia: true,
	supportsAudio:         true,
}

var leoVideoModelSpecs = map[string]leoVideoModelSpec{
	"seedance-2.0": {
		resolutions:       []string{"480p", "720p", "1080p"},
		defaultResolution: "720p",
		maxDurationByResolution: map[string]int{
			"1080p": 12,
		},
		aspects: map[string][]string{
			"480p":  {"16:9", "9:16", "1:1", "4:3", "3:4", "21:9", "9:21"},
			"720p":  {"16:9", "9:16", "1:1", "4:3", "3:4", "21:9"},
			"1080p": {"16:9", "9:16", "1:1", "4:3", "3:4", "21:9", "9:21"},
		},
		minDuration:           4,
		maxDuration:           15,
		defaultDuration:       8,
		maxPromptLength:       5000,
		maxStartFrames:        1,
		maxEndFrames:          1,
		maxImageRefs:          4,
		maxVideoRefs:          3,
		maxAudioRefs:          1,
		audioRefRequiresMedia: true,
		supportsAudio:         true,
		framesExcludeOtherRef: true,
	},
	"seedance-2.0-fast": {
		resolutions:       []string{"480p", "720p"},
		defaultResolution: "720p",
		aspects: map[string][]string{
			"480p": {"16:9", "9:16", "1:1", "4:3", "3:4", "21:9", "9:21"},
			"720p": {"16:9", "9:16", "1:1", "4:3", "3:4", "21:9"},
		},
		minDuration:           4,
		maxDuration:           15,
		defaultDuration:       8,
		maxPromptLength:       5000,
		maxStartFrames:        1,
		maxEndFrames:          1,
		maxImageRefs:          4,
		maxVideoRefs:          3,
		maxAudioRefs:          1,
		frameImageTypes:       []string{"UPLOADED", "GENERATED"},
		audioRefRequiresMedia: true,
		supportsAudio:         true,
		framesExcludeOtherRef: true,
	},
	"seedance-2.0-mini": {
		resolutions:       []string{"480p", "720p"},
		defaultResolution: "720p",
		aspects: map[string][]string{
			"480p": {"16:9", "1:1", "9:16"},
			"720p": {"16:9", "1:1", "9:16"},
		},
		minDuration:           4,
		maxDuration:           15,
		defaultDuration:       8,
		maxPromptLength:       5000,
		maxStartFrames:        1,
		maxEndFrames:          1,
		maxImageRefs:          4,
		maxVideoRefs:          3,
		maxAudioRefs:          1,
		audioRefRequiresMedia: true,
		supportsAudio:         true,
		framesExcludeOtherRef: true,
	},
	"bytedance/seedance-2.5": leoSeedance25ModelSpec,
	"seedance-2.5":           leoSeedance25ModelSpec,
	"happy-horse-1.1": {
		resolutions:       []string{"720p", "1080p"},
		defaultResolution: "1080p",
		aspects: map[string][]string{
			"720p":  {"16:9", "4:3", "1:1", "3:4", "9:16"},
			"1080p": {"16:9", "4:3", "1:1", "3:4", "9:16"},
		},
		minDuration:                      3,
		maxDuration:                      15,
		defaultDuration:                  5,
		maxPromptLength:                  2500,
		maxStartFrames:                   1,
		maxImageRefs:                     9,
		supportsAudio:                    true,
		supportsPromptEnhance:            true,
		rejectsPromptEnhanceOnStartFrame: true,
		framesExcludeOtherRef:            true,
	},
	"grok-imagine-1.5": {
		resolutions:       []string{"auto", "400p", "544p", "720p", "960p"},
		defaultResolution: "720p",
		aspects: map[string][]string{
			"auto": {"auto"},
			"400p": {"16:9", "9:16"},
			"544p": {"1:1"},
			"720p": {"16:9", "9:16"},
			"960p": {"1:1"},
		},
		minDuration:           3,
		maxDuration:           15,
		defaultDuration:       6,
		maxPromptLength:       5000,
		maxStartFrames:        1,
		requiresStartFrame:    true,
		supportsAudio:         true,
		framesExcludeOtherRef: true,
	},
	LeoLTX23ProModelID: {
		resolutions:       []string{"1080p", "1440p", "2160p"},
		defaultResolution: "1080p",
		aspects: map[string][]string{
			"1080p": {"16:9"},
			"1440p": {"16:9"},
			"2160p": {"16:9"},
		},
		allowedDurations:      []int{6, 8, 10},
		minDuration:           6,
		maxDuration:           10,
		defaultDuration:       6,
		maxPromptLength:       5000,
		maxStartFrames:        1,
		maxEndFrames:          1,
		maxImageRefs:          0,
		maxVideoRefs:          0,
		maxAudioRefs:          0,
		supportsAudio:         true,
		supportsPromptEnhance: true,
		rejectsSeedAndMode:    true,
	},
	LeoLTX23FastModelID: {
		resolutions:       []string{"1080p", "1440p", "2160p"},
		defaultResolution: "1080p",
		aspects: map[string][]string{
			"1080p": {"16:9"},
			"1440p": {"16:9"},
			"2160p": {"16:9"},
		},
		allowedDurations:      []int{6, 8, 10, 12, 14, 16, 18, 20},
		minDuration:           6,
		maxDuration:           20,
		defaultDuration:       6,
		maxPromptLength:       5000,
		maxStartFrames:        1,
		maxEndFrames:          1,
		maxImageRefs:          0,
		maxVideoRefs:          0,
		maxAudioRefs:          0,
		supportsAudio:         true,
		supportsPromptEnhance: true,
		rejectsSeedAndMode:    true,
	},
	"hailuo-03": {
		resolutions:       []string{"1440p"},
		defaultResolution: "1440p",
		aspects: map[string][]string{
			"1440p": {"16:9", "1:1", "9:16"},
		},
		minDuration:           5,
		maxDuration:           15,
		defaultDuration:       5,
		maxPromptLength:       2000,
		maxStartFrames:        1,
		maxEndFrames:          1,
		maxImageRefs:          5,
		maxAudioRefs:          3,
		maxAudioRefSeconds:    15,
		endFrameRequiresStart: true,
		framesExcludeOtherRef: true,
		audioRefRequiresMedia: true,
		supportsAudio:         true,
	},
	"gemini-omni-flash": {
		resolutions:       []string{"720p"},
		defaultResolution: "720p",
		aspects: map[string][]string{
			"720p": {"16:9", "9:16"},
		},
		minDuration:     3,
		maxDuration:     10,
		defaultDuration: 5,
		maxPromptLength: 2500,
		maxImageRefs:    5,
	},
	"kling-2.1": {
		resolutions:       []string{"1080p"},
		defaultResolution: "1080p",
		aspects: map[string][]string{
			"1080p": {"16:9", "1:1", "9:16"},
		},
		allowedDurations:      []int{5, 10},
		minDuration:           5,
		maxDuration:           10,
		defaultDuration:       5,
		maxPromptLength:       2500,
		maxStartFrames:        1,
		maxEndFrames:          1,
		requiresStartFrame:    true,
		endFrameRequiresStart: true,
		supportsPromptEnhance: true,
	},
	"kling-2.5": {
		resolutions:       []string{"720p", "1080p"},
		defaultResolution: "1080p",
		aspects: map[string][]string{
			"720p":  {"16:9", "1:1", "9:16"},
			"1080p": {"16:9", "1:1", "9:16"},
		},
		allowedDurations:      []int{5, 10},
		minDuration:           5,
		maxDuration:           10,
		defaultDuration:       5,
		maxPromptLength:       2500,
		maxStartFrames:        1,
		maxEndFrames:          1,
		endFrameRequiresStart: true,
		supportsPromptEnhance: true,
	},
	"kling-2.5-turbo-standard": {
		resolutions:       []string{"720p"},
		defaultResolution: "720p",
		aspects: map[string][]string{
			"720p": {"16:9", "1:1", "9:16"},
		},
		allowedDurations:      []int{5, 10},
		minDuration:           5,
		maxDuration:           10,
		defaultDuration:       5,
		maxPromptLength:       2500,
		maxStartFrames:        1,
		requiresStartFrame:    true,
		supportsPromptEnhance: true,
	},
	"kling-2.6": {
		resolutions:       []string{"auto", "1080p"},
		defaultResolution: "1080p",
		aspects: map[string][]string{
			"auto":  {"auto"},
			"1080p": {"16:9", "1:1", "9:16"},
		},
		allowedDurations: []int{5, 10},
		minDuration:      5,
		maxDuration:      10,
		defaultDuration:  5,
		maxPromptLength:  2500,
		maxStartFrames:   1,
		supportsAudio:    true,
	},
	"kling-video-o-1": {
		resolutions:       []string{"1080p"},
		defaultResolution: "1080p",
		aspects: map[string][]string{
			"1080p": {"16:9", "1:1", "9:16"},
		},
		minDuration:           3,
		maxDuration:           10,
		defaultDuration:       5,
		maxPromptLength:       2500,
		maxStartFrames:        1,
		maxEndFrames:          1,
		maxImageRefs:          5,
		maxVideoRefs:          1,
		videoReferenceTypes:   []string{"GENERATED"},
		endFrameRequiresStart: true,
		framesExcludeOtherRef: true,
	},
	"kling-3.0": {
		resolutions:       []string{"auto", "720p", "1080p", "2160p"},
		defaultResolution: "1080p",
		aspects: map[string][]string{
			"auto":  {"auto"},
			"720p":  {"16:9", "1:1", "9:16"},
			"1080p": {"16:9", "1:1", "9:16"},
			"2160p": {"16:9", "1:1", "9:16"},
		},
		minDuration:           3,
		maxDuration:           15,
		defaultDuration:       5,
		maxPromptLength:       2500,
		maxStartFrames:        1,
		maxEndFrames:          1,
		endFrameRequiresStart: true,
		supportsAudio:         true,
	},
	"kling-3.0-turbo": {
		resolutions:       []string{"auto", "720p", "1080p"},
		defaultResolution: "1080p",
		aspects: map[string][]string{
			"auto":  {"auto"},
			"720p":  {"16:9", "1:1", "9:16"},
			"1080p": {"16:9", "1:1", "9:16"},
		},
		minDuration:     3,
		maxDuration:     15,
		defaultDuration: 5,
		maxPromptLength: 2500,
		maxStartFrames:  1,
		supportsAudio:   true,
	},
	"kling-video-o-3": {
		resolutions:       []string{"720p", "1080p", "2160p"},
		defaultResolution: "1080p",
		aspects: map[string][]string{
			"720p":  {"16:9", "1:1", "9:16"},
			"1080p": {"16:9", "1:1", "9:16"},
			"2160p": {"16:9", "1:1", "9:16"},
		},
		minDuration:           3,
		maxDuration:           15,
		defaultDuration:       5,
		maxPromptLength:       2500,
		maxStartFrames:        1,
		maxEndFrames:          1,
		maxImageRefs:          7,
		maxVideoRefs:          1,
		maxImageRefsWithVideo: 4,
		maxDurationWithVideo:  10,
		videoReferenceTypes:   []string{"GENERATED"},
		endFrameRequiresStart: true,
		framesExcludeOtherRef: true,
		supportsAudio:         true,
	},
	"veo-3.1-generate-001": {
		resolutions:       []string{"720p", "1080p", "2160p"},
		defaultResolution: "720p",
		aspects: map[string][]string{
			"720p":  {"16:9", "9:16"},
			"1080p": {"16:9", "9:16"},
			"2160p": {"16:9", "9:16"},
		},
		allowedDurations:      []int{4, 6, 8},
		minDuration:           4,
		maxDuration:           8,
		defaultDuration:       8,
		maxPromptLength:       9999,
		maxStartFrames:        1,
		maxEndFrames:          1,
		maxImageRefs:          3,
		endFrameRequiresStart: true,
		supportsAudio:         true,
	},
	"veo-3.1-fast-generate-001": {
		resolutions:       []string{"720p", "1080p", "2160p"},
		defaultResolution: "720p",
		aspects: map[string][]string{
			"720p":  {"16:9", "9:16"},
			"1080p": {"16:9", "9:16"},
			"2160p": {"16:9", "9:16"},
		},
		allowedDurations:      []int{4, 6, 8},
		minDuration:           4,
		maxDuration:           8,
		defaultDuration:       8,
		maxPromptLength:       9999,
		maxStartFrames:        1,
		maxEndFrames:          1,
		endFrameRequiresStart: true,
		supportsAudio:         true,
	},
	"veo-3.1-lite": {
		resolutions:       []string{"720p", "1080p"},
		defaultResolution: "720p",
		aspects: map[string][]string{
			"720p":  {"16:9", "9:16"},
			"1080p": {"16:9", "9:16"},
		},
		allowedDurations:      []int{4, 6, 8},
		minDuration:           4,
		maxDuration:           8,
		defaultDuration:       8,
		maxPromptLength:       9999,
		maxStartFrames:        1,
		maxEndFrames:          1,
		endFrameRequiresStart: true,
		supportsAudio:         true,
	},
}

type LeoVideoValidationError struct {
	message string
}

func (e *LeoVideoValidationError) Error() string {
	if e == nil {
		return "invalid video request"
	}
	return e.message
}

func newLeoVideoValidationError(format string, args ...any) error {
	return &LeoVideoValidationError{message: fmt.Sprintf(format, args...)}
}

func lookupLeoVideoModelSpec(model string) (leoVideoModelSpec, bool) {
	model = normalizeLeoVideoModelID(model)
	if model == "seedance" {
		model = "seedance-2.0"
	}
	spec, ok := leoVideoModelSpecs[model]
	return spec, ok
}

func normalizeLeoVideoResolution(resolution string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(resolution)) {
	case "auto", "resolution_auto":
		return "auto", true
	case "480", "480p", "sd", "resolution_480":
		return "480p", true
	case "400", "400p", "resolution_400":
		return "400p", true
	case "544", "544p", "resolution_544":
		return "544p", true
	case "720", "720p", "hd", "resolution_720":
		return "720p", true
	case "960", "960p", "resolution_960":
		return "960p", true
	case "1080", "1080p", "full_hd", "full-hd", "fhd", "resolution_1080":
		return "1080p", true
	case "1440", "1440p", "resolution_1440":
		return "1440p", true
	case "2160", "2160p", "4k", "uhd", "resolution_2160":
		return "2160p", true
	default:
		return "", false
	}
}

func leoVideoModelSupportsResolution(model, resolution string) bool {
	spec, ok := lookupLeoVideoModelSpec(model)
	if !ok {
		return false
	}
	normalized, ok := normalizeLeoVideoResolution(resolution)
	if !ok {
		return false
	}
	return containsLeoVideoValue(spec.resolutions, normalized)
}

func ValidateLeoVideoRequest(body []byte) (LeoVideoRequestInfo, error) {
	return validateLeoVideoRequest(body, "")
}

func ValidateLeoVideoRequestForModel(body []byte, effectiveModel string) (LeoVideoRequestInfo, error) {
	return validateLeoVideoRequest(body, effectiveModel)
}

func validateLeoVideoRequest(body []byte, effectiveModel string) (LeoVideoRequestInfo, error) {
	info, err := ParseLeoVideoRequest(body)
	if err != nil {
		return LeoVideoRequestInfo{}, err
	}
	if info.Model == "" {
		return LeoVideoRequestInfo{}, newLeoVideoValidationError("model is required")
	}
	if info.Prompt == "" {
		return LeoVideoRequestInfo{}, newLeoVideoValidationError("prompt is required")
	}
	modelForSpec := info.Model
	if strings.TrimSpace(effectiveModel) != "" {
		modelForSpec = effectiveModel
	}
	spec, knownModel := lookupLeoVideoModelSpec(modelForSpec)
	if !knownModel {
		spec = defaultLeoVideoModelSpec
	}
	if strings.TrimSpace(effectiveModel) != "" && !knownModel {
		return LeoVideoRequestInfo{}, newLeoVideoValidationError("video model %q is not supported", effectiveModel)
	}
	duration := gjson.GetBytes(body, "duration")
	durationProvided := duration.Exists()
	if durationProvided {
		if duration.Type != gjson.Number || duration.Float() != float64(duration.Int()) || duration.Int() < 1 {
			return LeoVideoRequestInfo{}, newLeoVideoValidationError("%s", durationValidationMessage(spec))
		}
		info.DurationSeconds = int(duration.Int())
	} else {
		info.DurationSeconds = spec.defaultDuration
	}
	if knownModel {
		if spec.maxPromptLength > 0 && utf8.RuneCountInString(info.Prompt) > spec.maxPromptLength {
			return LeoVideoRequestInfo{}, newLeoVideoValidationError("prompt supports at most %d characters for the selected video model", spec.maxPromptLength)
		}
		if !durationProvided {
			info.DurationSeconds = spec.defaultDuration
		} else if info.DurationSeconds < spec.minDuration || info.DurationSeconds > spec.maxDuration || !containsLeoVideoInt(spec.allowedDurations, info.DurationSeconds) {
			return LeoVideoRequestInfo{}, newLeoVideoValidationError("%s for the selected video model", durationValidationMessage(spec))
		}
		resolution := info.Resolution
		if strings.TrimSpace(resolution) == "" {
			resolution = spec.defaultResolution
			if resolution == "" {
				resolution = "720p"
			}
		}
		normalizedResolution, ok := normalizeLeoVideoResolution(resolution)
		if !ok || !containsLeoVideoValue(spec.resolutions, normalizedResolution) {
			return LeoVideoRequestInfo{}, newLeoVideoValidationError("resolution is not supported by the selected video model")
		}
		info.Resolution = normalizedResolution
		if info.AspectRatio == "" {
			info.AspectRatio = leoVideoDefaultAspect(spec.aspects[normalizedResolution])
		}
		if !containsLeoVideoValue(spec.aspects[normalizedResolution], info.AspectRatio) {
			return LeoVideoRequestInfo{}, newLeoVideoValidationError("aspect_ratio is not supported by the selected video model and resolution")
		}
		if audio := gjson.GetBytes(body, "audio"); audio.Exists() && audio.Bool() && !spec.supportsAudio {
			return LeoVideoRequestInfo{}, newLeoVideoValidationError("audio is not supported by the selected video model")
		}
		if maxDuration, limited := spec.maxDurationByResolution[normalizedResolution]; limited && info.DurationSeconds > maxDuration {
			return LeoVideoRequestInfo{}, newLeoVideoValidationError("duration must be a whole number from 4 through %d seconds for the selected video model and resolution", maxDuration)
		}
		if promptEnhance := strings.TrimSpace(gjson.GetBytes(body, "prompt_enhance").String()); promptEnhance != "" {
			if !spec.supportsPromptEnhance {
				return LeoVideoRequestInfo{}, newLeoVideoValidationError("prompt_enhance is not supported by the selected video model")
			}
			switch strings.ToUpper(promptEnhance) {
			case "AUTO", "ON", "OFF":
			default:
				return LeoVideoRequestInfo{}, newLeoVideoValidationError("prompt_enhance must be AUTO, ON, or OFF for the selected video model")
			}
			if spec.rejectsPromptEnhanceOnStartFrame && strings.EqualFold(promptEnhance, "ON") && leoVideoStartFrameCount(body) > 0 {
				return LeoVideoRequestInfo{}, newLeoVideoValidationError("prompt_enhance ON is not supported with start_frame by happy-horse-1.1")
			}
		}
		if spec.rejectsSeedAndMode && (gjson.GetBytes(body, "seed").Exists() || gjson.GetBytes(body, "mode").Exists()) {
			return LeoVideoRequestInfo{}, newLeoVideoValidationError("seed and mode are not supported by the selected video model")
		}
		if spec.requiresStartFrame && leoVideoStartFrameCount(body) == 0 {
			return LeoVideoRequestInfo{}, newLeoVideoValidationError("start frame is required by %s", normalizeLeoVideoModelID(modelForSpec))
		}
	}
	if !knownModel && utf8.RuneCountInString(info.Prompt) > spec.maxPromptLength {
		return LeoVideoRequestInfo{}, newLeoVideoValidationError("prompt supports at most %d characters", spec.maxPromptLength)
	}
	if !knownModel && durationProvided && info.DurationSeconds < spec.minDuration {
		return LeoVideoRequestInfo{}, newLeoVideoValidationError("duration must be a whole number from %d through %d seconds", spec.minDuration, spec.maxDuration)
	}

	if err := validateLeoVideoGuidanceCounts(body, spec); err != nil {
		return LeoVideoRequestInfo{}, err
	}
	if err := validateLeoVideoMediaGuidances(body, spec); err != nil {
		return LeoVideoRequestInfo{}, err
	}
	return info, nil
}

func validateLeoVideoGuidanceCounts(body []byte, spec leoVideoModelSpec) error {
	startFrames, err := leoVideoArrayLength(body, "guidances.start_frame")
	if err != nil {
		return err
	}
	if strings.TrimSpace(gjson.GetBytes(body, "image_url").String()) != "" {
		startFrames++
	}
	if strings.TrimSpace(gjson.GetBytes(body, "start_frame_url").String()) != "" {
		startFrames++
	}
	endFrames, err := leoVideoArrayLength(body, "guidances.end_frame")
	if err != nil {
		return err
	}
	if strings.TrimSpace(gjson.GetBytes(body, "end_frame_url").String()) != "" {
		endFrames++
	}
	imageReferences := 0
	for _, path := range []string{"image_urls", "guidances.image_reference"} {
		length, err := leoVideoArrayLength(body, path)
		if err != nil {
			return err
		}
		imageReferences += length
	}
	if spec.framesExcludeOtherRef && imageReferences > 0 && startFrames+endFrames > 0 {
		return newLeoVideoValidationError("reference images cannot be combined with start or end frames")
	}
	if startFrames > spec.maxStartFrames {
		return newLeoVideoValidationError("start frame must be supplied only once")
	}
	if endFrames > 0 && endFrames > spec.maxEndFrames {
		return newLeoVideoValidationError("guidances.end_frame supports at most %d item(s)", spec.maxEndFrames)
	}
	if spec.endFrameRequiresStart && endFrames > 0 && startFrames == 0 {
		return newLeoVideoValidationError("end frame requires a start frame for the selected video model")
	}

	limits := []struct {
		paths []string
		name  string
		max   int
	}{
		{paths: []string{"guidances.end_frame"}, name: "guidances.end_frame", max: spec.maxEndFrames},
		{paths: []string{"image_urls", "guidances.image_reference"}, name: "guidances.image_reference", max: spec.maxImageRefs},
		{paths: []string{"guidances.video_reference_base"}, name: "guidances.video_reference_base", max: spec.maxVideoRefs},
		{paths: []string{"guidances.audio_reference"}, name: "guidances.audio_reference", max: spec.maxAudioRefs},
	}
	for _, limit := range limits {
		count := 0
		for _, path := range limit.paths {
			length, err := leoVideoArrayLength(body, path)
			if err != nil {
				return err
			}
			count += length
		}
		if limit.name == "guidances.end_frame" && strings.TrimSpace(gjson.GetBytes(body, "end_frame_url").String()) != "" {
			count++
		}
		if count > limit.max {
			return newLeoVideoValidationError("%s supports at most %d item(s)", limit.name, limit.max)
		}
	}
	if spec.maxImageRefsWithVideo > 0 && videoReferenceCount(body) > 0 && imageReferences > spec.maxImageRefsWithVideo {
		return newLeoVideoValidationError("guidances.image_reference supports at most %d item(s) with video_reference_base", spec.maxImageRefsWithVideo)
	}
	if spec.maxDurationWithVideo > 0 && videoReferenceCount(body) > 0 {
		duration := gjson.GetBytes(body, "duration")
		if duration.Exists() && duration.Type == gjson.Number && int(duration.Int()) > spec.maxDurationWithVideo {
			return newLeoVideoValidationError("duration supports at most %d seconds with video_reference_base", spec.maxDurationWithVideo)
		}
	}
	return nil
}

func videoReferenceCount(body []byte) int {
	value, _ := leoVideoArrayLength(body, "guidances.video_reference_base")
	return value
}

func leoVideoStartFrameCount(body []byte) int {
	count, _ := leoVideoArrayLength(body, "guidances.start_frame")
	if strings.TrimSpace(gjson.GetBytes(body, "image_url").String()) != "" {
		count++
	}
	if strings.TrimSpace(gjson.GetBytes(body, "start_frame_url").String()) != "" {
		count++
	}
	return count
}

func leoVideoDefaultAspect(aspects []string) string {
	for _, preferred := range []string{"16:9", "9:16", "1:1", "4:3", "3:4", "21:9", "9:21"} {
		if containsLeoVideoValue(aspects, preferred) {
			return preferred
		}
	}
	if len(aspects) > 0 {
		return aspects[0]
	}
	return "16:9"
}

func durationValidationMessage(spec leoVideoModelSpec) string {
	if len(spec.allowedDurations) > 0 {
		values := make([]string, 0, len(spec.allowedDurations))
		for _, duration := range spec.allowedDurations {
			values = append(values, fmt.Sprintf("%d", duration))
		}
		return fmt.Sprintf("duration must be one of %s seconds", strings.Join(values, ", "))
	}
	return fmt.Sprintf("duration must be a whole number from %d through %d seconds", spec.minDuration, spec.maxDuration)
}

func validateLeoVideoMediaGuidances(body []byte, spec leoVideoModelSpec) error {
	for _, path := range []string{"image_url", "start_frame_url", "end_frame_url"} {
		if rawURL := strings.TrimSpace(gjson.GetBytes(body, path).String()); rawURL != "" {
			if !isAbsoluteHTTPURL(rawURL) {
				return newLeoVideoValidationError("%s must be an absolute HTTP(S) URL", path)
			}
		}
	}
	for _, path := range []string{"guidances.start_frame", "guidances.end_frame"} {
		frames := gjson.GetBytes(body, path)
		for index, item := range frames.Array() {
			asset := item.Get("image")
			if err := validateLeoVideoImageAsset(asset, "image", spec.frameImageTypes); err != nil {
				return newLeoVideoValidationError("%s[%d] %s", path, index, err.Error())
			}
		}
	}
	for index, item := range gjson.GetBytes(body, "image_urls").Array() {
		if !isAbsoluteHTTPURL(strings.TrimSpace(item.String())) {
			return newLeoVideoValidationError("image_urls[%d] must be an absolute HTTP(S) URL", index)
		}
	}
	for index, item := range gjson.GetBytes(body, "guidances.image_reference").Array() {
		asset := item.Get("image")
		if err := validateLeoVideoImageAsset(asset, "image", nil); err != nil {
			return newLeoVideoValidationError("guidances.image_reference[%d] %s", index, err.Error())
		}
	}

	videoReferences := gjson.GetBytes(body, "guidances.video_reference_base")
	for index, item := range videoReferences.Array() {
		asset := item.Get("video")
		if err := validateLeoVideoMediaAsset(asset, "video", spec.videoReferenceTypes); err != nil {
			return newLeoVideoValidationError("guidances.video_reference_base[%d] %s", index, err.Error())
		}
	}

	audioReferences := gjson.GetBytes(body, "guidances.audio_reference")
	for index, item := range audioReferences.Array() {
		asset := item.Get("audio")
		if err := validateLeoVideoAudioAsset(asset); err != nil {
			return newLeoVideoValidationError("guidances.audio_reference[%d] %s", index, err.Error())
		}
	}
	if len(audioReferences.Array()) > 0 {
		imageReferences := len(gjson.GetBytes(body, "image_urls").Array()) + len(gjson.GetBytes(body, "guidances.image_reference").Array())
		videoReferenceCount := len(videoReferences.Array())
		if spec.audioRefRequiresMedia && imageReferences == 0 && videoReferenceCount == 0 {
			return newLeoVideoValidationError("guidances.audio_reference requires an image_reference or video_reference_base")
		}
		if spec.maxAudioRefSeconds > 0 {
			totalDuration := 0.0
			for index, item := range audioReferences.Array() {
				audio := item.Get("audio")
				if audio.Get("id").String() != "" && !audio.Get("duration").Exists() {
					return newLeoVideoValidationError("guidances.audio_reference[%d] audio.duration is required with audio.id", index)
				}
				if duration := audio.Get("duration"); duration.Exists() && duration.Type == gjson.Number {
					totalDuration += duration.Float()
				}
			}
			if totalDuration > spec.maxAudioRefSeconds {
				return newLeoVideoValidationError("guidances.audio_reference supports at most %.0f seconds total", spec.maxAudioRefSeconds)
			}
		}
	}
	return nil
}

func validateLeoVideoMediaAsset(value gjson.Result, kind string, allowedTypes []string) error {
	if !value.Exists() || !value.IsObject() {
		return fmt.Errorf("requires %s.id or %s.url", kind, kind)
	}
	id := strings.TrimSpace(value.Get("id").String())
	rawURL := strings.TrimSpace(value.Get("url").String())
	if id == "" && rawURL == "" {
		return fmt.Errorf("requires %s.id or %s.url", kind, kind)
	}
	if id != "" && rawURL != "" {
		return fmt.Errorf("%s.id and %s.url cannot both be set", kind, kind)
	}
	if value.Get("duration").Exists() {
		return fmt.Errorf("%s.duration is not supported", kind)
	}
	typeName := strings.ToUpper(strings.TrimSpace(value.Get("type").String()))
	if id != "" {
		if !leoVideoAssetIDPattern.MatchString(id) {
			return fmt.Errorf("%s.id must be a UUID", kind)
		}
		if typeName == "" {
			typeName = "UPLOADED"
		}
		if !leoVideoTypeAllowed(typeName, allowedTypes) {
			return fmt.Errorf("%s.type is invalid", kind)
		}
		return nil
	}
	if typeName == "" {
		typeName = "UPLOADED"
	}
	if typeName != "UPLOADED" || (len(allowedTypes) > 0 && !leoVideoTypeAllowed("UPLOADED", allowedTypes)) {
		return fmt.Errorf("%s.url requires type UPLOADED", kind)
	}
	if !isAbsoluteHTTPURL(rawURL) {
		return fmt.Errorf("%s.url must be an absolute HTTP(S) URL", kind)
	}
	return nil
}

func validateLeoVideoImageAsset(value gjson.Result, kind string, allowedTypes []string) error {
	if !value.Exists() || !value.IsObject() {
		return fmt.Errorf("requires %s.id or %s.url", kind, kind)
	}
	id := strings.TrimSpace(value.Get("id").String())
	rawURL := strings.TrimSpace(value.Get("url").String())
	if id == "" && rawURL == "" {
		return fmt.Errorf("requires %s.id or %s.url", kind, kind)
	}
	if id != "" && rawURL != "" {
		return fmt.Errorf("%s.id and %s.url cannot both be set", kind, kind)
	}
	typeName := strings.ToUpper(strings.TrimSpace(value.Get("type").String()))
	if typeName == "" {
		typeName = "UPLOADED"
	}
	if id != "" && !leoVideoAssetIDPattern.MatchString(id) {
		return fmt.Errorf("%s.id must be a UUID", kind)
	}
	if !leoVideoTypeAllowed(typeName, allowedTypes) {
		return fmt.Errorf("%s.type is invalid", kind)
	}
	if rawURL != "" {
		if typeName != "UPLOADED" || (len(allowedTypes) > 0 && !leoVideoTypeAllowed("UPLOADED", allowedTypes)) {
			return fmt.Errorf("%s.url requires type UPLOADED", kind)
		}
		if !isAbsoluteHTTPURL(rawURL) {
			return fmt.Errorf("%s.url must be an absolute HTTP(S) URL", kind)
		}
	}
	return nil
}

func leoVideoTypeAllowed(typeName string, allowedTypes []string) bool {
	if len(allowedTypes) == 0 {
		return typeName == "UPLOADED" || typeName == "GENERATED"
	}
	for _, allowed := range allowedTypes {
		if typeName == strings.ToUpper(strings.TrimSpace(allowed)) {
			return true
		}
	}
	return false
}

func validateLeoVideoAudioAsset(value gjson.Result) error {
	if !value.Exists() || !value.IsObject() {
		return fmt.Errorf("requires audio.id or audio.url")
	}
	id := strings.TrimSpace(value.Get("id").String())
	rawURL := strings.TrimSpace(value.Get("url").String())
	if id == "" && rawURL == "" {
		return fmt.Errorf("requires audio.id or audio.url")
	}
	if id != "" && rawURL != "" {
		return fmt.Errorf("audio.id and audio.url cannot both be set")
	}
	duration := value.Get("duration")
	if rawURL != "" && duration.Exists() {
		return fmt.Errorf("audio.duration must be omitted when audio.url is used")
	}
	if duration.Exists() {
		if duration.Type != gjson.Number || duration.Float() < 2 || duration.Float() > 30 {
			return fmt.Errorf("audio.duration must be between 2 and 30 seconds")
		}
	}
	typeName := strings.ToUpper(strings.TrimSpace(value.Get("type").String()))
	if typeName == "" {
		typeName = "UPLOADED"
	}
	if typeName != "UPLOADED" {
		return fmt.Errorf("audio.type must be UPLOADED")
	}
	if id != "" && !leoVideoAssetIDPattern.MatchString(id) {
		return fmt.Errorf("audio.id must be a UUID")
	}
	if rawURL != "" && !isAbsoluteHTTPURL(rawURL) {
		return fmt.Errorf("audio.url must be an absolute HTTP(S) URL")
	}
	return nil
}

func isAbsoluteHTTPURL(raw string) bool {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	return err == nil && parsed.Host != "" && (parsed.Scheme == "http" || parsed.Scheme == "https")
}

func leoVideoArrayLength(body []byte, path string) (int, error) {
	value := gjson.GetBytes(body, path)
	if !value.Exists() {
		return 0, nil
	}
	if !value.IsArray() {
		return 0, newLeoVideoValidationError("%s must be an array", path)
	}
	return len(value.Array()), nil
}

func containsLeoVideoValue(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func containsLeoVideoInt(values []int, target int) bool {
	if len(values) == 0 {
		return true
	}
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

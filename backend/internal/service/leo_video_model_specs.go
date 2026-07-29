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
	requiresStartFrame               bool
	supportsPromptEnhance            bool
	rejectsPromptEnhanceOnStartFrame bool
	rejectsSeedAndMode               bool
}

var defaultLeoVideoModelSpec = leoVideoModelSpec{
	defaultResolution: "720p",
	minDuration:       4,
	maxDuration:       15,
	defaultDuration:   8,
	maxPromptLength:   leoVideoMaxPromptLength,
	maxStartFrames:    1,
	maxEndFrames:      1,
	maxImageRefs:      4,
	maxVideoRefs:      3,
	maxAudioRefs:      1,
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
		minDuration:     4,
		maxDuration:     15,
		defaultDuration: 8,
		maxPromptLength: 5000,
		maxStartFrames:  1,
		maxEndFrames:    1,
		maxImageRefs:    4,
		maxVideoRefs:    3,
		maxAudioRefs:    1,
	},
	"seedance-2.0-fast": {
		resolutions:       []string{"480p", "720p"},
		defaultResolution: "720p",
		aspects: map[string][]string{
			"480p": {"16:9", "9:16", "1:1", "4:3", "3:4", "21:9", "9:21"},
			"720p": {"16:9", "9:16", "1:1", "4:3", "3:4", "21:9"},
		},
		minDuration:     4,
		maxDuration:     15,
		defaultDuration: 8,
		maxPromptLength: 5000,
		maxStartFrames:  1,
		maxEndFrames:    1,
		maxImageRefs:    4,
		maxVideoRefs:    3,
		maxAudioRefs:    1,
	},
	"seedance-2.0-mini": {
		resolutions:       []string{"480p", "720p"},
		defaultResolution: "720p",
		aspects: map[string][]string{
			"480p": {"16:9", "1:1", "9:16"},
			"720p": {"16:9", "1:1", "9:16"},
		},
		minDuration:     4,
		maxDuration:     15,
		defaultDuration: 8,
		maxPromptLength: 5000,
		maxStartFrames:  1,
		maxEndFrames:    1,
		maxImageRefs:    4,
		maxVideoRefs:    3,
		maxAudioRefs:    1,
	},
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
		supportsPromptEnhance:            true,
		rejectsPromptEnhanceOnStartFrame: true,
	},
	"grok-imagine-1.5": {
		resolutions:       []string{"400p", "544p", "720p", "960p"},
		defaultResolution: "720p",
		aspects: map[string][]string{
			"400p": {"16:9", "9:16"},
			"544p": {"1:1"},
			"720p": {"16:9", "9:16"},
			"960p": {"1:1"},
		},
		minDuration:        3,
		maxDuration:        15,
		defaultDuration:    6,
		maxPromptLength:    5000,
		maxStartFrames:     1,
		requiresStartFrame: true,
	},
	"ltxv-2.3-pro": {
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
		supportsPromptEnhance: true,
		rejectsSeedAndMode:    true,
	},
	"ltxv-2.3-fast": {
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
		supportsPromptEnhance: true,
		rejectsSeedAndMode:    true,
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
	model = strings.ToLower(strings.TrimSpace(model))
	if model == "seedance" {
		model = "seedance-2.0"
	}
	spec, ok := leoVideoModelSpecs[model]
	return spec, ok
}

func normalizeLeoVideoResolution(resolution string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(resolution)) {
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
			return LeoVideoRequestInfo{}, newLeoVideoValidationError("start frame is required by grok-imagine-1.5")
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
	if err := validateLeoVideoMediaGuidances(body); err != nil {
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
	if imageReferences > 0 && startFrames+endFrames > 0 {
		return newLeoVideoValidationError("reference images cannot be combined with start or end frames")
	}
	if startFrames > spec.maxStartFrames {
		return newLeoVideoValidationError("start frame must be supplied only once")
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
	return nil
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

func validateLeoVideoMediaGuidances(body []byte) error {
	videoReferences := gjson.GetBytes(body, "guidances.video_reference_base")
	for index, item := range videoReferences.Array() {
		asset := item.Get("video")
		if err := validateLeoVideoMediaAsset(asset, "video"); err != nil {
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
		if imageReferences == 0 && videoReferenceCount == 0 {
			return newLeoVideoValidationError("guidances.audio_reference requires an image_reference or video_reference_base")
		}
	}
	return nil
}

func validateLeoVideoMediaAsset(value gjson.Result, kind string) error {
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
		if typeName != "UPLOADED" && typeName != "GENERATED" {
			return fmt.Errorf("%s.type is invalid", kind)
		}
		return nil
	}
	if typeName == "" {
		typeName = "UPLOADED"
	}
	if typeName != "UPLOADED" {
		return fmt.Errorf("%s.url requires type UPLOADED", kind)
	}
	if !isAbsoluteHTTPURL(rawURL) {
		return fmt.Errorf("%s.url must be an absolute HTTP(S) URL", kind)
	}
	return nil
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

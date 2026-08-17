package service

import "strings"

const (
	VideoBillingResolution400P  = "400p"
	VideoBillingResolution480P  = "480p"
	VideoBillingResolution544P  = "544p"
	VideoBillingResolution720P  = "720p"
	VideoBillingResolution960P  = "960p"
	VideoBillingResolution1080P = "1080p"
	VideoBillingResolution1440P = "1440p"
	VideoBillingResolution2160P = "2160p"
)

var leoVideoPricingResolutions = []string{
	VideoBillingResolution480P,
	VideoBillingResolution720P,
	VideoBillingResolution1080P,
}

// LeoVideoPricingResolutions returns the resolution tiers supported by a Leo
// video model. The returned slice is a copy so callers cannot mutate the
// shared capability table.
func LeoVideoPricingResolutions(model string) []string {
	model = strings.ToLower(strings.TrimSpace(model))
	if model == "bytedance/seedance-2.5" || model == "seedance-2.5" {
		return []string{VideoBillingResolution480P, VideoBillingResolution720P}
	}
	if model == "seedance-2.0" {
		return []string{VideoBillingResolution480P, VideoBillingResolution720P, VideoBillingResolution1080P, VideoBillingResolution2160P}
	}
	if model == "seedance-2.0-mini" {
		return []string{VideoBillingResolution720P}
	}
	if model == "happy-horse-1.1" {
		return []string{VideoBillingResolution720P, VideoBillingResolution1080P}
	}
	if model == "grok-imagine-1.5" {
		return []string{
			VideoBillingResolution400P,
			VideoBillingResolution544P,
			VideoBillingResolution720P,
			VideoBillingResolution960P,
		}
	}
	if isLeoLTX23Model(model) {
		return []string{
			VideoBillingResolution1080P,
			VideoBillingResolution1440P,
			VideoBillingResolution2160P,
		}
	}
	switch model {
	case "hailuo-03":
		return []string{VideoBillingResolution1440P}
	case "gemini-omni-flash", "kling-2.5-turbo-standard":
		return []string{VideoBillingResolution720P}
	case "kling-2.1", "kling-2.6", "kling-video-o-1":
		return []string{VideoBillingResolution1080P}
	case "kling-2.5", "kling-3.0-turbo":
		return []string{VideoBillingResolution720P, VideoBillingResolution1080P}
	case "kling-3.0", "kling-video-o-3", "veo-3.1-generate-001", "veo-3.1-fast-generate-001":
		return []string{VideoBillingResolution720P, VideoBillingResolution1080P, VideoBillingResolution2160P}
	case "veo-3.1-lite":
		return []string{VideoBillingResolution720P, VideoBillingResolution1080P}
	}
	return append([]string(nil), leoVideoPricingResolutions...)
}

func DefaultLeoVideoResolution(model string) string {
	if spec, ok := lookupLeoVideoModelSpec(model); ok && strings.TrimSpace(spec.defaultResolution) != "" {
		return spec.defaultResolution
	}
	resolutions := LeoVideoPricingResolutions(model)
	if len(resolutions) == 1 {
		return resolutions[0]
	}
	return VideoBillingResolution720P
}

func LeoVideoModelSupportsResolution(model, resolution string) bool {
	return leoVideoModelSupportsResolution(model, resolution)
}

// 视频生成按秒计费；通用视频路径沿用 15 秒上限。
// 计费时长必须与实际消耗对齐，否则用户可通过拉长 duration 套利（提交时长由用户控制）。
const (
	VideoBillingMinDurationSeconds     = 1
	VideoBillingMaxDurationSeconds     = 15
	VideoBillingDefaultDurationSeconds = 8
)

// NormalizeVideoBillingDurationSecondsOrDefault 归一化计费用视频时长：
// 未指定（<=0）按上游默认 8 秒计，超出上游允许区间按边界收敛。
func NormalizeVideoBillingDurationSecondsOrDefault(durationSeconds int) int {
	return normalizeVideoBillingDurationSecondsOrDefault(durationSeconds, VideoBillingMaxDurationSeconds)
}

// NormalizeLeoVideoBillingDurationSecondsOrDefault widens the shared billing
// limit only for Leo models whose upstream contract explicitly allows it.
func NormalizeLeoVideoBillingDurationSecondsOrDefault(model string, durationSeconds int) int {
	maxDuration := VideoBillingMaxDurationSeconds
	if spec, ok := lookupLeoVideoModelSpec(model); ok && spec.maxDuration > maxDuration {
		maxDuration = spec.maxDuration
	}
	return normalizeVideoBillingDurationSecondsOrDefault(durationSeconds, maxDuration)
}

func normalizeVideoBillingDurationSecondsOrDefault(durationSeconds, maxDuration int) int {
	if durationSeconds <= 0 {
		return VideoBillingDefaultDurationSeconds
	}
	if durationSeconds < VideoBillingMinDurationSeconds {
		return VideoBillingMinDurationSeconds
	}
	if durationSeconds > maxDuration {
		return maxDuration
	}
	return durationSeconds
}

// LookupVideoBillingResolution 归一化分辨率并报告是否为已知档位。
// 配置解析路径必须用它而不是 OrDefault：把无法识别的档位（如 "4k"、拼错的
// "1080i"）静默折算成 480p，会让管理员配的高分辨率单价被挂到低分辨率档上。
func LookupVideoBillingResolution(resolution string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(resolution)) {
	case "480", "480p", "sd", "resolution_480":
		return VideoBillingResolution480P, true
	case "720", "720p", "hd", "resolution_720":
		return VideoBillingResolution720P, true
	case "1080", "1080p", "full_hd", "full-hd", "fhd", "resolution_1080":
		return VideoBillingResolution1080P, true
	default:
		return "", false
	}
}

// NormalizeVideoBillingResolutionOrDefault is the generic three-tier billing
// normalization. Leo-specific model tiers are handled by the model-aware helper.
func NormalizeVideoBillingResolutionOrDefault(resolution string) string {
	if normalized, ok := LookupVideoBillingResolution(resolution); ok {
		return normalized
	}
	switch strings.ToLower(strings.TrimSpace(resolution)) {
	case "400", "400p", "resolution_400", "544", "544p", "resolution_544":
		return VideoBillingResolution480P
	case "960", "960p", "resolution_960", "1440", "1440p", "resolution_1440", "2160", "2160p":
		return VideoBillingResolution1080P
	default:
		return VideoBillingResolution480P
	}
}

func NormalizeLeoVideoBillingResolutionOrDefault(model, resolution string) string {
	model = normalizeLeoVideoModelID(model)
	resolution = strings.ToLower(strings.TrimSpace(resolution))
	switch model {
	case "hailuo-03":
		return VideoBillingResolution1440P
	case "gemini-omni-flash", "kling-2.5-turbo-standard":
		return VideoBillingResolution720P
	case "kling-2.1", "kling-2.6", "kling-video-o-1":
		return VideoBillingResolution1080P
	case "kling-2.5", "kling-3.0-turbo", "veo-3.1-lite":
		if strings.HasPrefix(resolution, "720") {
			return VideoBillingResolution720P
		}
		return VideoBillingResolution1080P
	case "kling-3.0", "kling-video-o-3", "veo-3.1-generate-001", "veo-3.1-fast-generate-001":
		switch resolution {
		case "720", "720p", "hd", "resolution_720":
			return VideoBillingResolution720P
		case "2160", "2160p", "4k", "uhd", "resolution_2160":
			return VideoBillingResolution2160P
		default:
			return VideoBillingResolution1080P
		}
	case "happy-horse-1.1":
		switch resolution {
		case "720", "720p", "hd", "resolution_720":
			return VideoBillingResolution720P
		case "1080", "1080p", "full_hd", "full-hd", "fhd", "resolution_1080":
			return VideoBillingResolution1080P
		default:
			return VideoBillingResolution1080P
		}
	case "grok-imagine-1.5":
		switch resolution {
		case "400", "400p", "resolution_400":
			return VideoBillingResolution400P
		case "544", "544p", "resolution_544":
			return VideoBillingResolution544P
		case "720", "720p", "hd", "resolution_720":
			return VideoBillingResolution720P
		case "960", "960p", "resolution_960":
			return VideoBillingResolution960P
		default:
			return VideoBillingResolution720P
		}
	}
	if model == "seedance-2.0" {
		switch resolution {
		case "480", "480p", "sd", "resolution_480":
			return VideoBillingResolution480P
		case "720", "720p", "hd", "resolution_720":
			return VideoBillingResolution720P
		case "1080", "1080p", "full_hd", "full-hd", "fhd", "resolution_1080":
			return VideoBillingResolution1080P
		case "2160", "2160p", "4k", "uhd", "resolution_2160":
			return VideoBillingResolution2160P
		default:
			return VideoBillingResolution720P
		}
	}
	if !isLeoLTX23Model(model) {
		return NormalizeVideoBillingResolutionOrDefault(resolution)
	}

	switch resolution {
	case "1080", "1080p", "full_hd", "full-hd", "fhd", "resolution_1080":
		return VideoBillingResolution1080P
	case "1440", "1440p", "resolution_1440":
		return VideoBillingResolution1440P
	case "2160", "2160p", "4k", "uhd", "resolution_2160":
		return VideoBillingResolution2160P
	default:
		return VideoBillingResolution1080P
	}
}

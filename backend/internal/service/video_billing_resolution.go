package service

import "strings"

const (
	VideoBillingResolution480P  = "480p"
	VideoBillingResolution720P  = "720p"
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
	if model == "seedance-2.0-mini" {
		return []string{VideoBillingResolution720P}
	}
	if model == "ltxv-2.3-pro" || model == "ltxv-2.3-fast" {
		return []string{
			VideoBillingResolution1080P,
			VideoBillingResolution1440P,
			VideoBillingResolution2160P,
		}
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

func NormalizeVideoBillingResolutionOrDefault(resolution string) string {
	switch strings.ToLower(strings.TrimSpace(resolution)) {
	case "400", "400p", "resolution_400", "480", "480p", "sd", "resolution_480", "544", "544p", "resolution_544":
		return VideoBillingResolution480P
	case "720", "720p", "hd", "resolution_720":
		return VideoBillingResolution720P
	case "960", "960p", "resolution_960", "1080", "1080p", "full_hd", "full-hd", "fhd", "resolution_1080", "1440", "1440p", "resolution_1440", "2160", "2160p", "4k", "uhd", "resolution_2160":
		return VideoBillingResolution1080P
	default:
		return VideoBillingResolution480P
	}
}

func NormalizeLeoVideoBillingResolutionOrDefault(model, resolution string) string {
	model = strings.ToLower(strings.TrimSpace(model))
	if model != "ltxv-2.3-pro" && model != "ltxv-2.3-fast" {
		return NormalizeVideoBillingResolutionOrDefault(resolution)
	}

	switch strings.ToLower(strings.TrimSpace(resolution)) {
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

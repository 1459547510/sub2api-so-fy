package service

import "regexp"

var publicVendorNamePattern = regexp.MustCompile(`(?i)\b(?:leonardo(?:\.ai| ai)?|leo\s*studio|krea|leo)\b`)

func sanitizePublicVendorMessage(message, replacement string) string {
	return publicVendorNamePattern.ReplaceAllString(message, replacement)
}

func SanitizeVideoProviderMessage(message string) string {
	return sanitizePublicVendorMessage(message, "Video service")
}

func SanitizeImageProviderMessage(message string) string {
	return sanitizePublicVendorMessage(message, "Image service")
}

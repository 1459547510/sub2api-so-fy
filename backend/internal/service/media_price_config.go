package service

func imagePriceConfigFromAPIKey(apiKey *APIKey) *ImagePriceConfig {
	if apiKey == nil || apiKey.Group == nil {
		return nil
	}
	return &ImagePriceConfig{
		Price1K: apiKey.Group.ImagePrice1K,
		Price2K: apiKey.Group.ImagePrice2K,
		Price4K: apiKey.Group.ImagePrice4K,
	}
}

func apiKeyHasConfiguredImagePrice(apiKey *APIKey, imageSize string) bool {
	return apiKey != nil && apiKey.Group != nil && apiKey.Group.GetImagePrice(imageSize) != nil
}

func videoPriceConfigFromAPIKey(apiKey *APIKey) *VideoPriceConfig {
	if apiKey == nil || apiKey.Group == nil {
		return nil
	}
	return &VideoPriceConfig{
		Price480P:   apiKey.Group.VideoPrice480P,
		Price720P:   apiKey.Group.VideoPrice720P,
		Price1080P:  apiKey.Group.VideoPrice1080P,
		ModelPrices: apiKey.Group.VideoModelPrices,
	}
}

func cloneVideoPriceConfig(in *VideoPriceConfig) *VideoPriceConfig {
	if in == nil {
		return &VideoPriceConfig{}
	}
	out := *in
	return &out
}

// overlayGroupVideoModelPrices copies group video_model_prices onto a channel
// (or empty) config. Flat group video_price_* columns are not stamped onto
// models that already have per-model or channel tiers.
func overlayGroupVideoModelPrices(config *VideoPriceConfig, group *Group, model string) *VideoPriceConfig {
	out := cloneVideoPriceConfig(config)
	if group == nil {
		return out
	}
	prices := NormalizeVideoModelPrices(group.VideoModelPrices)
	if len(prices) == 0 {
		return out
	}
	if out.ModelPrices == nil {
		out.ModelPrices = prices
	} else {
		merged := make(map[string]map[string]float64, len(out.ModelPrices)+len(prices))
		for family, tiers := range out.ModelPrices {
			merged[family] = tiers
		}
		for family, tiers := range prices {
			merged[family] = tiers
		}
		out.ModelPrices = merged
	}
	if price := LookupVideoModelPrice(prices, model, VideoBillingResolution480P); price != nil {
		out.Price480P = price
	}
	if price := LookupVideoModelPrice(prices, model, VideoBillingResolution720P); price != nil {
		out.Price720P = price
	}
	if price := LookupVideoModelPrice(prices, model, VideoBillingResolution1080P); price != nil {
		out.Price1080P = price
	}
	return out
}

func groupVideoModelPrices(apiKey *APIKey) map[string]map[string]float64 {
	if apiKey == nil || apiKey.Group == nil {
		return nil
	}
	return apiKey.Group.VideoModelPrices
}

func apiKeyHasConfiguredVideoPrice(apiKey *APIKey, model, resolution string) bool {
	return apiKey != nil && apiKey.Group != nil && apiKey.Group.GetVideoPriceForModel(model, resolution) != nil
}

func webSearchPricePerCallFromAPIKey(apiKey *APIKey) *float64 {
	if apiKey == nil || apiKey.Group == nil {
		return nil
	}
	return apiKey.Group.WebSearchPricePerCall
}

func groupSearchPricePer1kFromAPIKey(apiKey *APIKey) *float64 {
	if apiKey == nil || apiKey.Group == nil {
		return nil
	}
	return apiKey.Group.GetSearchPricePer1k()
}

func groupAudioPriceConfigFromAPIKey(apiKey *APIKey) *audioPriceConfig {
	if apiKey == nil || apiKey.Group == nil {
		return nil
	}
	g := apiKey.Group
	return &audioPriceConfig{
		RealtimePerMin: g.AudioRealtimePricePerMin,
		TTSPerMChars:   g.AudioTTSPricePerMillionChars,
		STTPerHour:     g.AudioSTTPricePerHour,
	}
}

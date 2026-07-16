import { describe, expect, it } from "vitest";

import {
  getDefaultImagePreviewPrice,
  getDefaultVideoPreviewPrice,
  getImagePricePlaceholder,
  getVideoPricePlaceholder,
  hasCompleteLeoVideoPrices,
  imagePricingPlatforms,
  imagePricingI18nKey,
  supportsImagePricingPlatform,
  supportsVideoPricingPlatform,
  videoPricingI18nKey,
} from "../groupsImagePricing";

describe("groups image pricing platform support", () => {
  it("includes Grok image groups", () => {
    expect(supportsImagePricingPlatform("grok")).toBe(true);
    expect(imagePricingPlatforms.has("grok")).toBe(true);
  });

  it("enables video pricing controls for Grok only", () => {
    expect(supportsVideoPricingPlatform("grok")).toBe(true);
    expect(supportsVideoPricingPlatform("leo")).toBe(true);
    expect(supportsVideoPricingPlatform("openai")).toBe(false);
  });

  it("requires operator-supplied Leo video prices without enabling image pricing", () => {
    expect(supportsImagePricingPlatform("leo")).toBe(false);
    expect(imagePricingPlatforms.has("leo")).toBe(false);
    expect(getVideoPricePlaceholder("leo", "video_price_480p")).toBe("");
    expect(getDefaultVideoPreviewPrice("leo", "video_price_480p")).toBeNull();
  });

  it("accepts only complete non-negative Leo video prices", () => {
    expect(hasCompleteLeoVideoPrices({
      platform: "leo",
      video_price_480p: 0,
      video_price_720p: "0.1",
      video_price_1080p: 0.2,
    })).toBe(true);
    expect(hasCompleteLeoVideoPrices({
      platform: "leo",
      video_price_480p: 0,
      video_price_720p: null,
      video_price_1080p: 0.2,
    })).toBe(false);
    expect(hasCompleteLeoVideoPrices({
      platform: "leo",
      video_price_480p: 0,
      video_price_720p: -1,
      video_price_1080p: 0.2,
    })).toBe(false);
    expect(hasCompleteLeoVideoPrices({
      platform: "grok",
      video_price_480p: null,
      video_price_720p: null,
      video_price_1080p: null,
    })).toBe(true);
  });

  it("keeps non-media group platforms out of the image pricing controls", () => {
    expect(supportsImagePricingPlatform("anthropic")).toBe(false);
  });

  it("keeps image and video pricing copy separate", () => {
    expect(imagePricingI18nKey("grok", "title")).toBe(
      "admin.groups.imagePricing.title",
    );
    expect(videoPricingI18nKey("title")).toBe("admin.groups.videoPricing.title");
  });

  it("uses Grok media defaults instead of generic image fallback placeholders", () => {
    expect(getImagePricePlaceholder("grok", "image_price_1k")).toBe("0.02");
    expect(getImagePricePlaceholder("grok", "image_price_2k")).toBe("0.02");
    // 视频 placeholder 为每秒单价：480p/720p 取 grok-imagine-video 官方每秒价，
    // 1080p 仅 video-1.5 支持、取 1.5 每秒价。
    expect(getVideoPricePlaceholder("grok", "video_price_480p")).toBe("0.05");
    expect(getVideoPricePlaceholder("grok", "video_price_720p")).toBe("0.07");
    expect(getVideoPricePlaceholder("grok", "video_price_1080p")).toBe("0.25");
  });

  it("keeps non-Grok image placeholders on the generic image card", () => {
    expect(getImagePricePlaceholder("openai", "image_price_1k")).toBe("0.134");
    expect(getDefaultImagePreviewPrice("openai", "image_price_2k")).toBe(0.201);
    expect(getDefaultVideoPreviewPrice("openai", "video_price_480p")).toBeNull();
  });
});

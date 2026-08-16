# Media Pricing Page

Date: 2026-08-16

## Goal

Move media prices out of the V2 video/image API docs and into a third video-module tab. Show the administrator-configured unit prices, one card per model. Image and video use different resolution systems. Do not apply any rate multipliers.

## Navigation

- Add `/video-generation/pricing` (`VideoPricing`) next to the existing workbench and API docs routes.
- Add a third `VideoSectionTabs` item: 生成工作台 / API 对接文档 / 价格.
- Keep the same auth, layout, and `video_generation_enabled` visibility as the other two tabs.
- Do not add a sidebar or header entry. `/video-generation` already stays active for child paths.
- Remove the `#pricing` section and its TOC item from `VideoApiDocsView`. Do not leave a link or summary.

## Page

Title, one-line source note, and a refresh button. Then one section per visible media group. Inside a group: image cards first, then video cards.

Empty, loading, and fetch-failure states stay on this page. If `GET /groups/available` fails, show unavailable even when plaza or channels succeed. Plaza or channels failure only drops those model names; it does not make the page unavailable. A failed refresh keeps previously loaded groups and model names. If the user has no media group with any configured media price, show empty.

## Cards

Each card shows model name, type (image/video), and unit (per image / per second).

### Image

Read `image_price_1k`, `image_price_2k`, `image_price_4k`.

- If two or three configured tiers share one value, collapse to a single unlabeled price.
- If only one tier is configured, show that price with its 1K / 2K / 4K label.
- If configured values differ, split by 1K / 2K / 4K.
- Omit unconfigured tiers. Do not render a dash.

### Video

Unit is USD/s. Always list configured resolutions; do not collapse even when values match.

- Display only 480p / 720p / 1080p. Ignore other resolution keys.
- Prefer `video_model_prices[model][resolution]` for those three.
- Fall back to `video_price_480p` / `video_price_720p` / `video_price_1080p`.
- Omit unconfigured resolutions.

### Price math

Display the raw configured USD amounts. Do not multiply by user exclusive rate, group `rate_multiplier`, image/video independent rates, or peak rate. Do not fetch `/groups/rates` for this page. Do not show multiplier text.

Treat `null`, `undefined`, and `""` as unconfigured before calling `Number()`. `Number(null)` and `Number("")` are `0` in JavaScript and must not count as a price. After that guard, a tier is configured when `Number(value)` is finite and `>= 0`. Numeric strings such as `"0.08"` count. `0` is a real price.

## Data

No new backend endpoint.

1. `GET /groups/available` is the only price source. Keep a group if its platform is `leo`, `openai_media`, `video`, `composite`, or `grok`, and it has at least one configured media price: any `image_price_*`, any flat `video_price_*`, or any `video_model_prices` entry with a configured 480p / 720p / 1080p override. Nested model prices count; a group with only `video_model_prices` is kept.
2. Model names are dynamic. Union `GET /model-plaza` models for the group with `GET /channels/available` models whose section includes that group. Ignore plaza/channel token, per-request, and official prices. Do not hard-code a frontend model catalog.
3. Video cards for a group:
   - Keep a `video_model_prices` key only when it has at least one configured 480p / 720p / 1080p override. Empty or all-null keys do not become cards.
   - Also create a card for each dynamic model that is video-billed, or whose platform is a media platform and is not image-billed.
   - Each card lists 480p / 720p / 1080p. A resolution uses the override if configured, otherwise the matching flat `video_price_*`. Omit a dynamic/override model that still has no configured tier.
   - If that union is empty and the group has at least one configured flat video price, render one fallback video card using only the flat prices.
4. Image cards for a group:
   - Only if the group has at least one configured image price.
   - If plaza or available channels return image-billed models for that group, one card per such model, all using the group image prices.
   - Otherwise render one fallback image card for the group.
   - Do not render an image card that has no configured image tier.
5. Fallback cards use the group name as the title, not a fake model ID.
6. Keep API group order. Inside a group, image cards then video cards, each sorted by title.
7. Load on mount. Refresh is manual. No 60-second polling.

## Out of scope

- Changing how billing is calculated at request time
- Model plaza UI
- Token / search / audio prices
- New public pricing API

## Tests

- `VideoSectionTabs` includes the pricing route and marks it active.
- New pricing view: collapse identical image tiers; keep video tiers split; omit empty tiers; show raw prices with no multiplier; one card per dynamic plaza/channel/override model; image models or group-name fallback cards; group-fetch failure is unavailable; empty when no configured media prices.
- `VideoApiDocsView` no longer renders `#pricing` or loads pricing data.

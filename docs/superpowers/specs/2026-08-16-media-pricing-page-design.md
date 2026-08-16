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

Empty, loading, and fetch-failure states stay on this page. If `GET /groups/available` fails, show unavailable even when plaza succeeds. Plaza failure only drops image model names; it does not make the page unavailable. If the user has no media group with any configured media price, show empty.

## Cards

Each card shows model name, type (image/video), and unit (per image / per second).

### Image

Read `image_price_1k`, `image_price_2k`, `image_price_4k`.

- Configured tiers that share one value collapse to a single price. Do not list 1K/2K/4K.
- Different values stay split by 1K / 2K / 4K.
- Omit unconfigured tiers. Do not render a dash.

### Video

Unit is USD/s. Always list configured resolutions; do not collapse even when values match.

- Prefer `video_model_prices[model][resolution]`.
- Fall back to `video_price_480p` / `video_price_720p` / `video_price_1080p`.
- Omit unconfigured resolutions.

### Price math

Display the raw configured USD amounts. Do not multiply by user exclusive rate, group `rate_multiplier`, image/video independent rates, or peak rate. Do not fetch `/groups/rates` for this page. Do not show multiplier text.

A tier is configured when the value is a finite number `>= 0`. `null`, missing, empty string, and non-numeric values are unconfigured. `0` is a real price.

## Data

No new backend endpoint.

1. `GET /groups/available` is the only price source. Keep a group if its platform is `leo`, `openai_media`, `video`, `composite`, or `grok`, and it has at least one configured image or video price field.
2. Optional `GET /model-plaza` supplies image model names only. Ignore plaza token, per-request, and official prices. Do not use plaza to decide video cards.
3. Video cards for a group:
   - One card per `video_model_prices` key that has at least one configured resolution.
   - That card’s price for a resolution is the model override if configured, otherwise the group flat `video_price_*` for 480p / 720p / 1080p.
   - If the group has no such keys but has at least one configured flat video price, render one fallback video card for the group (same idea as the image fallback).
   - Do not expand the workbench catalog onto a group. Workbench IDs are not a second source of cards.
4. Image cards for a group:
   - If plaza returns models for that group whose billing mode is image, one card per such model, all using the group image prices.
   - Otherwise, if the group has at least one configured image price, render one fallback image card for the group.
5. Fallback cards use the group name as the title, not a fake model ID.
6. Load on mount. Refresh is manual. No 60-second polling.

## Out of scope

- Changing how billing is calculated at request time
- Model plaza UI
- Token / search / audio prices
- New public pricing API

## Tests

- `VideoSectionTabs` includes the pricing route and marks it active.
- New pricing view: collapse identical image tiers; keep video tiers split; omit empty tiers; show raw prices with no multiplier; one card per `video_model_prices` key; plaza image models or group-name fallback cards; group-fetch failure is unavailable; empty when no configured media prices.
- `VideoApiDocsView` no longer renders `#pricing` or loads pricing data.

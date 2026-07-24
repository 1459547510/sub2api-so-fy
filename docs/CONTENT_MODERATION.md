# Content Moderation

Content moderation applies to authenticated non-administrator API requests in the groups and modes selected in **Security Audit > Content Moderation**.

Requests owned by an authenticated user whose trusted role is `admin` bypass both the legacy content-moderation engine and the prompt-audit engine. Those requests are not blocked by local keywords or moderation results and do not create moderation violations, send hit notifications, or contribute to automatic account disabling.

The bypass requires the authenticated user ID, API key owner ID, and loaded user record ID to match. A username, email address, API key name, or user-supplied request field cannot enable the bypass. Requests for users with the `user` role continue through the configured moderation path unchanged.

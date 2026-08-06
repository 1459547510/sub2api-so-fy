# Content Moderation

Content moderation applies to authenticated non-administrator API requests in the groups and modes selected in **Security Audit > Content Moderation**.

Requests owned by an authenticated user whose trusted role is `admin` bypass both the legacy content-moderation engine and the prompt-audit engine. Those requests are not blocked by local keywords or moderation results and do not create moderation violations, send hit notifications, or contribute to automatic account disabling.

The bypass requires the authenticated user ID, API key owner ID, and loaded user record ID to match. A username, email address, API key name, or user-supplied request field cannot enable the bypass. Requests for users with the `user` role continue through the configured moderation path unchanged.

For upstream `cyber_policy` records, administrators can optionally select an active OpenAI exclusive group to revoke from regular users after a hit. The backend performs the role check and group removal server-side, invalidates the user's auth cache, and leaves administrators, other groups, API keys, user status, and account-pool records unchanged. Leave the selector disabled to keep the rule off.

# Account Error Email Alerts

## Overview

The account error alert periodically scans account availability and sends an
email only when a new account error is detected. It reuses the existing Ops
email recipients, SMTP configuration, scheduler, and multi-instance leader
lock.

The feature is disabled by default. Its default schedule is every five minutes:

```text
*/5 * * * *
```

Cron expressions use the application timezone and have one-minute precision.

## Configuration

1. Configure and test SMTP in the admin settings.
2. Enable Ops monitoring.
3. Open the Ops email notification configuration.
4. Enable report emails and configure at least one report recipient.
5. Enable **New account error alert** and set its Cron schedule.

If no report recipient is configured, the scheduler falls back to the first
administrator email when available.

## Detection Rules

- Only accounts whose status is `error` are included.
- The first scheduled scan establishes a baseline and does not email existing
  historical errors.
- A normal account entering `error` status is sent once.
- An account that remains in `error` with the same reason is not sent again.
- A changed error reason is treated as a new error.
- An account that recovers and later enters `error` again is sent again.

The email includes the account ID, name, platform, group, and error reason.
The scan baseline is stored in the settings table, so service restarts do not
resend unchanged errors. A new baseline is committed only after at least one
email is delivered successfully; failed deliveries are retried on a later
scheduled scan.

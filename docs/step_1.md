# Video 01 – Why Background Jobs Exist

## Branch

video-01-why-background-jobs

## Goal

Show why background jobs are needed before building any abstraction.

## What We Build

- HTTP `/register` endpoint
- Fake email sending
- Artificial delay using `time.Sleep`

## Demo Flow

1. User hits `/register`
2. API waits 3 seconds to “send email”
3. Response is slow

## Concepts Explained

- Blocking vs non-blocking code
- Why slow requests hurt UX
- Examples of background jobs:
  - Emails
  - Payments
  - Reports
  - Notifications

## What to Say

- “We are not building a queue yet”
- “We must feel the pain first”
- “This is how real problems start”

## No Architecture Yet

- No goroutines
- No channels
- No queue

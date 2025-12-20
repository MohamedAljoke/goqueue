# Video 04 – Error Handling & Job Status

## Branch

video-04-error-handling

## Goal

Teach how failures should be handled correctly.

## What We Build

- Job statuses:
  - pending
  - processing
  - completed
  - failed
- Error propagation

## Concepts

- `error` in Go
- State transitions
- Testing failure paths

## Teaching Moment

Always test success AND failure.

## Visual State Flow

pending → processing → completed
→ failed

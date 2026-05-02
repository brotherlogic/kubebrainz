# Proposal: Artist Resolution Failure Visibility

## Problem
Currently, artist lookup failures in `kubebrainz` are only visible via gRPC error codes (`codes.NotFound`). This provides no long-term observability or actionable data for developers to investigate resolution failures.

## Solution
Implement automated GitHub issue creation for failed artist resolutions using `githubridge`.

## Design
- **Integration**: Utilize `githubridge` to create and manage GitHub issues.
- **Deduplication**: Implement server-side tracking to ensure only one active GitHub issue exists per unique failed artist.
- **Reporting**: When `GetArtist` returns a `NotFound` error, check if an issue for that artist is already open. If not, post a new one.
- **Content**: Issues will include:
    - Failed artist name.
    - Timestamp of failure.
    - Frequency counter of failures for this artist.

## Implementation Details
1. Add tracking mechanism to `Server` struct.
2. Update `GetArtist` in `api.go` to handle failure logging and `githubridge` integration.
3. Include artist name, timestamp, and count in issue body.

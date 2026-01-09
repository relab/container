// Package container provides a small client for interacting with a Docker daemon
// over its HTTP API. It exposes helpers for common operations such as building,
// pulling and removing images, and managing networks.
//
// The package is intentionally kept minimal, with no dependencies, mirroring the
// Docker API where useful, but may deviate from the existing API structure and
// naming based on my own preferences.
package container

package container

import (
	"net/url"
	"strconv"
)

// The struct definitions in this file are largely copied from the Docker API
// types, but may have been simplified for our use case.
// See the original [RemoveOptions], [LogsOptions], and [StopOptions].
//
// [RemoveOptions]: https://github.com/moby/moby/blob/master/api/types/container/options.go#L34
// [LogsOptions]: https://github.com/moby/moby/blob/master/api/types/container/options.go#L58
// [StopOptions]: https://github.com/moby/moby/blob/master/api/types/container/config.go#L18

// RemoveOptions holds parameters to remove containers.
type RemoveOptions struct {
	RemoveVolumes bool
	RemoveLinks   bool
	Force         bool
}

func (o RemoveOptions) url(containerID string) string {
	query := url.Values{}
	if o.RemoveVolumes {
		query.Set("v", "1")
	}
	if o.RemoveLinks {
		query.Set("link", "1")
	}
	if o.Force {
		query.Set("force", "1")
	}
	u := url.URL{Scheme: "http", Host: "localhost", Path: "/containers/" + containerID, RawQuery: query.Encode()}
	return u.String()
}

// LogsOptions holds parameters to filter logs with.
type LogsOptions struct {
	ShowStdout bool
	ShowStderr bool
	Timestamps bool
	Follow     bool
	Tail       string
	Details    bool
}

func (o LogsOptions) url(containerID string) string {
	query := url.Values{}
	if o.ShowStdout {
		query.Set("stdout", "1")
	}
	if o.ShowStderr {
		query.Set("stderr", "1")
	}
	if o.Timestamps {
		query.Set("timestamps", "1")
	}
	if o.Follow {
		query.Set("follow", "1")
	}
	if o.Details {
		query.Set("details", "1")
	}
	query.Set("tail", o.Tail)

	u := url.URL{Scheme: "http", Host: "localhost", Path: "/containers/" + containerID + "/logs", RawQuery: query.Encode()}
	return u.String()
}

// StopOptions holds the options to stop or restart a container.
type StopOptions struct {
	// Signal (optional) is the signal to send to the container to (gracefully)
	// stop it before forcibly terminating the container with SIGKILL after the
	// timeout expires. If not value is set, the default (SIGTERM) is used.
	Signal string `json:",omitempty"`

	// Timeout (optional) is the timeout (in seconds) to wait for the container
	// to stop gracefully before forcibly terminating it with SIGKILL.
	//
	// - Use nil to use the default timeout (10 seconds).
	// - Use '-1' to wait indefinitely.
	// - Use '0' to not wait for the container to exit gracefully, and
	//   immediately proceeds to forcibly terminating the container.
	// - Other positive values are used as timeout (in seconds).
	Timeout *int `json:",omitempty"`
}

func (o StopOptions) url() string {
	query := url.Values{}
	if o.Timeout != nil {
		query.Set("t", strconv.Itoa(*o.Timeout))
	}
	if o.Signal != "" {
		query.Set("signal", o.Signal)
	}
	u := url.URL{Scheme: "http", Host: "localhost", Path: "/stop", RawQuery: query.Encode()}
	return u.String()
}

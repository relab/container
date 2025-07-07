package container

import "net/url"

// The struct definitions in this file are largely copied from the Docker API
// types, but may have been simplified for our use case.
// See the originals [WaitCondition], [WaitResponse], and [WaitExitError].
//
// [WaitCondition]: https://github.com/moby/moby/blob/master/api/types/container/waitcondition.go
// [WaitResponse]: https://github.com/moby/moby/blob/master/api/types/container/wait_response.go
// [WaitExitError]: https://github.com/moby/moby/blob/master/api/types/container/wait_exit_error.go

// WaitCondition is a type used to specify a container state for which
// to wait.
type WaitCondition string

// Possible WaitCondition Values.
//
// WaitConditionNotRunning (default) is used to wait for any of the non-running
// states: "created", "exited", "dead", "removing", or "removed".
//
// WaitConditionNextExit is used to wait for the next time the state changes
// to a non-running state. If the state is currently "created" or "exited",
// this would cause Wait() to block until either the container runs and exits
// or is removed.
//
// WaitConditionRemoved is used to wait for the container to be removed.
const (
	WaitConditionNotRunning WaitCondition = "not-running"
	WaitConditionNextExit   WaitCondition = "next-exit"
	WaitConditionRemoved    WaitCondition = "removed"
)

func (w WaitCondition) url(containerID string) string {
	query := url.Values{}
	if w != "" {
		query.Set("condition", string(w))
	}
	u := url.URL{Scheme: "http", Host: "localhost", Path: "/containers/" + containerID + "/wait", RawQuery: query.Encode()}
	return u.String()
}

// WaitExitError container waiting error, if any
// swagger:model WaitExitError
type WaitExitError struct {
	// Details of an error
	Message string `json:"Message,omitempty"`
}

// WaitResponse ContainerWaitResponse
//
// OK response to ContainerWait operation
// swagger:model WaitResponse
type WaitResponse struct {
	// error
	Error *WaitExitError `json:"Error,omitempty"`

	// Exit code of the container
	// Required: true
	StatusCode int64 `json:"StatusCode"`
}

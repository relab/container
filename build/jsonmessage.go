package build

// These message types are originally defined in the Docker API under the jsonmessage package.
// We define them here to avoid adding an extra package.

// JSONError wraps a concrete Code and Message, Code is
// an integer error code, Message is the error message.
type JSONError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *JSONError) Error() string {
	return e.Message
}

// JSONMessage defines a message struct for docker events.
//
// This is a simplified version of the Docker API's [JSONMessage].
// See the full version in case we want to extend this in the future.
//
// [JSONMessage]: https://github.com/moby/moby/blob/v28.5.1/pkg/jsonmessage/jsonmessage.go#L144
type JSONMessage struct {
	Stream string     `json:"stream,omitempty"`
	Status string     `json:"status,omitempty"`
	ID     string     `json:"id,omitempty"`
	Error  *JSONError `json:"errorDetail,omitempty"`
}

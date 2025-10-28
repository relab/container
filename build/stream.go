package build

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ConsumeStream reads the Docker image build JSON stream from r, writes any
// human-readable lines to w, and returns an error if the stream reports a
// build error. It drains the entire stream before returning to ensure the
// build completes.
func ConsumeStream(r io.Reader, w io.Writer) error {
	dec := json.NewDecoder(r)
	var buildErr error
	for {
		var jsonLine JSONMessage
		if err := dec.Decode(&jsonLine); err != nil {
			if err == io.EOF {
				break
			}
			// If decoding fails mid-stream, do a best-effort drain of any buffered
			// content to w and stop parsing. This mirrors previous raw streaming
			// behavior without claiming strict JSON compatibility.
			if buf := dec.Buffered(); buf != nil {
				_, _ = io.Copy(w, buf)
			}
			return buildErr
		}

		if jsonLine.Stream != "" {
			_, _ = io.WriteString(w, jsonLine.Stream)
		}
		if jsonLine.Error != nil {
			buildErr = fmt.Errorf("docker build error: %s", strings.TrimSpace(jsonLine.Error.Error()))
			// Do not return immediately; drain the rest so the build completes.
		}
	}
	return buildErr
}

package container

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/relab/container/build"
	"github.com/relab/container/image"
	"github.com/relab/container/network"
)

const (
	defaultTimeout = 10 * time.Second
)

type Container struct {
	client *http.Client
}

// NewContainer creates a new Container instance that can be used to interact with the Docker daemon.
// Currently we support only the default Docker host URL.
func NewContainer() (*Container, error) {
	u, err := url.Parse(DefaultDockerHost)
	if err != nil {
		return nil, err
	}

	return &Container{
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:    10,
				IdleConnTimeout: 30 * time.Second,
				DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
					return net.DialTimeout(u.Scheme, u.Path, defaultTimeout)
				},
			},
		},
	}, nil
}

// Ping checks if the Docker daemon is reachable and responds to a ping request.
func (c *Container) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/_ping", nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer close(resp)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return nil
}

// ImagePull requests the docker host to pull an image from a remote registry.
// It executes the privileged function if the operation is unauthorized
// and it tries one more time.
// It's up to the caller to handle the io.ReadCloser and close it properly.
func (c *Container) ImagePull(ctx context.Context, refStr string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost/images/create?fromImage="+refStr, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// ImageBuild builds a Docker image from the provided build context.
// The build context should be a tar archive containing the Dockerfile and
// any other files needed for the build.
func (c *Container) ImageBuild(ctx context.Context, buildCtx io.Reader, options build.ImageBuildOptions) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, options.URL(), buildCtx)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-tar")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("image build failed: %s", resp.Status)
	}
	return resp.Body, nil
}

// ImageRemove removes an image from the docker host.
func (c *Container) ImageRemove(ctx context.Context, imageID string, options image.RemoveOptions) ([]image.DeleteResponse, error) {
	imageID = strings.TrimSpace(imageID)
	if imageID == "" {
		return nil, fmt.Errorf("image ID cannot be empty")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, options.URL(imageID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer close(resp)

	if resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("image removal failed: %s", resp.Status)
	}
	var response []image.DeleteResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

// NetworkCreate creates a new network in the docker host.
func (c *Container) NetworkCreate(ctx context.Context, options network.CreateOptions) (network.CreateResponse, error) {
	body, err := encodeBody(options)
	if err != nil {
		return network.CreateResponse{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost/networks/create", body)
	if err != nil {
		return network.CreateResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return network.CreateResponse{}, err
	}
	defer close(resp)

	if resp.StatusCode != http.StatusCreated {
		return network.CreateResponse{}, fmt.Errorf("network creation failed: %s", resp.Status)
	}
	var response network.CreateResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

// NetworkRemove removes an existing network from the docker host.
func (c *Container) NetworkRemove(ctx context.Context, networkID string) error {
	networkID = strings.TrimSpace(networkID)
	if networkID == "" {
		return fmt.Errorf("network ID cannot be empty")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, "http://localhost/networks/"+networkID, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer close(resp)

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("network removal failed: %s", resp.Status)
	}
	return nil
}

// NetworkConnect connects a container to an existent network in the docker host.
func (c *Container) NetworkConnect(ctx context.Context, networkID, containerID string, config *network.EndpointSettings) error {
	networkID = strings.TrimSpace(networkID)
	if networkID == "" {
		return fmt.Errorf("network ID cannot be empty")
	}
	containerID = strings.TrimSpace(containerID)
	if containerID == "" {
		return fmt.Errorf("container ID cannot be empty")
	}

	body, err := encodeBody(network.ConnectOptions{
		Container:      containerID,
		EndpointConfig: config,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost/networks/"+networkID+"/connect", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer close(resp)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("network connect failed: %s", resp.Status)
	}
	return nil
}

// NetworkDisconnect disconnects a container from an existent network in the docker host.
func (c *Container) NetworkDisconnect(ctx context.Context, networkID, containerID string, force bool) error {
	networkID = strings.TrimSpace(networkID)
	if networkID == "" {
		return fmt.Errorf("network ID cannot be empty")
	}
	containerID = strings.TrimSpace(containerID)
	if containerID == "" {
		return fmt.Errorf("container ID cannot be empty")
	}

	body, err := encodeBody(network.DisconnectOptions{
		Container: containerID,
		Force:     force,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost/networks/"+networkID+"/disconnect", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer close(resp)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("network disconnect failed: %s", resp.Status)
	}
	return nil
}

// ContainerCreate creates a new container based on the given configuration.
// It can be associated with a name, but it's not mandatory.
func (c *Container) ContainerCreate(ctx context.Context, config *Config, hostConfig *HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (CreateResponse, error) {
	body, err := encodeBody(CreateRequest{
		Config:           config,
		HostConfig:       hostConfig,
		NetworkingConfig: networkingConfig,
	})
	if err != nil {
		return CreateResponse{}, err
	}
	query := url.Values{}
	if containerName != "" {
		query.Set("name", containerName)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost/containers/create?"+query.Encode(), body)
	if err != nil {
		return CreateResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return CreateResponse{}, err
	}
	defer close(resp)

	if resp.StatusCode != http.StatusCreated {
		return CreateResponse{}, fmt.Errorf("container creation failed: %s", resp.Status)
	}

	var response CreateResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

// ContainerRemove kills and removes a container from the docker host.
func (c *Container) ContainerRemove(ctx context.Context, containerID string, options RemoveOptions) error {
	containerID = strings.TrimSpace(containerID)
	if containerID == "" {
		return fmt.Errorf("container ID cannot be empty")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, options.url(containerID), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer close(resp)

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("container removal failed: %s", resp.Status)
	}
	return nil
}

// ContainerStart sends a request to the docker daemon to start a container.
func (c *Container) ContainerStart(ctx context.Context, containerID string) error {
	containerID = strings.TrimSpace(containerID)
	if containerID == "" {
		return fmt.Errorf("container ID cannot be empty")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost/containers/"+containerID+"/start", nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer close(resp)

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("container start failed: %s", resp.Status)
	}
	return nil
}

// ContainerStop stops a container. In case the container fails to stop
// gracefully within a time frame specified by the timeout argument,
// it is forcefully terminated (killed).
//
// If the timeout is nil, the container's StopTimeout value is used, if set,
// otherwise the engine default. A negative timeout value can be specified,
// meaning no timeout, i.e. no forceful termination is performed.
func (c *Container) ContainerStop(ctx context.Context, containerID string, options StopOptions) error {
	containerID = strings.TrimSpace(containerID)
	if containerID == "" {
		return fmt.Errorf("container ID cannot be empty")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, options.url(), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer close(resp)

	return nil
}

const containerWaitErrorMsgLimit = 2 * 1024 // 2KiB

// ContainerWait waits until the specified container is in a certain state
// indicated by the given condition, either "not-running" (default),
// "next-exit", or "removed".
//
// This is a blocking call that will return when the container is in the specified state.
//
// This is different from the original Docker API, which returns a channel
// that can be used to wait for the container state change.
// See the original [ContainerWait] API documentation for more details.
//
// [ContainerWait]: https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerWait
func (c *Container) ContainerWait(ctx context.Context, containerID string, condition WaitCondition) (WaitResponse, error) {
	containerID = strings.TrimSpace(containerID)
	if containerID == "" {
		return WaitResponse{}, fmt.Errorf("container ID cannot be empty")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, condition.url(containerID), nil)
	if err != nil {
		return WaitResponse{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return WaitResponse{}, err
	}
	defer close(resp)

	var buf bytes.Buffer
	stream := io.TeeReader(resp.Body, &buf)

	var result WaitResponse
	if err := json.NewDecoder(stream).Decode(&result); err != nil {
		// Try to extract plaintext error message (e.g. from proxy)
		if errors.As(err, new(*json.SyntaxError)) {
			body, _ := io.ReadAll(io.LimitReader(stream, containerWaitErrorMsgLimit))
			return WaitResponse{}, fmt.Errorf("malformed response: %s%s", buf.String(), string(body))
		}
		return WaitResponse{}, err
	}

	return result, nil
}

// ContainerInspect returns the container information.
func (c *Container) ContainerInspect(ctx context.Context, containerID string) (InspectResponse, error) {
	containerID = strings.TrimSpace(containerID)
	if containerID == "" {
		return InspectResponse{}, fmt.Errorf("container ID cannot be empty")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/containers/"+containerID+"/json", nil)
	if err != nil {
		return InspectResponse{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return InspectResponse{}, err
	}
	defer close(resp)

	if resp.StatusCode != http.StatusOK {
		return InspectResponse{}, fmt.Errorf("container inspect failed: %s", resp.Status)
	}

	var response InspectResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

// ContainerLogs returns the logs generated by a container in an io.ReadCloser.
// It's up to the caller to close the stream.
//
// The stream format on the response will be in one of two formats:
//
// If the container is using a TTY, there is only a single stream (stdout), and
// data is copied directly from the container output stream, no extra
// multiplexing or headers.
//
// If the container is *not* using a TTY, streams for stdout and stderr are
// multiplexed.
// The format of the multiplexed stream is as follows:
//
//	[8]byte{STREAM_TYPE, 0, 0, 0, SIZE1, SIZE2, SIZE3, SIZE4}[]byte{OUTPUT}
//
// STREAM_TYPE can be 1 for stdout and 2 for stderr
//
// SIZE1, SIZE2, SIZE3, and SIZE4 are four bytes of uint32 encoded as big endian.
// This is the size of OUTPUT.
//
// You can use github.com/docker/docker/pkg/stdcopy.StdCopy to demultiplex this
// stream.
func (c *Container) ContainerLogs(ctx context.Context, containerID string, options LogsOptions) (io.ReadCloser, error) {
	containerID = strings.TrimSpace(containerID)
	if containerID == "" {
		return nil, fmt.Errorf("container ID cannot be empty")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, options.url(containerID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func encodeBody(obj any) (*bytes.Buffer, error) {
	if obj == nil {
		return nil, nil
	}
	body := bytes.NewBuffer(nil)
	if err := json.NewEncoder(body).Encode(obj); err != nil {
		return nil, err
	}
	return body, nil
}

func close(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
}

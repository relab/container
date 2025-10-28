package container_test

import (
	"archive/tar"
	"bytes"
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/relab/container"
	"github.com/relab/container/build"
	"github.com/relab/container/network"
)

const containerTestTag = "container-test"

// TestMain builds the test image once before running tests that depend on it.
func TestMain(m *testing.M) {
	// Parse flags early so testing.Verbose() works in buildTestImage.
	flag.Parse()
	if err := buildTestImage(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to build test image: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestPing(t *testing.T) {
	c, err := container.NewContainer()
	if err != nil {
		t.Fatalf("Failed to create container client: %v", err)
	}

	if err := c.Ping(t.Context()); err != nil {
		t.Fatalf("Failed to ping Docker daemon: %v", err)
	}
	t.Log("Ping successful")
}

func TestNetworkCreateAndNetworkRemove(t *testing.T) {
	c, err := container.NewContainer()
	if err != nil {
		t.Fatalf("Failed to create container client: %v", err)
	}

	resp, err := c.NetworkCreate(t.Context(), network.CreateOptions{
		Name:   "container-" + rand.Text()[:8],
		Driver: "bridge",
	})
	if err != nil {
		t.Fatalf("Failed to create network: %v", err)
	}
	t.Cleanup(func() {
		// cannot use t.Context() here, since it may be canceled before cleanup runs
		if err := c.NetworkRemove(context.Background(), resp.ID); err != nil {
			t.Errorf("Failed to remove network: %v", err)
		}
		t.Logf("Network removed: %s", resp.ID)
	})

	t.Logf("Network created: %s", resp.ID)
}

func TestContainerCreateAndStartAndInspectAndStop(t *testing.T) {
	c, err := container.NewContainer()
	if err != nil {
		t.Fatalf("Failed to create container client: %v", err)
	}

	net, err := c.NetworkCreate(t.Context(), network.CreateOptions{
		Name:   "container-" + rand.Text()[:8],
		Driver: "bridge",
	})
	if err != nil {
		t.Fatalf("Failed to create network: %v", err)
	}
	t.Logf("Network created: %s", net.ID)

	resp, err := c.ContainerCreate(context.Background(), &container.Config{
		Env:   []string{"AUTHORIZED_KEYS=xyz"},
		Image: containerTestTag,
		ExposedPorts: container.PortSet{
			"22/tcp": struct{}{},
		},
	}, &container.HostConfig{
		PortBindings: container.PortMap{"22/tcp": {{}}}, // map ssh port to ephemeral port
		AutoRemove:   true,
	}, nil, "")
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}

	t.Cleanup(func() {
		timeout := 1 // seconds to wait before forcefully killing the container
		opts := container.StopOptions{Timeout: &timeout}

		// cannot use t.Context() here, since it may be canceled before cleanup runs
		ctx := context.Background()

		if err := c.ContainerStop(ctx, resp.ID, opts); err != nil {
			t.Errorf("Failed to stop container '%s': %v", resp.ID, err)
		} else {
			t.Logf("Container stopped: %s", resp.ID)
		}

		if err := c.NetworkDisconnect(ctx, net.ID, resp.ID, true); err != nil {
			t.Errorf("Failed to disconnect container from network '%s': %v", net.ID, err)
		} else {
			t.Logf("Container disconnected from network: %s", net.ID)
		}

		if err := c.NetworkRemove(ctx, net.ID); err != nil {
			t.Errorf("Failed to remove network: %v", err)
		} else {
			t.Logf("Network removed: %s", net.ID)
		}
	})
	t.Logf("Container created: %s", resp.ID)

	err = c.ContainerStart(t.Context(), resp.ID)
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}
	t.Logf("Container started: %s", resp.ID)

	insp, err := c.ContainerInspect(t.Context(), resp.ID)
	if err != nil {
		t.Fatalf("Failed to inspect container: %v", err)
	}
	t.Logf("Container inspected: %+v", insp)

	name := strings.TrimPrefix(insp.Name, "/")
	err = c.NetworkConnect(context.Background(), net.ID, resp.ID, &network.EndpointSettings{
		Aliases: []string{name},
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Container name: %s", name)
	t.Logf("Container network settings: %+v", insp.NetworkSettings.Ports)
}

func TestBuild(t *testing.T) {
	resp, err := startImageBuild(t.Context(), []string{containerTestTag})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := resp.Close(); err != nil {
			t.Error(err)
		}
	})
	if err := build.ConsumeStream(resp, t.Output()); err != nil {
		t.Error(err)
	}
}

// startImageBuild initializes a Docker image build and returns the build output stream.
// The caller owns the returned ReadCloser and must Close it after consuming the stream.
func startImageBuild(ctx context.Context, tags []string) (io.ReadCloser, error) {
	c, err := container.NewContainer()
	if err != nil {
		return nil, fmt.Errorf("create container failed: %w", err)
	}
	buildCtx, err := prepareBuildContext()
	if err != nil {
		return nil, fmt.Errorf("prepare build context failed: %w", err)
	}
	return c.ImageBuild(ctx, buildCtx, build.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       tags,
	})
}

// buildTestImage builds the container image needed for tests.
func buildTestImage() error {
	resp, err := startImageBuild(context.Background(), []string{containerTestTag})
	if err != nil {
		return err
	}
	defer func() { _ = resp.Close() }()

	// Capture logs. If tests are verbose (-v), also echo to stderr live.
	var buf bytes.Buffer
	out := io.Writer(&buf)
	if testing.Verbose() {
		out = io.MultiWriter(&buf, os.Stderr)
	}

	if err := build.ConsumeStream(resp, out); err != nil {
		// Ensure logs are visible on failure even if logging not enabled.
		if !testing.Verbose() {
			_, _ = io.Copy(os.Stderr, &buf)
		}
		return err
	}
	return nil
}

var (
	dockerfile = `FROM alpine:latest

RUN apk add --no-cache openssh lsb-release && \
    ssh-keygen -A && \
    mkdir -p /root/.ssh && \
    chmod 700 /root/.ssh

ADD entrypoint.sh /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh", "sleep", "infinity" ]
`

	entrypoint = `#!/bin/sh

mkdir "$HOME/.ssh"
echo "$AUTHORIZED_KEYS" > "$HOME/.ssh/authorized_keys"
/usr/sbin/sshd
exec "$@"
`
)

func prepareBuildContext() (r io.ReadCloser, err error) {
	var buf bytes.Buffer
	tarWriter := tar.NewWriter(&buf)

	err = tarWriter.WriteHeader(&tar.Header{
		Name:   "Dockerfile",
		Size:   int64(len(dockerfile)),
		Mode:   0o644,
		Format: tar.FormatUSTAR,
	})
	if err != nil {
		return nil, err
	}

	_, err = tarWriter.Write([]byte(dockerfile))
	if err != nil {
		return nil, err
	}

	err = tarWriter.WriteHeader(&tar.Header{
		Name:   "entrypoint.sh",
		Size:   int64(len(entrypoint)),
		Mode:   0o755,
		Format: tar.FormatUSTAR,
	})
	if err != nil {
		return nil, err
	}

	_, err = tarWriter.Write([]byte(entrypoint))
	if err != nil {
		return nil, err
	}
	return io.NopCloser(&buf), nil
}

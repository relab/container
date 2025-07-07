package container

// The struct definitions in this file is largely copied from the Docker API
// types, but may have been simplified for our use case.
// See the original [InspectResponse].
//
// [InspectResponse]: https://github.com/moby/moby/blob/master/api/types/container/container.go

// InspectResponse is the response for the GET "/containers/{name:.*}/json"
// endpoint.
type InspectResponse struct {
	ID              string `json:"Id"`
	Created         string
	Path            string
	Args            []string
	Image           string
	ResolvConfPath  string
	HostnamePath    string
	HostsPath       string
	LogPath         string
	Name            string
	RestartCount    int
	Driver          string
	Platform        string
	MountLabel      string
	ProcessLabel    string
	AppArmorProfile string
	ExecIDs         []string
	HostConfig      *HostConfig
	Config          *Config
	NetworkSettings *NetworkSettings
}

// NetworkSettings exposes the network settings in the api
type NetworkSettings struct {
	Ports PortMap // Ports is a collection of PortBinding indexed by Port
}

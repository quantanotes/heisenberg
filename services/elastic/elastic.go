package elastic

import (
	"services/utils"

	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// Services to handle elastic compute containers

var sw *scw.Client     // Scaleway client instance
var api *container.API // Scaleway container api

// Create one heisenberg container per user
func createContainer(usr string) error {
	// No idea why the fuck Go makes me do this????
	maxScale := new(uint32)
	*maxScale = uint32(1)
	memoryLimit := new(uint32)
	*memoryLimit = uint32(128)
	port := new(uint32)
	*port = 420

	// New container configuration
	req := &container.CreateContainerRequest{
		NamespaceID: "heisenberg",
		Name:        usr, // Set id of container to username 1-1 user-container (will change)
		MaxScale:    maxScale,
		MemoryLimit: memoryLimit, // Container mem allocation in Mb
		Port:        port,
		//RegistryImage: ,
	}

	_, err := api.CreateContainer(req)
	return err
}

// Wipe that bitch
func deleteContainer(usr string) error {
	req := &container.DeleteContainerRequest{
		ContainerID: usr,
	}
	_, err := api.DeleteContainer(req)
	return err
}

func getContainersForUser() error {
	return nil
}

func getAllContainers() error {
	return nil
}

func init() {
	var err error
	sw, err = scw.NewClient(
		scw.WithDefaultOrganizationID(utils.Env["SCW_DEFAULT_ORGANIZATION_ID"]),
		scw.WithAuth(utils.Env["SCW_ACCESS_KEY"], utils.Env["SCW_SECRET_KEY"]),
		scw.WithDefaultRegion(scw.Region(scw.ZoneFrPar1)),
	)
	if err != nil {
		panic(err)
	}

	api = container.NewAPI(sw)
}

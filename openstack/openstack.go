package openstack

import (
	"errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

func CreateServerSnapShot(params map[string]string) (string, error) {
	provider, err := auth(params["username"], params["password"], params["project"], params["endpoint"])
	if err != nil {
		return "error", err
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	serverId := params["server"]

	opts := servers.CreateImageOpts{Name: params["imageName"]}
	imageId, err := servers.CreateImage(client, serverId, opts).ExtractImageID()

	if err != nil {
		return "error", err
	}

	return imageId, nil
}

func DeleteServerSnapshot(params map[string]string) gophercloud.ErrResult {
	provider, err := auth(params["username"], params["password"], params["project"], params["endpoint"])
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err != nil {
		return err
	}

	res := images.Delete(client, params["imageId"])
	return res
}

func CreateVolumeSnapshot() {

}

func DeleteVolumeSnapshot() {

}

func auth(_username string, _password string, _project string, _endPoint string) (*gophercloud.ProviderClient, error) {
	sc := &gophercloud.ProviderClient{}

	if _endPoint == "" {
		return sc, errors.New("End point is missing.")
	}

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "" + _endPoint + ":5000/v2.0",
		Username:         "" + _username + "",
		Password:         "" + _password + "",
		TenantID:         "" + _project + "",
	}

	provider, err := openstack.AuthenticatedClient(opts)

	if err != nil {
		return provider, err
	}
	return provider, err
}

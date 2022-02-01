package openstack

import (
	"encoding/json"
	"errors"

	validator "github.com/go-playground/validator/v10"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v2/snapshots"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/pagination"
)

type imageDetails struct {
	ID       string `json:"ID"`
	Created  string `json:"Created"`
	MinDisk  int    `json:"MinDisk"`
	MinRAM   int    `json:"MinRAM"`
	Name     string `json:"Name"`
	Progress int    `json:"Progress"`
	Status   string `json:"Status"`
	Updated  string `json:"Updated"`
	Metadata struct {
		BdmV2              string `json:"bdmV2"`
		RootDeviceName     string `json:"root_device_name"`
		BlockDeviceMapping []struct {
			DeviceName          string `json:"device_name"`
			SourceType          string `json:"source_type"`
			VolumeSize          int    `json:"volume_size"`
			BootIndex           int    `json:"boot_index"`
			DeleteOnTermination bool   `json:"delete_on_termination"`
			DestinationType     string `json:"destination_type"`
			SnapshotID          string `json:"snapshot_id"`
			DeviceType          string `json:"device_type"`
			DiskBus             string `json:"disk_bus"`
		} `json:"block_device_mapping"`
	} `json:"Metadata"`
}

type OpenstackAuth struct {
	Username string `json:"username" validate:"required"` //
	Password string `json:"password" validate:"required"` //
	Project  string `json:"project" validate:"required"`  //
	EndPoint string `json:"endPoint" validate:"required"` //
}

/*
Creating Server Image
*/
func CreateServerSnapShot(authCredential OpenstackAuth, serverId string, snapShotName string) (string, error) {
	validate := validator.New()
	errVal := validate.Var(serverId, "required")
	if errVal != nil {
		return "error", errors.New("Server Id missing.")
	}

	errVal = validate.Var(snapShotName, "required")

	if errVal != nil {
		return "error", errors.New("Snapshot name missing.")
	}

	provider, err := auth(authCredential)
	if err != nil {
		return "error", err
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	opts := servers.CreateImageOpts{Name: snapShotName}
	imageId, err := servers.CreateImage(client, serverId, opts).ExtractImageID()

	if err != nil {
		return "error", err
	}

	return imageId, nil
}

/*
Deleting Server Image
*/
func DeleteServerSnapshot(authCredential OpenstackAuth, snapShotId string) error {
	validate := validator.New()
	err := validate.Var(snapShotId, "required")

	if err != nil {
		return errors.New("Schedule Id missing.")
	}
	provider, err := auth(authCredential)
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err != nil {
		return err
	}

	//Get the List of connecting Snapshot of volumes
	_imageDetails, _ := ListServerSnapshotDetails(authCredential, snapShotId)
	jsonStr, err := json.Marshal(_imageDetails)
	if err != nil {
		return err
	}

	var m imageDetails
	if err := json.Unmarshal(jsonStr, &m); err != nil {
		return err
	}
	for _, x := range m.Metadata.BlockDeviceMapping {
		//Deleting the connecting Block Volume
		DeleteVolumeSnapshot(authCredential, x.SnapshotID)
	}

	res := images.Delete(client, snapShotId)
	if res.Err != nil {
		return res.Err
	}
	return nil
}

/*
List Volume Snapshots
*/
func ListServerSnapshotDetails(authCredential OpenstackAuth, snapShotId string) (*images.Image, error) {
	result := &images.Image{}

	validate := validator.New()
	err := validate.Var(snapShotId, "required")

	if err != nil {
		return result, errors.New("Snapshot Id missing.")
	}

	provider, err := auth(authCredential)
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err != nil {
		return result, err
	}

	imageDetails, err := images.Get(client, snapShotId).Extract()

	if err != nil {
		return result, err
	}

	return imageDetails, nil
}

/*
Create Volume Snapshot
*/
func CreateVolumeSnapshot(authCredential OpenstackAuth, volumeId string, snapShotName string) (string, error) {

	validate := validator.New()
	err := validate.Var(volumeId, "required")

	if err != nil {
		return "error", errors.New("Volume Id missing.")
	}

	err = validate.Var(snapShotName, "required")
	if err != nil {
		return "error", errors.New("Snapshot Name missing.")
	}

	provider, err := auth(authCredential)
	client, err := openstack.NewBlockStorageV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err != nil {
		return "error", err
	}

	options := snapshots.CreateOpts{VolumeID: volumeId, Name: snapShotName}
	volumeSnapShot, err := snapshots.Create(client, options).Extract()

	if err != nil {
		return "error", err
	}

	return volumeSnapShot.ID, nil
}

/*
Delete Volume Snapshot
*/
func DeleteVolumeSnapshot(authCredential OpenstackAuth, snapShotId string) error {
	validate := validator.New()
	err := validate.Var(snapShotId, "required")

	if err != nil {
		return errors.New("SnapShot Id missing.")
	}
	provider, err := auth(authCredential)
	client, err := openstack.NewBlockStorageV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err != nil {
		return err
	}

	res := snapshots.Delete(client, snapShotId)
	if res.Err != nil {
		return res.Err
	}
	return nil
}

/*
List Images
*/
func ListImages(authCredentials OpenstackAuth) ([]map[string]interface{}, error) {
	result := make(map[string]interface{})
	result1 := []map[string]interface{}{}

	provider, err := auth(authCredentials)
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err != nil {
		return result1, err
	}
	opt := &images.ListOpts{}
	pager := images.ListDetail(client, opt)
	pager.EachPage(func(page pagination.Page) (bool, error) {
		imageList, _ := images.ExtractImages(page)
		for _, images := range imageList {
			result["ID"] = images.ID
			result["Name"] = images.Name
			result1 = append(result1, result)
		}
		return true, nil
	})
	return result1, nil
}

/*
Auth Function
*/
func auth(authCredential OpenstackAuth) (*gophercloud.ProviderClient, error) {
	sc := &gophercloud.ProviderClient{}

	validate := validator.New()
	err := validate.Struct(&authCredential)

	if err != nil {
		return sc, errors.New("Validation failure, not all auth parameters provided.")
	}

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "" + authCredential.EndPoint + ":5000/v2.0",
		Username:         "" + authCredential.Username + "",
		Password:         "" + authCredential.Password + "",
		TenantID:         "" + authCredential.Project + "",
	}

	provider, err := openstack.AuthenticatedClient(opts)

	if err != nil {
		return provider, err
	}
	return provider, err
}

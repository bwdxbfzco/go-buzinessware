package openstack

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	conv "github.com/cstockton/go-conv"
	validator "github.com/go-playground/validator/v10"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v2/snapshots"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/bootfromvolume"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	imageservices "github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
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

/*
Delete Port
*/
func DeletePort(authCredentials OpenstackAuth, ipAddress string) error {
	provider, err := auth(authCredentials)

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err != nil {
		return err
	}

	port, err := ListPorts(authCredentials, ipAddress)

	if err != nil {
		return err
	}

	res := ports.Delete(client, port.ID)
	if res.ErrResult.Err != nil {
		return err
	}

	return nil
}

/*
List Ports
*/
func ListPorts(authCredentials OpenstackAuth, ipaddress string) (ports.Port, error) {
	var result ports.Port
	provider, err := auth(authCredentials)
	if err != nil {
		return result, err
	}
	client, err1 := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err1 != nil {
		return result, err1
	}

	var fixedIP ports.FixedIPOpts
	var fixedIPs []ports.FixedIPOpts

	fixedIP.IPAddress = ipaddress
	fixedIPs = append(fixedIPs, fixedIP)

	opts := ports.ListOpts{
		FixedIPs: fixedIPs,
	}

	errPort := ports.List(client, opts).EachPage(func(page pagination.Page) (bool, error) {
		actual, _ := ports.ExtractPorts(page)
		for _, b := range actual {
			result = b
		}
		return true, nil
	})
	if errPort != nil {
		return result, errPort
	}

	return result, nil
}

/*
Create Port
*/
func CreatePort(authCredentials OpenstackAuth, adminCredentials *gophercloud.ProviderClient, params map[string]string) (string, error) {
	provider := adminCredentials
	client, err1 := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err1 != nil {
		return "", err1
	}

	port, _ := ListPorts(authCredentials, params["ipAddress"])

	var ip ports.IP
	var ips []ports.IP

	ip.SubnetID = params["subnetId"]
	ip.IPAddress = params["ipAddress"]
	ips = append(ips, ip)

	if port.ID == "" {
		asu := true
		options := ports.CreateOpts{
			AdminStateUp: &asu,
			NetworkID:    params["networkId"],
			TenantID:     params["tenantId"],
			FixedIPs:     ips,
		}
		n, err := ports.Create(client, options).Extract()
		if err != nil {
			return "", err
		}
		return n.ID, nil
	} else {
		if port.TenantID == authCredentials.Project {
			return port.ID, nil
		} else {
			return "", errors.New("Port already in use")
		}
	}
	return "", nil
}

/*
Reinstall Server
*/
func Reinstall(authCredentials OpenstackAuth, adminCredentials *gophercloud.ProviderClient, serverParams map[string]string) (string, error) {
	var reinstall int
	port, _ := ListPorts(authCredentials, serverParams["ipaddress"])
	if port.ID != "" {
		//Terminate the server
		DeleteServer(authCredentials, serverParams["serverId"])
		//Sleep for sometime
		time.Sleep(15 * time.Second)
		//Verify if the port is in use
		portVerify, _ := ListPorts(authCredentials, serverParams["ipaddress"])
		if portVerify.ID == "" {
			reinstall = 1
		} else {
			//Delete Port forcefully
			DeletePort(authCredentials, serverParams["ipaddress"])
			time.Sleep(5 * time.Second)
			reinstall = 1
		}
	} else {
		reinstall = 1
	}
	if reinstall == 1 {
		//Create Port with same IP Address
		params := make(map[string]string)
		params["networkId"] = port.NetworkID
		params["subnetId"] = port.FixedIPs[0].SubnetID
		params["ipAddress"] = serverParams["ipaddress"]
		params["tenantId"] = port.TenantID
		if port.SecurityGroups[0] != "" {
			params["securitygroup"] = port.SecurityGroups[0]
		}
		newPort, portErr := CreatePort(authCredentials, adminCredentials, params)

		if portErr != nil {
			log.Println(portErr)
			return "", portErr
		}
		//Create server with new Port
		serverParams["port"] = newPort
		if port.SecurityGroups[0] != "" {
			serverParams["securitygroup"] = port.SecurityGroups[0]
		}
		result, serverErr := CreateServerFromImage(authCredentials, params["region"], serverParams)
		if serverErr != nil {
			return "", serverErr
		}
		time.Sleep(10 * time.Second)
		return result, nil
	}
	return "", nil
}

/*
Create Server From Image
*/
func CreateServerFromImage(authCredentials OpenstackAuth, region string, params map[string]string) (string, error) {
	provider, err := auth(authCredentials)

	if err != nil {
		return "error", err
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err != nil {
		return "error", err
	}

	var _sg1 []string
	var _block1 []bootfromvolume.BlockDevice
	_userdata := []byte(params["password"])
	_metadata := make(map[string]string)
	_net1 := []servers.Network{}
	_net := servers.Network{}

	if params["port"] == "" {
		_net = servers.Network{UUID: params["networkId"]}
	} else {
		_net = servers.Network{Port: params["port"]}
	}

	_net1 = append(_net1, _net)
	_sg := params["securitygroup"]
	_sg1 = append(_sg1, _sg)

	_metadata["admin_pass"] = params["password"]
	_imageRef := params["image"]
	_flavorRef := params["size"]
	_serverName := params["serverName"]
	_size, _ := conv.Int(params["disksize"])

	_block := bootfromvolume.BlockDevice{
		SourceType:          "image",
		UUID:                strings.TrimSpace(_imageRef),
		BootIndex:           0,
		DeleteOnTermination: true,
		DestinationType:     "volume",
		VolumeSize:          _size,
	}

	_block1 = append(_block1, _block)

	opts := bootfromvolume.CreateOptsExt{
		servers.CreateOpts{
			Name:           strings.TrimSpace(_serverName),
			ImageRef:       strings.TrimSpace(_imageRef),
			FlavorRef:      strings.TrimSpace(_flavorRef),
			Networks:       _net1,
			SecurityGroups: _sg1,
			UserData:       _userdata,
			Metadata:       _metadata,
		},
		_block1,
	}

	actual, err := bootfromvolume.Create(client, opts).Extract()

	if err != nil {
		return "error", err
	}

	return actual.ID, nil
}

/*
Delete Server
*/
func DeleteServer(authCredentials OpenstackAuth, id string) error {
	provider, _ := auth(authCredentials)
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err != nil {
		return err
	}
	res := servers.Delete(client, id)

	if res.ErrResult.Err != nil {
		return res.ErrResult.Err
	}

	return nil
}

/*
Server Info
*/

func ServerDetails(authCredentials OpenstackAuth, serverId string) (*servers.Server, error) {
	var _serverInfo *servers.Server
	provider, _ := auth(authCredentials)
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err != nil {
		return _serverInfo, err
	}

	_serverInfo, errInfo := servers.Get(client, serverId).Extract()

	if errInfo != nil {
		return _serverInfo, errInfo
	}

	return _serverInfo, nil

}

/*
List OS Images
*/
func ListOSImages(authCredentials OpenstackAuth) ([]map[string]interface{}, error) {
	result1 := []map[string]interface{}{}

	provider, _ := auth(authCredentials)
	client, err := openstack.NewImageServiceV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err != nil {
		return result1, err
	}
	opt := &imageservices.ListOpts{Visibility: "public"}
	pager := imageservices.List(client, opt)
	pager.EachPage(func(page pagination.Page) (bool, error) {
		imageList, _ := imageservices.ExtractImages(page)
		for _, images := range imageList {
			result := make(map[string]interface{})
			result["ID"] = images.ID
			result["Name"] = images.Name
			result1 = append(result1, result)
		}
		return true, nil
	})

	return result1, nil
}

/*
List OS Flavors
*/
func ListOSFlavors(authCredentials OpenstackAuth) ([]map[string]interface{}, error) {
	result1 := []map[string]interface{}{}

	provider, _ := auth(authCredentials)
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err != nil {
		return result1, err
	}
	opt := &flavors.ListOpts{AccessType: flavors.PublicAccess}
	pager := flavors.ListDetail(client, opt)
	pager.EachPage(func(page pagination.Page) (bool, error) {
		flavorList, _ := flavors.ExtractFlavors(page)
		for _, flavor := range flavorList {
			result := make(map[string]interface{})
			result["ID"] = flavor.ID
			result["Name"] = flavor.Name
			result["Vcpu"] = flavor.VCPUs
			result["Ram"] = flavor.RAM
			result1 = append(result1, result)
		}
		return true, nil
	})

	return result1, err
}

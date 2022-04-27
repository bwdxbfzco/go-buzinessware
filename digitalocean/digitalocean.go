package digitalocean

import (
	"context"

	godo "github.com/digitalocean/godo"
)

var PERPAGE = 100

var Ctx = context.TODO()

func authDo(apiKey string) *godo.Client {
	client := godo.NewFromToken(apiKey)
	return client
}

func ListImages(apiKey string) ([]map[string]interface{}, error) {
	result1 := []map[string]interface{}{}
	client := authDo(apiKey)

	opts := &godo.ListOptions{PerPage: 200}

	result, _, err := client.Images.ListDistribution(Ctx, opts)

	if err != nil {
		return result1, err
	}

	for _, a := range result {
		result := make(map[string]interface{})
		result["ID"] = a.Slug
		result["Name"] = a.Description
		result["Regions"] = a.Regions
		result1 = append(result1, result)
	}

	return result1, nil

}

func ListRegions(apiKey string) ([]map[string]interface{}, error) {
	result1 := []map[string]interface{}{}

	client := authDo(apiKey)

	opts := &godo.ListOptions{}

	result, _, err := client.Regions.List(Ctx, opts)

	if err != nil {
		return result1, err
	}

	for _, a := range result {
		result := make(map[string]interface{})
		result["Region"] = a.Name
		result["Slug"] = a.Slug
		result1 = append(result1, result)
	}

	return result1, nil
}

func ListSizes(apiKey string) ([]map[string]interface{}, error) {
	result1 := []map[string]interface{}{}

	client := authDo(apiKey)

	opts := &godo.ListOptions{PerPage: 200}

	result, _, err := client.Sizes.List(Ctx, opts)

	if err != nil {
		return result1, err
	}

	for _, a := range result {
		result := make(map[string]interface{})
		result["ID"] = a.Slug
		result["Name"] = a.Slug
		result["Vcpu"] = a.Vcpus
		result["Ram"] = a.Memory
		result["Disk"] = a.Disk
		result["Regions"] = a.Regions

		result1 = append(result1, result)
	}

	return result1, nil
}

func ListDroplets(apiKey string) ([]map[string]interface{}, error) {
	result1 := []map[string]interface{}{}

	client := authDo(apiKey)

	opt := &godo.ListOptions{PerPage: PERPAGE}

	result, _, err := client.Droplets.List(Ctx, opt)

	if err != nil {
		return result1, err
	}

	for _, a := range result {
		result := make(map[string]interface{})
		result["ID"] = a.ID
		result["name"] = a.Name
		result["region"] = a.Region.Slug
		result["image"] = a.Image.Slug
		result["flavor"] = a.SizeSlug
		result["status"] = a.Status
		if len(a.Networks.V4) == 2 {
			result["ipaddress"] = a.Networks.V4[1].IPAddress
		} else {
			result["ipaddress"] = a.Networks.V4[0].IPAddress
		}

		result1 = append(result1, result)
	}

	//TODO we have to use resp.Meta to check the total records returned and then loop through till we get all the records.

	return result1, nil
}

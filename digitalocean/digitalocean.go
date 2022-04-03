package digitalocean

import (
	"context"

	godo "github.com/digitalocean/godo"
)

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

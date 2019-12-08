package resourceservice

import (
	"focus/cfg"
	"focus/types/resource"
)

func InitServiceResource() error {
	rows, err := cfg.FocusCtx.DB.Table("resource").Select("path, service_name").Joins(
		`left join service_resource sr on resource.id = sr.resource_id
				left join service on sr.service_id = service.id`).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()
	var resources []*resourcetype.Resource
	for rows.Next() {
		resource := &resourcetype.Resource{}
		rows.Scan(&resource.Path, &resource.ServiceName)
		resources = append(resources, resource)
	}
	cfg.FocusCtx.ServiceResource = resources
	return nil
}

func FilterResource(filter resourcetype.ResourceFilter) []*resourcetype.Resource {
	var resources []*resourcetype.Resource
	for _, resource := range cfg.FocusCtx.ServiceResource {
		if filter(resource) {
			resources = append(resources, resource)
		}
	}
	return resources
}

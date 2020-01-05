package resourceserv

import (
	"focus/cfg"
	"focus/types/resource"
	dbutil "focus/util/db"
)

func InitServiceResource() error {
	var resources []*resourcetype.Resource
	dbutil.NewDbExecutor(cfg.FocusCtx.DB.Table("resource").Select("service_id, path, service_name").Joins(
		`left join service_resource sr on resource.id = sr.resource_id
				left join service on sr.service_id = service.id`).Find(&resources))
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

func FilterSingleResource(filter resourcetype.ResourceFilter) *resourcetype.Resource {
	for _, resource := range cfg.FocusCtx.ServiceResource {
		if filter(resource) {
			return resource
		}
	}
	return nil
}

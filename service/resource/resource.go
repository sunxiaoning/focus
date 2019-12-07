package resourceservice

import (
	"context"
	"fmt"
	"focus/cfg"
	"focus/types/resource"
)

func QueryServiceResource(ctx context.Context) error {
	rows, err := cfg.FocusCtx.DB.Table("resource").Select("path, service_name").Joins(`left join service_resource sr on resource.id = sr.resource_id 
		left join service on sr.service_id = service.id`).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()
	var resources []*resourcetype.Resource
	for rows.Next() {
		resource := &resourcetype.Resource{}
		rows.Scan(&resource.Path, &resource.ServiceName)
		fmt.Println(resource.Path, resource.ServiceName)
		resources = append(resources, resource)
	}
	return nil
}

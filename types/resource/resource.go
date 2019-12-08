package resourcetype

type Resource struct {
	Path        string
	ServiceName string
}

type ResourceFilter func(resource *Resource) bool

package resourcetype

type Resource struct {
	ServiceId   int
	Path        string
	ServiceName string
}

type ResourceWithLimit struct {
	Resource          *Resource
	ConcurrencyNumber int
}

type ResourceFilter func(resource *Resource) bool

package pagetype

type PageInfo struct {
	PageIndex int
	PageSize  int
}

func NewPage(pageIndex int, pageSize int) *PageInfo {
	return &PageInfo{
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}
}

type PageQuery struct {
	Page   *PageInfo
	Params map[string]interface{}
}

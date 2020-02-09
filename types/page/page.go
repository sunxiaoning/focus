package pagetype

import "focus/types"

type PageInfo struct {
	PageIndex int
	PageSize  int
}

func NewPage(pageIndex int, pageSize int) *PageInfo {
	if pageIndex < 1 {
		types.InvalidParamPanic("pageIndex is invalid!")
	}
	if pageSize < 1 || pageSize > 1000 {
		types.InvalidParamPanic("pageSize is invalid!")
	}
	return &PageInfo{
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}
}

type PageQuery struct {
	Page   *PageInfo
	Params map[string]interface{}
}

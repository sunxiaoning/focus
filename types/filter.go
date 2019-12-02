package types

type Filter struct {

	// 顺序
	Order int

	// 过滤路径
	Paths []string

	// 排除路径
	ExculdePaths []string

	// 处理器
	Process Handle
}

type FilterComparable []*Filter

func (filterComparable FilterComparable) Len() int {
	return len(filterComparable)
}

func (filterComparable FilterComparable) Less(i, j int) bool {
	return filterComparable[i].Order < filterComparable[j].Order
}

func (filterComparable FilterComparable) Swap(i, j int) {
	filterComparable[i], filterComparable[j] = filterComparable[j], filterComparable[i]
}

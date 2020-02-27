package dto

// 分页查询
type PaginationParam struct {
	Total     int `json:"total"`     // 统计
	PageIndex int `json:"pageIndex"` // 页索引
	PageSize  int `json:"pageSize"`  // 页大小
}

type ResponseList struct {
	List       interface{}      `json:"list"`
	Pagination *PaginationParam `json:"pagination,omitempty"`
}

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Error   string `json:"error,omitempty"`
}

type QueryParams struct {
	PageNum   int // 分页计算
	PageSize  int // 分页条数
	PageIndex int // 当前页
}

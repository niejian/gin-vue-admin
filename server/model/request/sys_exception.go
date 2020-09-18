package request

// 根据索引名称，字段名称获取异常详情请求参数
type GetExceptionDetailStruct struct {
	IndexName    string `json:"indexName"`
	ExceptionTag string `json:"exceptionTag"`
	CreateDate   string `json:"createDate"`
}

// 根据索引名获取异常信息
type GetExceptionOverviewByIndexNameStruct struct {
	IndexName  string `json:"indexName"`
	CreateDate string `json:"createDate"`
	Days       string `json:"days"`
}

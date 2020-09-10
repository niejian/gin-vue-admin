package request

// 获取异常详情请求参数
type GetExceptionDetailStruct struct {
	IndexName    string `json:"indexName"`
	ExceptionTag string `json:"exceptionTag"`
}

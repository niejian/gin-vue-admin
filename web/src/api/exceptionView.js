import service from '@/utils/request'

// 
// @Summary 预览异常统计信息
// @Produce  application/json
// @Router /exception/viewException [post]
export const getExceptionView = () => {
  return service({
    url: "/exception/viewException",
    method: 'post',
  })
}

// @Summary 根据索引名称、字段名获取异常详细信息
// @Produce  application/json
// @Param {
//  indexName     string
//	exceptionTag  string
// }
// @Router /exception/exceptionDetails [post]
export const exceptionDetails = (data) => {
  return service({
    url: "/exception/exceptionDetails",
    method: 'post',
    data
  })
}

// @Summary 根据索引名称详细信息
// @Produce  application/json
// @Param {
//  indexName     string
// }
// @Router /exception/indexException [post]
export const indexException = (data) => {
  return service({
    url: "/exception/indexException",
    method: 'post',
    data
  })
}
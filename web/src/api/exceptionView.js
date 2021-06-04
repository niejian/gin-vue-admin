import service from '@/utils/request'

// 
// @Summary 异常汇总列表
// @Produce  application/json
// @Router /exception/exceptionOverview [post]
export const getExceptionOverview = (data) => {
  return service({
    url: "/exception/exceptionOverview",
    method: 'post',
    data
  })
}
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

// @Summary 根据索引名称和Id详细信息
// @Produce  application/json
// @Param {
//  indexName     string
// }
// @Router /exception/indexException [post]
export const getExceptionById = (indexName, id) => {
  return service({
    url: "/exception/getExceptionById/"+ indexName + "/" + id,
    method: 'get'
  })
}

// 添加或更新告警提醒
export const addOrUpdateErrorWarn = (data) => {
  return service({
    url: '/errorWarn/addOrUpdate',
    method: 'POST',
    data
  })
}

export const getConfInfoByIndexName = (indexName) => {
  return service({
    url: '/errorWarn/getConfInfoByIndexName/'+indexName,
    method: 'get',
    
  })
}

export const getUserInfo = (userId) => {
  return service({
    url: '/errorWarn/getUserInfo/'+userId,
    method: 'get',
    
  })
}
import service from '@/utils/request'

// @Summary 获取所有namespace
// @Produce  application/json
// @Router/k8sapi/listNs [post]
export const listNs = () => {
  return service({
    url: "/k8sapi/nsList",
    method: 'GET'
 
  })
}

export const deploys = (ns) => {
  return service({
    url: "/k8sapi/deploys/"+ns,
    method: 'GET'
 
  })
}

// 添加配置信息
export const addOrUpdate = (data) => {
  return service({
    url: "/k8sapi/watchConf/addOrUpdate",
    method: "POST",
    data
  })
}

// 根据命名空间、应用名称信息查询配置信息
export const getConfByNsAndAppName = (ns, appName) => {
  return service({
    url: "/k8sapi/getConfByNsAndAppName/"+ns+"/"+appName,
    method: "GET"
  })
}

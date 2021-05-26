// v1 doc

package v1

import (
	"fmt"
	"gin-vue-admin/global/response"
	"gin-vue-admin/model/k8s"
	"gin-vue-admin/model/request"
	resp "gin-vue-admin/model/response"
	"gin-vue-admin/service"
	"gin-vue-admin/utils"
	"github.com/gin-gonic/gin"
)

//ListNs doc
//@Description: 获取所有namespace
//@Author niejian
//@Date 2021-04-29 16:55:59
//@param clientSet
//@return []string
//@return error
func ListNs(c *gin.Context) {
	listNs, err := service.ListNs()
	if err != nil {
		response.FailWithMessage("获取命名空间失败", c)
		return
	}

	response.OkDetailed(&listNs, "请求成功", c)

}

//ListDeploy doc
//@Description: 获取k8s中所有的deploy资源
//@Author niejian
//@Date 2021-04-30 09:35:22
//@param c
//@return *v1.DeploymentList
//@return error
func ListDeploy(c *gin.Context) {
	// xxx/ns 获取reset请求参数信息
	ns := c.Param("ns")
	if "" == ns {
		response.FailWithMessage("请选择命名空间", c)
		return
	}
	deploymentList, err := service.ListDeploy(ns)
	if err != nil {
		response.FailWithMessage("获取deploy信息失败", c)
		return
	}
	response.OkDetailed(deploymentList, "请求成功", c)
}

//ListPodsByLabels doc
//@Description: 通过标签获取pod信息
//@Author niejian
//@Date 2021-05-07 14:07:52
//@param c
func ListPodsByLabels(c *gin.Context) {
	ns := c.Param("ns")
	labels := c.Param("labels")
	if "" == ns {
		response.FailWithMessage("请选择命名空间", c)
		return
	}
	podNamesByLabels, err := service.ListPodNamesByLabels(ns, labels)
	if err != nil {
		response.FailWithMessage("获取deploy信息失败", c)
		return
	}
	response.OkDetailed(podNamesByLabels, "请求成功", c)
}

//AddWatchdogConf doc
//@Description: 添加看门狗信息
//@Author niejian
//@Date 2021-05-07 14:08:47
//@param c
func AddWatchdogConf(c *gin.Context) {

	var requestData request.ErrorLogAlterConfig
	// 绑定参数
	c.ShouldBindJSON(&requestData)
	// 校验请求参数是否为空
	rules := utils.Rules{
		"ns":        {utils.NotEmpty()},
		"appName":   {utils.NotEmpty()},
		"errs":      {utils.NotEmpty()},
		"ignores":   {utils.NotEmpty()},
		"toUserIds": {utils.NotEmpty()},
	}

	verifyErr := utils.Verify(requestData, rules)
	if nil != verifyErr {
		response.FailWithMessage(verifyErr.Error(), c)
		return
	}
	// 组装成结构体
	logAlter := &k8s.ErrorLogAlterConfig{
		Errs:        requestData.Errs,
		AppName:     requestData.AppName,
		Namespace:   requestData.Namespace,
		Ignores:     requestData.Ignores,
		ToUserIds:   requestData.ToUserIds,
		EnableStore: requestData.EnableStore,
	}

	// 判断新增还是更新
	id := requestData.Id
	if id <= 0 {
		err, _ := service.AddWatchdogConf(logAlter)
		if err != nil {
			response.FailWithDetailed(response.ERROR, "添加失败", fmt.Sprintf("%v", err), c)
		} else {
			response.OkDetailed("", "添加成功", c)
		}
	} else {
		// 更新
		logAlter.ID = uint(id)
		service.UpdateAlterConfig(logAlter)
		response.OkDetailed("", "更新成功", c)
	}

}

//GetConfByNsAndAppName doc
//@Description: 通过命名空间和应用名称获取配置信息
//@Author niejian
//@Date 2021-05-07 16:45:52
//@param c
func GetConfByNsAndAppName(c *gin.Context) {
	ns := c.Param("ns")
	appName := c.Param("appName")

	if "" == ns {
		response.FailWithMessage("请选择命名空间", c)
		return
	}

	if "" == appName {
		response.FailWithMessage("请选择应用名称", c)
		return
	}
	configs := service.GetLogAlterConfByNsAndAppName(ns, appName)
	response.OkDetailed(&configs, "请求成功", c)
}

//GetConfigList doc
//@Description: 分页获取配置列表
//@Author niejian
//@Date 2021-05-26 11:16:02
//@param c
func GetConfigList(c *gin.Context) {

	// 此结构体仅本方法使用
	var sp request.SearchWatchdogConfParam
	_ = c.ShouldBindJSON(&sp)
	PageVerifyErr := utils.Verify(sp.PageInfo, utils.CustomizeMap["PageVerify"])
	if PageVerifyErr != nil {
		response.FailWithMessage(PageVerifyErr.Error(), c)
		return
	}

	list, total, err := service.GetWatchDogConfList(sp.ErrorLogAlterConfig, sp.PageInfo)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取数据失败，%v", err), c)
	} else {
		response.OkWithData(resp.PageResult{
			List:     &list,
			Total:    total,
			Page:     sp.PageInfo.Page,
			PageSize: sp.PageInfo.PageSize,
		}, c)
	}

}

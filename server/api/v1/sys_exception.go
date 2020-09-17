package v1

import (
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/global/response"
	"gin-vue-admin/model/request"
	resp "gin-vue-admin/model/response"
	"gin-vue-admin/service"
	"gin-vue-admin/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

const (
	INDEX_PREFIX    = "watchdog_store_"
	FormatStartTime = "20060102"
)

// @Tags GetExceptionOverview
// @Summary 获取异常的x天内的报表信息
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body request.RegisterAndLoginStruct true "可以什么都不填"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /exception/view [post]
func GetExceptionOverview(c *gin.Context) {
	var requestData request.GetExceptionOverviewByIndexNameStruct
	// 获取所有相关的索引信息
	names, err := service.FindAppIndex(INDEX_PREFIX)
	c.ShouldBindJSON(&requestData)
	indexName := requestData.IndexName
	var initIndexName string = ""
	if err != nil {
		global.GVA_LOG.Error("获取索引失败：", zap.Any("err", err))
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
	} else {

		if len(names) > 0 && "" != indexName {
			// 初始化页面
			initIndexName = indexName
		} else if len(names) > 0 && "" == indexName {
			initIndexName = names[0]
		}

		// 初始化页面
		global.GVA_LOG.Info("初始化异常报表页面：", zap.Any("indexName", initIndexName))
		// 聚合查询错误列表信息
		aggIndexs := service.GetExceptionOverview(initIndexName, 30)
		response.OkWithData(resp.IndexNameResponse{IndexNames: names, AggIndexs: aggIndexs}, c)

	}
}

// @Tags GetExceptionView
// @Summary 获取异常视图信息
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body request.RegisterAndLoginStruct true "可以什么都不填"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /exception/view [post]
func GetExceptionView(c *gin.Context) {
	names, err := service.FindAppIndex(INDEX_PREFIX)

	if err != nil {
		global.GVA_LOG.Error("获取索引失败：", zap.Any("err", err))
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
	} else {
		if len(names) > 0 {
			// 初始化页面
			initIndexName := names[0]
			fmt.Printf("initIndexName %v \n", initIndexName)
			// 聚合查询错误列表信息
			aggIndexs := service.IndexOverview(initIndexName, time.Now().Format(FormatStartTime))

			response.OkWithData(resp.IndexNameResponse{IndexNames: names, AggIndexs: aggIndexs}, c)

		}
	}

}

func GetExceptionDailyViewByIndexName(c *gin.Context) {
	var requestData request.GetExceptionOverviewByIndexNameStruct
	c.ShouldBindJSON(&requestData)
	// 校验请求参数是否为空
	rules := utils.Rules{
		"IndexName":  {utils.NotEmpty()},
		"CreateDate": {utils.NotEmpty()},
	}

	verifyErr := utils.Verify(requestData, rules)
	if nil != verifyErr {
		response.FailWithMessage(verifyErr.Error(), c)
		return
	}
	// 聚合查询错误列表信息
	aggIndexs := service.IndexOverview(requestData.IndexName, requestData.CreateDate)

	response.OkWithData(resp.AggIndexResponse{AggIndexs: aggIndexs}, c)

}

// @Tags GetExceptionDetails
// @Summary 获取异常详细信息
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body request.RegisterAndLoginStruct true "可以什么都不填"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /exception/view [post]
func GetExceptionDetails(c *gin.Context) {
	// 获取请求参数
	var exceptionRequest request.GetExceptionDetailStruct
	c.ShouldBindJSON(&exceptionRequest)
	// 校验
	exceptionValid := utils.Rules{
		"IndexName":    {utils.NotEmpty()},
		"ExceptionTag": {utils.NotEmpty()},
		"CreateDate":   {utils.NotEmpty()},
	}

	exceptionVerifyErr := utils.Verify(exceptionRequest, exceptionValid)
	if exceptionVerifyErr != nil {
		response.FailWithMessage(exceptionVerifyErr.Error(), c)
		return
	}

	// 查询索引详细信息
	datas, err := service.FindFDatasByIndiceName(&exceptionRequest)
	if nil != err {
		global.GVA_LOG.Info("查询异常详细信息失败，请重试：", zap.Any("err", err))
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)

	} else {

		response.OkWithData(datas, c)
	}

}

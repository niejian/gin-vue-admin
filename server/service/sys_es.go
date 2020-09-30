package service

import (
	"context"
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/model/es"
	"gin-vue-admin/model/request"
	"gin-vue-admin/utils"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"reflect"
	"strings"
	"time"
)

var (
	EXCEPTION_AGGS    = "exceptionTag_aggs"
	EXCEPTION_TAG     = "exceptionTag"
	CREATEDATE        = "createDate"
	FormatStartTime   = "20060102"
	GroupByCreateDate = "groupByCreateDate"
)

// 根据索引前缀获取所有相关的索引列表
func FindAppIndex(indexPrefix string) ([]string, error) {
	var indexNames []string
	names, err := global.GVA_ES.IndexNames()
	if err != nil {
		global.GVA_LOG.Error("es获取索引失败：", zap.Any("err", err))
		return nil, err
	}
	for _, name := range names {
		if strings.Index(name, indexPrefix) == 0 {
			if strings.Contains(name, ".") {
				continue
			}
			indexNames = append(indexNames, name)
		}
	}

	return indexNames, nil
}

// 聚合查询改索引近{days}天内的信息
func GetExceptionOverview(indexName string, days int) []es.AggIndex {
	var result []es.AggIndex
	ctx := context.Background()
	if days <= 0 {
		days = 30
	}
	duration := -24 * days
	durationHour, err := time.ParseDuration(fmt.Sprintf("%d%s", duration, "h"))
	if err != nil {
		global.GVA_LOG.Error("时间转换失败", zap.Any("err", err))
	}

	from := time.Now().Add(durationHour).Format(FormatStartTime)
	now := time.Now().Format(FormatStartTime)

	query := elastic.NewRangeQuery(CREATEDATE).Gte(from).Lte(now)

	// 根据结果key（插入日期：20200916）升序排序
	agg := elastic.NewTermsAggregation().Field(CREATEDATE).Order("_key", true)
	searchResult, err := global.GVA_ES.
		Search(indexName).
		Query(query).
		Aggregation(GroupByCreateDate, agg).
		Do(ctx)
	if err != nil {
		global.GVA_LOG.Error("GetExceptionOverview失败", zap.Any("err", err))
	}

	// 获取聚合结果
	terms, exists := searchResult.Aggregations.Terms(GroupByCreateDate)
	if exists {
		for _, bucket := range terms.Buckets {
			var aggType es.AggIndex
			aggType.Key = bucket.Key.(string)
			aggType.DocCount = bucket.DocCount
			result = append(result, aggType)
		}
	}
	return result
}

// 聚合查询，获取错误统计信息
func IndexOverview(indexName, createDate string) []es.AggIndex {
	var result []es.AggIndex
	ctx := context.Background()
	query := elastic.NewTermQuery(CREATEDATE, createDate)
	agg := elastic.NewTermsAggregation().Field(EXCEPTION_TAG)
	searchResult, err := global.GVA_ES.
		Search(indexName).
		Size(0).
		Query(query).
		Aggregation(EXCEPTION_AGGS, agg).
		Do(ctx)

	if nil != err {
		global.GVA_LOG.Error("聚合查询失败:", zap.Any("err", err))
	}

	// 获取聚合结果
	terms, exists := searchResult.Aggregations.Terms(EXCEPTION_AGGS)
	if exists {
		for _, bucket := range terms.Buckets {
			var aggType es.AggIndex
			aggType.Key = bucket.Key.(string)
			aggType.DocCount = bucket.DocCount
			result = append(result, aggType)
		}
	}

	return result
}

// 通过索引名称、字段名获取到具体信息
func FindFDatasByIndiceName(queryExceptionRequest *request.GetExceptionDetailStruct) ([]es.Exception, error) {
	var exs []es.Exception
	ctx := context.Background()

	queryExceptionTag := elastic.NewTermQuery(EXCEPTION_TAG, queryExceptionRequest.ExceptionTag)

	queryCreateDate := elastic.NewTermQuery(CREATEDATE, queryExceptionRequest.CreateDate)
	// 组合条件查询
	query := elastic.NewBoolQuery().Must(queryExceptionTag, queryCreateDate)
	// 数据总数
	count, err := global.GVA_ES.Count(queryExceptionRequest.IndexName).
		Query(query).Do(ctx)

	if nil != err {
		count = 0
	}
	global.GVA_LOG.Info("数据总量：", zap.Any("total", count))
	// 翻页获取数据信息
	pageNum := 0
	pageSize := 100
	index := 0

	for {
		datas, err := getExceptionDetailDataByPage(queryExceptionRequest.IndexName,
			query, pageNum, pageSize)
		if err != nil {
			global.GVA_LOG.Error("查询失败，", zap.Any("err", err))
		} else {

			// 先追加数据
			for _, data := range datas {
				exs = append(exs, data)
			}

			// 再判断是否退出循环
			if len(datas) < pageSize {
				break
			}

			if pageNum > int(count) {
				break
			}
			index = index + 1
			pageNum = index * pageSize

		}
	}

	fmt.Printf("======》总数：%v\n", len(exs))

	return exs, nil
}

// 分页获取错误信息
func getExceptionDetailDataByPage(indexName string, query elastic.Query,
	pageNum, pageSize int) ([]es.Exception, error) {
	var exs []es.Exception
	ctx := context.Background()

	global.GVA_LOG.Info("开始分页获取异常详情",
		zap.Int("from", pageNum),
		zap.Int("size", pageSize))

	searchResult, err := global.
		GVA_ES.Search(indexName).
		Query(query).
		Sort("id", false).
		From(pageNum).
		Size(pageSize).
		FilterPath("hits.total", "hits.hits", "hits.hits._source&_source=id,ip,createDate,appName,exceptionTag,createTime"). // 排除详情字段，懒加载
		Do(ctx)
	if nil != err {
		global.GVA_LOG.Info("分页查询失败：",
			zap.Any("err", err),
			zap.Any("pageNum", pageNum))
		return exs, err
	}

	var ex es.Exception
	for _, item := range searchResult.Each(reflect.TypeOf(ex)) {
		i := item.(es.Exception)
		//msg := "<div>" + strings.ReplaceAll(i.Msg, "\n", "</div><div>") + "</div>"
		i.Msg = ""
		// 时间格式化
		createTime := i.CreateTime
		i.CreateDate = utils.FormatTimeByTimestamp(createTime)
		exs = append(exs, i)
	}

	return exs, nil
}

// 根据Id获取某条索引下面的具体数据信息
func GetExceptionDetailById(indexName, id string) (es.Exception, error) {
	var ex es.Exception
	ctx := context.Background()
	query := elastic.NewTermQuery("id", id)
	searchResult, err := global.GVA_ES.Search(indexName).
		Query(query).
		Do(ctx)

	if nil != err {
		global.GVA_LOG.Info("根据Id查询异常信息失败：",
			zap.Any("index", indexName),
			zap.Any("id", id))
		return ex, err
	}

	for _, item := range searchResult.Each(reflect.TypeOf(ex)) {
		i := item.(es.Exception)
		//msg := i.Id + "<br/>" + "<div>" + strings.ReplaceAll(i.Msg, "\n", "</div><div>") + "</div>"sc
		msg := "<div>" + strings.ReplaceAll(i.Msg, "\n", "</div><div>") + "</div>"
		i.Msg = msg
		ex = i
		break
	}

	return ex, nil
}

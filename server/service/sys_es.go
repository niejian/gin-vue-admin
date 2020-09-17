package service

import (
	"context"
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/model/es"
	"gin-vue-admin/model/request"
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

	searchResult, err := global.
		GVA_ES.Search(queryExceptionRequest.IndexName).
		Query(query).
		Sort("id", true).
		From(0).
		Size(10000).
		Do(ctx)
	if nil != err {
		global.GVA_LOG.Info("查询失败：", zap.Any("err", err))
		return nil, err
	}

	var ex es.Exception
	for _, item := range searchResult.Each(reflect.TypeOf(ex)) {
		i := item.(es.Exception)
		msg := "<div>" + strings.ReplaceAll(i.Msg, "\n", "</div><div>") + "</div>"
		i.Msg = msg
		exs = append(exs, i)
	}

	return exs, nil
}

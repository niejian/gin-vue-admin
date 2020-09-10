package service

import (
	"context"
	"gin-vue-admin/global"
	"gin-vue-admin/model/es"
	"gin-vue-admin/model/request"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"reflect"
	"strings"
)

var (
	EXCEPTION_AGGS = "exceptionTag_aggs"
	EXCEPTION_TAG  = "exceptionTag"
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
			indexNames = append(indexNames, name)
		}
	}

	return indexNames, nil
}

// 聚合查询，获取错误统计信息
func IndexOverview(indexName string) []es.AggIndex {
	var result []es.AggIndex
	ctx := context.Background()
	agg := elastic.NewTermsAggregation().Field(EXCEPTION_TAG)
	searchResult, err := global.GVA_ES.Search(indexName).Size(0).
		Aggregation(EXCEPTION_AGGS, agg).Do(ctx)

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
	query := elastic.NewTermQuery(EXCEPTION_TAG, queryExceptionRequest.ExceptionTag)
	searchResult, err := global.
		GVA_ES.Search(queryExceptionRequest.IndexName).
		Query(query).
		Sort("id", true).
		From(0).
		Do(ctx)
	if nil != err {
		global.GVA_LOG.Info("查询失败：", zap.Any("err", err))
		return nil, err
	}

	var ex es.Exception
	for _, item := range searchResult.Each(reflect.TypeOf(ex)) {
		i := item.(es.Exception)
		exs = append(exs, i)
	}

	return exs, nil
}

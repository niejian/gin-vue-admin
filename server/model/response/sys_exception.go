package response

import "gin-vue-admin/model/es"

type IndexNameResponse struct {
	IndexNames []string      `json:"indexNames"`
	AggIndexs  []es.AggIndex `json:"aggIndexs"`
}

type AggIndexResponse struct {
	AggIndexs []es.AggIndex `json:"aggIndexs"`
}

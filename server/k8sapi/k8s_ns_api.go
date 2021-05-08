// k8sapi doc

package k8sapi

import (
	"context"
	"gin-vue-admin/global"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//ListNs doc
//@Description: 获取所有namespace
//@Author niejian
//@Date 2021-04-29 16:55:59
//@param clientSet
//@return []string
//@return error
func ListNs(clientSet *kubernetes.Clientset) ([]string, error) {
	nsList := make([]string, 10)
	namespaces, err := clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		global.GVA_LOG.Error("获取命名空间失败")
		return nil, err
	}

	items := namespaces.Items
	for _, ns := range items {
		nsList = append(nsList, ns.Name)
	}
	return nsList, nil
}

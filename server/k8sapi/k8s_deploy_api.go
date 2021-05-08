// k8sapi doc

package k8sapi

import (
	"context"
	"gin-vue-admin/global"
	"go.uber.org/zap"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//ListDeploy doc
//@Description: 获取k8s中所有的deploy资源
//@Author niejian
//@Date 2021-04-29 16:57:31
//@param ns
//@param clientSet
//@return []string
//@return error
func ListDeploy(ns string, clientSet *kubernetes.Clientset) (*v1.DeploymentList, error) {
	deploymentList, err := clientSet.AppsV1().Deployments(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		global.GVA_LOG.Error("获取deployment资源失败", zap.String("ns", ns))
		return nil, err
	}
	return deploymentList, nil
}

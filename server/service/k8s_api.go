// service doc

package service

import (
	"context"
	"gin-vue-admin/global"
	"gin-vue-admin/k8sapi"
	"gin-vue-admin/model/k8s"
	"gin-vue-admin/model/request"
	"go.uber.org/zap"
	v1 "k8s.io/api/apps/v1"
	vCore "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

var (
	clientSet  = k8sapi.InitK8s()
	anno       = "kubesphere.io/creator"
	kubesphere = "kubesphere"
	kubeSystem = "kube-system"
)

//var db = global.GVA_DB.Table("log_alter_conf")

//ListNs doc
//@Description: 获取所有namespace
//@Author niejian
//@Date 2021-04-29 16:55:59
//@param clientSet
//@return []string
//@return error
func ListNs() (*[]string, error) {
	var nsList []string
	namespaces, err := clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		global.GVA_LOG.Error("获取命名空间失败")
		return nil, err
	}

	items := namespaces.Items
	for _, ns := range items {
		namespace := ns.Name
		// kubesphere开头，过滤
		if strings.HasPrefix(namespace, kubesphere) {
			continue
		}

		// kube-system过滤
		if namespace == kubeSystem {
			continue
		}
		// 获取普通用户创建的命名空间(剔除admin创建的namespace)
		annotations := ns.Annotations
		if creator, ok := annotations[anno]; ok && creator == "admin" {
			continue
		}
		isContinue := true
		for _, ns := range global.GVA_CONFIG.ExcludeNs {
			if ns == namespace || strings.HasPrefix(namespace, ns) {
				global.GVA_LOG.Info(" ns=" + namespace + "， 为忽略的命名空间，不处理")
				isContinue = false
				break
			}
		}

		if isContinue {
			nsList = append(nsList, namespace)
		}
	}
	return &nsList, nil
}

//ListDeploy doc
//@Description: 获取k8s中所有的deploy资源
//@Author niejian
//@Date 2021-04-29 16:57:31
//@param ns
//@param clientSet
//@return []string
//@return error
func ListDeploy(ns string) (*v1.DeploymentList, error) {
	deploymentList, err := clientSet.AppsV1().Deployments(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		global.GVA_LOG.Error("获取deployment资源失败", zap.String("ns", ns))
		return nil, err
	}
	return deploymentList, nil
}

//ListPodNamesByLabels doc
//@Description: 通过标签获取pod名称信息，多个标签 key1=value1,key2=value2
//@Author niejian
//@Date 2021-05-06 11:40:18
//@param namespace
//@param labels
//@return []string
//@return error
func ListPodNamesByLabels(namespace, labels string) ([]string, error) {
	podInterface := clientSet.CoreV1().Pods(namespace)
	var list *vCore.PodList
	list, err := podInterface.List(context.TODO(), metav1.ListOptions{LabelSelector: labels})
	if err != nil {
		return nil, err
	}

	//if "" != labels {
	//	list1, err := podInterface.List(context.TODO(), metav1.ListOptions{LabelSelector: labels})
	//	if err != nil {
	//		return nil, err
	//	}
	//	list = list1
	//} else {
	//	list1, err := podInterface.List(context.TODO(), metav1.ListOptions{})
	//	if err != nil {
	//		return nil, err
	//	}
	//	list = list1
	//}

	var podNames []string
	for _, pod := range list.Items {
		//podNames[index] = pod.Name
		podNames = append(podNames, pod.Name)
	}
	return podNames, nil
}

//AddWatchdogConf doc
//@Description: 添加配置信息
//@Author niejian
//@Date 2021-05-07 14:20:11
//@param config
//@return error
//@return *k8s.ErrorLogAlterConfig
func AddWatchdogConf(config *k8s.ErrorLogAlterConfig) (error, *k8s.ErrorLogAlterConfig) {
	err := global.GVA_DB.Table("log_alter_conf").Create(&config).Error
	return err, config
}

//GetLogAlterConfByNsAndAppName doc
//@Description: 通过命名空间和应用名称获取告警配置信息
//@Author niejian
//@Date 2021-05-07 15:15:07
//@param ns 命名空间
//@param appName 应用名称data
func GetLogAlterConfByNsAndAppName(ns, appName string) []*k8s.ErrorLogAlterConfig {
	var data []*k8s.ErrorLogAlterConfig
	global.GVA_DB.Table("log_alter_conf").Where("namespace = ? and app_name = ?", ns, appName).Find(&data)
	//if err != nil {
	//	global.GVA_LOG.Error("GetLogAlterConfByNsAndAppName 查询失败。", zap.String("err", err.Error()))
	//	return nil
	//}
	return data
}

func UpdateAlterConfig(config *k8s.ErrorLogAlterConfig) {
	// 查询该记录的创建时间等，然后赋值
	var data k8s.ErrorLogAlterConfig
	// 按住键查询
	global.GVA_DB.Table("log_alter_conf").Find(&data, config.ID)
	// 查到数据
	if data.ID > 0 {
		config.CreatedAt = data.CreatedAt
		config.DeletedAt = data.DeletedAt
	}
	global.GVA_DB.Table("log_alter_conf").Save(config)
}

//GetWatchDogConfList doc
//@Description: 分页获取配置列表信息
//@Author niejian
//@Date 2021-05-26 11:20:33
//@param api
//@param info
//@return list
//@return total
//@return err
func GetWatchDogConfList(conf k8s.ErrorLogAlterConfig, info request.PageInfo) (list *[]k8s.ErrorLogAlterConfig, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&k8s.ErrorLogAlterConfig{})
	var confList []k8s.ErrorLogAlterConfig

	if conf.AppName != "" {
		db = db.Debug().Where("app_name LIKE ?", "%"+conf.AppName+"%")
	}

	if conf.Namespace != "" {
		db = db.Where("namespace LIKE ?", "%"+conf.Namespace+"%")
	}

	err = db.Count(&total).Error

	if err == nil {
		// 分页查询
		err = db.Debug().Limit(limit).Offset(offset).Find(&confList).Error
	}
	return &confList, total, err

}

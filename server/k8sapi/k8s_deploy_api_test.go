// k8sapi doc

package k8sapi

import (
	"fmt"
	"testing"
)

func TestListDeploy(t *testing.T) {
	t.Run("获取deploy信息", func(t *testing.T) {
		deploy, err := ListDeploy("sso", clientSet)
		if err != nil {
			fmt.Printf("get deploy error %v \n", err)
			return
		}
		fmt.Printf("%v \n", deploy)
	})
}

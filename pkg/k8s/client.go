package k8s

import (
	"os"
	"path/filepath"

	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	clientcmd "k8s.io/client-go/tools/clientcmd"
	homedir "k8s.io/client-go/util/homedir"
	"k8s.io/klog"
)

var KubeClient *KubeClientContext

type KubeClientContext struct {
	Config *restclient.Config
	Client *clientset.Clientset
}

// NewKubeClient 创建 kube client, bind 接口需要用此客户端建立 Pod 与 Node 的绑定关系.
func NewKubeClient() *KubeClientContext {
	// 先尝试从 ~/.kube 目录下获取配置, 如果没有, 则尝试寻找 Pod 内置的认证配置
	var kubeconfig string
	kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
	if _, err := os.Stat(kubeconfig); err != nil {
		klog.Warningf("kube config %s doesn't exist, buid config from InCluster", kubeconfig)
		kubeconfig = ""
	}
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		klog.Errorf("Error building kubeconfig: %s", err.Error())
	}

	// kubeClient 用于集群内资源操作, crdClient 用于操作 crd 资源本身.
	// 具体区别目前还不清楚, 不过示例中大多都是这么做的.
	kubeClient, err := clientset.NewForConfig(kubeConfig)
	if err != nil {
		klog.Errorf("Error building kubernetes clientset: %s", err.Error())
	}

	return &KubeClientContext{
		Config: kubeConfig,
		Client: kubeClient,
	}
}

func InitKubeClient() {
	KubeClient = NewKubeClient()
}

package response

type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type AKSInfo struct {
	NodeNum    int    `json:"node_num"`
	K8sVersion string `json:"k8s_version"`
	KubeConfig string `json:"kube_config"`
}

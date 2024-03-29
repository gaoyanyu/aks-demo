package action

import (
	"aks-demo/model/response"
	"aks-demo/pkg/util"
	"errors"
	"fmt"
	"os"
	"strings"

	"k8s.io/klog/v2"
)

func CreateAksShort(master, version string) error {
	// bash script to init k8s cluster
	initAction := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s bash /root/yanyu/install-k8s.sh", master)
	res, err, output := util.ExecShortCMD("bash", "-c", initAction)
	klog.Infof("init k8s output: %s", output)
	if err != nil {
		klog.Error(err)
	}
	if res != 0 {
		klog.Error(fmt.Sprintf("Fail to create aks, code:%d, err:%+v, output: %s", res, err, output))
	}

	return nil
}

func CreateAks(master, version string) error {
	infoFile, err := os.Create(fmt.Sprintf("%s/%s.create.info.log", fmt.Sprintf("./logs"), master))
	if err != nil {
		klog.Errorf("create master %s info file failed", master)
		infoFile = nil
	}
	errFile, err := os.Create(fmt.Sprintf("%s/%s.create.err.log", fmt.Sprintf("./logs"), master))
	if err != nil {
		klog.Errorf("create master %s error file failed", master)
		errFile = nil
	}

	// init k8s cluster
	initAction := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s kubeadm init --kubernetes-version=%s --image-repository=registry.sensetime.com/diamond --apiserver-advertise-address=%s --service-cidr=10.96.0.0/12 --pod-network-cidr=10.244.0.0/16 -v=10", master, version, master)
	res, err, output := util.ExecCMD(infoFile, errFile, "bash", "-c", initAction)
	klog.Infof("init k8s output: %s", output)
	if err != nil {
		klog.Error(err)
		return err
	}
	if res != 0 {
		return errors.New(fmt.Sprintf("Fail to create aks, code:%d, err:%+v, output: %s", res, err, output))
	}

	// mv k8s config
	changeConfigPath := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s 'mkdir -p $HOME/.kube && sudo cp -rf /etc/kubernetes/admin.conf $HOME/.kube/config && sudo chown $(id -u):$(id -g) $HOME/.kube/config'", master)
	res, err, output = util.ExecCMD(infoFile, errFile, "bash", "-c", changeConfigPath)
	klog.Infof("changeConfigPath output: %s", output)
	if err != nil {
		klog.Error(err)
		return err
	}
	if res != 0 {
		return errors.New(fmt.Sprintf("Fail to mv k8s config, code:%d, err:%+v, output: %s", res, err, output))
	}

	// install calico cni
	installCalico := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s 'kubectl create -f /root/yanyu/calico/tigera-operator.yaml && kubectl create -f /root/yanyu/calico/custom-resources.yaml'", master)
	res, err, output = util.ExecCMD(infoFile, errFile, "bash", "-c", installCalico)
	klog.Infof("installCalico output: %s", output)
	if err != nil {
		klog.Error(err)
		//return err
	}
	//if res != 0 {
	//	return errors.New(fmt.Sprintf("Fail to install cni, code:%d, err:%+v, output: %s", res, err, output))
	//}

	return nil
}

func GetAks(master string) (error, response.AKSInfo) {
	var aksInfo response.AKSInfo
	node := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s kubectl get node -owide | grep Ready", master)
	res, err, output := util.ExecShortCMD("bash", "-c", node)
	if err != nil {
		klog.Error(err)
		return err, aksInfo
	}
	if res != 0 {
		return errors.New(fmt.Sprintf("Fail to check k8s node, code:%d, err:%+v, output: %s", res, err, output)), aksInfo
	}

	nodeNum := strings.Count(output, "Ready")
	if nodeNum < 1 {
		return errors.New("not found"), response.AKSInfo{}
	}
	aksInfo.NodeNum = nodeNum
	aksInfo.K8sVersion = "v1.21.5"

	//kubeVersion := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s kubectl version | grep Server | grep GitVersion", master)
	//res, err, output = util.ExecShortCMD("bash", "-c", kubeVersion)
	//if err != nil {
	//	klog.Error(err)
	//	return err, aksInfo
	//}
	//if res != 0 {
	//	return errors.New(fmt.Sprintf("Fail to check cni, code:%d, err:%+v, output: %s", res, err, output)), aksInfo
	//}

	kubeConfig := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s cat $HOME/.kube/config", master)
	res, err, output = util.ExecShortCMD("bash", "-c", kubeConfig)
	if err != nil {
		klog.Error(err)
		return err, aksInfo
	}
	if res != 0 {
		return errors.New(fmt.Sprintf("Fail to check cni, code:%d, err:%+v, output: %s", res, err, output)), aksInfo
	}
	cut1Str := "Warning: Permanently added '10.119.250.16' (ED25519) to the list of known hosts."
	cut2Str := "Authorized uses only. All activity may be monitored and reported."

	kube := strings.TrimLeft(strings.TrimLeft(strings.TrimLeft(strings.TrimLeft(strings.TrimLeft(output, cut1Str), "\r"), "\n"), cut2Str), "\n")
	aksInfo.KubeConfig = kube

	return nil, aksInfo
}

func DeleteAks(master string) error {
	infoFile, err := os.Create(fmt.Sprintf("%s/%s.delete.info.log", fmt.Sprintf("./logs"), master))
	if err != nil {
		klog.Errorf("create master %s info file failed", master)
		infoFile = nil
	}
	errFile, err := os.Create(fmt.Sprintf("%s/%s.delete.err.log", fmt.Sprintf("./logs"), master))
	if err != nil {
		klog.Errorf("create master %s error file failed", master)
		errFile = nil
	}

	deleteAction := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s kubeadm reset -f", master)
	res, err, output := util.ExecCMD(infoFile, errFile, "bash", "-c", deleteAction)
	klog.Infof("delete output: %s", output)
	if err != nil {
		klog.Error(err)
		return err
	}
	if res != 0 {
		return errors.New(fmt.Sprintf("Fail to delete k8s, code:%d, err:%+v, output: %s", res, err, output))
	}

	return nil
}

func UpdateAks(master string) error {
	infoFile, err := os.Create(fmt.Sprintf("%s/%s.update.info.log", fmt.Sprintf("./logs"), master))
	if err != nil {
		klog.Errorf("create master %s info file failed", master)
		infoFile = nil
	}
	errFile, err := os.Create(fmt.Sprintf("%s/%s.update.err.log", fmt.Sprintf("./logs"), master))
	if err != nil {
		klog.Errorf("create master %s error file failed", master)
		errFile = nil
	}

	update := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s kubeadm upgrade plan", master)
	res, err, output := util.ExecCMD(infoFile, errFile, "bash", "-c", update)
	klog.Infof("update k8s output: %s", output)
	if err != nil {
		klog.Error(err)
		return err
	}
	if res != 0 {
		return errors.New(fmt.Sprintf("Fail to update aks, code:%d, err:%+v, output: %s", res, err, output))
	}

	return nil
}

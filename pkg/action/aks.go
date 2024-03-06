package action

import (
	"aks-demo/pkg/util"
	"errors"
	"fmt"
	"os"

	"k8s.io/klog/v2"
)

func CreateAks(master string) error {
	// init k8s cluster
	createAction := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s kubeadm init --kubernetes-version=v1.21.5 --image-repository=registry.sensetime.com/diamond --apiserver-advertise-address=%s --service-cidr=10.96.0.0/12 --pod-network-cidr=10.244.0.0/16 -v=10", master, master)
	res, err, output := util.ExecCMD(nil, nil, "bash", "-c", createAction)
	if err != nil {
		klog.Error(err)
		return err
	}
	if res != 0 {
		return errors.New(fmt.Sprintf("Fail to create aks, code:%d, err:%+v, output: %s", res, err, output))
	}

	// mv k8s config
	changeConfigPath := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s mkdir -p $HOME/.kube && sudo cp -rf /etc/kubernetes/admin.conf $HOME/.kube/config && sudo chown $(id -u):$(id -g) $HOME/.kube/config", master)
	res, err, output = util.ExecCMD(nil, nil, "bash", "-c", changeConfigPath)
	if err != nil {
		klog.Error(err)
		return err
	}
	if res != 0 {
		return errors.New(fmt.Sprintf("Fail to mv k8s config, code:%d, err:%+v, output: %s", res, err, output))
	}

	// install calico cni
	installCalico := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s kubectl create -f /root/yanyu/calico/tigera-operator.yaml && kubectl create -f /root/yanyu/calico/custom-resources.yaml", master)
	res, err, output = util.ExecCMD(nil, nil, "bash", "-c", installCalico)
	if err != nil {
		klog.Error(err)
		return err
	}
	if res != 0 {
		return errors.New(fmt.Sprintf("Fail to install cni, code:%d, err:%+v, output: %s", res, err, output))
	}

	return nil
}

func GetAks(master string) (error, string) {
	infoFile, err := os.Create(fmt.Sprintf("%s/%s.get.info.log", fmt.Sprintf("./logs"), master))
	if err != nil {
		klog.Errorf("create master %s info file failed", master)
		infoFile = nil
	}
	errFile, err := os.Create(fmt.Sprintf("%s/%s.get.err.log", fmt.Sprintf("./logs"), master))
	if err != nil {
		klog.Errorf("create master %s error file failed", master)
		errFile = nil
	}

	node := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s kubectl get node -owide | grep Ready |wc -l", master)
	klog.Infof(node)
	res, err, output := util.ExecCMD(infoFile, errFile, "bash", "-c", node)
	klog.Infof("node output: %s", output)
	if err != nil {
		klog.Error(err)
		return err, ""
	}
	if res != 0 {
		return errors.New(fmt.Sprintf("Fail to check k8s, code:%d, err:%+v, output: %s", res, err, output)), ""
	}

	kubeVersion := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s kubectl version | grep Server | grep GitVersion | awk '{print $5}'", master)
	klog.Infof(kubeVersion)
	res, err, output = util.ExecCMD(infoFile, errFile, "bash", "-c", kubeVersion)
	klog.Infof("kubeVersion output: %s", output)
	if err != nil {
		klog.Error(err)
		return err, ""
	}
	if res != 0 {
		return errors.New(fmt.Sprintf("Fail to check cni, code:%d, err:%+v, output: %s", res, err, output)), ""
	}

	return nil, output
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

	delete := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s kubeadm reset -f", master)
	klog.Infof(delete)
	res, err, output := util.ExecCMD(infoFile, errFile, "bash", "-c", delete)
	klog.Infof("delete output: %s", output)
	if err != nil {
		klog.Error(err)
		return err
	}
	if res != 0 {
		return errors.New(fmt.Sprintf("Fail to delete aks, code:%d, err:%+v, output: %s", res, err, output))
	}

	return nil
}

func UpdateAks(master string) error {
	action := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s kubeadm upgrade ", master)
	res, err, output := util.ExecCMD(nil, nil, "bash", "-c", action)
	if err != nil {
		klog.Error(err)
		return err
	}
	if res != 0 {
		return errors.New(fmt.Sprintf("Fail to update aks, code:%d, err:%+v, output: %s", res, err, output))
	}

	return nil
}

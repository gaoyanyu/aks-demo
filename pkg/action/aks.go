package action

import (
	"aks-demo/pkg/util"
	"errors"
	"fmt"
)

func CreateAks(master string) error {
	// init k8s cluster
	createAction := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s kubeadm init --kubernetes-version=v1.21.5 --image-repository=registry.sensetime.com/diamond --apiserver-advertise-address=%s --service-cidr=10.96.0.0/12 --pod-network-cidr=10.244.0.0/16 -v=10", master, master)
	res, err, output := util.ExecCMD(nil, nil, "bash", "-c", createAction)
	if err != nil {
		return err
	}
	if res != 0 {
		return errors.New(fmt.Sprintf("Fail to create aks, code:%d, err:%+v, output: %s", res, err, output))
	}

	// mv k8s config
	changeConfigPath := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s mkdir -p $HOME/.kube && sudo cp -rf /etc/kubernetes/admin.conf $HOME/.kube/config && sudo chown $(id -u):$(id -g) $HOME/.kube/config", master)
	res, err, output = util.ExecCMD(nil, nil, "bash", "-c", changeConfigPath)
	if err != nil {
		return err
	}
	if res != 0 {
		return errors.New(fmt.Sprintf("Fail to mv k8s config, code:%d, err:%+v, output: %s", res, err, output))
	}

	// install calico cni
	installCalico := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s kubectl create -f /root/yanyu/calico/tigera-operator.yaml && kubectl create -f /root/yanyu/calico/custom-resources.yaml", master)
	res, err, output = util.ExecCMD(nil, nil, "bash", "-c", installCalico)
	if err != nil {
		return err
	}
	if res != 0 {
		return errors.New(fmt.Sprintf("Fail to install cni, code:%d, err:%+v, output: %s", res, err, output))
	}

	return nil
}

func GetAks(master string) error {
	k8s := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s kubectl get node -owide | grep Ready |wc -l", master)
	res, err, output := util.ExecCMD(nil, nil, "bash", "-c", k8s)
	if err != nil {
		return err
	}
	if res != 0 {
		return errors.New(fmt.Sprintf("Fail to check k8s, code:%d, err:%+v, output: %s", res, err, output))
	}

	cni := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s kubectl get pod -n calico-system -owide | grep -v Running | grep -v NAME |wc -l", master)
	res, err, output = util.ExecCMD(nil, nil, "bash", "-c", cni)
	if err != nil {
		return err
	}
	if res != 0 {
		return errors.New(fmt.Sprintf("Fail to check cni, code:%d, err:%+v, output: %s", res, err, output))
	}

	return nil
}

func DeleteAks(master string) error {
	action := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s kubeadm reset -f", master)
	res, err, output := util.ExecCMD(nil, nil, "bash", "-c", action)
	if err != nil {
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
		return err
	}
	if res != 0 {
		return errors.New(fmt.Sprintf("Fail to update aks, code:%d, err:%+v, output: %s", res, err, output))
	}

	return nil
}

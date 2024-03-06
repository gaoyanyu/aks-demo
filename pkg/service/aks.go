package service

import (
	"aks-demo/model/request"
	"aks-demo/model/response"
	"aks-demo/pkg/action"
	"aks-demo/pkg/util"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"k8s.io/klog/v2"
)

func Version(c *gin.Context) {
	c.JSON(http.StatusOK, response.Result{Code: http.StatusOK, Message: "success", Data: nil})
}

func CreateAks(c *gin.Context) {
	var createInfo request.CreateBody
	if err := c.BindJSON(&createInfo); err != nil {
		c.JSON(http.StatusBadRequest, response.Result{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	cmd := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s ps aux| grep 'kubeadm init' | grep -v grep|wc -l", createInfo.Master)
	if exist, err := util.EnsureProcessExist(cmd); err == nil && exist {
		c.JSON(http.StatusConflict, response.Result{Code: http.StatusConflict, Message: "aks is initing"})
		return
	}

	err := action.CreateAks(createInfo.Master)
	if err != nil {
		klog.Error(err)
		c.JSON(http.StatusInternalServerError, response.Result{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.Result{Code: http.StatusOK, Message: "Success to create aks"})
}

func GetAks(c *gin.Context) {
	master := c.GetHeader("master")
	cmd := fmt.Sprintf("sshpass -p 235659YANyy@ ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@%s ps aux| grep 'kubeadm init' | grep -v grep|wc -l", master)
	if exist, err := util.EnsureProcessExist(cmd); err == nil && exist {
		c.JSON(http.StatusInternalServerError, response.Result{Code: http.StatusInternalServerError, Message: "aks is initing"})
		return
	}

	err, kubeVersion := action.GetAks(master)
	if err != nil {
		klog.Error(err)
		c.JSON(http.StatusInternalServerError, response.Result{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.Result{Code: http.StatusOK, Message: "Success to get aks", Data: []string{kubeVersion}})
}

func DeleteAks(c *gin.Context) {
	master := c.GetHeader("master")
	err := action.DeleteAks(master)
	if err != nil {
		klog.Error(err)
		c.JSON(http.StatusInternalServerError, response.Result{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.Result{Code: http.StatusOK, Message: "Success to delete aks"})
}

func UpdateAks(c *gin.Context) {
	var updateInfo request.UpdateBody
	if err := c.BindJSON(&updateInfo); err != nil {
		c.JSON(http.StatusBadRequest, response.Result{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	err := action.UpdateAks(updateInfo.Master)
	if err != nil {
		klog.Error(err)
		c.JSON(http.StatusInternalServerError, response.Result{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.Result{Code: http.StatusOK, Message: "Success to delete aks"})
}

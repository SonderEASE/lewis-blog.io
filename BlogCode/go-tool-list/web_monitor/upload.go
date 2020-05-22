package web_monitor

import (
	"io"
	"net/http"
	"os"
	"os/exec"

	"wx-gitlab.xunlei.cn/scdn/x/logger"
)

func handleUploadFile(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for {
		part, err := reader.NextPart()
		if err == io.EOF || part == nil {
			break
		}

		if part.FileName() == "" || opt.serviceName+".tar" != part.FileName() {
			logger.Error("[web_monitor] error file name %s, need %s", part.FileName(), opt.serviceName+".tar")
			return
		}
		logger.Info("[web_monitor] start recv file", part.FileName())
		path := "./" + part.FileName()
		_ = os.Remove(path)
		dst, err := os.Create(path)
		if err != nil {
			logger.Errorf("[web_monitor] %v", err)
			return
		}
		_, err = io.Copy(dst, part)
		if err != nil {
			logger.Errorf("[web_monitor] %v", err)
			return
		}
		err = dst.Close()
		if err != nil {
			logger.Errorf("[web_monitor] %v", err)
			return
		}

		err = os.RemoveAll("./bin")
		if err != nil {
			logger.Errorf("[web_monitor] %v", err)
			return
		}
		cmd := exec.Command("tar", "-xzvf", "./"+opt.serviceName+".tar")
		err = cmd.Run()
		if err != nil {
			logger.Errorf("[web_monitor] %v", err)
			return
		}
		err = os.Remove("./" + opt.serviceName)
		if err != nil {
			logger.Errorf("[web_monitor] %v", err)
			return
		}
		err = os.Rename("./bin/"+opt.serviceName, "./"+opt.serviceName)
		if err != nil {
			logger.Errorf("[web_monitor] %v", err)
			return
		}

		cmd = exec.Command("chmod", "a+rwx", "./"+opt.serviceName)
		err = cmd.Run()
		if err != nil {
			logger.Errorf("[web_monitor] %v", err)
			return
		}

		_ = os.RemoveAll("./bin")
		_ = os.RemoveAll("./" + opt.serviceName + ".tar")

		logger.Infof("[web_monitor] upload %s file success", opt.serviceName)
	}
}

package web_monitor

import (
	"encoding/base64"
	"strings"
	"time"

	"wx-gitlab.xunlei.cn/scdn/x/config_manager"
	"wx-gitlab.xunlei.cn/scdn/x/logger"
)

type AccountConfirm func(username, pwd string) bool
type Encryption func([]byte) []byte
type Decrypt func([]byte) ([]byte, error)

type webMonitorOptions struct {
	port           int
	startTime      time.Time
	serviceName    string
	version        string
	pProf          bool
	accountConfirm AccountConfirm
	encryption     Encryption
	decrypt        Decrypt
}

func newDefaultWebMonitorOptions() *webMonitorOptions {
	return &webMonitorOptions{
		startTime: time.Now(),
		encryption: func(src []byte) []byte {
			return []byte(base64.StdEncoding.EncodeToString(src))
		},
		decrypt: func(src []byte) ([]byte, error) {
			return base64.StdEncoding.DecodeString(string(src))
		},
	}
}

func (w *webMonitorOptions) Init(opts []WebMonitorOption) {
	if w == nil {
		return
	}

	for _, opt := range opts {
		opt(w)
	}
}

// WebMonitorOption web_monitor的初始化参数
type WebMonitorOption func(*webMonitorOptions)

// WithVersion 设置版本号
func WithVersion(version string) WebMonitorOption {
	return func(o *webMonitorOptions) {
		o.version = version
	}
}

// WithServiceName 设置服务名称
func WithServiceName(serviceName string) WebMonitorOption {
	return func(o *webMonitorOptions) {
		o.serviceName = serviceName
	}
}

// WithPProf 设置是否开启pprof
func WithPProf(startPProf bool) WebMonitorOption {
	return func(o *webMonitorOptions) {
		o.pProf = startPProf
	}
}

// WithAccount 设置用户认证机制
func WithAccount(accountConfirm AccountConfirm) WebMonitorOption {
	return func(o *webMonitorOptions) {
		o.accountConfirm = accountConfirm
	}
}

// WithNoAccount 设置关闭用户认证机制(默认开启)
func WithNoAccount() WebMonitorOption {
	return func(o *webMonitorOptions) {
		o.accountConfirm = func(_, _ string) bool { return true }
	}
}

// WithAccountConfig 设置config_manager的用户认证方式
// 控制台通用配置要设置成 "monitor_account":"admin,888;user,123", 这种形式
func WithAccountConfig(conf, key string) WebMonitorOption {
	return WithAccount(func(username, pwd string) bool {
		str, ok := config_manager.GetString(conf, key)
		if !ok {
			logger.Errorf("[web_monitor] undefined %s.%s from config_manager", conf, key)
			return false
		}
		for _, v := range strings.Split(str, ";") {
			if v == username+","+pwd {
				return true
			}
		}
		return false
	})
}

// WithEncryption 加密算法
// 默认base64
func WithEncryption(encryption Encryption) WebMonitorOption {
	return func(o *webMonitorOptions) {
		o.encryption = encryption
	}
}

// WithDecrypt 解密算法
// 默认base64
func WithDecrypt(decrypt Decrypt) WebMonitorOption {
	return func(o *webMonitorOptions) {
		o.decrypt = decrypt
	}
}

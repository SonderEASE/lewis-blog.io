package web_monitor

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"time"

	"wx-gitlab.xunlei.cn/scdn/x/logger"
)

//获取登陆界面字节流
func GetLoginByte() []byte {
	return []byte(`<html><head><meta charset="utf-8"></head><body>
<form style="text-align: center;margin-top: 150px" name="input" action="/user/login" method="POST">
				Username: <input type="text" name="user" style="font-size:14pt">
				Password: <input type="password" name="pwd" style="font-size:14pt">
				<input type="submit" value="Login" style="color:sienna;margin-left:20px;font-size:16pt">
		</form></body></html>`)
}

//获得请求头中的ip
func GetRemoteIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if opt.accountConfirm == nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`404`))
		return
	}
	path := r.URL.Path

	if r.Method == "GET" {
		_, _ = w.Write(GetLoginByte())
		return
	} else if r.Method != "POST" {
		return
	}

	err := r.ParseForm()
	if err != nil {
		logger.Infof("[web_monitor] Request ParseForm is error! path:%s err:%v ", path, err)
		return
	}
	user := r.Form.Get("user")
	pwd := r.Form.Get("pwd")
	if opt.accountConfirm(user, pwd) {
		sessionID := opt.encryption([]byte(user + "," + pwd))
		dd, _ := time.ParseDuration("24h")
		cookie := http.Cookie{Name: "token", Value: string(sessionID), Path: "/", Expires: time.Now().Add(30 * dd)}
		http.SetCookie(w, &cookie)
		u, _ := url.Parse(r.Header.Get("Referer"))
		redirect, err := base64.URLEncoding.DecodeString(u.Query().Get("r"))
		if err != nil || strings.TrimSpace(string(redirect)) == "" {
			redirect = []byte("/monitor")
		}
		http.Redirect(w, r, string(redirect), http.StatusFound)
		logger.Infof("[web_monitor] Login success user:%s password:%s path:%s IP:%s request head:%v ", user, pwd, path, GetRemoteIP(r), r.Header)
		return
	} else {
		write(w, []byte("Username or Password is error "))
		logger.Warnf("[web_monitor] Login Failed user:%s password:%s path:%s IP:%s request head:%v ", user, pwd, path, GetRemoteIP(r), r.Header)
		return
	}
}

func checkAccount(r *http.Request) bool {
	if opt.accountConfirm == nil {
		return false
	}
	c, err := r.Cookie("token")
	if err != nil {
		return false
	}

	deSrc, err := opt.decrypt([]byte(c.Value))
	account := strings.Split(string(deSrc), ",")
	if err != nil || len(account) != 2 || !opt.accountConfirm(account[0], account[1]) {
		logger.Warnf("[web_monitor] token does not exist or is wrong! path:%s IP:%s request head:%v ", r.URL.Path, GetRemoteIP(r), r.Header)
		return false
	}

	return true
}

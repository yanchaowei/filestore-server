package handler

import (
	dblayer "filestore-server/db"
	"filestore-server/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	pwd_salt = "*#890"
)

//SigninHandler： 登陆接口
func SigninHandler(w http.ResponseWriter, r *http.Request)  {

	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encPasswd := util.Sha1([]byte(password + pwd_salt))
	// 1.校验用户名以及密码
	pwdChecked := dblayer.UserSignin(username, encPasswd)
	if !pwdChecked {
		w.Write([]byte("FAILDED"))
		return
	}
	// 2.生成访问凭证（token）
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("FAILED"))
	}
	// 3、登陆成功后重定向到首页
	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	//http.Redirect(w, r, "http://" + r.Host + "/static/view/home.html", http.StatusFound)
	resp := util.RespMsg{
		Code: 0,
		Msg: "OK",
		Data: struct {
			Location string
			Username string
			Token string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token: token,
		},
	}
	w.Write(resp.JSONBytes())
}


// UserInfoHandler: 查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request)  {

	// 1、解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	//token := r.Form.Get("token")
	// 2、验证Token是否有效
	//isValidToken := IsTokenValid(token)
	//if !isValidToken {
	//	w.WriteHeader(http.StatusForbidden)
	//	return
	//}
	// 3、查询用户信息
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// 4、组装并且相应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg: "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

//
func IsTokenValid(token string) bool {
	// 1、判断token时效性，是否过期
	// 2、从数据库表tbl_token查询username对应的token信息
	// 3、对比两个token是否一致
	return true

}

func GenToken(username string) string {
	// md5(username+timestamp+token_salt+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

// SignupHandler: 处理用户注册请求
func SignupHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if len(username) < 3 || len(password) < 5 {
		w.Write([]byte("invalid parameter"))
		return
	}

	encPasswd := util.Sha1([]byte(password + pwd_salt))
	suc := dblayer.UserSignup(username, encPasswd)
	if suc {
		w.Write([]byte("SUCCESS"))
	}else {
		w.Write([]byte("FAILED"))
	}
}

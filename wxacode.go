package wxacode

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const baseUrl = "https://api.weixin.qq.com"

type Api struct {
	appid  string
	secret string
}

func New(appid, secret string) *Api {
	return &Api{
		appid:  appid,
		secret: secret,
	}
}

func (a *Api) Code2Session(code string) (*AuthInfo, error) {
	url := fmt.Sprintf(
		"%s/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		baseUrl,
		a.appid,
		a.secret,
		code,
	)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	_ = resp.Body.Close()

	result := &struct {
		SessionKey string `json:"session_key"`
		Unionid    string `json:"unionid"`
		Errmsg     string `json:"errmsg"`
		Openid     string `json:"openid"`
		Errcode    int32  `json:"errcode"`
	}{}

	err = json.Unmarshal(bs, result)
	if err != nil {
		return nil, err
	}

	if result.Errcode != 0 {
		dist := map[int32]string{
			40029: "js_code无效",
			45011: "API 调用太频繁，请稍候再试",
			40226: "高风险等级用户，小程序登录拦截",
			-1:    "系统繁忙，此时请开发者稍候再试",
		}

		msg, ok := dist[result.Errcode]
		if !ok {
			msg = result.Errmsg
		}

		return nil, errors.New(msg)
	}

	return &AuthInfo{SessionKey: result.SessionKey, Unionid: result.Unionid, Openid: result.Openid}, nil
}

func Test(a, b int64) int64 {
	return a + b + 12
}

package controller

import (
	"InterviewPush/dao/redis"
	"crypto/tls"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"time"
)

func GetAccessToken() (access_token string, err error) {
	var accessToken string
	expire, err := redis.Client.TTL("AccessToken").Result()
	if err != nil {
		zap.L().Error("判断token剩余生存时间失败", zap.Error(err))
	}
	if expire == -time.Second*2 {
		//申请新的accesstoken
		accessToken, _, err = GetAccessTokenDing()
		if err != nil {
			zap.L().Error("申请新的token失败", zap.Error(err))
			return
		}
		err = redis.Client.Set("AccessToken", accessToken, time.Second*7200).Err()
		if err != nil {
			zap.L().Error("重新设置token和token的过期时间失败", zap.Error(err))
			return
		}
		result, err := redis.Client.Get("AccessToken").Result()
		if err != nil {
			zap.L().Error("重新申请后，获取token失败", zap.Error(err))
		}
		access_token = result
		if access_token == "" {
			zap.L().Error("重新申请后，获取token失败")
		}
	} else {
		access_token, err = redis.Client.Get("AccessToken").Result()
	}
	// 如果err 是key不存在的话，应该重新申请一遍
	if err != nil {
		zap.L().Error("从redis从取access_token失败", zap.Error(err))
		return
	}
	return
}

// 定时调用获得access_token
func GetAccessTokenDing() (accessToken string, expireTime int64, err error) {
	//发消息
	var client *http.Client
	var request *http.Request
	var resp *http.Response
	URL := "https://api.dingtalk.com/v1.0/oauth2/accessToken"
	client = &http.Client{Transport: &http.Transport{ //对客户端进行一些配置
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}, Timeout: time.Duration(time.Second * 5)}
	//此处是post请求的请求体，我们先初始化一个对象
	b := struct {
		AppKey    string `json:"appKey"`
		AppSecret string `json:"appSecret"`
	}{
		AppKey:    "钉钉开放平台企业的key",
		AppSecret: "钉钉开放平台企业的AppSecret",
	}
	//然后把结构体对象序列化一下
	bodymarshal, err := json.Marshal(&b)
	if err != nil {
		return
	}
	//再处理一下
	reqBody := strings.NewReader(string(bodymarshal))
	//然后就可以放入具体的request中的
	request, _ = http.NewRequest(http.MethodPost, URL, reqBody)
	request.Header.Set("Content-Type", "application/json")
	resp, _ = client.Do(request)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body) //把请求到的body转化成byte[]
	h := struct {
		ExpireIn    int64  `json:"expireIn"`
		AccessToken string `json:"accessToken"`
	}{}
	//把请求到的结构反序列化到专门接受返回值的对象上面
	err = json.Unmarshal(body, &h)
	return h.AccessToken, h.ExpireIn, err
}

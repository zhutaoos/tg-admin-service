package service

import (
	"app/internal/config"
	"app/internal/model"
	"app/tools/conv"
	"app/tools/jwt"
	"app/tools/logger"
	"app/tools/resp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type TokenLogic struct {
}

var TokenLogicInstance *TokenLogic

func init() {
	TokenLogicInstance = &TokenLogic{}
}

func (tl *TokenLogic) GenerateJwt(uid uint, exTime int64) (string, *jwt.UserJwt) {
	j, userJwt := jwt.CreateJwt(uid, exTime)

	tokenModel := new(model.Token)
	tokenModel.UserId = userJwt.UserId
	tokenModel.Token = userJwt.Token
	tokenModel.ExpireTime = userJwt.ExpireTime
	tokenModel.DelToken() // 删除这个用户的 token

	tokenKey := GetTokenKey(userJwt.Token)
	uidTokenKey := GetUidToken(int(uid))

	// 删除旧的 key
	get := config.Redis.Get(context.Background(), uidTokenKey)
	if get.Val() != "" {
		config.Redis.Del(context.Background(), get.Val())
	}

	tokenModel.CreateToken()

	var m, err = json.Marshal(userJwt)
	_, err = config.Redis.Set(context.Background(), tokenKey, m, -1).Result()
	if err != nil {
		logger.Error("jwt 缓存失败", err)
		(&resp.JsonResp{Code: resp.ReFail, Msg: "jwt 缓存失败", Data: nil}).Response()
	}
	_, err = config.Redis.Set(context.Background(), uidTokenKey, tokenKey, -1).Result()
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "uidTokenKey 缓存失败", Data: nil}).Response()
	}
	if exTime > 0 {
		config.Redis.Expire(context.Background(), tokenKey, time.Duration(exTime)*time.Second)
		config.Redis.Expire(context.Background(), uidTokenKey, time.Duration(exTime)*time.Second)
	}
	return j, userJwt
}

func (tl *TokenLogic) CheckJwt(j string) (*jwt.UserJwt, error) {
	userJwt, err := jwt.ParseJwt(j)
	if err != nil {
		return nil, err
	}

	cacheKey := GetTokenKey(userJwt.Token)
	r, err := config.Redis.HGetAll(context.Background(), cacheKey).Result()
	if err != nil {
		return nil, err
	}
	i, ok := r["uid"]
	if !ok {
		return nil, errors.New("账户已经在其他终端上登录")
	}

	uid, _ := conv.Conv[uint](i)
	fmt.Println(r)
	fmt.Println(userJwt)
	fmt.Println(uid)
	if uid != userJwt.UserId || r["token"] != userJwt.Token {
		return nil, errors.New("账户已经在其他终端上登录[1]")
	}

	if userJwt.ExpireTime < time.Now().Unix() {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "token 过期", Data: nil}).Response()
	}
	return userJwt, nil
}

// GetTokenKey 存放 userJwt struct
func GetTokenKey(token string) string {
	return "token:" + token
}

func GetUidToken(uid int) string {
	return "uid-token:" + strconv.Itoa(uid)
}

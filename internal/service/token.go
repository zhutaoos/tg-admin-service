package service

import (
	"app/internal/model"
	"app/tools/conv"
	"app/tools/jwt"
	"app/tools/logger"
	"app/tools/resp"
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// TokenService Token服务接口
type TokenService interface {
	GenerateJwt(uid uint, exTime int64) (string, *jwt.UserJwt)
	CheckJwt(j string) (*jwt.UserJwt, error)
}

type TokenServiceImpl struct {
	redis *redis.Client
	db    *gorm.DB
}

// NewTokenLogic 创建TokenLogic实例
func NewTokenLogic(redis *redis.Client, db *gorm.DB) TokenService {
	return &TokenServiceImpl{
		redis: redis,
		db:    db,
	}
}

func (tl *TokenServiceImpl) GenerateJwt(uid uint, exTime int64) (string, *jwt.UserJwt) {
	j, userJwt := jwt.CreateJwt(uid, exTime)

	// 创建Token模型
	tokenModel := &model.Token{
		UserId:     userJwt.UserId,
		Token:      userJwt.Token,
		ExpireTime: userJwt.ExpireTime,
	}

	// 删除这个用户的旧token
	delToken(tokenModel, tl.db)

	tokenKey := GetTokenKey(userJwt.Token)
	uidTokenKey := GetUidToken(int(uid))

	// 删除旧的 Redis key
	get := tl.redis.Get(context.Background(), uidTokenKey)
	if get.Val() != "" {
		tl.redis.Del(context.Background(), get.Val())
	}

	// 创建新的Token记录
	err := createToken(tokenModel, tl.db)
	if err != nil {
		logger.Error("创建token记录失败", "error", err)
		(&resp.JsonResp{Code: resp.ReFail, Msg: "创建token记录失败", Data: nil}).Response()
	}

	// 缓存到Redis
	m, err := json.Marshal(userJwt)
	if err != nil {
		logger.Error("序列化JWT失败", "error", err)
		(&resp.JsonResp{Code: resp.ReFail, Msg: "序列化JWT失败", Data: nil}).Response()
	}

	_, err = tl.redis.Set(context.Background(), tokenKey, m, -1).Result()
	if err != nil {
		logger.Error("jwt 缓存失败", "error", err)
		(&resp.JsonResp{Code: resp.ReFail, Msg: "jwt 缓存失败", Data: nil}).Response()
	}

	_, err = tl.redis.Set(context.Background(), uidTokenKey, tokenKey, -1).Result()
	if err != nil {
		logger.Error("uidTokenKey 缓存失败", "error", err)
		(&resp.JsonResp{Code: resp.ReFail, Msg: "uidTokenKey 缓存失败", Data: nil}).Response()
	}

	// 设置过期时间
	if exTime > 0 {
		tl.redis.Expire(context.Background(), tokenKey, time.Duration(exTime)*time.Second)
		tl.redis.Expire(context.Background(), uidTokenKey, time.Duration(exTime)*time.Second)
	}

	return j, userJwt
}

// DelToken 删除Token记录
func delToken(token *model.Token, db *gorm.DB) error {
	where := make(map[string]interface{})

	if token.UserId > 0 {
		where["user_id"] = token.UserId
	}
	if token.Id > 0 {
		where["id"] = token.Id
	}
	if token.Token != "" {
		where["token"] = token.Token
	}

	return db.Where(where).Delete(&model.Token{}).Error
}

// CreateToken 创建Token记录
func createToken(token *model.Token, db *gorm.DB) error {
	return db.Create(token).Error
}

func (tl *TokenServiceImpl) CheckJwt(j string) (*jwt.UserJwt, error) {
	userJwt, err := jwt.ParseJwt(j)
	if err != nil {
		return nil, err
	}

	cacheKey := GetTokenKey(userJwt.Token)
	r, err := tl.redis.Get(context.Background(), cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("token 不存在")
		}
		return nil, err
	}

	var userJwtCache jwt.UserJwt
	err = json.Unmarshal([]byte(r), &userJwtCache)
	if err != nil {
		return nil, err
	}

	uid, _ := conv.Conv[uint](userJwtCache.UserId)
	if uid != userJwt.UserId || userJwtCache.Token != userJwt.Token {
		return nil, errors.New("账户已经在其他终端上登录")
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

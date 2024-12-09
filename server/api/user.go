package api

import (
	"log"
	"strconv"

	"crm/models"
	"crm/response"
	"crm/service"

	"github.com/gin-gonic/gin"
)

type UserApi struct {
	userService *service.UserService
}

func NewUserApi() *UserApi {
	userApi := UserApi{
		userService: service.NewUserService(),
	}
	return &userApi
}

// Register 用户注册
func (u *UserApi) Register(context *gin.Context) {
	var param models.UserCreateParam
	if err := context.ShouldBind(&param); err != nil {
		response.Result(response.ErrCodeParamInvalid, nil, context)
		log.Printf("[error]UserApi:Register:%s", err)
		return
	}
	errCode := u.userService.Register(&param)
	response.Result(errCode, nil, context)
}

// Login 用户登录
func (u *UserApi) Login(context *gin.Context) {
	var param models.UserLoginParam
	if err := context.ShouldBind(&param); err != nil {
		response.Result(response.ErrCodeParamInvalid, nil, context)
		return
	}
	userInfo, errCode := u.userService.Login(&param)
	if userInfo == nil {
		response.Result(errCode, nil, context)
		return
	}
	response.Result(errCode, userInfo, context)
}

// GetVerifyCode 获取验证码
func (u *UserApi) GetVerifyCode(context *gin.Context) {
	var param models.UserVerifyCodeParam
	if err := context.ShouldBind(&param); err != nil {
		response.Result(response.ErrCodeParamInvalid, nil, context)
		return
	}
	errCode := u.userService.GetVerifyCode(param.Email)
	response.Result(errCode, nil, context)
}

// ForgotPass 忘记密码
func (u *UserApi) ForgotPass(context *gin.Context) {
	var param models.UserPassParam
	if err := context.ShouldBind(&param); err != nil {
		response.Result(response.ErrCodeParamInvalid, nil, context)
		return
	}
	errCode := u.userService.ForgotPass(&param)
	response.Result(errCode, nil, context)
}

// Delete 注销账号
func (u *UserApi) Delete(context *gin.Context) {
	var param models.UserDeleteParam
	uid, _ := strconv.Atoi(context.Request.Header.Get("uid"))
	err := context.ShouldBind(&param)
	if uid <= 0 || err != nil {
		response.Result(response.ErrCodeParamInvalid, nil, context)
		return
	}
	param.Id = int64(uid)
	errCode := u.userService.Delete(param)
	response.Result(errCode, nil, context)
}

// GetInfo 获取用户信息
func (u *UserApi) GetInfo(context *gin.Context) {
	uid, _ := strconv.Atoi(context.Request.Header.Get("uid"))
	if uid <= 0 {
		response.Result(response.ErrCodeParamInvalid, nil, context)
		return
	}
	userInfo, errCode := u.userService.GetInfo(int64(uid))
	response.Result(errCode, userInfo, context)
}

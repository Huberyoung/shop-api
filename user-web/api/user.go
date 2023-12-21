package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"shop-api/user-web/forms"
	"shop-api/user-web/global"
	"shop-api/user-web/global/response"
	"shop-api/user-web/proto/user"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err == nil {
		return
	}

	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"msg": e.Message()})
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "内部错误"})
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "其他错误"})
		}
	}
}

func TranslateErr(err error) any {
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	if !ok {
		return err.Error()
	}
	return removeTopStruct(errs.Translate(global.Trans))
}

func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

func HandleValidatorError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"msg": TranslateErr(err)})
}

func getClient() user.UserClient {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", global.ServerSetting.RemoteConfig.Host, global.ServerSetting.RemoteConfig.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[GetUserList 连接【用户服务失败】", "msg", err.Error())
	}
	return user.NewUserClient(conn)
}

// GetUserList 获得用户列表
func GetUserList(ctx *gin.Context) {
	var param forms.GetUserListForm
	if err := ctx.ShouldBind(&param); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	client := getClient()
	rsp, err := client.GetUserList(context.Background(), &user.PageInfo{PageNum: param.PageNum, PageSize: param.PageSize})
	if err != nil {
		zap.S().Errorw("[GetUserList 查询【用户列表】失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]any, 0)
	for _, value := range rsp.Data {
		data := response.UserResponse{
			Id:       value.Id,
			Mobile:   value.Mobile,
			NickName: value.NikeName,
			Birthday: response.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			Gender:   int(value.Gender.Number()),
		}
		result = append(result, data)
	}
	ctx.JSON(http.StatusOK, result)
	return
}

// PasswordLogin 密码登录
func PasswordLogin(ctx *gin.Context) {
	var param forms.PasswordLoginForm
	if err := ctx.ShouldBind(&param); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	cbg := context.Background()

	client := getClient()
	rsp, err := client.GetUserByMobile(cbg, &user.MobileRequest{Mobile: param.Mobile})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, gin.H{"mobile": "用户不存在", "msg": err.Error()})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"mobile": "登录失败", "msg": err.Error()})
			}
			return
		}
		zap.S().Errorw("[GetUserByMobile 查询【用户信息】失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	if passRsp, err := client.CheckPassword(cbg, &user.PasswordCheckRequest{Password: param.Password, EncryptedPassword: rsp.Password}); err != nil {
		zap.S().Errorw("[CheckPassword 验证【密码信息】失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	} else {
		if !passRsp.Success {
			ctx.JSON(http.StatusForbidden, gin.H{"msg": "密码验证失败～"})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "登录成功～"})
	return
}

// CreateUser 创建用户
func CreateUser(ctx *gin.Context) {
	var param forms.CreateUserForm
	if err := ctx.ShouldBind(&param); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	birthday, err := time.ParseInLocation("2006-01-02 15:04:05", param.BirthDay, global.Location)
	if err != nil {
		zap.S().Errorw("[ParseInLocation 失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	client := getClient()
	crq := user.CreateUserRequest{
		Nickname: param.Nickname,
		Password: param.Password,
		Mobile:   param.Mobile,
		Gender:   user.CreateUserRequest_Gender(param.Gender).Enum(),
		BirthDay: uint64(birthday.Unix()),
	}
	createUser, err := client.CreateUser(context.Background(), &crq)
	if err != nil {
		zap.S().Errorw("[CreateUser【创建用户】失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"createUser": createUser})
	return
}

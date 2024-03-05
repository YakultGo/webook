package web

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/service"
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// UserHandler 与用户有关的路由
type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}
func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.GET("/profile", u.Profile)
	ug.POST("/signup", u.Signup)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
}

func (u *UserHandler) Signup(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}
	var req SignUpReq
	// Bind方法会根据Content-Type自动选择绑定方式
	// 解析错误会直接返回一个400的错误响应
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "邮箱格式错误")
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次密码不一致")
		return
	}

	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，且包含数字和字母")
		return
	}
	// 调用svc的方法
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserDuplicateEmail) {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	ctx.String(http.StatusOK, "注册成功")
	fmt.Printf("%v", req)
}
func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "账号/密码错误")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	// 设置session
	sess := sessions.Default(ctx)
	sess.Set("userId", user.Id)
	sess.Save()
	ctx.String(http.StatusOK, "登录成功")
}
func (u *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		Name        string `json:"name"`
		birthday    string `json:"birthday"`
		description string `json:"desciption"`
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 从session中获取userId
	sess := sessions.Default(ctx)
	userId := sess.Get("userId")
	birth, err := time.Parse("2006-01-02", req.birthday)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	err = u.svc.Edit(ctx, domain.User{
		Id:          userId.(int64),
		NickName:    req.Name,
		Birthday:    birth,
		Description: req.description,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "修改成功")
}
func (u *UserHandler) Profile(ctx *gin.Context) {
	// 从session中获取userId
	sess := sessions.Default(ctx)
	userId := sess.Get("userId")
	user, err := u.svc.Profile(ctx, userId.(int64))
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, user)
}

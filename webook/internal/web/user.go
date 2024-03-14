package web

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/service"
	myJwt "basic-go/webook/internal/web/jwt"
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

const (
	biz = "login"
)

var _ handler = (*UserHandler)(nil)

// UserHandler 与用户有关的路由
type UserHandler struct {
	svc         service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	codeSvc     service.CodeService
	cmd         redis.Cmdable
	myJwt.Handler
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService,
	cmd redis.Cmdable, jwtHandler myJwt.Handler) *UserHandler {
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
		codeSvc:     codeSvc,
		Handler:     jwtHandler,
		cmd:         cmd,
	}
}
func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.GET("/profile", u.ProfileJWT)
	ug.POST("/signup", u.Signup)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/edit", u.Edit)
	ug.POST("/logout", u.LogoutJWT)
	ug.POST("/login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("/login_sms", u.LoginSMS)
	ug.POST("/refresh_token", u.RefreshToken)
}
func (u *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := u.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码错误",
		})
		return
	}
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if err = u.SetLoginToken(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 4,
		Msg:  "验证码通过",
	})
}
func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "手机号有误",
		})
		return
	}
	err := u.codeSvc.Send(ctx, biz, req.Phone)
	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case errors.Is(err, service.ErrCodeSendTooMany):
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送太频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
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
	if errors.Is(err, service.ErrUserDuplicate) {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "注册成功")
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
	sess.Options(sessions.Options{
		MaxAge: 60 * 60,
	})
	sess.Save()
	ctx.String(http.StatusOK, "登录成功")
}
func (u *UserHandler) LoginJWT(ctx *gin.Context) {
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
	// 在这里使用JWT生成token
	if err = u.SetLoginToken(ctx, user.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "登录成功")
}

func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	sess.Save()
	sess.Clear()
	ctx.String(http.StatusOK, "退出成功")
}

func (u *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := u.Handler.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.String(http.StatusOK, "退出成功")
}
func (u *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		Name        string `json:"name"`
		Birthday    string `json:"birthday"`
		Description string `json:"description"`
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	birth, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	c, _ := ctx.Get("claims")
	claims, ok := c.(*myJwt.UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	err = u.svc.Edit(ctx, domain.User{
		Id:          claims.UserId,
		NickName:    req.Name,
		Birthday:    birth,
		Description: req.Description,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "修改成功")
}
func (u *UserHandler) Profile(ctx *gin.Context) {
	// JWT中获取ID
	c, _ := ctx.Get("claims")
	claims, ok := c.(*myJwt.UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	user, err := u.svc.Profile(ctx, claims.UserId)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, user)
}
func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	c, _ := ctx.Get("claims")
	claims, ok := c.(*myJwt.UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	user, err := u.svc.Profile(ctx, claims.UserId)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (u *UserHandler) RefreshToken(ctx *gin.Context) {
	// 拿出长token
	refreshToken := u.ExtractToken(ctx)
	var rc myJwt.RefreshClaims
	token, err := jwt.ParseWithClaims(refreshToken, &rc, func(token *jwt.Token) (interface{}, error) {
		return myJwt.RefreshTokenKey, nil
	})
	if err != nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = u.SetRefreshToken(ctx, rc.UserId, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "刷新成功",
	})
}

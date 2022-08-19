package middleware

import (
	"github.com/XM-GO/PandaKit/biz"
	"github.com/XM-GO/PandaKit/casbin"
	"github.com/XM-GO/PandaKit/ginx"
	"github.com/XM-GO/PandaKit/token"
	"github.com/dgrijalva/jwt-go"
	"pandax/pkg/global"
	"strconv"
)

func PermissionHandler(rc *ginx.ReqCtx) error {
	permission := rc.RequiredPermission
	// 如果需要的权限信息不为空，并且不需要token，则不返回错误，继续后续逻辑
	if permission != nil && !permission.NeedToken {
		return nil
	}
	tokenStr := rc.GinCtx.Request.Header.Get("X-TOKEN")
	// header不存在则从查询参数token中获取
	if tokenStr == "" {
		tokenStr = rc.GinCtx.Query("token")
	}
	if tokenStr == "" {
		return biz.PermissionErr
	}
	j := token.NewJWT("", []byte("PandaX"), jwt.SigningMethodHS256)
	loginAccount, err := j.ParseToken(tokenStr)
	if err != nil || loginAccount == nil {
		return biz.PermissionErr
	}
	rc.LoginAccount = loginAccount

	if !permission.NeedCasbin {
		return nil
	}
	ca := casbin.CasbinS{ModelPath: global.Conf.Casbin.ModelPath}
	e := ca.Casbin()
	// 判断策略中是否存在
	tenantId := strconv.Itoa(int(rc.LoginAccount.TenantId))
	success, err := e.Enforce(tenantId, loginAccount.RoleKey, rc.GinCtx.Request.URL.Path, rc.GinCtx.Request.Method)
	if !success {
		return biz.CasbinErr
	}

	return nil
}

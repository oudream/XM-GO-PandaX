package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/XM-GO/PandaKit/biz"
	"github.com/XM-GO/PandaKit/ginx"
	"github.com/XM-GO/PandaKit/logger"
	"github.com/XM-GO/PandaKit/utils"
	"github.com/sirupsen/logrus"
	"reflect"
	"runtime/debug"
)

func LogHandler(rc *ginx.ReqCtx) error {
	li := rc.LogInfo
	if li == nil {
		return nil
	}

	lfs := logrus.Fields{}
	if la := rc.LoginAccount; la != nil {
		lfs["uid"] = la.UserId
		lfs["uname"] = la.UserName
	}

	req := rc.GinCtx.Request
	lfs[req.Method] = req.URL.Path

	if err := rc.Err; err != nil {
		logger.Log.WithFields(lfs).Error(getErrMsg(rc, err))
		return nil
	}
	logger.Log.WithFields(lfs).Info(getLogMsg(rc))
	return nil
}

func getLogMsg(rc *ginx.ReqCtx) string {
	msg := rc.LogInfo.Description + fmt.Sprintf(" ->%dms", rc.Timed)
	if !utils.IsBlank(reflect.ValueOf(rc.ReqParam)) {
		rb, _ := json.Marshal(rc.ReqParam)
		msg = msg + fmt.Sprintf("\n--> %s", string(rb))
	}

	// 返回结果不为空，则记录返回结果
	if rc.LogInfo.LogResp && !utils.IsBlank(reflect.ValueOf(rc.ResData)) {
		respB, _ := json.Marshal(rc.ResData)
		msg = msg + fmt.Sprintf("\n<-- %s", string(respB))
	}
	return msg
}

func getErrMsg(rc *ginx.ReqCtx, err any) string {
	msg := rc.LogInfo.Description
	if !utils.IsBlank(reflect.ValueOf(rc.ReqParam)) {
		rb, _ := json.Marshal(rc.ReqParam)
		msg = msg + fmt.Sprintf("\n--> %s", string(rb))
	}

	var errMsg string
	switch t := err.(type) {
	case *biz.BizError:
		errMsg = fmt.Sprintf("\n<-e errCode: %d, errMsg: %s", t.Code(), t.Error())
	case error:
		errMsg = fmt.Sprintf("\n<-e errMsg: %s\n%s", t.Error(), string(debug.Stack()))
	case string:
		errMsg = fmt.Sprintf("\n<-e errMsg: %s\n%s", t, string(debug.Stack()))
	}
	return (msg + errMsg)
}

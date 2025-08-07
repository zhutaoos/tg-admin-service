package evaluate_controller

import (
	"app/internal/request"
	"app/internal/service"
	"app/tools/resp"

	"github.com/gin-gonic/gin"
)

type EvaluateController struct {
	evaluateService service.EvaluateService
}

func NewEvaluateController(evaluateService service.EvaluateService) *EvaluateController {
	return &EvaluateController{
		evaluateService: evaluateService,
	}
}

func (ec *EvaluateController) GetEvaluateList(ctx *gin.Context) {
	var req request.EvaluateSearchRequest
	if err := ctx.ShouldBind(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}

	evaluates, err := ec.evaluateService.GetList(req)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReError, Msg: "查询失败: " + err.Error()}).Response()
		return
	}
	resp.Ok(evaluates)
}

func (ec *EvaluateController) UpdateEvaluate(ctx *gin.Context) {
	var param request.EvaluateUpdateParam
	if err := ctx.ShouldBind(&param); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}
	if err := param.Validate(); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}

	if err := ec.evaluateService.UpdateEvaluate(param); err != nil {
		(&resp.JsonResp{Code: resp.ReError, Msg: "更新失败: " + err.Error()}).Response()
		return
	}

	resp.Success()
}

func (ec *EvaluateController) DeleteEvaluate(ctx *gin.Context) {
	var param request.EvaluateDeleteParam
	if err := ctx.ShouldBind(&param); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}
	if err := param.Validate(); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}

	if err := ec.evaluateService.DeleteEvaluate(param); err != nil {
		(&resp.JsonResp{Code: resp.ReError, Msg: "删除失败: " + err.Error()}).Response()
		return
	}

	resp.Success()
}

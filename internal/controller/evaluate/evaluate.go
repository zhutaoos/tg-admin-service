package evaluate_controller

import (
	"app/internal/converter"
	"app/internal/request"
	"app/internal/response"
	"app/internal/service"
	"app/tools/resp"

	"github.com/gin-gonic/gin"
)

// EvaluateController 评价控制器 - 展示Go中类似Java VO层的实现
type EvaluateController struct {
	evaluateService service.EvaluateService
}

// NewEvaluateController 创建评价控制器实例
func NewEvaluateController(evaluateService service.EvaluateService) *EvaluateController {
	return &EvaluateController{
		evaluateService: evaluateService,
	}
}

// GetEvaluateList 获取评价列表 - 使用VO层返回格式化数据
// GET /api/evaluate/list?group_id=xxx&page=1&limit=10
func (ec *EvaluateController) GetEvaluateList(ctx *gin.Context) {
	// 1. 接收请求参数（类似Java的DTO）
	var req struct {
		GroupID string `form:"group_id" binding:"required"`
		Page    int    `form:"page" default:"1"`
		Limit   int    `form:"limit" default:"10"`
	}

	if err := ctx.ShouldBind(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}

	// 2. 调用Service层获取数据
	evaluates, total, err := ec.evaluateService.GetList(req.GroupID, req.Page, req.Limit)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReError, Msg: "查询失败: " + err.Error()}).Response()
		return
	}

	// 3. 使用Converter将Model转换为Response（VO层）
	responseData := converter.ToEvaluateListResponse(evaluates, total, req.Page, req.Limit)

	// 4. 返回格式化后的响应
	(&resp.JsonResp{
		Code: resp.ReSuccess,
		Msg:  "获取评价列表成功",
		Data: responseData,
	}).Response()
}

// GetEvaluateDetail 获取评价详情 - 单个对象的VO转换
// GET /api/evaluate/:id
func (ec *EvaluateController) GetEvaluateDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "评价ID不能为空"}).Response()
		return
	}

	// 调用Service获取数据
	evaluate, err := ec.evaluateService.GetByID(id)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReError, Msg: "查询失败: " + err.Error()}).Response()
		return
	}

	// 转换为Response VO
	responseData := converter.ToEvaluateResponse(evaluate)

	(&resp.JsonResp{
		Code: resp.ReSuccess,
		Msg:  "获取评价详情成功",
		Data: responseData,
	}).Response()
}

// GetEvaluateStats 获取评价统计 - 自定义VO示例
// GET /api/evaluate/stats?group_id=xxx
func (ec *EvaluateController) GetEvaluateStats(ctx *gin.Context) {
	groupID := ctx.Query("group_id")
	if groupID == "" {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "群组ID不能为空"}).Response()
		return
	}

	// 调用Service获取统计数据
	stats, err := ec.evaluateService.GetStats(groupID)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReError, Msg: "获取统计失败: " + err.Error()}).Response()
		return
	}

	// 构建统计响应VO
	statsResponse := &response.EvaluateStatResponse{
		TotalCount:    stats.TotalCount,
		AvgScore:      stats.AvgScore,
		PendingCount:  stats.PendingCount,
		ApprovedCount: stats.ApprovedCount,
		RejectedCount: stats.RejectedCount,
	}

	(&resp.JsonResp{
		Code: resp.ReSuccess,
		Msg:  "获取统计数据成功",
		Data: statsResponse,
	}).Response()
}

// CreateEvaluate 创建评价 - 展示Request到Model的转换
// POST /api/evaluate
func (ec *EvaluateController) CreateEvaluate(ctx *gin.Context) {
	// 接收请求DTO
	var req request.CreateEvaluateRequest // 假设已定义
	if err := ctx.ShouldBindJSON(&req); err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Msg: "参数错误: " + err.Error()}).Response()
		return
	}

	// 转换请求为业务对象
	evaluate := converter.ToEvaluateModel(&req) // 需要实现这个转换方法

	// 调用Service处理业务逻辑
	err := ec.evaluateService.Create(evaluate)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReError, Msg: "创建失败: " + err.Error()}).Response()
		return
	}

	// 返回创建后的数据（转换为VO）
	responseData := converter.ToEvaluateResponse(evaluate)

	(&resp.JsonResp{
		Code: resp.ReSuccess,
		Msg:  "创建评价成功",
		Data: responseData,
	}).Response()
}

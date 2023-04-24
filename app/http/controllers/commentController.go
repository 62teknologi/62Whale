package controllers

import (
	"net/http"
	"whale/62teknologi-golang-utility/utils"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	SingularName  string
	PluralName    string
	SingularLabel string
	PluralLabel   string
	Table         string
}

func (ctrl *CommentController) Init(ctx *gin.Context) {
	ctrl.SingularName = utils.Pluralize.Singular(ctx.Param("table"))
	ctrl.PluralName = utils.Pluralize.Plural(ctx.Param("table"))
	ctrl.SingularLabel = ctrl.SingularName + " comment"
	ctrl.PluralLabel = ctrl.SingularName + " comments"
	ctrl.Table = ctrl.SingularName + "_comments"
}

func (ctrl *CommentController) Find(ctx *gin.Context) {
	ctrl.Init(ctx)

	value := map[string]any{}
	columns := []string{ctrl.Table + ".*"}
	transformer, _ := utils.JsonFileParser("setting/transformers/response/" + ctrl.Table + "/find.json")
	query := utils.DB.Table(ctrl.Table + "")

	utils.SetJoin(query, transformer, &columns)

	if err := query.Select(columns).Where(ctrl.Table+".id = ?", ctx.Param("id")).Take(&value).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", ctrl.SingularLabel+" not found", nil))
		return
	}

	utils.MapValuesShifter(transformer, value)
	utils.AttachJoin(transformer, value)

	if transformer["id"] != nil {
		transformer["childs"] = ctrl.FetchChild(transformer["id"].(int32))
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "find "+ctrl.SingularLabel+" success", transformer))
}

func (ctrl *CommentController) FindAll(ctx *gin.Context) {
	ctrl.Init(ctx)

	values := []map[string]any{}
	columns := []string{ctrl.Table + ".*"}
	transformer, _ := utils.JsonFileParser("setting/transformers/response/" + ctrl.Table + "/find.json")
	query := utils.DB.Table(ctrl.Table + "")

	utils.SetJoin(query, transformer, &columns)

	filterable, _ := utils.JsonFileParser("setting/filter/" + ctrl.Table + "/find.json")
	filter := utils.FilterByQueries(query, filterable, ctx)
	pagination := utils.SetPagination(query, ctx)

	if err := query.Select(columns).Find(&values).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", ctrl.PluralLabel+" not found", nil))
		return
	}

	customResponses := utils.MultiMapValuesShifter(values, transformer)

	if ctx.Query("include_childs") != "" {
		for _, value := range customResponses {
			if value["id"] != nil {
				value["childs"] = ctrl.FetchChild(value["id"].(int32))
			}
		}
	}

	ctx.JSON(http.StatusOK, utils.ResponseDataPaginate("success", "find "+ctrl.PluralLabel+" success", customResponses, pagination, filter))
}

func (ctrl *CommentController) Create(ctx *gin.Context) {
	ctrl.Init(ctx)

	transformer, _ := utils.JsonFileParser("setting/transformers/request/" + ctrl.Table + "/create.json")
	var input map[string]any

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	if validation, err := utils.Validate(input, transformer); err {
		ctx.JSON(http.StatusOK, utils.ResponseData("failed", "validation", validation.Errors))
		return
	}

	utils.MapValuesShifter(transformer, input)
	utils.MapNullValuesRemover(transformer)

	if err := utils.DB.Table(ctrl.Table + "").Create(&transformer).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "create "+ctrl.SingularLabel+" success", transformer))
}

func (ctrl *CommentController) Update(ctx *gin.Context) {
	ctrl.Init(ctx)

	transformer, _ := utils.JsonFileParser("setting/transformers/request/" + ctrl.Table + "/update.json")
	var input map[string]any

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	if validation, err := utils.Validate(input, transformer); err {
		ctx.JSON(http.StatusOK, utils.ResponseData("failed", "validation", validation.Errors))
		return
	}

	utils.MapValuesShifter(transformer, input)
	utils.MapNullValuesRemover(transformer)

	if err := utils.DB.Table(ctrl.Table).Where("id = ?", ctx.Param("id")).Updates(&transformer).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "update "+ctrl.SingularLabel+" success", transformer))
}

// todo : need to check constraint error
func (ctrl *CommentController) Delete(ctx *gin.Context) {
	ctrl.Init(ctx)

	if err := utils.DB.Table(ctrl.Table).Where("id = ?", ctx.Param("id")).Delete(map[string]any{}).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "delete "+ctrl.SingularLabel+" success", nil))
}

func (ctrl *CommentController) FetchChild(id int32) []map[string]any {
	var values []map[string]any

	if err := utils.DB.Table(ctrl.Table).Where("parent_id = ?", id).Find(&values).Error; err != nil {
		return values
	}

	transformer, _ := utils.JsonFileParser("setting/transformers/response/" + ctrl.Table + "/find.json")
	customResponses := utils.MultiMapValuesShifter(values, transformer)

	for _, value := range customResponses {
		value["childs"] = ctrl.FetchChild(value["id"].(int32))
	}

	return customResponses
}

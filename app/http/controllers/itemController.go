package controllers

import (
	"net/http"
	"whale/62teknologi-golang-utility/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type ItemController struct {
	SingularName  string
	PluralName    string
	SingularLabel string
	PluralLabel   string
	Table         string
}

func (ctrl *ItemController) Init(ctx *gin.Context) {
	ctrl.SingularName = utils.Pluralize.Singular(ctx.Param("table"))
	ctrl.PluralName = utils.Pluralize.Plural(ctx.Param("table"))
	ctrl.SingularLabel = ctrl.SingularName + " item"
	ctrl.PluralLabel = ctrl.SingularName + " items"
	ctrl.Table = ctrl.SingularName + "_items"
}

func (ctrl *ItemController) Find(ctx *gin.Context) {
	ctrl.Init(ctx)

	value := map[string]any{}
	columns := []string{ctrl.Table + ".*"}
	order := "id desc"
	transformer, _ := utils.JsonFileParser("setting/transformers/response/" + ctrl.Table + "/find.json")
	query := utils.DB.Table(ctrl.Table + "")

	utils.SetBelongsTo(query, transformer, &columns)
	delete(transformer, "filterable")

	if err := query.Select(columns).Order(order).Where(ctrl.Table+"."+"id = ?", ctx.Param("id")).Take(&value).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", ctrl.SingularLabel+" not found", nil))
		return
	}

	utils.MapValuesShifter(transformer, value)
	utils.AttachBelongsTo(transformer, value)

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "find "+ctrl.SingularLabel+" success", transformer))
}

func (ctrl *ItemController) FindAll(ctx *gin.Context) {
	ctrl.Init(ctx)

	values := []map[string]any{}
	columns := []string{ctrl.Table + ".*"}
	transformer, _ := utils.JsonFileParser("setting/transformers/response/" + ctrl.Table + "/find.json")
	query := utils.DB.Table(ctrl.Table + "")
	filter := utils.SetFilterByQuery(query, transformer, ctx)
	filter["search"] = utils.SetGlobalSearch(query, transformer, ctx)
	utils.SetBelongsTo(query, transformer, &columns)

	if err := query.Select(columns).Find(&values).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", ctrl.PluralLabel+" not found", nil))
		return
	}

	pagination := utils.SetPagination(query, ctx)
	customResponses := utils.MultiMapValuesShifter(transformer, values)

	ctx.JSON(http.StatusOK, utils.ResponseDataPaginate("success", "find "+ctrl.PluralLabel+" success", customResponses, pagination, filter))
}

func (ctrl *ItemController) Create(ctx *gin.Context) {
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

	var name string

	if transformer["name"] != nil {
		name, _ = transformer["name"].(string)
		transformer["slug"] = slug.Make(name)
	} else {
		transformer["slug"] = uuid.New()
	}

	if err := utils.DB.Table(ctrl.Table + "").Create(&transformer).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "create "+ctrl.SingularLabel+" success", transformer))
}

func (ctrl *ItemController) Update(ctx *gin.Context) {
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

	var name string

	if transformer["name"] != nil {
		name, _ = transformer["name"].(string)
		// not sure is it needed or not, may confusing if slug changes
		transformer["slug"] = slug.Make(name)
	}

	if err := utils.DB.Table(ctrl.Table+"").Where("id = ?", ctx.Param("id")).Updates(&transformer).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "update "+ctrl.SingularLabel+" success", transformer))
}

// todo : need to check constraint error
func (ctrl *ItemController) Delete(ctx *gin.Context) {
	ctrl.Init(ctx)

	if err := utils.DB.Table(ctrl.Table+"").Where("id = ?", ctx.Param("id")).Delete(map[string]any{}).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "delete "+ctrl.SingularLabel+" success", nil))
}

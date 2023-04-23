package controllers

import (
	"fmt"
	"net/http"
	"whale/62teknologi-golang-utility/utils"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

type GroupController struct{}

func (ctrl GroupController) Find(ctx *gin.Context) {
	value := map[string]any{}
	columns := []string{utils.SingularName + "_groups.*"}
	transformer, _ := utils.JsonFileParser("setting/transformers/response/" + utils.SingularName + "_groups/find.json")
	query := utils.DB.Table(utils.SingularName + "_groups")

	utils.SetJoin(query, transformer, &columns)

	if err := query.Select(columns).Where(utils.SingularName+"_groups."+"id = ?", ctx.Param("id")).Take(&value).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", utils.SingularName+" not found", nil))
		return
	}

	utils.MapValuesShifter(transformer, value)
	utils.AttachJoin(transformer, value)

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "find "+utils.SingularName+" success", transformer))
}

func (ctrl GroupController) FindAll(ctx *gin.Context) {
	values := []map[string]any{}
	columns := []string{utils.SingularName + "_groups.*"}
	transformer, _ := utils.JsonFileParser("setting/transformers/response/" + utils.SingularName + "_groups/find.json")
	query := utils.DB.Table(utils.SingularName + "_groups")

	utils.SetJoin(query, transformer, &columns)

	filterable, _ := utils.JsonFileParser("setting/filter/" + utils.SingularName + "_groups/find.json")
	filter := utils.FilterByQueries(query, filterable, ctx)
	pagination := utils.SetPagination(query, ctx)

	if err := query.Select(columns).Find(&values).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	customResponses := utils.MultiMapValuesShifter(values, transformer)

	ctx.JSON(http.StatusOK, utils.ResponseDataPaginate("success", "find "+utils.PluralName+" success", customResponses, pagination, filter))
}

func (ctrl GroupController) Create(ctx *gin.Context) {
	transformer, _ := utils.JsonFileParser("setting/transformers/request/" + utils.SingularName + "_groups/create.json")
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

	fmt.Println(input)
	fmt.Println(transformer)

	name, _ := transformer["name"].(string)
	transformer["slug"] = slug.Make(name)

	if err := utils.DB.Table(utils.SingularName + "_groups").Create(&transformer).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "create "+utils.SingularName+" success", transformer))
}

func (ctrl GroupController) Update(ctx *gin.Context) {
	transformer, _ := utils.JsonFileParser("setting/transformers/request/" + utils.SingularName + "_groups/update.json")
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

	name, _ := transformer["name"].(string)
	transformer["slug"] = slug.Make(name)

	if err := utils.DB.Table(utils.SingularName+"_groups").Where("id = ?", ctx.Param("id")).Updates(&transformer).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "update "+utils.SingularName+" success", transformer))
}

// todo : need to check constraint error
func (ctrl GroupController) Delete(ctx *gin.Context) {
	if err := utils.DB.Table(utils.SingularName+"_groups").Where("id = ?", ctx.Param("id")).Delete(map[string]any{}).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "delete "+utils.SingularName+" success", nil))
}

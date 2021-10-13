package controller

import (
	"gin_demo/common"
	"gin_demo/model"
	"gin_demo/response"
	"gin_demo/vo"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
)

type IPostController interface {
	RestController
	PageList(ctx *gin.Context)
}

type PostController struct {
	DB *gorm.DB
}

func NewPostController() IPostController {
	db := common.GetDB()
	db.AutoMigrate(model.Post{})
	return PostController{DB: db}
}

func (p PostController) Create(ctx *gin.Context) {
	var requestPost vo.CreatePostRequest

	if err := ctx.ShouldBind(&requestPost); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, "数据验证错误", nil)
		return
	}
	user, _ := ctx.Get("user")

	post := model.Post{
		UserId: user.(model.User).ID,
		CategoryId: requestPost.CategoryId,
		Title: requestPost.Title,
		HeadImg: requestPost.HeadImg,
		Content: requestPost.Content,
	}

	if err := p.DB.Create(&post).Error; err != nil {
		panic(err)
		return
	}
	response.Success(ctx, nil,"创建成功")
}

func (p PostController) Update(ctx *gin.Context) {
	var requestPost vo.CreatePostRequest

	if err := ctx.ShouldBind(&requestPost); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, "数据验证错误", nil)
		return
	}
	postId := ctx.Params.ByName("id")
	var post model.Post
	if p.DB.Where("id = ?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx,"文章不存在",nil)
		return
	}
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID

	if userId != post.UserId {
		response.Fail(ctx, "文章不属于你,请勿非法操作", nil)
		return
	}

	if err := p.DB.Model(&post).Update(requestPost).Error; err != nil {
		response.Fail(ctx, "更新失败", nil)
		return
	}
	response.Success(ctx, gin.H{"post":post},"更新成功")
}

func (p PostController) Show(ctx *gin.Context) {
	postId := ctx.Params.ByName("id")

	var post model.Post
	if p.DB.Preload("Category").Where("id = ?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx, "文章不存在", nil)
		return
	}
	response.Success(ctx, gin.H{"post":post},"成功")
}

func (p PostController) Delete(ctx *gin.Context) {
	postId := ctx.Params.ByName("id")

	var post model.Post
	if p.DB.Where("id = ?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx, "文章不存在", nil)
		return
	}
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID
	if userId != post.UserId {
		response.Fail(ctx, "文章不属于你,请勿非法操作", nil)
		return
	}
	p.DB.Delete(&post)
	response.Success(ctx, gin.H{"post":post},"成功")
}

func (p PostController) PageList(ctx *gin.Context) {
	pageNum,_ := strconv.Atoi(ctx.DefaultQuery("pageNum","1"))
	pageSize,_:= strconv.Atoi(ctx.DefaultQuery("pageSize","20"))

	var posts []model.Post
	p.DB.Order("created_at desc").Offset((pageNum - 1)* pageSize).Limit(pageSize).Find(&posts)

	var total int
	p.DB.Model(model.Post{}).Count(&total)
	response.Success(ctx, gin.H{"data":posts,"total":total},"成功")
}


package controllers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"meetingmanage/database"
	"meetingmanage/models"
	"meetingmanage/utils"

	"gorm.io/gorm"
)

type CarouselController struct{}

func (cc *CarouselController) GetAll(w http.ResponseWriter, r *http.Request) {
	var carousels []models.Carousel
	if err := database.DB.Order("`order` ASC, created_at DESC").
		Find(&carousels).Error; err != nil {
		utils.InternalServerErrorResponse(w, "获取轮播图列表失败")
		return
	}

	utils.SuccessResponse(w, "获取成功", carousels)
}

func (cc *CarouselController) GetActive(w http.ResponseWriter, r *http.Request) {
	var carousels []models.Carousel
	if err := database.DB.Where("status = ?", "active").
		Order("`order` ASC, created_at DESC").
		Find(&carousels).Error; err != nil {
		utils.InternalServerErrorResponse(w, "获取轮播图列表失败")
		return
	}

	utils.SuccessResponse(w, "获取成功", carousels)
}

func (cc *CarouselController) GetByID(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取轮播图ID
	carouselIDStr := r.URL.Query().Get("id")
	if carouselIDStr == "" {
		utils.BadRequestResponse(w, "缺少轮播图ID参数")
		return
	}

	carouselID, err := strconv.ParseUint(carouselIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "轮播图ID格式无效")
		return
	}

	var carousel models.Carousel
	if err := database.DB.First(&carousel, uint(carouselID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "轮播图不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取轮播图失败")
		}
		return
	}

	utils.SuccessResponse(w, "获取成功", carousel)
}

func (cc *CarouselController) Create(w http.ResponseWriter, r *http.Request) {
	// 解析表单数据
	err := r.ParseMultipartForm(10 << 20) // 10MB限制
	if err != nil {
		utils.BadRequestResponse(w, "表单解析失败")
		return
	}

	// 获取标题
	title := r.FormValue("title")
	if title == "" {
		utils.BadRequestResponse(w, "缺少标题参数")
		return
	}

	// 获取链接URL
	linkURL := r.FormValue("link_url")

	// 获取排序
	orderStr := r.FormValue("order")
	order := 0
	if orderStr != "" {
		order, _ = strconv.Atoi(orderStr)
	}

	// 获取状态
	status := r.FormValue("status")
	if status == "" {
		status = "active"
	}

	// 获取图片文件
	file, handler, err := r.FormFile("image")
	if err != nil {
		utils.BadRequestResponse(w, "获取图片失败")
		return
	}
	defer file.Close()

	// 创建上传目录
	uploadDir := "./uploads/carousels"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	// 生成文件名
	ext := filepath.Ext(handler.Filename)
	fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
	filePath := filepath.Join(uploadDir, fileName)

	// 保存文件
	dst, err := os.Create(filePath)
	if err != nil {
		utils.InternalServerErrorResponse(w, "创建文件失败")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		utils.InternalServerErrorResponse(w, "保存文件失败")
		return
	}

	// 创建轮播图记录
	carousel := models.Carousel{
		Title:    title,
		ImageURL: filePath,
		LinkURL:  linkURL,
		Order:    order,
		Status:   status,
	}

	if err := database.DB.Create(&carousel).Error; err != nil {
		// 删除已上传的文件
		os.Remove(filePath)
		utils.InternalServerErrorResponse(w, "创建轮播图失败")
		return
	}

	utils.SuccessResponse(w, "创建成功", carousel)
}

func (cc *CarouselController) Update(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取轮播图ID
	carouselIDStr := r.URL.Query().Get("id")
	if carouselIDStr == "" {
		utils.BadRequestResponse(w, "缺少轮播图ID参数")
		return
	}

	carouselID, err := strconv.ParseUint(carouselIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "轮播图ID格式无效")
		return
	}

	// 解析表单数据
	err = r.ParseMultipartForm(10 << 20) // 10MB限制
	if err != nil {
		utils.BadRequestResponse(w, "表单解析失败")
		return
	}

	var carousel models.Carousel
	if err := database.DB.First(&carousel, uint(carouselID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "轮播图不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取轮播图失败")
		}
		return
	}

	// 获取标题
	title := r.FormValue("title")
	if title != "" {
		carousel.Title = title
	}

	// 获取链接URL
	linkURL := r.FormValue("link_url")
	if linkURL != "" {
		carousel.LinkURL = linkURL
	}

	// 获取排序
	orderStr := r.FormValue("order")
	if orderStr != "" {
		order, _ := strconv.Atoi(orderStr)
		carousel.Order = order
	}

	// 获取状态
	status := r.FormValue("status")
	if status != "" {
		carousel.Status = status
	}

	// 获取图片文件（可选）
	file, handler, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		// 删除旧图片
		oldImagePath := carousel.ImageURL
		if _, err := os.Stat(oldImagePath); err == nil {
			os.Remove(oldImagePath)
		}

		// 创建上传目录
		uploadDir := "./uploads/carousels"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.MkdirAll(uploadDir, 0755)
		}

		// 生成文件名
		ext := filepath.Ext(handler.Filename)
		fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
		filePath := filepath.Join(uploadDir, fileName)

		// 保存文件
		dst, err := os.Create(filePath)
		if err != nil {
			utils.InternalServerErrorResponse(w, "创建文件失败")
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			utils.InternalServerErrorResponse(w, "保存文件失败")
			return
		}

		carousel.ImageURL = filePath
	}

	if err := database.DB.Save(&carousel).Error; err != nil {
		utils.InternalServerErrorResponse(w, "更新轮播图失败")
		return
	}

	utils.SuccessResponse(w, "更新成功", carousel)
}

func (cc *CarouselController) Delete(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取轮播图ID
	carouselIDStr := r.URL.Query().Get("id")
	if carouselIDStr == "" {
		utils.BadRequestResponse(w, "缺少轮播图ID参数")
		return
	}

	carouselID, err := strconv.ParseUint(carouselIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "轮播图ID格式无效")
		return
	}

	var carousel models.Carousel
	if err := database.DB.First(&carousel, uint(carouselID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "轮播图不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取轮播图失败")
		}
		return
	}

	// 获取图片路径
	imagePath := carousel.ImageURL

	// 删除数据库记录
	if err := database.DB.Delete(&carousel).Error; err != nil {
		utils.InternalServerErrorResponse(w, "删除轮播图失败")
		return
	}

	// 删除图片文件
	if _, err := os.Stat(imagePath); err == nil {
		os.Remove(imagePath)
	}

	utils.SuccessResponse(w, "删除成功", nil)
}

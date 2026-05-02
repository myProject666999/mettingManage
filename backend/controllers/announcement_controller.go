package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"meetingmanage/database"
	"meetingmanage/middleware"
	"meetingmanage/models"
	"meetingmanage/utils"

	"gorm.io/gorm"
)

type AnnouncementController struct{}

func (ac *AnnouncementController) GetAll(w http.ResponseWriter, r *http.Request) {
	var announcements []models.Announcement
	if err := database.DB.Preload("Author").
		Order("is_top DESC, created_at DESC").
		Find(&announcements).Error; err != nil {
		utils.InternalServerErrorResponse(w, "获取公告列表失败")
		return
	}

	utils.SuccessResponse(w, "获取成功", announcements)
}

func (ac *AnnouncementController) GetPublished(w http.ResponseWriter, r *http.Request) {
	var announcements []models.Announcement
	if err := database.DB.Where("status = ?", "published").
		Preload("Author").
		Order("is_top DESC, created_at DESC").
		Find(&announcements).Error; err != nil {
		utils.InternalServerErrorResponse(w, "获取公告列表失败")
		return
	}

	utils.SuccessResponse(w, "获取成功", announcements)
}

func (ac *AnnouncementController) GetByID(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取公告ID
	announcementIDStr := r.URL.Query().Get("id")
	if announcementIDStr == "" {
		utils.BadRequestResponse(w, "缺少公告ID参数")
		return
	}

	announcementID, err := strconv.ParseUint(announcementIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "公告ID格式无效")
		return
	}

	var announcement models.Announcement
	if err := database.DB.Preload("Author").
		First(&announcement, uint(announcementID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "公告不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取公告失败")
		}
		return
	}

	utils.SuccessResponse(w, "获取成功", announcement)
}

func (ac *AnnouncementController) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	if !ok {
		utils.UnauthorizedResponse(w, "无法获取用户信息")
		return
	}

	var announcement models.Announcement
	if err := json.NewDecoder(r.Body).Decode(&announcement); err != nil {
		utils.BadRequestResponse(w, "请求体解析失败")
		return
	}

	// 设置作者ID
	announcement.AuthorID = claims.UserID

	// 创建公告
	if err := database.DB.Create(&announcement).Error; err != nil {
		utils.InternalServerErrorResponse(w, "创建公告失败")
		return
	}

	// 加载关联信息
	database.DB.Preload("Author").First(&announcement, announcement.ID)

	utils.SuccessResponse(w, "创建成功", announcement)
}

func (ac *AnnouncementController) Update(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取公告ID
	announcementIDStr := r.URL.Query().Get("id")
	if announcementIDStr == "" {
		utils.BadRequestResponse(w, "缺少公告ID参数")
		return
	}

	announcementID, err := strconv.ParseUint(announcementIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "公告ID格式无效")
		return
	}

	var updateData struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		IsTop   bool   `json:"is_top"`
		Status  string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		utils.BadRequestResponse(w, "请求体解析失败")
		return
	}

	var announcement models.Announcement
	if err := database.DB.First(&announcement, uint(announcementID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "公告不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取公告失败")
		}
		return
	}

	if updateData.Title != "" {
		announcement.Title = updateData.Title
	}

	if updateData.Content != "" {
		announcement.Content = updateData.Content
	}

	announcement.IsTop = updateData.IsTop

	if updateData.Status != "" {
		announcement.Status = updateData.Status
	}

	if err := database.DB.Save(&announcement).Error; err != nil {
		utils.InternalServerErrorResponse(w, "更新公告失败")
		return
	}

	// 加载关联信息
	database.DB.Preload("Author").First(&announcement, announcement.ID)

	utils.SuccessResponse(w, "更新成功", announcement)
}

func (ac *AnnouncementController) Delete(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取公告ID
	announcementIDStr := r.URL.Query().Get("id")
	if announcementIDStr == "" {
		utils.BadRequestResponse(w, "缺少公告ID参数")
		return
	}

	announcementID, err := strconv.ParseUint(announcementIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "公告ID格式无效")
		return
	}

	var announcement models.Announcement
	if err := database.DB.First(&announcement, uint(announcementID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "公告不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取公告失败")
		}
		return
	}

	if err := database.DB.Delete(&announcement).Error; err != nil {
		utils.InternalServerErrorResponse(w, "删除公告失败")
		return
	}

	utils.SuccessResponse(w, "删除成功", nil)
}

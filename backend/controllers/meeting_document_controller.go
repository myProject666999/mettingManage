package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"meetingmanage/database"
	"meetingmanage/middleware"
	"meetingmanage/models"
	"meetingmanage/utils"

	"gorm.io/gorm"
)

type MeetingDocumentController struct{}

func (mdc *MeetingDocumentController) GetAll(w http.ResponseWriter, r *http.Request) {
	var documents []models.MeetingDocument
	if err := database.DB.Preload("Meeting").Preload("Meeting.Room").Preload("Meeting.Organizer").
		Find(&documents).Error; err != nil {
		utils.InternalServerErrorResponse(w, "获取会议资料列表失败")
		return
	}

	utils.SuccessResponse(w, "获取成功", documents)
}

func (mdc *MeetingDocumentController) GetByMeeting(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取会议ID
	meetingIDStr := r.URL.Query().Get("meeting_id")
	if meetingIDStr == "" {
		utils.BadRequestResponse(w, "缺少会议ID参数")
		return
	}

	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "会议ID格式无效")
		return
	}

	var documents []models.MeetingDocument
	if err := database.DB.Where("meeting_id = ?", uint(meetingID)).
		Preload("Meeting").Preload("Meeting.Room").Preload("Meeting.Organizer").
		Find(&documents).Error; err != nil {
		utils.InternalServerErrorResponse(w, "获取会议资料列表失败")
		return
	}

	utils.SuccessResponse(w, "获取成功", documents)
}

func (mdc *MeetingDocumentController) GetByID(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取会议资料ID
	documentIDStr := r.URL.Query().Get("id")
	if documentIDStr == "" {
		utils.BadRequestResponse(w, "缺少会议资料ID参数")
		return
	}

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "会议资料ID格式无效")
		return
	}

	var document models.MeetingDocument
	if err := database.DB.Preload("Meeting").Preload("Meeting.Room").Preload("Meeting.Organizer").
		First(&document, uint(documentID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "会议资料不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取会议资料失败")
		}
		return
	}

	utils.SuccessResponse(w, "获取成功", document)
}

func (mdc *MeetingDocumentController) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	if !ok {
		utils.UnauthorizedResponse(w, "无法获取用户信息")
		return
	}

	// 解析表单数据
	err := r.ParseMultipartForm(10 << 20) // 10MB限制
	if err != nil {
		utils.BadRequestResponse(w, "表单解析失败")
		return
	}

	// 获取会议ID
	meetingIDStr := r.FormValue("meeting_id")
	if meetingIDStr == "" {
		utils.BadRequestResponse(w, "缺少会议ID参数")
		return
	}

	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "会议ID格式无效")
		return
	}

	// 检查会议是否存在
	var meeting models.Meeting
	if err := database.DB.First(&meeting, uint(meetingID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "会议不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取会议信息失败")
		}
		return
	}

	// 获取文件
	file, handler, err := r.FormFile("file")
	if err != nil {
		utils.BadRequestResponse(w, "获取文件失败")
		return
	}
	defer file.Close()

	// 获取标题和描述
	title := r.FormValue("title")
	if title == "" {
		title = handler.Filename
	}

	description := r.FormValue("description")

	// 创建上传目录
	uploadDir := "./uploads/documents"
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

	// 获取文件信息
	fileInfo, err := dst.Stat()
	if err != nil {
		utils.InternalServerErrorResponse(w, "获取文件信息失败")
		return
	}

	// 创建会议资料记录
	document := models.MeetingDocument{
		MeetingID:     uint(meetingID),
		Title:         title,
		Description:   description,
		FileType:      strings.TrimPrefix(ext, "."),
		FileSize:      fileInfo.Size(),
		FilePath:      filePath,
		UploaderID:    claims.UserID,
		DownloadCount: 0,
	}

	if err := database.DB.Create(&document).Error; err != nil {
		// 删除已上传的文件
		os.Remove(filePath)
		utils.InternalServerErrorResponse(w, "创建会议资料失败")
		return
	}

	// 加载关联信息
	database.DB.Preload("Meeting").Preload("Meeting.Room").Preload("Meeting.Organizer").
		First(&document, document.ID)

	utils.SuccessResponse(w, "创建成功", document)
}

func (mdc *MeetingDocumentController) Update(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取会议资料ID
	documentIDStr := r.URL.Query().Get("id")
	if documentIDStr == "" {
		utils.BadRequestResponse(w, "缺少会议资料ID参数")
		return
	}

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "会议资料ID格式无效")
		return
	}

	var updateData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		utils.BadRequestResponse(w, "请求体解析失败")
		return
	}

	var document models.MeetingDocument
	if err := database.DB.First(&document, uint(documentID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "会议资料不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取会议资料失败")
		}
		return
	}

	if updateData.Title != "" {
		document.Title = updateData.Title
	}

	if updateData.Description != "" {
		document.Description = updateData.Description
	}

	if err := database.DB.Save(&document).Error; err != nil {
		utils.InternalServerErrorResponse(w, "更新会议资料失败")
		return
	}

	// 加载关联信息
	database.DB.Preload("Meeting").Preload("Meeting.Room").Preload("Meeting.Organizer").
		First(&document, document.ID)

	utils.SuccessResponse(w, "更新成功", document)
}

func (mdc *MeetingDocumentController) Delete(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取会议资料ID
	documentIDStr := r.URL.Query().Get("id")
	if documentIDStr == "" {
		utils.BadRequestResponse(w, "缺少会议资料ID参数")
		return
	}

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "会议资料ID格式无效")
		return
	}

	var document models.MeetingDocument
	if err := database.DB.First(&document, uint(documentID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "会议资料不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取会议资料失败")
		}
		return
	}

	// 获取文件路径
	filePath := document.FilePath

	// 删除数据库记录
	if err := database.DB.Delete(&document).Error; err != nil {
		utils.InternalServerErrorResponse(w, "删除会议资料失败")
		return
	}

	// 删除文件
	if _, err := os.Stat(filePath); err == nil {
		os.Remove(filePath)
	}

	utils.SuccessResponse(w, "删除成功", nil)
}

func (mdc *MeetingDocumentController) Download(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取会议资料ID
	documentIDStr := r.URL.Query().Get("id")
	if documentIDStr == "" {
		utils.BadRequestResponse(w, "缺少会议资料ID参数")
		return
	}

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "会议资料ID格式无效")
		return
	}

	var document models.MeetingDocument
	if err := database.DB.First(&document, uint(documentID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "会议资料不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取会议资料失败")
		}
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(document.FilePath); os.IsNotExist(err) {
		utils.NotFoundResponse(w, "文件不存在")
		return
	}

	// 更新下载次数
	document.DownloadCount++
	database.DB.Save(&document)

	// 设置下载头
	w.Header().Set("Content-Disposition", "attachment; filename="+document.Title+"."+document.FileType)
	w.Header().Set("Content-Type", "application/octet-stream")

	// 读取文件并发送
	file, err := os.Open(document.FilePath)
	if err != nil {
		utils.InternalServerErrorResponse(w, "打开文件失败")
		return
	}
	defer file.Close()

	io.Copy(w, file)
}

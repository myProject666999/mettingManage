package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"meetingmanage/database"
	"meetingmanage/models"
	"meetingmanage/utils"

	"gorm.io/gorm"
)

type MeetingRoomController struct{}

func (mrc *MeetingRoomController) GetAll(w http.ResponseWriter, r *http.Request) {
	var meetingRooms []models.MeetingRoom
	if err := database.DB.Find(&meetingRooms).Error; err != nil {
		utils.InternalServerErrorResponse(w, "获取会议室列表失败")
		return
	}

	utils.SuccessResponse(w, "获取成功", meetingRooms)
}

func (mrc *MeetingRoomController) GetAvailable(w http.ResponseWriter, r *http.Request) {
	var meetingRooms []models.MeetingRoom
	if err := database.DB.Where("status = ?", "available").Find(&meetingRooms).Error; err != nil {
		utils.InternalServerErrorResponse(w, "获取可用会议室列表失败")
		return
	}

	utils.SuccessResponse(w, "获取成功", meetingRooms)
}

func (mrc *MeetingRoomController) GetByID(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取会议室ID
	roomIDStr := r.URL.Query().Get("id")
	if roomIDStr == "" {
		utils.BadRequestResponse(w, "缺少会议室ID参数")
		return
	}

	roomID, err := strconv.ParseUint(roomIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "会议室ID格式无效")
		return
	}

	var meetingRoom models.MeetingRoom
	if err := database.DB.First(&meetingRoom, uint(roomID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "会议室不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取会议室信息失败")
		}
		return
	}

	utils.SuccessResponse(w, "获取成功", meetingRoom)
}

func (mrc *MeetingRoomController) Create(w http.ResponseWriter, r *http.Request) {
	var meetingRoom models.MeetingRoom
	if err := json.NewDecoder(r.Body).Decode(&meetingRoom); err != nil {
		utils.BadRequestResponse(w, "请求体解析失败")
		return
	}

	// 创建会议室
	if err := database.DB.Create(&meetingRoom).Error; err != nil {
		utils.InternalServerErrorResponse(w, "会议室创建失败")
		return
	}

	utils.SuccessResponse(w, "创建成功", meetingRoom)
}

func (mrc *MeetingRoomController) Update(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取会议室ID
	roomIDStr := r.URL.Query().Get("id")
	if roomIDStr == "" {
		utils.BadRequestResponse(w, "缺少会议室ID参数")
		return
	}

	roomID, err := strconv.ParseUint(roomIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "会议室ID格式无效")
		return
	}

	var updateData struct {
		Name        string `json:"name"`
		Location    string `json:"location"`
		Capacity    int    `json:"capacity"`
		Description string `json:"description"`
		Equipment   string `json:"equipment"`
		Status      string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		utils.BadRequestResponse(w, "请求体解析失败")
		return
	}

	var meetingRoom models.MeetingRoom
	if err := database.DB.First(&meetingRoom, uint(roomID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "会议室不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取会议室信息失败")
		}
		return
	}

	if updateData.Name != "" {
		meetingRoom.Name = updateData.Name
	}

	if updateData.Location != "" {
		meetingRoom.Location = updateData.Location
	}

	if updateData.Capacity > 0 {
		meetingRoom.Capacity = updateData.Capacity
	}

	if updateData.Description != "" {
		meetingRoom.Description = updateData.Description
	}

	if updateData.Equipment != "" {
		meetingRoom.Equipment = updateData.Equipment
	}

	if updateData.Status != "" {
		meetingRoom.Status = updateData.Status
	}

	if err := database.DB.Save(&meetingRoom).Error; err != nil {
		utils.InternalServerErrorResponse(w, "更新会议室信息失败")
		return
	}

	utils.SuccessResponse(w, "更新成功", meetingRoom)
}

func (mrc *MeetingRoomController) Delete(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取会议室ID
	roomIDStr := r.URL.Query().Get("id")
	if roomIDStr == "" {
		utils.BadRequestResponse(w, "缺少会议室ID参数")
		return
	}

	roomID, err := strconv.ParseUint(roomIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "会议室ID格式无效")
		return
	}

	var meetingRoom models.MeetingRoom
	if err := database.DB.First(&meetingRoom, uint(roomID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "会议室不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取会议室信息失败")
		}
		return
	}

	if err := database.DB.Delete(&meetingRoom).Error; err != nil {
		utils.InternalServerErrorResponse(w, "删除会议室失败")
		return
	}

	utils.SuccessResponse(w, "删除成功", nil)
}

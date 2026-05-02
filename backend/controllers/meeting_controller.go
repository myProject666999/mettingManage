package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"meetingmanage/database"
	"meetingmanage/middleware"
	"meetingmanage/models"
	"meetingmanage/utils"

	"gorm.io/gorm"
)

type MeetingController struct{}

func (mc *MeetingController) GetAll(w http.ResponseWriter, r *http.Request) {
	var meetings []models.Meeting
	if err := database.DB.Preload("Room").Preload("Organizer").Find(&meetings).Error; err != nil {
		utils.InternalServerErrorResponse(w, "获取会议列表失败")
		return
	}

	utils.SuccessResponse(w, "获取成功", meetings)
}

func (mc *MeetingController) GetUserMeetings(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	if !ok {
		utils.UnauthorizedResponse(w, "无法获取用户信息")
		return
	}

	// 获取用户组织的会议
	var organizedMeetings []models.Meeting
	if err := database.DB.Where("organizer_id = ?", claims.UserID).
		Preload("Room").Preload("Organizer").
		Find(&organizedMeetings).Error; err != nil {
		utils.InternalServerErrorResponse(w, "获取会议列表失败")
		return
	}

	// 获取用户参加的会议
	var participatedMeetings []models.Meeting
	var participants []models.Participant
	if err := database.DB.Where("user_id = ?", claims.UserID).
		Preload("Meeting").Preload("Meeting.Room").Preload("Meeting.Organizer").
		Find(&participants).Error; err != nil {
		utils.InternalServerErrorResponse(w, "获取参会会议列表失败")
		return
	}

	for _, p := range participants {
		participatedMeetings = append(participatedMeetings, p.Meeting)
	}

	utils.SuccessResponse(w, "获取成功", map[string]interface{}{
		"organized":   organizedMeetings,
		"participated": participatedMeetings,
	})
}

func (mc *MeetingController) GetByID(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取会议ID
	meetingIDStr := r.URL.Query().Get("id")
	if meetingIDStr == "" {
		utils.BadRequestResponse(w, "缺少会议ID参数")
		return
	}

	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "会议ID格式无效")
		return
	}

	var meeting models.Meeting
	if err := database.DB.Preload("Room").Preload("Organizer").
		First(&meeting, uint(meetingID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "会议不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取会议信息失败")
		}
		return
	}

	utils.SuccessResponse(w, "获取成功", meeting)
}

func (mc *MeetingController) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	if !ok {
		utils.UnauthorizedResponse(w, "无法获取用户信息")
		return
	}

	var meeting models.Meeting
	if err := json.NewDecoder(r.Body).Decode(&meeting); err != nil {
		utils.BadRequestResponse(w, "请求体解析失败")
		return
	}

	// 设置组织者ID
	meeting.OrganizerID = claims.UserID

	// 检查会议室是否存在
	var meetingRoom models.MeetingRoom
	if err := database.DB.First(&meetingRoom, meeting.MeetingRoomID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "会议室不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取会议室信息失败")
		}
		return
	}

	// 检查会议室是否可用
	if meetingRoom.Status != "available" {
		utils.BadRequestResponse(w, "会议室不可用")
		return
	}

	// 检查会议时间是否有效
	if meeting.StartTime.After(meeting.EndTime) {
		utils.BadRequestResponse(w, "开始时间不能晚于结束时间")
		return
	}

	if meeting.StartTime.Before(time.Now()) {
		utils.BadRequestResponse(w, "开始时间不能早于当前时间")
		return
	}

	// 检查会议室是否在该时间段已被占用
	var existingMeetings []models.Meeting
	if err := database.DB.Where("meeting_room_id = ? AND status != ? AND ((start_time <= ? AND end_time >= ?) OR (start_time <= ? AND end_time >= ?) OR (start_time >= ? AND end_time <= ?))",
		meeting.MeetingRoomID, "cancelled",
		meeting.StartTime, meeting.StartTime,
		meeting.EndTime, meeting.EndTime,
		meeting.StartTime, meeting.EndTime).Find(&existingMeetings).Error; err != nil {
		utils.InternalServerErrorResponse(w, "检查会议室可用性失败")
		return
	}

	if len(existingMeetings) > 0 {
		utils.BadRequestResponse(w, "会议室在该时间段已被占用")
		return
	}

	// 创建会议
	if err := database.DB.Create(&meeting).Error; err != nil {
		utils.InternalServerErrorResponse(w, "会议创建失败")
		return
	}

	// 加载关联信息
	database.DB.Preload("Room").Preload("Organizer").First(&meeting, meeting.ID)

	utils.SuccessResponse(w, "创建成功", meeting)
}

func (mc *MeetingController) Update(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取会议ID
	meetingIDStr := r.URL.Query().Get("id")
	if meetingIDStr == "" {
		utils.BadRequestResponse(w, "缺少会议ID参数")
		return
	}

	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "会议ID格式无效")
		return
	}

	var updateData struct {
		Title         string    `json:"title"`
		Description   string    `json:"description"`
		MeetingRoomID uint      `json:"meeting_room_id"`
		StartTime     time.Time `json:"start_time"`
		EndTime       time.Time `json:"end_time"`
		Status        string    `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		utils.BadRequestResponse(w, "请求体解析失败")
		return
	}

	var meeting models.Meeting
	if err := database.DB.First(&meeting, uint(meetingID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "会议不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取会议信息失败")
		}
		return
	}

	if updateData.Title != "" {
		meeting.Title = updateData.Title
	}

	if updateData.Description != "" {
		meeting.Description = updateData.Description
	}

	if updateData.MeetingRoomID > 0 {
		// 检查会议室是否存在
		var meetingRoom models.MeetingRoom
		if err := database.DB.First(&meetingRoom, updateData.MeetingRoomID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				utils.NotFoundResponse(w, "会议室不存在")
			} else {
				utils.InternalServerErrorResponse(w, "获取会议室信息失败")
			}
			return
		}

		// 检查会议室是否可用
		if meetingRoom.Status != "available" {
			utils.BadRequestResponse(w, "会议室不可用")
			return
		}

		meeting.MeetingRoomID = updateData.MeetingRoomID
	}

	if !updateData.StartTime.IsZero() {
		meeting.StartTime = updateData.StartTime
	}

	if !updateData.EndTime.IsZero() {
		meeting.EndTime = updateData.EndTime
	}

	if updateData.Status != "" {
		meeting.Status = updateData.Status
	}

	// 检查会议时间是否有效
	if meeting.StartTime.After(meeting.EndTime) {
		utils.BadRequestResponse(w, "开始时间不能晚于结束时间")
		return
	}

	// 检查会议室是否在该时间段已被占用
	var existingMeetings []models.Meeting
	if err := database.DB.Where("id != ? AND meeting_room_id = ? AND status != ? AND ((start_time <= ? AND end_time >= ?) OR (start_time <= ? AND end_time >= ?) OR (start_time >= ? AND end_time <= ?))",
		meeting.ID, meeting.MeetingRoomID, "cancelled",
		meeting.StartTime, meeting.StartTime,
		meeting.EndTime, meeting.EndTime,
		meeting.StartTime, meeting.EndTime).Find(&existingMeetings).Error; err != nil {
		utils.InternalServerErrorResponse(w, "检查会议室可用性失败")
		return
	}

	if len(existingMeetings) > 0 {
		utils.BadRequestResponse(w, "会议室在该时间段已被占用")
		return
	}

	if err := database.DB.Save(&meeting).Error; err != nil {
		utils.InternalServerErrorResponse(w, "更新会议信息失败")
		return
	}

	// 加载关联信息
	database.DB.Preload("Room").Preload("Organizer").First(&meeting, meeting.ID)

	utils.SuccessResponse(w, "更新成功", meeting)
}

func (mc *MeetingController) Cancel(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取会议ID
	meetingIDStr := r.URL.Query().Get("id")
	if meetingIDStr == "" {
		utils.BadRequestResponse(w, "缺少会议ID参数")
		return
	}

	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "会议ID格式无效")
		return
	}

	var meeting models.Meeting
	if err := database.DB.First(&meeting, uint(meetingID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "会议不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取会议信息失败")
		}
		return
	}

	// 取消会议
	meeting.Status = "cancelled"
	if err := database.DB.Save(&meeting).Error; err != nil {
		utils.InternalServerErrorResponse(w, "取消会议失败")
		return
	}

	utils.SuccessResponse(w, "会议已取消", nil)
}

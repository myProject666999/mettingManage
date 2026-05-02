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

type ParticipantController struct{}

func (pc *ParticipantController) GetAll(w http.ResponseWriter, r *http.Request) {
	var participants []models.Participant
	if err := database.DB.Preload("Meeting").Preload("Meeting.Room").Preload("Meeting.Organizer").
		Preload("User").Find(&participants).Error; err != nil {
		utils.InternalServerErrorResponse(w, "获取参会者列表失败")
		return
	}

	utils.SuccessResponse(w, "获取成功", participants)
}

func (pc *ParticipantController) GetUserParticipations(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	if !ok {
		utils.UnauthorizedResponse(w, "无法获取用户信息")
		return
	}

	var participants []models.Participant
	if err := database.DB.Where("user_id = ?", claims.UserID).
		Preload("Meeting").Preload("Meeting.Room").Preload("Meeting.Organizer").
		Find(&participants).Error; err != nil {
		utils.InternalServerErrorResponse(w, "获取参会列表失败")
		return
	}

	utils.SuccessResponse(w, "获取成功", participants)
}

func (pc *ParticipantController) GetByMeeting(w http.ResponseWriter, r *http.Request) {
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

	var participants []models.Participant
	if err := database.DB.Where("meeting_id = ?", uint(meetingID)).
		Preload("Meeting").Preload("Meeting.Room").Preload("Meeting.Organizer").
		Preload("User").Find(&participants).Error; err != nil {
		utils.InternalServerErrorResponse(w, "获取参会者列表失败")
		return
	}

	utils.SuccessResponse(w, "获取成功", participants)
}

func (pc *ParticipantController) Create(w http.ResponseWriter, r *http.Request) {
	var participant models.Participant
	if err := json.NewDecoder(r.Body).Decode(&participant); err != nil {
		utils.BadRequestResponse(w, "请求体解析失败")
		return
	}

	// 检查会议是否存在
	var meeting models.Meeting
	if err := database.DB.First(&meeting, participant.MeetingID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "会议不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取会议信息失败")
		}
		return
	}

	// 检查会议是否已取消
	if meeting.Status == "cancelled" {
		utils.BadRequestResponse(w, "会议已取消")
		return
	}

	// 检查用户是否存在
	var user models.User
	if err := database.DB.First(&user, participant.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "用户不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取用户信息失败")
		}
		return
	}

	// 检查用户是否已参会
	var existingParticipant models.Participant
	if err := database.DB.Where("meeting_id = ? AND user_id = ?", participant.MeetingID, participant.UserID).
		First(&existingParticipant).Error; err == nil {
		utils.BadRequestResponse(w, "用户已参会")
		return
	}

	// 创建参会记录
	if err := database.DB.Create(&participant).Error; err != nil {
		utils.InternalServerErrorResponse(w, "创建参会记录失败")
		return
	}

	// 加载关联信息
	database.DB.Preload("Meeting").Preload("Meeting.Room").Preload("Meeting.Organizer").
		Preload("User").First(&participant, participant.ID)

	utils.SuccessResponse(w, "创建成功", participant)
}

func (pc *ParticipantController) JoinMeeting(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	if !ok {
		utils.UnauthorizedResponse(w, "无法获取用户信息")
		return
	}

	var joinData struct {
		MeetingID uint `json:"meeting_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&joinData); err != nil {
		utils.BadRequestResponse(w, "请求体解析失败")
		return
	}

	// 检查会议是否存在
	var meeting models.Meeting
	if err := database.DB.First(&meeting, joinData.MeetingID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "会议不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取会议信息失败")
		}
		return
	}

	// 检查会议是否已取消
	if meeting.Status == "cancelled" {
		utils.BadRequestResponse(w, "会议已取消")
		return
	}

	// 检查用户是否已参会
	var existingParticipant models.Participant
	if err := database.DB.Where("meeting_id = ? AND user_id = ?", joinData.MeetingID, claims.UserID).
		First(&existingParticipant).Error; err == nil {
		utils.BadRequestResponse(w, "您已参会")
		return
	}

	// 创建参会记录
	participant := models.Participant{
		MeetingID: joinData.MeetingID,
		UserID:    claims.UserID,
		Status:    "invited",
	}

	if err := database.DB.Create(&participant).Error; err != nil {
		utils.InternalServerErrorResponse(w, "加入会议失败")
		return
	}

	// 加载关联信息
	database.DB.Preload("Meeting").Preload("Meeting.Room").Preload("Meeting.Organizer").
		Preload("User").First(&participant, participant.ID)

	utils.SuccessResponse(w, "加入会议成功", participant)
}

func (pc *ParticipantController) Update(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取参会记录ID
	participantIDStr := r.URL.Query().Get("id")
	if participantIDStr == "" {
		utils.BadRequestResponse(w, "缺少参会记录ID参数")
		return
	}

	participantID, err := strconv.ParseUint(participantIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "参会记录ID格式无效")
		return
	}

	var updateData struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		utils.BadRequestResponse(w, "请求体解析失败")
		return
	}

	var participant models.Participant
	if err := database.DB.First(&participant, uint(participantID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "参会记录不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取参会记录失败")
		}
		return
	}

	if updateData.Status != "" {
		participant.Status = updateData.Status
	}

	if err := database.DB.Save(&participant).Error; err != nil {
		utils.InternalServerErrorResponse(w, "更新参会记录失败")
		return
	}

	// 加载关联信息
	database.DB.Preload("Meeting").Preload("Meeting.Room").Preload("Meeting.Organizer").
		Preload("User").First(&participant, participant.ID)

	utils.SuccessResponse(w, "更新成功", participant)
}

func (pc *ParticipantController) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	if !ok {
		utils.UnauthorizedResponse(w, "无法获取用户信息")
		return
	}

	// 从URL参数中获取参会记录ID
	participantIDStr := r.URL.Query().Get("id")
	if participantIDStr == "" {
		utils.BadRequestResponse(w, "缺少参会记录ID参数")
		return
	}

	participantID, err := strconv.ParseUint(participantIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "参会记录ID格式无效")
		return
	}

	var statusData struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&statusData); err != nil {
		utils.BadRequestResponse(w, "请求体解析失败")
		return
	}

	var participant models.Participant
	if err := database.DB.Where("id = ? AND user_id = ?", uint(participantID), claims.UserID).
		First(&participant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "参会记录不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取参会记录失败")
		}
		return
	}

	participant.Status = statusData.Status

	// 如果状态是attended，设置加入时间
	if statusData.Status == "attended" {
		now := time.Now()
		participant.JoinTime = &now
	}

	if err := database.DB.Save(&participant).Error; err != nil {
		utils.InternalServerErrorResponse(w, "更新参会状态失败")
		return
	}

	utils.SuccessResponse(w, "更新成功", participant)
}

func (pc *ParticipantController) Delete(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取参会记录ID
	participantIDStr := r.URL.Query().Get("id")
	if participantIDStr == "" {
		utils.BadRequestResponse(w, "缺少参会记录ID参数")
		return
	}

	participantID, err := strconv.ParseUint(participantIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "参会记录ID格式无效")
		return
	}

	var participant models.Participant
	if err := database.DB.First(&participant, uint(participantID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "参会记录不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取参会记录失败")
		}
		return
	}

	if err := database.DB.Delete(&participant).Error; err != nil {
		utils.InternalServerErrorResponse(w, "删除参会记录失败")
		return
	}

	utils.SuccessResponse(w, "删除成功", nil)
}

func (pc *ParticipantController) LeaveMeeting(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	if !ok {
		utils.UnauthorizedResponse(w, "无法获取用户信息")
		return
	}

	// 从URL参数中获取参会记录ID
	participantIDStr := r.URL.Query().Get("id")
	if participantIDStr == "" {
		utils.BadRequestResponse(w, "缺少参会记录ID参数")
		return
	}

	participantID, err := strconv.ParseUint(participantIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "参会记录ID格式无效")
		return
	}

	var participant models.Participant
	if err := database.DB.Where("id = ? AND user_id = ?", uint(participantID), claims.UserID).
		First(&participant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "参会记录不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取参会记录失败")
		}
		return
	}

	// 设置离开时间
	now := time.Now()
	participant.LeaveTime = &now
	participant.Status = "declined"

	if err := database.DB.Save(&participant).Error; err != nil {
		utils.InternalServerErrorResponse(w, "离开会议失败")
		return
	}

	utils.SuccessResponse(w, "离开会议成功", nil)
}

import { request } from '@/utils/request'

// 用户信息类型
interface UserInfo {
  id: number
  username: string
  email: string
  phone?: string
  fullName?: string
  createdAt?: string
}

// 登录响应类型
interface LoginResponse {
  user: UserInfo
  token: string
}

// 注册响应类型
interface RegisterResponse {
  user: UserInfo
  token: string
}

export const userApi = {
  // 用户登录
  login(username: string, password: string): Promise<LoginResponse> {
    return request.post('/login', { username, password })
  },

  // 用户注册
  register(data: { username: string; password: string; email: string; phone?: string; fullName?: string }): Promise<RegisterResponse> {
    return request.post('/register', data)
  },

  // 获取用户信息
  getUserInfo(): Promise<UserInfo> {
    return request.get('/user/profile')
  },

  // 更新用户信息
  updateUserInfo(data: { email?: string; phone?: string; fullName?: string }): Promise<any> {
    return request.put('/user/profile', data)
  },

  // 修改密码
  changePassword(oldPassword: string, newPassword: string): Promise<any> {
    return request.post('/user/change-password', { old_password: oldPassword, new_password: newPassword })
  },

  // 获取用户会议列表
  getUserMeetings(): Promise<any> {
    return request.get('/user/meetings')
  },

  // 获取会议详情
  getMeetingById(id: number): Promise<any> {
    return request.get(`/user/meetings/${id}`)
  },

  // 创建会议
  createMeeting(data: {
    title: string
    description?: string
    meeting_room_id: number
    start_time: string
    end_time: string
  }): Promise<any> {
    return request.post('/user/meetings', data)
  },

  // 更新会议
  updateMeeting(id: number, data: {
    title?: string
    description?: string
    meeting_room_id?: number
    start_time?: string
    end_time?: string
  }): Promise<any> {
    return request.put(`/user/meetings/${id}`, data)
  },

  // 取消会议
  cancelMeeting(id: number): Promise<any> {
    return request.delete(`/user/meetings/${id}`)
  },

  // 获取用户参会列表
  getUserParticipations(): Promise<any> {
    return request.get('/user/participants')
  },

  // 加入会议
  joinMeeting(meetingId: number): Promise<any> {
    return request.post('/user/participants/join', { meeting_id: meetingId })
  },

  // 更新参会状态
  updateParticipationStatus(id: number, status: string): Promise<any> {
    return request.put(`/user/participants/${id}`, { status })
  },

  // 离开会议
  leaveMeeting(id: number): Promise<any> {
    return request.delete(`/user/participants/${id}`)
  },

  // 获取会议资料列表
  getMeetingDocuments(meetingId: number): Promise<any> {
    return request.get(`/user/meetings/${meetingId}/documents`)
  },

  // 下载会议资料
  downloadDocument(id: number): Promise<any> {
    return request.get(`/user/documents/${id}/download`, { responseType: 'blob' })
  },

  // 获取公告列表
  getAnnouncements(): Promise<any> {
    return request.get('/user/announcements')
  },

  // 获取公告详情
  getAnnouncementById(id: number): Promise<any> {
    return request.get(`/user/announcements/${id}`)
  },

  // 获取会议提醒列表
  getUserReminders(): Promise<any> {
    return request.get('/user/reminders')
  },

  // 获取会议提醒详情
  getReminderById(id: number): Promise<any> {
    return request.get(`/user/reminders/${id}`)
  },

  // 获取会议室列表
  getMeetingRooms(): Promise<any> {
    return request.get('/user/meeting-rooms')
  },

  // 获取会议室详情
  getMeetingRoomById(id: number): Promise<any> {
    return request.get(`/user/meeting-rooms/${id}`)
  }
}

export default userApi

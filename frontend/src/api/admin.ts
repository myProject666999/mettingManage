import { request } from '@/utils/request'

export const adminApi = {
  // 管理员登录
  login(username: string, password: string): Promise<any> {
    return request.post('/admin/login', { username, password })
  },

  // 获取管理员信息
  getAdminInfo(): Promise<any> {
    return request.get('/admin/profile')
  },

  // 更新管理员信息
  updateAdminInfo(data: { email?: string }): Promise<any> {
    return request.put('/admin/profile', data)
  },

  // 修改密码
  changePassword(oldPassword: string, newPassword: string): Promise<any> {
    return request.post('/admin/change-password', { old_password: oldPassword, new_password: newPassword })
  },

  // 管理员管理
  getAdmins(): Promise<any> {
    return request.get('/admin/admins')
  },

  getAdminById(id: number): Promise<any> {
    return request.get(`/admin/admins/${id}`)
  },

  createAdmin(data: {
    username: string
    password: string
    email: string
    role?: string
  }): Promise<any> {
    return request.post('/admin/admins', data)
  },

  updateAdmin(id: number, data: {
    username?: string
    email?: string
    role?: string
  }): Promise<any> {
    return request.put(`/admin/admins/${id}`, data)
  },

  deleteAdmin(id: number): Promise<any> {
    return request.delete(`/admin/admins/${id}`)
  },

  // 用户管理
  getUsers(): Promise<any> {
    return request.get('/admin/users')
  },

  getUserById(id: number): Promise<any> {
    return request.get(`/admin/users/${id}`)
  },

  createUser(data: {
    username: string
    password: string
    email: string
    phone?: string
    fullName?: string
  }): Promise<any> {
    return request.post('/admin/users', data)
  },

  updateUser(id: number, data: {
    username?: string
    email?: string
    phone?: string
    fullName?: string
  }): Promise<any> {
    return request.put(`/admin/users/${id}`, data)
  },

  deleteUser(id: number): Promise<any> {
    return request.delete(`/admin/users/${id}`)
  },

  // 会议室管理
  getMeetingRooms(): Promise<any> {
    return request.get('/admin/meeting-rooms')
  },

  getMeetingRoomById(id: number): Promise<any> {
    return request.get(`/admin/meeting-rooms/${id}`)
  },

  createMeetingRoom(data: {
    name: string
    location?: string
    capacity: number
    description?: string
    equipment?: string
    status?: string
  }): Promise<any> {
    return request.post('/admin/meeting-rooms', data)
  },

  updateMeetingRoom(id: number, data: {
    name?: string
    location?: string
    capacity?: number
    description?: string
    equipment?: string
    status?: string
  }): Promise<any> {
    return request.put(`/admin/meeting-rooms/${id}`, data)
  },

  deleteMeetingRoom(id: number): Promise<any> {
    return request.delete(`/admin/meeting-rooms/${id}`)
  },

  // 会议管理
  getMeetings(): Promise<any> {
    return request.get('/admin/meetings')
  },

  getMeetingById(id: number): Promise<any> {
    return request.get(`/admin/meetings/${id}`)
  },

  createMeeting(data: {
    title: string
    description?: string
    meeting_room_id: number
    organizer_id: number
    start_time: string
    end_time: string
  }): Promise<any> {
    return request.post('/admin/meetings', data)
  },

  updateMeeting(id: number, data: {
    title?: string
    description?: string
    meeting_room_id?: number
    start_time?: string
    end_time?: string
    status?: string
  }): Promise<any> {
    return request.put(`/admin/meetings/${id}`, data)
  },

  cancelMeeting(id: number): Promise<any> {
    return request.delete(`/admin/meetings/${id}`)
  },

  // 参会信息管理
  getParticipants(): Promise<any> {
    return request.get('/admin/participants')
  },

  getParticipantsByMeeting(meetingId: number): Promise<any> {
    return request.get(`/admin/meetings/${meetingId}/participants`)
  },

  createParticipant(data: {
    meeting_id: number
    user_id: number
    status?: string
  }): Promise<any> {
    return request.post('/admin/participants', data)
  },

  updateParticipant(id: number, data: {
    status?: string
  }): Promise<any> {
    return request.put(`/admin/participants/${id}`, data)
  },

  deleteParticipant(id: number): Promise<any> {
    return request.delete(`/admin/participants/${id}`)
  },

  // 会议资料管理
  getDocuments(): Promise<any> {
    return request.get('/admin/documents')
  },

  getDocumentsByMeeting(meetingId: number): Promise<any> {
    return request.get(`/admin/meetings/${meetingId}/documents`)
  },

  createDocument(data: FormData): Promise<any> {
    return request.post('/admin/documents', data, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  updateDocument(id: number, data: {
    title?: string
    description?: string
  }): Promise<any> {
    return request.put(`/admin/documents/${id}`, data)
  },

  deleteDocument(id: number): Promise<any> {
    return request.delete(`/admin/documents/${id}`)
  },

  // 公告管理
  getAnnouncements(): Promise<any> {
    return request.get('/admin/announcements')
  },

  getAnnouncementById(id: number): Promise<any> {
    return request.get(`/admin/announcements/${id}`)
  },

  createAnnouncement(data: {
    title: string
    content: string
    is_top?: boolean
    status?: string
  }): Promise<any> {
    return request.post('/admin/announcements', data)
  },

  updateAnnouncement(id: number, data: {
    title?: string
    content?: string
    is_top?: boolean
    status?: string
  }): Promise<any> {
    return request.put(`/admin/announcements/${id}`, data)
  },

  deleteAnnouncement(id: number): Promise<any> {
    return request.delete(`/admin/announcements/${id}`)
  },

  // 轮播图管理
  getCarousels(): Promise<any> {
    return request.get('/admin/carousels')
  },

  getCarouselById(id: number): Promise<any> {
    return request.get(`/admin/carousels/${id}`)
  },

  createCarousel(data: FormData): Promise<any> {
    return request.post('/admin/carousels', data, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  updateCarousel(id: number, data: FormData): Promise<any> {
    return request.put(`/admin/carousels/${id}`, data, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  deleteCarousel(id: number): Promise<any> {
    return request.delete(`/admin/carousels/${id}`)
  },

  // 会议提醒管理
  getReminders(): Promise<any> {
    return request.get('/admin/reminders')
  },

  getRemindersByMeeting(meetingId: number): Promise<any> {
    return request.get(`/admin/meetings/${meetingId}/reminders`)
  },

  createReminder(data: {
    meeting_id: number
    reminder_type: string
    reminder_time: string
    message?: string
  }): Promise<any> {
    return request.post('/admin/reminders', data)
  },

  updateReminder(id: number, data: {
    meeting_id?: number
    reminder_type?: string
    reminder_time?: string
    message?: string
  }): Promise<any> {
    return request.put(`/admin/reminders/${id}`, data)
  },

  deleteReminder(id: number): Promise<any> {
    return request.delete(`/admin/reminders/${id}`)
  }
}

export default adminApi

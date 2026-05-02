<template>
  <div class="user-home">
    <!-- 轮播图 -->
    <el-carousel height="300px" indicator-position="outside" v-if="carousels.length > 0">
      <el-carousel-item v-for="carousel in carousels" :key="carousel.id">
        <div class="carousel-item">
          <img :src="carousel.imageUrl" :alt="carousel.title" class="carousel-image" />
          <div class="carousel-content">
            <h3>{{ carousel.title }}</h3>
          </div>
        </div>
      </el-carousel-item>
    </el-carousel>
    
    <!-- 公告 -->
    <el-card class="announcement-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="title">
            <el-icon><Bell /></el-icon>
            最新公告
          </span>
        </div>
      </template>
      <el-empty v-if="announcements.length === 0" description="暂无公告" />
      <div v-else class="announcement-list">
        <div
          v-for="announcement in announcements.slice(0, 5)"
          :key="announcement.id"
          class="announcement-item"
          @click="viewAnnouncement(announcement)"
        >
          <div class="announcement-title">
            <el-tag v-if="announcement.isTop" type="danger" size="small">置顶</el-tag>
            <span>{{ announcement.title }}</span>
          </div>
          <span class="announcement-time">{{ formatDate(announcement.createdAt) }}</span>
        </div>
      </div>
    </el-card>
    
    <!-- 快捷操作 -->
    <el-row :gutter="20" class="quick-actions">
      <el-col :span="6">
        <el-card shadow="hover" class="action-card" @click="navigateTo('/user/meeting-rooms')">
          <div class="action-content">
            <el-icon size="40" color="#409eff"><OfficeBuilding /></el-icon>
            <span class="action-text">查看会议室</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="action-card" @click="navigateTo('/user/meetings')">
          <div class="action-content">
            <el-icon size="40" color="#67c23a"><Calendar /></el-icon>
            <span class="action-text">我的会议</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="action-card" @click="navigateTo('/user/participations')">
          <div class="action-content">
            <el-icon size="40" color="#e6a23c"><User /></el-icon>
            <span class="action-text">参会信息</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="action-card" @click="navigateTo('/user/reminders')">
          <div class="action-content">
            <el-icon size="40" color="#f56c6c"><AlarmClock /></el-icon>
            <span class="action-text">会议提醒</span>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 统计信息 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="8">
        <el-card shadow="hover" class="stats-card">
          <div class="stats-content">
            <div class="stats-icon blue">
              <el-icon><Calendar /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-value">{{ stats.meetings }}</div>
              <div class="stats-label">我的会议</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover" class="stats-card">
          <div class="stats-content">
            <div class="stats-icon green">
              <el-icon><User /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-value">{{ stats.participations }}</div>
              <div class="stats-label">参会次数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover" class="stats-card">
          <div class="stats-content">
            <div class="stats-icon orange">
              <el-icon><AlarmClock /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-value">{{ stats.reminders }}</div>
              <div class="stats-label">待提醒</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 公告详情弹窗 -->
    <el-dialog
      v-model="announcementDialogVisible"
      title="公告详情"
      width="600px"
    >
      <div v-if="currentAnnouncement" class="announcement-detail">
        <h2 class="detail-title">{{ currentAnnouncement.title }}</h2>
        <div class="detail-meta">
          <span>发布时间：{{ formatDate(currentAnnouncement.createdAt) }}</span>
        </div>
        <div class="detail-content">
          {{ currentAnnouncement.content }}
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import userApi from '@/api/user'
import dayjs from 'dayjs'

const router = useRouter()

const carousels = ref<any[]>([])
const announcements = ref<any[]>([])
const stats = reactive({
  meetings: 0,
  participations: 0,
  reminders: 0
})

const announcementDialogVisible = ref(false)
const currentAnnouncement = ref<any>(null)

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const viewAnnouncement = (announcement: any) => {
  currentAnnouncement.value = announcement
  announcementDialogVisible.value = true
}

const navigateTo = (path: string) => {
  router.push(path)
}

const fetchData = async () => {
  try {
    // 获取轮播图
    await userApi.getAnnouncements()
    // 这里应该获取轮播图，但当前API没有直接的轮播图接口
    // 暂时使用公告数据作为示例
    
    // 获取公告
    const announcementsRes = await userApi.getAnnouncements()
    announcements.value = announcementsRes || []
    
    // 获取用户会议
    const meetingsRes = await userApi.getUserMeetings()
    stats.meetings = (meetingsRes.organized?.length || 0) + (meetingsRes.participated?.length || 0)
    
    // 获取用户参会记录
    const participationsRes = await userApi.getUserParticipations()
    stats.participations = participationsRes?.length || 0
    
    // 获取用户提醒
    const remindersRes = await userApi.getUserReminders()
    stats.reminders = remindersRes?.length || 0
  } catch (error) {
    console.error('获取数据失败:', error)
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.user-home {
  max-width: 1200px;
  margin: 0 auto;
}

.carousel-item {
  position: relative;
  height: 100%;
  overflow: hidden;
  border-radius: 8px;
}

.carousel-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.carousel-content {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.6));
  padding: 20px;
  color: #fff;
}

.carousel-content h3 {
  margin: 0;
  font-size: 18px;
}

.announcement-card {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header .title {
  font-size: 16px;
  font-weight: bold;
  display: flex;
  align-items: center;
  gap: 8px;
}

.announcement-list {
  max-height: 300px;
  overflow-y: auto;
}

.announcement-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #ebeef5;
  cursor: pointer;
  transition: background-color 0.2s;
}

.announcement-item:hover {
  background-color: #f5f7fa;
  margin: 0 -10px;
  padding: 12px 10px;
}

.announcement-title {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
}

.announcement-title span {
  font-size: 14px;
  color: #303133;
}

.announcement-time {
  font-size: 12px;
  color: #909399;
  flex-shrink: 0;
}

.quick-actions {
  margin-top: 20px;
}

.action-card {
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.action-card:hover {
  transform: translateY(-5px);
}

.action-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.action-text {
  margin-top: 10px;
  font-size: 14px;
  color: #606266;
}

.stats-row {
  margin-top: 20px;
}

.stats-card {
  cursor: pointer;
}

.stats-content {
  display: flex;
  align-items: center;
}

.stats-icon {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
}

.stats-icon.blue {
  background-color: rgba(64, 158, 255, 0.1);
  color: #409eff;
}

.stats-icon.green {
  background-color: rgba(103, 194, 58, 0.1);
  color: #67c23a;
}

.stats-icon.orange {
  background-color: rgba(230, 162, 60, 0.1);
  color: #e6a23c;
}

.stats-icon .el-icon {
  font-size: 30px;
}

.stats-info {
  flex: 1;
}

.stats-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
}

.stats-label {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}

.announcement-detail {
  padding: 10px;
}

.detail-title {
  font-size: 20px;
  color: #303133;
  margin-bottom: 10px;
}

.detail-meta {
  font-size: 12px;
  color: #909399;
  margin-bottom: 20px;
}

.detail-content {
  font-size: 14px;
  line-height: 1.8;
  color: #606266;
}
</style>

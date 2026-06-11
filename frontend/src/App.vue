<template>
  <el-container style="min-height: 100vh">
    <el-header style="background: linear-gradient(135deg, #409eff, #337ecc); color: #fff; display: flex; align-items: center; padding: 0 20px; justify-content: space-between">
      <span style="font-size: 20px; font-weight: bold; display: flex; align-items: center; gap: 8px">
        <svg viewBox="0 0 24 24" width="22" height="22" style="fill: currentColor">
          <path d="M19 9.3V4h-3v2.6L12 3 2 12h3v8h5v-6h4v6h5v-8h3l-3-2.7z"/>
        </svg>
        出租房收租管理系统
      </span>
      <el-button
        :type="isDashboard ? 'default' : 'primary'"
        :icon="null"
        class="dash-btn"
        :class="{ 'dash-btn-active': isDashboard }"
        @click="$router.push('/dashboard')"
      >
        <svg viewBox="0 0 24 24" width="16" height="16" class="dash-btn-icon" style="fill: currentColor">
          <path d="M3 13h8V3H3v10zm0 8h8v-6H3v6zm10 0h8V11h-8v10zm0-18v6h8V3h-8z"/>
        </svg>
        <span>数据大屏</span>
      </el-button>
    </el-header>
    <el-container>
      <el-aside width="200px" style="background: #f0f4f8">
        <el-menu
          :default-active="activeMenu"
          router
          style="border-right: none; height: 100%"
        >
          <el-menu-item index="/">
            <el-icon><HomeFilled /></el-icon>
            <span>房间管理</span>
          </el-menu-item>
          <el-menu-item index="/bills">
            <el-icon><Document /></el-icon>
            <span>月度账单</span>
          </el-menu-item>
          <el-menu-item index="/history">
            <el-icon><Clock /></el-icon>
            <span>历史账单</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-main style="padding: 20px">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const activeMenu = computed(() => route.path)
const isDashboard = computed(() => route.path === '/dashboard')
</script>

<style>
:root {
  --el-bg-color-page: #f0f4f8;
}
body {
  background: #f0f4f8;
}
.dash-btn {
  display: inline-flex !important;
  align-items: center !important;
  gap: 6px !important;
  padding: 8px 18px !important;
  border-radius: 10px !important;
  font-weight: 600 !important;
  font-size: 14px !important;
  border: none !important;
  background: rgba(255, 255, 255, 0.2) !important;
  color: #fff !important;
  transition: all 0.25s ease !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1) !important;
}

.dash-btn:hover {
  background: rgba(255, 255, 255, 0.35) !important;
  transform: translateY(-1px) !important;
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.18) !important;
}

.dash-btn-active {
  background: #fff !important;
  color: #409eff !important;
  box-shadow: 0 2px 10px rgba(64, 158, 255, 0.35) !important;
}

.dash-btn-active:hover {
  background: #f0f5ff !important;
  color: #337ecc !important;
}

.dash-btn-icon {
  flex-shrink: 0;
}
</style>

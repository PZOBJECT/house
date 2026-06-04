<template>
  <el-container style="min-height: 100vh">
    <el-header style="background: #409eff; color: #fff; display: flex; align-items: center; padding: 0 20px; font-size: 20px; font-weight: bold">
      出租房收租管理系统
    </el-header>
    <el-container>
      <el-aside width="200px" style="background: #f0f2f5">
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
</script>

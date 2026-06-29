<script setup lang="ts">
import { useRoute } from 'vue-router'
import BottomNav from './components/BottomNav.vue'
import zhCn from 'element-plus/es/locale/lang/zh-cn'

const route = useRoute()
</script>

<template>
  <el-config-provider :locale="zhCn">
    <div class="app-background"></div>
    <div :class="['page-container', { 'login-page': route.name === 'Login' }]">
      <router-view v-slot="{ Component }">
        <component :is="Component" />
      </router-view>
    </div>
    <BottomNav v-if="route.name !== 'Login'" />
  </el-config-provider>
</template>

<style>
/* Global styles for App.vue if needed in future */
</style>

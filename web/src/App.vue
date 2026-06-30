<script setup lang="ts">
import { useRoute } from 'vue-router'
import BottomNav from './components/BottomNav.vue'
import zhCn from 'element-plus/es/locale/lang/zh-cn'

const route = useRoute()
</script>

<template>
  <el-config-provider :locale="zhCn">
    <div class="app-background"></div>
    <div :class="['page-container', { 'login-page': route.name === 'Login' }]" style="position: relative;">
      <router-view v-slot="{ Component }">
        <transition name="page">
          <keep-alive>
            <component :is="Component" :key="route.path" />
          </keep-alive>
        </transition>
      </router-view>
    </div>
    <BottomNav v-if="route.name !== 'Login'" />
  </el-config-provider>
</template>

<style>
.page-enter-active,
.page-leave-active {
  transition: opacity 0.15s ease;
}
.page-leave-active {
  position: absolute;
  top: 20px;
  left: 20px;
  right: 20px;
}
@media screen and (max-width: 600px) {
  .page-leave-active {
    top: 12px;
    left: 12px;
    right: 12px;
  }
}
.page-enter-from,
.page-leave-to {
  opacity: 0;
}
</style>

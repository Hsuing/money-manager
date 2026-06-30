<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const navItems = [
  { path: '/dashboard', label: '总览', id: 'home' },
  { path: '/ledger', label: '明细', id: 'list' },
  { path: '/analysis', label: '统计', id: 'chart' },
  { path: '/settings', label: '设置', id: 'setting' }
]

const activeIndex = computed(() => {
  const idx = navItems.findIndex(item => item.path === route.path)
  return idx === -1 ? 0 : idx
})
</script>

<template>
  <nav class="bottom-nav-container">
    <div class="bottom-nav">
      <!-- Sliding Pill Indicator -->
      <div 
        class="nav-indicator" 
        :style="{ transform: `translateX(${activeIndex * 100}%)` }"
      >
        <div class="indicator-pill"></div>
      </div>

      <router-link
        v-for="(item, index) in navItems"
        :key="item.path"
        :to="item.path"
        class="nav-item"
        :class="{ active: route.path === item.path }"
      >
        <div class="nav-content">
          <svg v-if="item.id === 'home'" class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
            <polyline points="9 22 9 12 15 12 15 22"></polyline>
          </svg>
          <svg v-else-if="item.id === 'list'" class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="8" y1="6" x2="21" y2="6"></line>
            <line x1="8" y1="12" x2="21" y2="12"></line>
            <line x1="8" y1="18" x2="21" y2="18"></line>
            <line x1="3" y1="6" x2="3.01" y2="6"></line>
            <line x1="3" y1="12" x2="3.01" y2="12"></line>
            <line x1="3" y1="18" x2="3.01" y2="18"></line>
          </svg>
          <svg v-else-if="item.id === 'chart'" class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="18" y1="20" x2="18" y2="10"></line>
            <line x1="12" y1="20" x2="12" y2="4"></line>
            <line x1="6" y1="20" x2="6" y2="14"></line>
          </svg>
          <svg v-else-if="item.id === 'setting'" class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="3"></circle>
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
          </svg>
          <span class="label">{{ item.label }}</span>
        </div>
      </router-link>
    </div>
  </nav>
</template>

<style scoped>
.bottom-nav-container {
  position: fixed;
  bottom: 0;
  left: 0;
  width: 100%;
  display: flex;
  justify-content: center;
  z-index: 1000;
  padding: 0 16px 16px;
  pointer-events: none;
}
.bottom-nav {
  position: relative;
  display: flex;
  width: 100%;
  max-width: 600px;
  background: rgba(22, 22, 24, 0.75);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 28px;
  padding: 8px 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.5);
  pointer-events: auto;
}
.nav-indicator {
  position: absolute;
  top: 8px;
  left: 12px;
  width: calc((100% - 24px) / 4);
  height: calc(100% - 16px);
  transition: transform 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
  z-index: 0;
}
.indicator-pill {
  width: 100%;
  height: 100%;
  background: rgba(106, 176, 255, 0.15); /* var(--color-primary) with opacity */
  border-radius: 20px;
}
.nav-item {
  flex: 1;
  position: relative;
  z-index: 1;
  text-decoration: none;
  color: #888;
  display: flex;
  justify-content: center;
  align-items: center;
  -webkit-tap-highlight-color: transparent;
  outline: none;
}
.nav-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 6px 0;
  width: 100%;
  transition: transform 0.2s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.nav-item:active .nav-content {
  transform: scale(0.85);
}
.nav-item.active {
  color: var(--color-primary);
}
.icon {
  width: 22px;
  height: 22px;
  transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.nav-item.active .icon {
  transform: translateY(-2px);
}
.label {
  font-size: 11px;
  font-weight: 500;
  transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.nav-item.active .label {
  font-weight: 700;
}
</style>

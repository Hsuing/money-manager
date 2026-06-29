<script setup lang="ts">
import { useRouter } from 'vue-router'
const router = useRouter()
</script>
<template>
  <div class="subpage">
    <div class="header">
      <button @click="router.back()" class="back-btn">‹</button>
      <h2>系统管理</h2>
      <div style="width:24px"></div>
    </div>
    <div class="content">
      <p style="color:var(--color-text-secondary);text-align:center;margin-top:50px">功能开发中...</p>
    </div>
  </div>
</template>
<style scoped>
.subpage { display: flex; flex-direction: column; gap: 20px; }
.header { display: flex; justify-content: space-between; align-items: center; padding-top: 10px; }
.back-btn { background: none; border: none; color: white; font-size: 30px; line-height: 1; cursor: pointer; padding: 0 10px; margin-left:-10px; }
</style>

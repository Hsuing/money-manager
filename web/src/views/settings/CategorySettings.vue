<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { getCategories, addCategory, updateCategory, deleteCategory, getTransactionTypes } from '../../service/api'

const router = useRouter()

const categories = ref<any[]>([])
const transactionTypes = ref<any[]>([])

const currentTab = ref('list') // 'list' or 'archive'
const expandedParents = ref<number[]>([])
const activeCategoryId = ref<number | null>(null)
const isEdit = ref(false)
const showForm = ref(false)

const currentCategory = ref({
  id: 0,
  name: '',
  type: '', // Bound TransactionType Name
  parent_id: 0,
  icon: '📁',
  is_exclude: false
})

const showTypeSelector = ref(false)

const fetchDependencies = async () => {
  try {
    const [cats, types] = await Promise.all([
      getCategories(),
      getTransactionTypes()
    ])
    categories.value = cats
    transactionTypes.value = types
    
    // Auto-expand is removed per user request. 
    // State of expandedParents will be preserved across saves.
  } catch (error) {
    console.error('Failed to fetch data', error)
  }
}

onMounted(() => {
  fetchDependencies()
})

const primaryCategories = computed(() => {
  return categories.value.filter(c => c.parent_id === 0)
})

const getChildren = (parentId: number) => {
  return categories.value.filter(c => c.parent_id === parentId)
}

const toggleExpand = (id: number) => {
  const index = expandedParents.value.indexOf(id)
  if (index > -1) {
    expandedParents.value.splice(index, 1)
  } else {
    expandedParents.value.push(id)
  }
}

const openAdd = () => {
  isEdit.value = false
  activeCategoryId.value = null
  showForm.value = true
  currentCategory.value = {
    id: 0,
    name: '',
    type: '',
    parent_id: 0,
    icon: '📁',
    is_exclude: false
  }
}

const selectCategory = (cat: any) => {
  isEdit.value = true
  activeCategoryId.value = cat.id
  currentCategory.value = { ...cat }
  showForm.value = true
}

const saveCategory = async () => {
  if (!currentCategory.value.name.trim()) {
    alert('请输入分类名称')
    return
  }
  
  const isDuplicate = categories.value.some(c => 
    c.name === currentCategory.value.name && 
    c.parent_id === currentCategory.value.parent_id &&
    c.id !== currentCategory.value.id
  )
  if (isDuplicate) {
    alert('该层级下已存在同名分类！')
    return
  }

  try {
    if (isEdit.value) {
      await updateCategory(currentCategory.value.id, currentCategory.value)
    } else {
      await addCategory(currentCategory.value)
    }
    await fetchDependencies()
    showForm.value = false
    isEdit.value = false
    activeCategoryId.value = null
  } catch (error) {
    console.error('Failed to save category', error)
  }
}

const handleArchive = () => {
  alert('归档功能暂未接入数据库 TODO')
}

const handleDisable = () => {
  alert('停用功能暂未接入数据库 TODO')
}

const selectType = (typeName: string) => {
  currentCategory.value.type = typeName
  showTypeSelector.value = false
}

const getFlowTag = (flow: string) => {
  if (flow === 'inflow') return { class: 'tag-income', text: '收入' }
  if (flow === 'outflow') return { class: 'tag-expense', text: '支出' }
  return { class: 'tag-transfer', text: '转账' }
}

const getTypeDisplay = (type: string) => {
  if (!type) return { text: '未知', class: 'tag-expense' }
  if (type.includes('收入') || type.includes('入')) return { text: '收入', class: 'tag-income' }
  if (type.includes('支出') || type.includes('出')) return { text: '支出', class: 'tag-expense' }
  return { text: type, class: 'tag-expense' }
}
</script>

<template>
  <div class="app-container">
    <!-- Header -->
    <div class="app-header">
      <h2 class="title">分类管理</h2>
      <div class="header-actions">
        <button class="btn-primary" @click="openAdd">
          <svg viewBox="0 0 24 24" width="16" height="16" stroke="currentColor" stroke-width="2" fill="none"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
          添加
        </button>
        <button class="close-btn" @click="router.back()">✕</button>
      </div>
    </div>

    <!-- Centered Tabs -->
    <div class="tabs-container">
      <div class="tabs">
        <div class="tab-item" :class="{active: currentTab === 'list'}" @click="currentTab = 'list'">
          <svg viewBox="0 0 24 24" width="16" height="16" stroke="currentColor" stroke-width="2" fill="none" style="margin-right: 6px"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path></svg>
          分类列表
        </div>
        <div class="tab-item" :class="{active: currentTab === 'archive'}" @click="currentTab = 'archive'">
          <svg viewBox="0 0 24 24" width="16" height="16" stroke="currentColor" stroke-width="2" fill="none" style="margin-right: 6px"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path></svg>
          归档
        </div>
      </div>
    </div>

    <!-- List -->
    <div class="category-list">
      <div class="primary-card" v-for="p in primaryCategories" :key="p.id">
        <!-- Primary Header -->
        <div class="card-header" @click="toggleExpand(p.id)">
          <div class="collapse-icon">
            <svg v-if="expandedParents.includes(p.id)" viewBox="0 0 24 24" width="16" height="16" stroke="currentColor" stroke-width="2" fill="none"><polyline points="6 9 12 15 18 9"></polyline></svg>
            <svg v-else viewBox="0 0 24 24" width="16" height="16" stroke="currentColor" stroke-width="2" fill="none"><polyline points="9 18 15 12 9 6"></polyline></svg>
          </div>
          <div class="cat-title">{{ p.name }}</div>
          <div class="cat-badge" v-if="p.type" :class="getTypeDisplay(p.type).class">{{ getTypeDisplay(p.type).text }}</div>
          <div class="spacer"></div>
          <div class="edit-icon" @click.stop="selectCategory(p)">
            <svg viewBox="0 0 24 24" width="16" height="16" stroke="currentColor" stroke-width="2" fill="none"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
          </div>
        </div>
        
        <!-- Children List -->
        <div class="children-container" v-show="expandedParents.includes(p.id)">
          <div 
            class="child-card" 
            v-for="c in getChildren(p.id)" 
            :key="c.id"
          >
            <div class="cat-title">{{ c.name }}</div>
            <div class="cat-badge" v-if="c.type" :class="getTypeDisplay(c.type).class">{{ getTypeDisplay(c.type).text }}</div>
            <div class="spacer"></div>
            <div class="edit-icon" @click.stop="selectCategory(c)">
              <svg viewBox="0 0 24 24" width="16" height="16" stroke="currentColor" stroke-width="2" fill="none"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit Form Overlay -->
    <div class="form-overlay" :class="{'show': showForm}">
      <div class="form-content">
        <div class="editor-header">
          <h3>{{ isEdit ? '编辑分类' : '添加分类' }}</h3>
          <button class="close-btn" @click="showForm = false">✕</button>
        </div>

        <div class="form-scroll-area">
          <div class="form-card">
            <div class="form-group">
              <label>分类名称 <span class="required">*</span></label>
              <input type="text" v-model="currentCategory.name" placeholder="请输入分类名称" />
            </div>

            <div class="form-group">
              <label>父级分类</label>
              <div class="select-box">
                <select v-model="currentCategory.parent_id">
                  <option :value="0">无 (一级分类)</option>
                  <template v-for="p in primaryCategories" :key="p.id">
                    <option :value="p.id" v-if="p.id !== currentCategory.id">{{ p.name }}</option>
                  </template>
                </select>
                <div class="arrow">›</div>
              </div>
            </div>

            <div class="form-group relative">
              <label>收支绑定</label>
              <div class="select-box" @click="showTypeSelector = !showTypeSelector">
                <span class="binding-text" :class="{'has-value': currentCategory.type}">
                  <span v-if="currentCategory.type" class="cat-badge" :class="getTypeDisplay(currentCategory.type).class" style="margin:0;">
                    {{ currentCategory.type }}
                  </span>
                  <span v-else>点击选择收支</span>
                </span>
                <div class="arrow">›</div>
              </div>
              <p class="help-text">绑定收支后，选择此分类时自动选择对应收支</p>

              <!-- Dropdown for Type Selection -->
              <div class="type-dropdown" v-if="showTypeSelector">
                <div class="type-item" v-for="t in transactionTypes" :key="t.id" @click.stop="selectType(t.name)">
                  {{ t.name }} <span class="cat-badge" :class="getFlowTag(t.flow).class" style="margin-left: 8px;">{{ getFlowTag(t.flow).text }}</span>
                </div>
              </div>
            </div>

            <div class="form-group switch-row">
              <div class="switch-left">
                <div class="switch-title">不计入分析</div>
                <div class="switch-desc">开启后此分类不计入统计分析</div>
              </div>
              <label class="toggle-switch">
                <input type="checkbox" v-model="currentCategory.is_exclude">
                <span class="slider"></span>
              </label>
            </div>
          </div>

          <div class="editor-footer">
            <div class="footer-actions-left">
              <button class="text-btn outline-yellow" @click="handleArchive">归档分类</button>
              <div class="divider"></div>
              <button class="text-btn outline-red" @click="handleDisable">停用分类</button>
            </div>
            <div class="footer-actions-right">
              <button class="btn-outline" @click="showForm = false">取消</button>
              <button class="btn-primary" @click="saveCategory">保存修改</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.app-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
  background-color: #1a1a1c; /* Match dark background */
  color: #eeeeee;
  overflow: hidden;
  position: relative;
}

/* Header */
.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background-color: #1a1a1c;
  border-bottom: 1px solid #2a2a2c;
}
.app-header .title {
  font-size: 20px;
  font-weight: 600;
  margin: 0;
  color: #ffffff;
}
.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}
.btn-primary {
  background-color: #3b82f6;
  color: #fff;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 14px;
  padding: 8px 16px;
  transition: opacity 0.2s;
}
.btn-primary:hover {
  opacity: 0.9;
}
.close-btn {
  background: none;
  border: none;
  font-size: 22px;
  color: #888;
  cursor: pointer;
  padding: 0;
  display: flex;
  align-items: center;
}
.close-btn:hover {
  color: #fff;
}

/* Tabs */
.tabs-container {
  padding: 24px 0;
  display: flex;
  justify-content: center;
}
.tabs {
  display: inline-flex;
  background-color: transparent;
  border-radius: 20px;
  border: 1px solid #3a3a3c;
  padding: 4px;
}
.tab-item {
  padding: 8px 24px;
  border-radius: 16px;
  font-size: 14px;
  color: #8e8e93;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}
.tab-item.active {
  background-color: #3a5c7c; /* Blue-ish bg for active tab based on screenshot */
  color: #86b9f2;
}
.tab-item.active svg {
  color: #86b9f2;
}

/* List */
.category-list {
  flex: 1;
  overflow-y: auto;
  padding: 0 24px 100px 24px; /* Added horizontal padding to shrink towards middle */
  max-width: 600px;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
}

/* Cards */
.primary-card {
  background-color: #2c2c2e;
  border-radius: 12px;
  margin-bottom: 16px;
  overflow: hidden;
}
.card-header {
  display: flex;
  align-items: center;
  padding: 20px 24px;
  cursor: pointer;
}
.collapse-icon {
  color: #888;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
}
.cat-title {
  font-size: 16px;
  font-weight: 500;
  color: #ffffff;
}
.cat-badge {
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 4px;
  margin-left: 12px;
}
/* Subdued dark mode tags */
.tag-expense {
  background-color: #3a2c2c;
  color: #ff6b6b;
}
.tag-income {
  background-color: #2c3a30;
  color: #67c23a;
}
.tag-transfer {
  background-color: #2c343a;
  color: #409eff;
}
.spacer {
  flex: 1;
}
.edit-icon {
  color: #888;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4px;
  cursor: pointer;
  transition: color 0.2s;
}
.edit-icon:hover {
  color: #fff;
}

/* Children */
.children-container {
  padding: 0 24px 24px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.child-card {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  background-color: #262628;
  border: 1px solid #3a3a3c;
  border-radius: 8px;
  cursor: pointer;
}
.child-card .cat-title {
  font-size: 15px;
}

/* Edit Form Overlay */
.form-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  height: 100vh;
  background-color: #1a1a1c;
  z-index: 2001;
  transform: translateX(100%);
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  justify-content: center;
}
.form-overlay.show {
  transform: translateX(0);
}
.form-content {
  width: 100%;
  max-width: 600px;
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #1a1a1c;
  position: relative;
}
.form-scroll-area {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

/* Form */
.form-card {
  background: #2c2c2e;
  border-radius: 12px;
  padding: 32px;
}
.form-group {
  margin-bottom: 24px;
}
.form-group label {
  display: block;
  font-size: 14px;
  color: #aaaaaa;
  margin-bottom: 8px;
  font-weight: 500;
}
.required {
  color: #f56c6c;
}
.form-group input[type="text"] {
  width: 100%;
  padding: 12px 16px;
  background-color: #1c1c1e;
  border: 1px solid #3a3a3c;
  border-radius: 6px;
  font-size: 14px;
  color: #eeeeee;
  transition: border-color 0.2s;
  box-sizing: border-box;
}
.form-group input[type="text"]:focus {
  outline: none;
  border-color: #3b82f6;
}
.select-box {
  position: relative;
  background-color: #1c1c1e;
  border: 1px solid #3a3a3c;
  border-radius: 6px;
  cursor: pointer;
}
.select-box select {
  width: 100%;
  padding: 12px 16px;
  appearance: none;
  border: none;
  background: transparent;
  font-size: 14px;
  color: #eeeeee;
  cursor: pointer;
}
.select-box select:focus {
  outline: none;
}
.select-box select option {
  background: #1c1c1e;
  color: #eee;
}
.binding-text {
  display: inline-block;
  padding: 12px 16px;
  font-size: 14px;
  color: #666666;
  width: 100%;
  box-sizing: border-box;
}
.binding-text.has-value {
  color: #eeeeee;
}
.arrow {
  position: absolute;
  right: 16px;
  top: 50%;
  transform: translateY(-50%);
  color: #666;
  pointer-events: none;
  font-family: monospace;
}
.help-text {
  font-size: 12px;
  color: #666666;
  margin-top: 8px;
}
.relative {
  position: relative;
}
.type-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: #2c2c2e;
  border: 1px solid #3a3a3c;
  border-radius: 6px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.5);
  z-index: 100;
  max-height: 200px;
  overflow-y: auto;
  margin-top: 4px;
}
.type-item {
  padding: 10px 16px;
  font-size: 14px;
  color: #eeeeee;
  cursor: pointer;
  display: flex;
  align-items: center;
}
.type-item:hover {
  background-color: #3a3a3c;
}

/* Switch */
.switch-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #2c2c2e;
  padding: 0; /* Reset for new layout */
  margin-bottom: 0;
}
.switch-left {
  flex: 1;
}
.switch-title {
  font-size: 15px;
  font-weight: 500;
  color: #eeeeee;
  margin-bottom: 4px;
}
.switch-desc {
  font-size: 12px;
  color: #888888;
}
.toggle-switch {
  position: relative;
  display: inline-block;
  width: 44px;
  height: 24px;
}
.toggle-switch input {
  opacity: 0;
  width: 0;
  height: 0;
}
.slider {
  position: absolute;
  cursor: pointer;
  top: 0; left: 0; right: 0; bottom: 0;
  background-color: #444;
  transition: .4s;
  border-radius: 24px;
}
.slider:before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background-color: #eeeeee;
  transition: .4s;
  border-radius: 50%;
}
input:checked + .slider {
  background-color: #3b82f6;
}
input:checked + .slider:before {
  transform: translateX(20px);
}

/* Footer Actions */
.editor-footer {
  padding: 16px 0 60px 0; /* Reduce top padding to bring it closer */
  margin-top: 16px;
  width: 100%;
  background-color: transparent;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.footer-actions-left {
  display: flex;
  align-items: center;
  gap: 16px;
}
.text-btn {
  background: none;
  border: none;
  font-size: 14px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0; /* Force no padding so it's perfectly flush */
}
.outline-yellow {
  color: #e6a23c;
}
.outline-red {
  color: #f56c6c;
}
.divider {
  width: 1px;
  height: 14px;
  background-color: #3a3a3c;
}
.footer-actions-right {
  display: flex;
  gap: 12px;
}
.btn-outline {
  background-color: transparent;
  border: 1px solid #444444;
  color: #dddddd;
  padding: 10px 24px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
}
.btn-outline:hover {
  background-color: rgba(255, 255, 255, 0.05);
}
</style>

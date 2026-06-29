<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getTransactionTypes, updateTransactionType, addTransactionType, deleteTransactionType } from '../../service/api'

const router = useRouter()

const transactionTypes = ref<any[]>([])

const fetchTypes = async () => {
  try {
    transactionTypes.value = await getTransactionTypes()
  } catch (error) {
    console.error('Failed to fetch transaction types', error)
  }
}

onMounted(() => {
  fetchTypes()
})

const showModal = ref(false)
const isEdit = ref(false)
const currentType = ref({
  id: 0,
  name: '',
  flow: 'inflow',
  is_exclude: false,
  is_disabled: false
})

const openEdit = (type: any) => {
  isEdit.value = true
  currentType.value = { ...type }
  showModal.value = true
}

const openAdd = () => {
  isEdit.value = false
  currentType.value = {
    id: 0,
    name: '',
    flow: 'inflow',
    is_exclude: false,
    is_disabled: false
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const saveType = async () => {
  if (!currentType.value.name.trim()) {
    alert('请输入收支名称')
    return
  }
  // Check for duplicate name
  const isDuplicate = transactionTypes.value.some(t => t.name === currentType.value.name && t.id !== currentType.value.id)
  if (isDuplicate) {
    alert('已存在相同名称的收支类型，请换一个名称！')
    return
  }

  try {
    if (isEdit.value) {
      await updateTransactionType(currentType.value.id, currentType.value)
    } else {
      await addTransactionType(currentType.value)
    }
    showModal.value = false
    fetchTypes()
  } catch (error) {
    console.error('Failed to save transaction type', error)
  }
}

const deleteType = async () => {
  if (!confirm('确定要删除这个收支类型吗？')) return
  try {
    await deleteTransactionType(currentType.value.id)
    showModal.value = false
    fetchTypes()
  } catch (error) {
    console.error('Failed to delete transaction type', error)
  }
}

const getFlowTagInfo = (flow: string) => {
  if (flow === 'inflow') return { text: '流入', class: 'tag-inflow' }
  if (flow === 'outflow') return { text: '流出', class: 'tag-outflow' }
  return { text: '转账', class: 'tag-transfer' }
}
</script>

<template>
  <div class="page-container subpage">
    <div class="header">
      <button @click="router.back()" class="back-btn">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="15 18 9 12 15 6"></polyline></svg>
      </button>
      <h2>收支管理</h2>
      <button class="add-btn" @click="openAdd">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
        添加
      </button>
    </div>

    <div class="list-container">
      <div 
        v-for="item in transactionTypes" 
        :key="item.id" 
        class="type-card"
        :class="{ 'disabled-card': item.is_disabled }"
        @click="openEdit(item)"
      >
        <div class="type-info">
          <div class="type-name">{{ item.name }} <span v-if="item.is_disabled" class="disabled-badge">已禁用</span></div>
          <div class="tags">
            <span class="tag" :class="getFlowTagInfo(item.flow).class">
              {{ getFlowTagInfo(item.flow).text }}
            </span>
            <span v-if="item.is_exclude" class="tag tag-exclude">不计入</span>
          </div>
        </div>
        <div class="arrow">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>
        </div>
      </div>
    </div>

    <!-- Edit Modal Overlay -->
    <div class="modal-overlay" v-if="showModal" @click.self="closeModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ isEdit ? '编辑收支' : '新建收支' }}</h3>
          <button class="close-btn" @click="closeModal">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
          </button>
        </div>

        <div class="form-group">
          <label>收支名称</label>
          <input type="text" v-model="currentType.name" class="form-input" placeholder="输入名称" />
        </div>

        <div class="form-group">
          <label>收支类型</label>
          <div class="radio-group">
            <div class="radio-item" :class="{ 'active inflow': currentType.flow === 'inflow' }" @click="currentType.flow = 'inflow'">
              <div class="radio-icon">
                <div class="radio-dot"></div>
              </div>
              <div class="radio-text">
                <span class="main">流入</span>
                <span class="sub">账户金额增加</span>
              </div>
            </div>

            <div class="radio-item" :class="{ 'active outflow': currentType.flow === 'outflow' }" @click="currentType.flow = 'outflow'">
              <div class="radio-icon">
                <div class="radio-dot"></div>
              </div>
              <div class="radio-text">
                <span class="main">流出</span>
                <span class="sub">账户金额减少</span>
              </div>
            </div>

            <div class="radio-item" :class="{ 'active transfer': currentType.flow === 'transfer' }" @click="currentType.flow = 'transfer'">
              <div class="radio-icon">
                <div class="radio-dot"></div>
              </div>
              <div class="radio-text">
                <span class="main">转账</span>
                <span class="sub">账户间流转</span>
              </div>
            </div>
          </div>
        </div>

        <div class="switch-group">
          <div class="switch-info">
            <div class="switch-title">不计入总金额</div>
            <div class="switch-desc">开启后该收支不计入统计</div>
          </div>
          <label class="switch">
            <input type="checkbox" v-model="currentType.is_exclude">
            <span class="slider"></span>
          </label>
        </div>

        <div class="switch-group">
          <div class="switch-info">
            <div class="switch-title">禁用此收支</div>
            <div class="switch-desc">禁用后不会出现在记账选项中</div>
          </div>
          <label class="switch">
            <input type="checkbox" v-model="currentType.is_disabled">
            <span class="slider"></span>
          </label>
        </div>

        <div class="modal-footer">
          <button v-if="isEdit" class="btn btn-delete" @click="deleteType">删除</button>
          <button class="btn btn-cancel" @click="closeModal">取消</button>
          <button class="btn btn-save" @click="saveType">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.subpage {
  display: flex;
  flex-direction: column;
  gap: 20px;
  min-height: 100vh;
  background-color: var(--color-bg);
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
}
.header h2 {
  font-size: 18px;
  font-weight: 600;
}
.back-btn {
  background: none;
  border: none;
  color: var(--color-text-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  padding: 4px;
  margin-left: -8px;
}
.add-btn {
  background: rgba(96, 165, 250, 0.2);
  color: #60a5fa;
  border: none;
  padding: 6px 12px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
}

.list-container {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.type-card {
  background: var(--color-card-bg);
  border-radius: var(--radius-md);
  padding: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
  transition: background-color 0.2s;
  border: 1px solid transparent;
}
.type-card:active {
  background: #222;
}
.disabled-card {
  opacity: 0.6;
}

.type-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.type-name {
  font-size: 16px;
  font-weight: 500;
  color: var(--color-text-primary);
  display: flex;
  align-items: center;
  gap: 8px;
}
.disabled-badge {
  font-size: 10px;
  background: #333;
  color: #999;
  padding: 2px 6px;
  border-radius: 4px;
}

.tags {
  display: flex;
  gap: 8px;
}

.tag {
  font-size: 11px;
  padding: 2px 6px;
  border-radius: 4px;
  font-weight: 500;
}
.tag-inflow {
  background: rgba(92, 242, 92, 0.15);
  color: var(--color-income);
}
.tag-outflow {
  background: rgba(255, 107, 107, 0.15);
  color: var(--color-expense);
}
.tag-transfer {
  background: rgba(96, 165, 250, 0.15);
  color: #60a5fa;
}
.tag-exclude {
  background: rgba(255, 255, 255, 0.1);
  color: var(--color-text-secondary);
}

.arrow {
  color: var(--color-text-secondary);
  display: flex;
  align-items: center;
}

/* Modal Overlay */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(4px);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}
.modal-content {
  background: #171717;
  width: 100%;
  max-width: 400px;
  border-radius: var(--radius-lg);
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
  border: 1px solid #2a2a2a;
}
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.modal-header h3 {
  font-size: 18px;
  font-weight: 600;
}
.close-btn {
  background: none;
  border: none;
  color: var(--color-text-secondary);
  cursor: pointer;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.form-group label {
  font-size: 13px;
  color: var(--color-text-secondary);
  font-weight: 500;
}
.form-input {
  background: #111;
  border: 1px solid #333;
  border-radius: 8px;
  padding: 12px 16px;
  color: white;
  font-size: 16px;
  outline: none;
}
.form-input:focus {
  border-color: #60a5fa;
}

.radio-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.radio-item {
  background: #111;
  border: 1px solid #333;
  border-radius: 12px;
  padding: 14px 16px;
  display: flex;
  align-items: center;
  gap: 16px;
  cursor: pointer;
  transition: all 0.2s;
}
.radio-item.active.inflow { border-color: var(--color-income); box-shadow: 0 0 10px rgba(92,242,92,0.1); }
.radio-item.active.outflow { border-color: var(--color-expense); box-shadow: 0 0 10px rgba(255,107,107,0.1); }
.radio-item.active.transfer { border-color: #60a5fa; box-shadow: 0 0 10px rgba(96,165,250,0.1); }

.radio-icon {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: 2px solid #555;
  display: flex;
  align-items: center;
  justify-content: center;
}
.radio-item.active.inflow .radio-icon { border-color: var(--color-income); }
.radio-item.active.outflow .radio-icon { border-color: var(--color-expense); }
.radio-item.active.transfer .radio-icon { border-color: #60a5fa; }

.radio-item.active.inflow .radio-dot { background: var(--color-income); }
.radio-item.active.outflow .radio-dot { background: var(--color-expense); }
.radio-item.active.transfer .radio-dot { background: #60a5fa; }
.radio-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: transparent;
}

.radio-text {
  display: flex;
  justify-content: space-between;
  flex: 1;
  align-items: center;
}
.radio-text .main {
  font-size: 15px;
  font-weight: 500;
  color: white;
}
.radio-item.active.inflow .main { color: var(--color-income); }
.radio-item.active.outflow .main { color: var(--color-expense); }
.radio-item.active.transfer .main { color: #60a5fa; }

.radio-text .sub {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.switch-group {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #111;
  padding: 16px;
  border-radius: 12px;
  border: 1px solid #333;
}
.switch-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.switch-title {
  font-size: 14px;
  font-weight: 500;
}
.switch-desc {
  font-size: 11px;
  color: var(--color-text-secondary);
}

/* Switch Component */
.switch {
  position: relative;
  display: inline-block;
  width: 44px;
  height: 24px;
}
.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}
.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #444;
  transition: .2s;
  border-radius: 24px;
}
.slider:before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: .2s;
  border-radius: 50%;
}
input:checked + .slider {
  background-color: #fff;
}
input:checked + .slider:before {
  transform: translateX(20px);
  background-color: #111;
}

.modal-footer {
  display: flex;
  gap: 12px;
  margin-top: 12px;
}
.btn {
  flex: 1;
  padding: 12px 0;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  border: none;
}
.btn-cancel {
  background: transparent;
  color: white;
  border: 1px solid #444;
}
.btn-delete {
  background: rgba(255, 107, 107, 0.15);
  color: var(--color-expense);
  border: 1px solid rgba(255, 107, 107, 0.3);
}
.btn-save {
  background: #3b82f6; /* Blue like screenshot */
  color: white;
}
</style>

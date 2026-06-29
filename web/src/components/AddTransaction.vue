<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { getCategories, getAccounts, getTransactionTypes } from '../service/api'

const props = defineProps<{
  show: boolean
}>()

const emit = defineEmits(['close', 'submit'])

const form = ref({
  type: '支出', // Default to Name of the outflow
  amount: '',
  categoryId: '',
  accountId: '',
  recordDate: new Date().toISOString().substring(0, 10),
  remark: ''
})

const categories = ref<any[]>([])
const accounts = ref<any[]>([])
const transactionTypes = ref<any[]>([])

const fetchFormDependencies = async () => {
  try {
    const [c, a, t] = await Promise.all([
      getCategories(),
      getAccounts(),
      getTransactionTypes()
    ])
    categories.value = c
    accounts.value = a
    transactionTypes.value = t.filter((x: any) => !x.is_disabled)
    
    // Set default type to first available type
    if (transactionTypes.value.length > 0 && !transactionTypes.value.find(t => t.name === form.value.type)) {
      form.value.type = transactionTypes.value[0].name
    }
  } catch (error) {
    console.error('Failed to load dependencies', error)
  }
}

onMounted(() => {
  fetchFormDependencies()
})

watch(() => props.show, (newVal) => {
  if (newVal) fetchFormDependencies()
})

const currentTypeObj = computed(() => {
  return transactionTypes.value.find(t => t.name === form.value.type) || { flow: 'outflow' }
})

const filteredCategories = computed(() => {
  const flow = currentTypeObj.value.flow
  // 匹配绑定的具体收支类型名称，或者匹配旧的'income'/'expense'，或者如果是支出/收入就全显
  return categories.value.filter(c => {
    if (c.type === form.value.type) return true
    if (flow === 'inflow' && (c.type === 'income' || c.type === '收入')) return true
    if (flow === 'outflow' && (c.type === 'expense' || c.type === '支出')) return true
    if (!c.type) return true // 如果没绑定收支类型，默认都可以选
    return false
  })
})

const groupedCategories = computed(() => {
  const filtered = filteredCategories.value
  const primaries = filtered.filter(c => c.parent_id === 0)
  const result: any[] = []
  
  primaries.forEach(p => {
    result.push({
      ...p,
      children: filtered.filter(c => c.parent_id === p.id)
    })
  })
  
  const handledIds = new Set()
  result.forEach(p => {
    handledIds.add(p.id)
    p.children.forEach((c: any) => handledIds.add(c.id))
  })
  
  const orphans = filtered.filter(c => !handledIds.has(c.id))
  if (orphans.length > 0) {
    result.push({
      id: -1,
      name: '其他',
      children: orphans
    })
  }
  
  return result
})

const submit = () => {
  const amt = parseFloat(form.value.amount)
  if (isNaN(amt) || amt <= 0) {
    alert('金额必须大于 0')
    return
  }
  if (!form.value.accountId) {
    alert('请选择账户')
    return
  }
  
  emit('submit', form.value)
  emit('close')
}
</script>

<template>
  <transition name="slide-up">
    <div class="add-tx-modal" v-if="show">
      <div class="modal-header">
        <button class="cancel-btn" @click="$emit('close')">取消</button>
        <h3>记一笔</h3>
        <button class="submit-btn" @click="submit">保存</button>
      </div>

      <div class="type-selector-wrapper">
        <div class="type-selector">
          <button 
            v-for="t in transactionTypes" 
            :key="t.id" 
            :class="{ active: form.type === t.name }" 
            @click="form.type = t.name"
          >
            {{ t.name }}
          </button>
        </div>
      </div>

      <div class="amount-input">
        <span class="currency">¥</span>
        <input type="number" v-model="form.amount" placeholder="0.00" />
      </div>

      <div class="form-group">
        <label>账户</label>
        <select v-model="form.accountId">
          <option value="">请选择账户</option>
          <option v-for="a in accounts" :key="a.id" :value="a.id">{{ a.name }}</option>
        </select>
      </div>

      <div class="form-group" v-if="groupedCategories.length > 0">
        <label>分类</label>
        <select v-model="form.categoryId">
          <option value="">请选择分类(可选)</option>
          <optgroup v-for="group in groupedCategories" :key="group.id" :label="group.name">
            <option v-if="group.id !== -1" :value="group.id">{{ group.name }}</option>
            <option v-for="child in group.children" :key="child.id" :value="child.id">
              |- {{ child.name }}
            </option>
          </optgroup>
        </select>
      </div>

      <div class="form-group">
        <label>日期</label>
        <input type="date" v-model="form.recordDate" />
      </div>

      <div class="form-group">
        <label>备注</label>
        <input type="text" v-model="form.remark" placeholder="写点什么..." />
      </div>
    </div>
  </transition>
  <transition name="fade">
    <div class="modal-overlay" v-if="show" @click="$emit('close')"></div>
  </transition>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.4);
  z-index: 1000;
}

.add-tx-modal {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: var(--color-bg);
  border-radius: var(--radius-lg) var(--radius-lg) 0 0;
  padding: 20px;
  z-index: 1001;
  box-shadow: var(--shadow-lg);
  max-width: 600px;
  margin: 0 auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}
.modal-header h3 {
  font-size: 18px;
  font-weight: 600;
}
.cancel-btn, .submit-btn {
  background: none;
  border: none;
  font-size: 16px;
  cursor: pointer;
}
.cancel-btn {
  color: var(--color-text-secondary);
}
.submit-btn {
  color: var(--color-primary);
  font-weight: bold;
}

.type-selector-wrapper {
  overflow-x: auto;
  margin-bottom: 24px;
}
.type-selector-wrapper::-webkit-scrollbar {
  display: none;
}
.type-selector {
  display: inline-flex;
  background: rgba(255, 255, 255, 0.08);
  border-radius: var(--radius-pill);
  padding: 4px;
  min-width: 100%;
}
.type-selector button {
  flex: 1;
  min-width: 70px;
  border: none;
  background: transparent;
  padding: 8px 12px;
  border-radius: var(--radius-pill);
  font-size: 14px;
  cursor: pointer;
  color: #888;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  white-space: nowrap;
}
.type-selector button.active {
  background: #3b82f6;
  color: white;
  font-weight: bold;
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.4);
}

.amount-input {
  display: flex;
  align-items: center;
  border-bottom: 2px solid var(--color-primary);
  padding-bottom: 8px;
  margin-bottom: 24px;
}
.currency {
  font-size: 32px;
  font-weight: bold;
  margin-right: 8px;
}
.amount-input input {
  border: none;
  background: transparent;
  font-size: 40px;
  font-weight: bold;
  width: 100%;
  outline: none;
  color: var(--color-text-primary);
}

.form-group {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  background: rgba(255, 255, 255, 0.05);
  padding: 12px 16px;
  border-radius: var(--radius-md);
  border: 1px solid rgba(255, 255, 255, 0.05);
}
.form-group label {
  width: 60px;
  color: #aaa;
}
.form-group input, .form-group select {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 16px;
  color: var(--color-text-primary);
  outline: none;
}
.form-group select option {
  background: #1e1e1e;
  color: white;
}

/* Animations */
.slide-up-enter-active, .slide-up-leave-active {
  transition: transform 0.3s ease;
}
.slide-up-enter-from, .slide-up-leave-to {
  transform: translateY(100%);
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>

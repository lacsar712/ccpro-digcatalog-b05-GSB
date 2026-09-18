<template>
  <div>
    <div class="toolbar">
      <div>
        <h2 class="page-title">出土文物</h2>
        <p class="page-sub">登记器物类型、材质、完整度与存放位置</p>
      </div>
      <button class="btn" @click="openCreate">新增文物</button>
    </div>

    <div class="card">
      <div class="filters">
        <label>
          探方筛选
          <select v-model="filterUnitId" @change="load">
            <option value="">全部探方</option>
            <option v-for="u in units" :key="u.id" :value="String(u.id)">
              {{ u.site?.name || '' }} / {{ u.code }}
            </option>
          </select>
        </label>
        <label>
          器物类型
          <select v-model="filterType" @change="load">
            <option value="">全部类型</option>
            <option v-for="t in artifactTypes" :key="t" :value="t">{{ t }}</option>
          </select>
        </label>
      </div>

      <table class="table">
        <thead>
          <tr>
            <th>登记号</th>
            <th>探方</th>
            <th>器物类型</th>
            <th>材质</th>
            <th>完整度</th>
            <th>出土日期</th>
            <th>存放位置</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in list" :key="item.id">
            <td>{{ item.registerNo }}</td>
            <td>{{ item.unit?.code || '-' }}</td>
            <td><span class="tag">{{ item.artifactType }}</span></td>
            <td>{{ item.materialName || item.material?.name || '-' }}</td>
            <td>{{ item.completeness || '-' }}</td>
            <td>{{ formatDate(item.findDate) }}</td>
            <td>{{ item.storageLoc || '-' }}</td>
            <td>
              <button class="btn secondary small" @click="openEdit(item)">编辑</button>
              <button class="btn secondary small" @click="openDrawer(item)">轨迹</button>
              <button class="btn danger small" @click="remove(item)">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
      <p v-if="!list.length" class="page-sub">暂无数据</p>
      <p v-if="error" class="error">{{ error }}</p>
    </div>

    <div v-if="showModal" class="modal-mask" @click.self="showModal = false">
      <div class="modal">
        <h3>{{ form.id ? '编辑文物' : '新增文物' }}</h3>
        <div class="form-grid">
          <label>
            所属探方
            <select v-model.number="form.unitId">
              <option :value="0" disabled>请选择</option>
              <option v-for="u in units" :key="u.id" :value="u.id">
                {{ u.site?.name || '' }} / {{ u.code }}
              </option>
            </select>
          </label>
          <label>
            登记号
            <input v-model="form.registerNo" />
          </label>
          <label>
            器物类型
            <select v-model="form.artifactType">
              <option v-for="t in artifactTypes" :key="t" :value="t">{{ t }}</option>
            </select>
          </label>
          <label>
            材质
            <select v-model="form.materialId">
              <option :value="null">未指定</option>
              <option v-for="m in materials" :key="m.id" :value="m.id">{{ m.name }}</option>
            </select>
          </label>
          <label>
            完整度
            <select v-model="form.completeness">
              <option>完整</option>
              <option>残缺</option>
              <option>碎片</option>
            </select>
          </label>
          <label>
            出土日期
            <input v-model="form.findDate" type="date" />
          </label>
          <label v-if="completenessChanged" class="full">
            完整度变更备注（可选）
            <input v-model="form.completenessNote" placeholder="例如：修复后由残缺改为完整" />
          </label>
          <label class="full">
            存放位置
            <input v-model="form.storageLoc" />
          </label>
          <label class="full">
            描述
            <textarea v-model="form.description" />
          </label>
        </div>
        <p v-if="formError" class="error">{{ formError }}</p>
        <div class="modal-actions">
          <button class="btn secondary" @click="showModal = false">取消</button>
          <button class="btn" @click="save">保存</button>
        </div>
      </div>
    </div>

    <div v-if="drawerOpen" class="drawer-mask" @click.self="closeDrawer">
      <aside class="drawer">
        <div class="drawer-head">
          <div>
            <h3>完整度变更轨迹</h3>
            <p class="page-sub">
              {{ drawerFind?.registerNo }} · 当前完整度：{{ drawerFind?.completeness || '-' }}
            </p>
          </div>
          <div class="drawer-head-actions">
            <button class="btn secondary small" @click="openEdit(drawerFind)">编辑</button>
            <button class="btn secondary small" @click="closeDrawer">关闭</button>
          </div>
        </div>
        <p v-if="logsLoading" class="page-sub">加载中…</p>
        <p v-else-if="logsError" class="error">{{ logsError }}</p>
        <p v-else-if="!logs.length" class="page-sub">暂无变更记录</p>
        <div v-else class="timeline">
          <div v-for="log in logs" :key="log.id" class="timeline-item">
            <span class="timeline-dot"></span>
            <div class="timeline-head">
              <span class="tag">{{ log.fromValue || '空' }} → {{ log.toValue || '空' }}</span>
              <span class="timeline-time">{{ formatDateTime(log.changedAt) }}</span>
            </div>
            <div class="timeline-meta">
              操作人：{{ log.operator?.username || `#${log.operatorId}` }}
            </div>
            <div v-if="log.note" class="timeline-note">{{ log.note }}</div>
          </div>
        </div>
      </aside>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import api from '../api/http'

const artifactTypes = ['陶片', '青铜器', '骨器', '玉器', '石器', '铁器', '其他']
const list = ref([])
const units = ref([])
const materials = ref([])
const filterUnitId = ref('')
const filterType = ref('')
const error = ref('')
const formError = ref('')
const showModal = ref(false)

const form = reactive({
  id: null,
  unitId: 0,
  materialId: null,
  registerNo: '',
  artifactType: '陶片',
  completeness: '完整',
  completenessNote: '',
  findDate: '',
  description: '',
  storageLoc: ''
})

const originalCompleteness = ref('')
const drawerOpen = ref(false)
const drawerFind = ref(null)
const logs = ref([])
const logsLoading = ref(false)
const logsError = ref('')

// 编辑态下完整度与打开时的原值不一致，才算一次变更（用于显示备注输入框）
const completenessChanged = computed(
  () => !!form.id && form.completeness !== originalCompleteness.value
)

function formatDate(v) {
  if (!v) return '-'
  return String(v).slice(0, 10)
}

function formatDateTime(v) {
  if (!v) return '-'
  const d = new Date(v)
  if (Number.isNaN(d.getTime())) return String(v)
  return d.toLocaleString('zh-CN', { hour12: false })
}

async function loadMeta() {
  const [u, m] = await Promise.all([api.get('/units'), api.get('/materials')])
  units.value = u.data
  materials.value = m.data
}

async function load() {
  error.value = ''
  try {
    const params = {}
    if (filterUnitId.value) params.unitId = filterUnitId.value
    if (filterType.value) params.artifactType = filterType.value
    const { data } = await api.get('/finds', { params })
    list.value = data
  } catch (e) {
    error.value = e.response?.data?.error || '加载失败'
  }
}

function openCreate() {
  Object.assign(form, {
    id: null,
    unitId: units.value[0]?.id || 0,
    materialId: materials.value[0]?.id ?? null,
    registerNo: '',
    artifactType: '陶片',
    completeness: '完整',
    completenessNote: '',
    findDate: '',
    description: '',
    storageLoc: ''
  })
  originalCompleteness.value = ''
  formError.value = ''
  showModal.value = true
}

function openEdit(item) {
  Object.assign(form, {
    id: item.id,
    unitId: item.unitId,
    materialId: item.materialId,
    registerNo: item.registerNo,
    artifactType: item.artifactType,
    completeness: item.completeness || '完整',
    completenessNote: '',
    findDate: formatDate(item.findDate) === '-' ? '' : formatDate(item.findDate),
    description: item.description || '',
    storageLoc: item.storageLoc || ''
  })
  originalCompleteness.value = item.completeness || '完整'
  formError.value = ''
  showModal.value = true
}

async function openDrawer(item) {
  drawerFind.value = item
  drawerOpen.value = true
  await loadLogs()
}

function closeDrawer() {
  drawerOpen.value = false
  drawerFind.value = null
  logs.value = []
  logsError.value = ''
}

async function loadLogs() {
  if (!drawerFind.value) return
  logsLoading.value = true
  logsError.value = ''
  try {
    const { data } = await api.get(`/finds/${drawerFind.value.id}/completeness-logs`)
    logs.value = data || []
  } catch (e) {
    logsError.value = e.response?.data?.error || '轨迹加载失败'
  } finally {
    logsLoading.value = false
  }
}

async function save() {
  formError.value = ''
  const savedId = form.id
  try {
    const payload = {
      unitId: form.unitId,
      materialId: form.materialId || null,
      registerNo: form.registerNo,
      artifactType: form.artifactType,
      completeness: form.completeness,
      findDate: form.findDate || null,
      description: form.description,
      storageLoc: form.storageLoc
    }
    if (savedId && completenessChanged.value) {
      payload.completenessNote = form.completenessNote || null
    }
    let saved = null
    if (savedId) {
      const { data } = await api.put(`/finds/${savedId}`, payload)
      saved = data
    } else {
      await api.post('/finds', payload)
    }
    showModal.value = false
    await load()
    // 抽屉打开时，保存成功后立即重新拉取轨迹，新记录即时出现在时间线上
    if (savedId && drawerOpen.value && drawerFind.value?.id === savedId) {
      drawerFind.value = saved
      await loadLogs()
    }
  } catch (e) {
    formError.value = e.response?.data?.error || '保存失败'
  }
}

async function remove(item) {
  if (!confirm(`确认删除文物「${item.registerNo}」？`)) return
  try {
    await api.delete(`/finds/${item.id}`)
    await load()
  } catch (e) {
    alert(e.response?.data?.error || '删除失败')
  }
}

onMounted(async () => {
  await loadMeta()
  await load()
})
</script>

<style scoped>
.filters {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}

.filters label {
  min-width: 200px;
}

.drawer-mask {
  position: fixed;
  inset: 0;
  background: rgba(30, 22, 14, 0.35);
  z-index: 40;
}

.drawer {
  position: absolute;
  top: 0;
  right: 0;
  height: 100%;
  width: min(420px, 92vw);
  background: var(--panel);
  border-left: 1px solid var(--border);
  box-shadow: var(--shadow);
  padding: 1.25rem;
  overflow-y: auto;
}

.drawer-head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

.drawer-head h3 {
  margin: 0 0 0.35rem;
}

.drawer-head-actions {
  display: flex;
  gap: 0.5rem;
  flex-shrink: 0;
}

.timeline {
  position: relative;
  margin-left: 0.5rem;
  padding-left: 1.25rem;
  border-left: 2px solid var(--border);
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.timeline-item {
  position: relative;
}

.timeline-dot {
  position: absolute;
  left: calc(-1.25rem - 6px);
  top: 0.3rem;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: var(--accent);
  border: 2px solid var(--panel);
}

.timeline-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.timeline-time {
  color: var(--muted);
  font-size: 0.82rem;
}

.timeline-meta {
  color: var(--muted);
  font-size: 0.85rem;
  margin-top: 0.3rem;
}

.timeline-note {
  margin-top: 0.35rem;
  padding: 0.45rem 0.6rem;
  background: #f6efe3;
  border-radius: 8px;
  font-size: 0.88rem;
}
</style>

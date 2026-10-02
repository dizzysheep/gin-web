<template>
  <div class="page">
    <div class="page-head">
      <div>
        <h1 class="page-title">友链管理</h1>
        <p class="page-desc">维护友情链接与外部站点</p>
      </div>
      <div class="page-actions">
        <el-button type="primary" @click="openDialog()"><el-icon><Plus /></el-icon>新建友链</el-button>
      </div>
    </div>

    <div class="card filter-card">
      <div class="toolbar">
        <el-input
          v-model="query.name"
          placeholder="按名称搜索"
          clearable
          style="width: 240px"
          @keyup.enter="load"
          @clear="load"
        />
        <el-button type="primary" @click="load"><el-icon><Search /></el-icon>查询</el-button>
      </div>
    </div>

    <div class="card table-card">
      <el-table v-loading="loading" :data="rows">
        <el-table-column prop="id" label="ID" width="72">
          <template #default="{ row }">
            <span class="id-text">#{{ row.id }}</span>
          </template>
        </el-table-column>
        <el-table-column label="站点" min-width="200">
          <template #default="{ row }">
            <div class="site-cell">
              <el-avatar :size="30" :src="row.logo" shape="square">{{ row.name?.charAt(0) }}</el-avatar>
              <div>
                <div class="cell-title">{{ row.name }}</div>
                <div class="cell-slug">{{ row.url }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="180" show-overflow-tooltip>
          <template #default="{ row }">
            <span class="id-text">{{ row.description || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="90" align="center">
          <template #default="{ row }">
            <span class="id-text">{{ row.sort }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <span :class="['badge', row.state === 1 ? 'badge-success' : 'badge-danger']">
              {{ stateText[row.state] }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="更新时间" width="170">
          <template #default="{ row }">
            <span class="id-text">{{ formatTime(row.update_time) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="148" fixed="right" align="right">
          <template #default="{ row }">
            <div class="row-actions">
              <el-button class="table-action table-action-edit" text aria-label="编辑友链" @click="openDialog(row)">
                <el-icon><EditPen /></el-icon>
              </el-button>
              <el-button class="table-action table-action-delete" text aria-label="删除友链" @click="onDelete(row)">
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑友链' : '新建友链'" width="480px">
      <el-form :model="form" label-width="70px">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" maxlength="100" />
        </el-form-item>
        <el-form-item label="URL" required>
          <el-input v-model="form.url" placeholder="https://example.com" />
        </el-form-item>
        <el-form-item label="Logo">
          <el-input v-model="form.logo" placeholder="图标地址，可留空" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="2" maxlength="255" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.state" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="submit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { EditPen, Delete, Plus, Search } from '@element-plus/icons-vue'
import { listLinks, addLink, editLink, deleteLink } from '../api/link'
import { formatTime, stateText } from '../utils/format'

const loading = ref(false)
const rows = ref([])
const dialogVisible = ref(false)
const saving = ref(false)
const query = reactive({ name: '' })

const form = reactive({
  id: 0,
  name: '',
  url: '',
  logo: '',
  description: '',
  sort: 0,
  state: 1
})

async function load() {
  loading.value = true
  try {
    const data = await listLinks(query.name ? { name: query.name } : {})
    rows.value = data.list || []
  } finally {
    loading.value = false
  }
}

function openDialog(row) {
  if (row) {
    Object.assign(form, {
      id: row.id,
      name: row.name,
      url: row.url,
      logo: row.logo,
      description: row.description,
      sort: row.sort,
      state: row.state
    })
  } else {
    Object.assign(form, { id: 0, name: '', url: '', logo: '', description: '', sort: 0, state: 1 })
  }
  dialogVisible.value = true
}

async function submit() {
  if (!form.name || !form.url) {
    ElMessage.warning('请填写名称和 URL')
    return
  }
  saving.value = true
  try {
    const payload = {
      name: form.name,
      url: form.url,
      logo: form.logo,
      description: form.description,
      sort: form.sort,
      state: form.state
    }
    if (form.id) {
      await editLink(form.id, payload)
    } else {
      await addLink(payload)
    }
    ElMessage.success('保存成功')
    dialogVisible.value = false
    load()
  } finally {
    saving.value = false
  }
}

async function onDelete(row) {
  await ElMessageBox.confirm(`确定删除友链「${row.name}」？`, '警告', { type: 'warning' })
  await deleteLink(row.id)
  ElMessage.success('已删除')
  load()
}

onMounted(load)
</script>

<style scoped>
.site-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}
</style>

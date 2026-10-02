<template>
  <div class="page">
    <div class="page-head">
      <div>
        <h1 class="page-title">标签管理</h1>
        <p class="page-desc">维护文章标签，便于内容归类</p>
      </div>
      <div class="page-actions">
        <el-button type="primary" @click="openDialog()"><el-icon><Plus /></el-icon>新建标签</el-button>
      </div>
    </div>

    <div class="card filter-card">
      <el-form inline class="filter-form" @submit.prevent>
        <el-form-item label="名称">
          <el-input
            v-model="query.name"
            placeholder="按名称搜索"
            clearable
            style="width: 200px"
            @keyup.enter="load"
            @clear="load"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="query.state" clearable placeholder="全部" style="width: 120px" @change="load">
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="load"><el-icon><Search /></el-icon>查询</el-button>
          <el-button @click="reset"><el-icon><Refresh /></el-icon>重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="card table-card">
      <el-table v-loading="loading" :data="rows">
        <el-table-column prop="id" label="ID" width="72">
          <template #default="{ row }">
            <span class="id-text">#{{ row.id }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称" min-width="180">
          <template #default="{ row }">
            <span class="tag-plain">{{ row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="article_count" label="文章数" width="100" align="center">
          <template #default="{ row }">
            <span class="id-text">{{ row.article_count }}</span>
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
        <el-table-column prop="update_user" label="更新人" width="120">
          <template #default="{ row }">
            <span class="id-text">{{ row.update_user || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="148" fixed="right" align="right">
          <template #default="{ row }">
            <div class="row-actions">
              <el-button class="table-action table-action-edit" text aria-label="编辑标签" @click="openDialog(row)">
                <el-icon><EditPen /></el-icon>
              </el-button>
              <el-button class="table-action table-action-delete" text aria-label="删除标签" @click="onDelete(row)">
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <div class="pager">
        <el-pagination
          v-model:current-page="query.page"
          v-model:page-size="query.page_size"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="load"
          @current-change="load"
        />
      </div>
    </div>

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑标签' : '新建标签'" width="400px">
      <el-form :model="form" label-width="70px">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" maxlength="100" />
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
import { EditPen, Delete, Plus, Search, Refresh } from '@element-plus/icons-vue'
import { listTags, addTag, editTag, deleteTag } from '../api/tag'
import { formatTime, stateText } from '../utils/format'

const loading = ref(false)
const rows = ref([])
const total = ref(0)
const dialogVisible = ref(false)
const saving = ref(false)

const query = reactive({ page: 1, page_size: 10, name: '', state: null })
const form = reactive({ id: 0, name: '', state: 1 })

async function load() {
  loading.value = true
  try {
    const params = { ...query }
    if (params.name === '') delete params.name
    if (params.state === null) delete params.state
    const data = await listTags(params)
    rows.value = data.list || []
    total.value = data.pager?.total || 0
  } finally {
    loading.value = false
  }
}

function reset() {
  query.name = ''
  query.state = null
  query.page = 1
  load()
}

function openDialog(row) {
  if (row) {
    Object.assign(form, { id: row.id, name: row.name, state: row.state })
  } else {
    Object.assign(form, { id: 0, name: '', state: 1 })
  }
  dialogVisible.value = true
}

async function submit() {
  if (!form.name) {
    ElMessage.warning('请输入标签名称')
    return
  }
  saving.value = true
  try {
    const payload = { name: form.name, state: form.state }
    if (form.id) {
      await editTag(form.id, payload)
    } else {
      await addTag(payload)
    }
    ElMessage.success('保存成功')
    dialogVisible.value = false
    load()
  } finally {
    saving.value = false
  }
}

async function onDelete(row) {
  await ElMessageBox.confirm(`确定删除标签「${row.name}」？`, '警告', { type: 'warning' })
  await deleteTag(row.id)
  ElMessage.success('已删除')
  load()
}

onMounted(load)
</script>

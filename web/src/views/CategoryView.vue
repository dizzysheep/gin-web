<template>
  <div class="page">
    <div class="page-head">
      <div>
        <h1 class="page-title">分类管理</h1>
        <p class="page-desc">维护博客文章的分类结构</p>
      </div>
      <div class="page-actions">
        <el-button type="primary" @click="openDialog()"><el-icon><Plus /></el-icon>新建分类</el-button>
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
        <el-table-column prop="name" label="名称" min-width="180">
          <template #default="{ row }">
            <span class="cell-title">{{ row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="90" align="center">
          <template #default="{ row }">
            <span class="id-text">{{ row.sort }}</span>
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
              <el-button class="table-action table-action-edit" text aria-label="编辑分类" @click="openDialog(row)">
                <el-icon><EditPen /></el-icon>
              </el-button>
              <el-button class="table-action table-action-delete" text aria-label="删除分类" @click="onDelete(row)">
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑分类' : '新建分类'" width="440px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" maxlength="100" />
        </el-form-item>
        <el-form-item label="父分类">
          <el-select v-model="form.parent_id" clearable placeholder="顶级分类" style="width: 100%">
            <el-option
              v-for="c in parentOptions"
              :key="c.id"
              :label="c.name"
              :value="c.id"
            />
          </el-select>
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
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { EditPen, Delete, Plus, Search } from '@element-plus/icons-vue'
import { listCategories, addCategory, editCategory, deleteCategory } from '../api/category'
import { formatTime, stateText } from '../utils/format'

const loading = ref(false)
const rows = ref([])
const dialogVisible = ref(false)
const saving = ref(false)
const query = reactive({ name: '' })

const form = reactive({ id: 0, name: '', parent_id: null, sort: 0, state: 1 })

// 父分类只能选顶级分类，且编辑时排除自身
const parentOptions = computed(() =>
  rows.value.filter((c) => c.parent_id === 0 && c.id !== form.id)
)

async function load() {
  loading.value = true
  try {
    const data = await listCategories(query.name ? { name: query.name } : {})
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
      parent_id: row.parent_id || null,
      sort: row.sort,
      state: row.state
    })
  } else {
    Object.assign(form, { id: 0, name: '', parent_id: null, sort: 0, state: 1 })
  }
  dialogVisible.value = true
}

async function submit() {
  if (!form.name) {
    ElMessage.warning('请输入分类名称')
    return
  }
  saving.value = true
  try {
    const payload = {
      name: form.name,
      parent_id: form.parent_id || 0,
      sort: form.sort,
      state: form.state
    }
    if (form.id) {
      await editCategory(form.id, payload)
    } else {
      await addCategory(payload)
    }
    ElMessage.success('保存成功')
    dialogVisible.value = false
    load()
  } finally {
    saving.value = false
  }
}

async function onDelete(row) {
  await ElMessageBox.confirm(`确定删除分类「${row.name}」？`, '警告', { type: 'warning' })
  await deleteCategory(row.id)
  ElMessage.success('已删除')
  load()
}

onMounted(load)
</script>

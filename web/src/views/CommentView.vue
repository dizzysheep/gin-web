<template>
  <div class="page">
    <div class="page-head">
      <div>
        <h1 class="page-title">评论管理</h1>
        <p class="page-desc">审核与管理读者评论</p>
      </div>
    </div>

    <div class="card filter-card">
      <el-form inline class="filter-form" @submit.prevent>
        <el-form-item label="状态">
          <el-select v-model="query.state" clearable placeholder="全部" style="width: 130px" @change="load">
            <el-option label="待审核" :value="0" />
            <el-option label="已通过" :value="1" />
            <el-option label="已拒绝" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="文章 ID">
          <el-input
            v-model.number="query.article_id"
            placeholder="按文章过滤"
            clearable
            style="width: 160px"
            @keyup.enter="load"
            @clear="load"
          />
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
        <el-table-column prop="article_id" label="文章 ID" width="90">
          <template #default="{ row }">
            <span class="id-text">#{{ row.article_id }}</span>
          </template>
        </el-table-column>
        <el-table-column label="评论内容" min-width="280">
          <template #default="{ row }">
            <div class="cell-title">{{ row.content }}</div>
            <div class="cell-slug">
              {{ row.nickname }}{{ row.email ? ` · ${row.email}` : '' }} · IP {{ row.ip || '-' }}
            </div>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="110">
          <template #default="{ row }">
            <span :class="['badge', commentBadge[row.state]]">{{ commentStateText[row.state] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="时间" width="170">
          <template #default="{ row }">
            <span class="id-text">{{ formatTime(row.create_time) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="148" fixed="right" align="right">
          <template #default="{ row }">
            <div class="row-actions">
              <el-button class="table-action table-action-add" text aria-label="通过评论" :disabled="row.state === 1" @click="onAudit(row, 1)">
                <el-icon><CircleCheck /></el-icon>
              </el-button>
              <el-dropdown trigger="click" @command="(cmd) => onCommand(cmd, row)">
                <el-button class="table-action table-action-more" text aria-label="评论更多操作">
                  <el-icon><MoreFilled /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="reject" :disabled="row.state === 2">拒绝</el-dropdown-item>
                    <el-dropdown-item command="delete" divided class="dropdown-danger">删除</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { CircleCheck, MoreFilled, Search, Refresh } from '@element-plus/icons-vue'
import { adminListComment, auditComment, deleteComment } from '../api/comment'
import { formatTime, commentStateText } from '../utils/format'

const loading = ref(false)
const rows = ref([])
const total = ref(0)

const query = reactive({ page: 1, page_size: 10, state: null, article_id: null })

const commentBadge = { 0: 'badge-warning', 1: 'badge-success', 2: 'badge-danger' }

async function load() {
  loading.value = true
  try {
    const params = { ...query }
    if (params.state === null) delete params.state
    if (!params.article_id) delete params.article_id
    const data = await adminListComment(params)
    rows.value = data.list || []
    total.value = data.pager?.total || 0
  } finally {
    loading.value = false
  }
}

function reset() {
  query.state = null
  query.article_id = null
  query.page = 1
  load()
}

function onCommand(cmd, row) {
  if (cmd === 'reject') return onAudit(row, 2)
  if (cmd === 'delete') return onDelete(row)
}

async function onAudit(row, state) {
  await auditComment(row.id, state)
  ElMessage.success(state === 1 ? '已通过' : '已拒绝')
  load()
}

async function onDelete(row) {
  await ElMessageBox.confirm('确定删除该评论？', '警告', { type: 'warning' })
  await deleteComment(row.id)
  ElMessage.success('已删除')
  load()
}

onMounted(load)
</script>

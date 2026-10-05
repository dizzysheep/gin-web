<template>
  <div class="page">
    <div class="page-head">
      <div>
        <h1 class="page-title">文章管理</h1>
        <p class="page-desc">管理和维护博客文章内容</p>
      </div>
      <div class="page-actions">
        <el-button type="primary" @click="$router.push('/article/add')"><el-icon><Plus /></el-icon>新建文章</el-button>
      </div>
    </div>

    <div class="card filter-card">
      <el-form inline class="filter-form" @submit.prevent>
        <el-form-item label="关键词">
          <el-input
            v-model="query.keyword"
            placeholder="标题 / 简述 / 正文"
            clearable
            style="width: 220px"
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
        <el-form-item class="filter-actions">
          <el-button type="primary" @click="load"><el-icon><Search /></el-icon>查询</el-button>
          <el-button @click="reset"><el-icon><Refresh /></el-icon>重置</el-button>
          <el-button text class="filter-toggle" @click="showAdvanced = !showAdvanced">
            <el-icon><component :is="showAdvanced ? ArrowUp : ArrowDown" /></el-icon>
            {{ showAdvanced ? '收起筛选' : '更多筛选' }}
          </el-button>
        </el-form-item>
        <el-collapse-transition>
          <div v-show="showAdvanced" class="advanced-filters">
            <el-form-item label="草稿">
              <el-select v-model="query.is_draft" clearable placeholder="全部" style="width: 120px" @change="load">
                <el-option label="草稿" :value="1" />
                <el-option label="非草稿" :value="0" />
              </el-select>
            </el-form-item>
            <el-form-item label="分类">
              <el-select v-model="query.category_id" clearable filterable placeholder="全部分类" style="width: 160px" @change="load">
                <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id" />
              </el-select>
            </el-form-item>
            <el-form-item label="发布时间">
              <el-date-picker
                v-model="query.published_range"
                type="daterange"
                value-format="YYYY-MM-DD"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                clearable
                unlink-panels
                style="width: 250px"
              />
            </el-form-item>
          </div>
        </el-collapse-transition>
      </el-form>
    </div>

    <div class="card table-card">
      <el-table v-loading="loading" :data="rows">
        <el-table-column prop="id" label="ID" width="72">
          <template #default="{ row }">
            <span class="id-text">#{{ row.id }}</span>
          </template>
        </el-table-column>

        <el-table-column label="标题" min-width="260">
          <template #default="{ row }">
            <div class="cell-title">
              {{ row.title }}
              <span v-if="row.is_top === 1" class="badge badge-info" style="margin-left: 6px">置顶</span>
            </div>
            <div class="cell-slug">{{ row.slug || '#' + row.id }}</div>
          </template>
        </el-table-column>

        <el-table-column label="分类" width="120">
          <template #default="{ row }">
            <span>{{ row.category?.name || '-' }}</span>
          </template>
        </el-table-column>

        <el-table-column label="标签" width="200">
          <template #default="{ row }">
            <span v-if="row.tags?.length" class="tag-list">
              <span v-for="t in row.tags" :key="t.id" class="tag-plain">{{ t.name }}</span>
            </span>
            <span v-else class="id-text">-</span>
          </template>
        </el-table-column>

        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <span v-if="row.is_draft === 1" class="badge badge-warning">草稿</span>
            <span v-else-if="row.state === 1" class="badge badge-success">启用</span>
            <span v-else class="badge badge-danger">禁用</span>
          </template>
        </el-table-column>

        <el-table-column prop="view_count" label="浏览" width="90" align="right">
          <template #default="{ row }">
            <span class="id-text">{{ row.view_count }}</span>
          </template>
        </el-table-column>

        <el-table-column label="发布时间" width="160">
          <template #default="{ row }">
            <span class="id-text">{{ formatTime(row.published_on) }}</span>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="148" fixed="right" align="right">
          <template #default="{ row }">
            <div class="row-actions">
              <el-button class="table-action table-action-edit" text aria-label="编辑文章" @click="$router.push(`/article/edit/${row.id}`)">
                <el-icon><EditPen /></el-icon>
              </el-button>
              <el-dropdown trigger="click" @command="(cmd) => onCommand(cmd, row)">
                <el-button class="table-action table-action-more" text aria-label="文章更多操作">
                  <el-icon><MoreFilled /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="publish">
                      {{ row.is_draft === 1 ? '发布' : '转草稿' }}
                    </el-dropdown-item>
                    <el-dropdown-item command="top">
                      {{ row.is_top === 1 ? '取消置顶' : '置顶' }}
                    </el-dropdown-item>
                    <el-dropdown-item command="state">
                      {{ row.state === 1 ? '禁用' : '启用' }}
                    </el-dropdown-item>
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
import { ArrowDown, ArrowUp, EditPen, MoreFilled, Plus, Search, Refresh } from '@element-plus/icons-vue'
import {
  adminListArticle, publishArticle, setArticleState, setArticleTop, deleteArticle
} from '../api/article'
import { listCategories } from '../api/category'
import { formatTime } from '../utils/format'

const loading = ref(false)
const rows = ref([])
const total = ref(0)
const categories = ref([])
const showAdvanced = ref(false)

const query = reactive({
  page: 1,
  page_size: 10,
  keyword: '',
  state: null,
  is_draft: null,
  category_id: null,
  published_range: null
})

async function load() {
  loading.value = true
  try {
    const params = { ...query }
    // 空筛选不传，保持后端 omitempty 语义
    if (params.keyword === '') delete params.keyword
    if (params.state === null) delete params.state
    if (params.is_draft === null) delete params.is_draft
    if (!params.category_id) delete params.category_id
    if (params.published_range?.length === 2) {
      params.published_from = params.published_range[0]
      params.published_to = params.published_range[1]
    }
    delete params.published_range
    const data = await adminListArticle(params)
    rows.value = data.list || []
    total.value = data.pager?.total || 0
  } finally {
    loading.value = false
  }
}

function reset() {
  query.keyword = ''
  query.state = null
  query.is_draft = null
  query.category_id = null
  query.published_range = null
  showAdvanced.value = false
  query.page = 1
  load()
}

async function onCommand(cmd, row) {
  if (cmd === 'publish') return onPublish(row)
  if (cmd === 'top') return onTop(row)
  if (cmd === 'state') return onState(row)
  if (cmd === 'delete') return onDelete(row)
}

async function onPublish(row) {
  await publishArticle(row.id, row.is_draft === 1)
  ElMessage.success(row.is_draft === 1 ? '已发布' : '已转为草稿')
  await load()
}

async function onTop(row) {
  await setArticleTop(row.id, row.is_top === 1 ? 0 : 1)
  ElMessage.success(row.is_top === 1 ? '已取消置顶' : '已置顶')
  load()
}

async function onState(row) {
  await setArticleState(row.id, row.state === 1 ? 0 : 1)
  ElMessage.success(row.state === 1 ? '已禁用' : '已启用')
  load()
}

async function onDelete(row) {
  await ElMessageBox.confirm(`确定删除文章「${row.title}」？`, '警告', { type: 'warning' })
  await deleteArticle(row.id)
  ElMessage.success('已删除')
  load()
}

onMounted(async () => {
  load()
  listCategories({}).then((d) => { categories.value = d.list || [] }).catch(() => {})
})
</script>

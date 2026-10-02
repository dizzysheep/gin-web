<template>
  <div class="page">
    <div class="page-head">
      <div>
        <h1 class="page-title">{{ id ? '编辑文章' : '新建文章' }}</h1>
        <p class="page-desc">填写文章内容并保存，支持 Markdown 编辑与实时预览</p>
      </div>
      <div class="page-actions">
        <el-button @click="$router.push('/article')">返回</el-button>
      </div>
    </div>

    <div class="card form-card">
    <el-form :model="form" label-width="80px">
      <el-row :gutter="16">
        <el-col :span="16">
          <el-form-item label="标题" required>
            <el-input v-model="form.title" maxlength="100" show-word-limit placeholder="文章标题" />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="Slug">
            <el-input v-model="form.slug" maxlength="150" placeholder="url 路径，留空自动" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="16">
        <el-col :span="8">
          <el-form-item label="分类">
            <el-select v-model="form.category_id" filterable clearable placeholder="选择分类" style="width: 100%">
              <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="标签">
            <el-select v-model="form.tag_ids" multiple filterable placeholder="选择标签" style="width: 100%">
              <el-option v-for="t in tags" :key="t.id" :label="t.name" :value="t.id" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="4">
          <el-form-item label="草稿">
            <el-switch v-model="form.is_draft" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item label="简述">
        <el-input v-model="form.desc" type="textarea" :rows="2" maxlength="255" show-word-limit placeholder="文章简述" />
      </el-form-item>

      <el-form-item label="封面">
        <el-upload :show-file-list="false" :http-request="onCoverUpload" accept="image/*">
          <img v-if="form.cover" :src="form.cover" class="cover-preview" alt="cover" />
          <el-button v-else>上传封面图</el-button>
        </el-upload>
        <el-button v-if="form.cover" link type="danger" style="margin-left: 8px" @click="form.cover = ''">移除</el-button>
      </el-form-item>

      <el-form-item label="正文" required>
        <div style="width: 100%">
          <el-radio-group v-model="editMode" size="small" style="margin-bottom: 8px">
            <el-radio-button label="edit">编辑</el-radio-button>
            <el-radio-button label="preview">预览</el-radio-button>
          </el-radio-group>
          <el-input
            v-if="editMode === 'edit'"
            v-model="form.content_md"
            type="textarea"
            :rows="18"
            placeholder="支持 Markdown"
            class="md-editor"
          />
          <div v-else class="md-preview" v-html="html" />
        </div>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" :loading="saving" @click="submit">保存</el-button>
        <el-button @click="$router.push('/article')">取消</el-button>
      </el-form-item>
    </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { marked } from 'marked'
import { adminGetArticle, addArticle, editArticle } from '../api/article'
import { listCategories } from '../api/category'
import { listTags } from '../api/tag'
import { uploadImage } from '../api/upload'

const route = useRoute()
const router = useRouter()
const id = Number(route.params.id) || 0

const categories = ref([])
const tags = ref([])
const saving = ref(false)
const editMode = ref('edit')

const form = reactive({
  title: '',
  slug: '',
  desc: '',
  content_md: '',
  cover: '',
  category_id: null,
  tag_ids: [],
  is_draft: false
})

const html = computed(() => marked.parse(form.content_md || ''))

async function onCoverUpload(opt) {
  const fd = new FormData()
  fd.append('file', opt.file)
  const data = await uploadImage(fd)
  form.cover = data.url
}

async function submit() {
  if (!form.title) {
    ElMessage.warning('请输入标题')
    return
  }
  if (!form.content_md) {
    ElMessage.warning('请输入正文内容')
    return
  }
  saving.value = true
  try {
    const payload = {
      title: form.title,
      slug: form.slug,
      desc: form.desc,
      content_md: form.content_md,
      cover: form.cover,
      category_id: form.category_id || 0,
      tag_ids: form.tag_ids,
      is_draft: form.is_draft
    }
    if (id) {
      await editArticle(id, payload)
    } else {
      await addArticle(payload)
    }
    ElMessage.success('保存成功')
    router.push('/article')
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  const [catRes, tagRes] = await Promise.all([
    listCategories({}),
    listTags({ page: 1, page_size: 100 })
  ])
  categories.value = catRes.list || []
  tags.value = tagRes.list || []

  if (id) {
    const d = await adminGetArticle(id)
    Object.assign(form, {
      title: d.title,
      slug: d.slug,
      desc: d.desc,
      content_md: d.content_md,
      cover: d.cover,
      category_id: d.category_id || null,
      tag_ids: (d.tags || []).map((t) => t.id),
      is_draft: d.is_draft === 1
    })
  }
})
</script>

<style scoped>
.md-editor :deep(textarea) {
  font-family: 'SF Mono', Menlo, Consolas, Courier New, monospace;
  font-size: 13px;
  line-height: 1.7;
}
</style>

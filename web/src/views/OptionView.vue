<template>
  <div class="page">
    <div class="page-head">
      <div>
        <h1 class="page-title">站点配置</h1>
        <p class="page-desc">站点配置项由后端 options 表驱动，动态渲染，修改后立即对前台生效</p>
      </div>
    </div>

    <div class="card form-card">
      <el-form v-loading="loading" label-width="160px" class="option-form">
        <el-form-item v-for="item in items" :key="item.key" :label="item.key">
          <el-input v-model="item.value" :placeholder="`请输入 ${item.key}`" />
        </el-form-item>

        <el-empty v-if="items.length === 0 && !loading" description="暂无配置项" :image-size="80" />

        <el-form-item v-if="items.length > 0">
          <el-button type="primary" :loading="saving" :disabled="items.length === 0" @click="submit">
            保存配置
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getOptions, saveOptions } from '../api/option'

const loading = ref(false)
const saving = ref(false)
const items = ref([])

async function load() {
  loading.value = true
  try {
    const data = await getOptions()
    const options = data?.options || {}
    items.value = Object.keys(options)
      .sort()
      .map((key) => ({ key, value: options[key] ?? '' }))
  } finally {
    loading.value = false
  }
}

async function submit() {
  saving.value = true
  try {
    await saveOptions(items.value.map(({ key, value }) => ({ key, value: value ?? '' })))
    ElMessage.success('保存成功')
    load()
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.option-form {
  max-width: 720px;
}
</style>

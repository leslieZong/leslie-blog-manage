<template>
  <el-dialog
    :title="$t('searchDialog.title')"
    :model-value="modelValue"
    @update:model-value="updateModelValue"
    width="600px"
    :close-on-click-modal="false"
    destroy-on-close
    :modal="false"
    modal-penetrable
  >
    <div class="search-dialog-content">
      <el-input v-model="searchText" :placeholder="$t('searchDialog.placeholder')" />
    </div>
    <template #footer>
      <el-button type="primary" @click="handleSearch">
        {{ $t('searchDialog.search') }}
      </el-button>
    </template>
  </el-dialog>
</template>
<script setup lang="ts">
import { ref } from 'vue'

// 定义 props
defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
})

// 定义 emit
const emit = defineEmits(['update:model-value'])

// 定义 updateModelValue 函数
const updateModelValue = (val: boolean) => {
  emit('update:model-value', val)
}
const searchText = ref('')

// 定义 handleSearch 函数
const handleSearch = () => {
  console.log(searchText.value)
  updateModelValue(false)
}
</script>
<style scoped lang="scss">
.search-dialog {
  width: 400px;
}
</style>

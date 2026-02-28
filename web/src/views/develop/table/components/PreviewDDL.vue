<template>
  <n-modal
    v-model:show="showModal"
    preset="card"
    title="DDL 预览"
    :style="{ width: '70%' }"
    :mask-closable="false"
  >
    <div class="ddl-preview">
      <n-code :code="ddlContent" language="sql" :hljs="hljs" word-wrap />
    </div>
    <template #action>
      <n-space justify="end">
        <n-button @click="handleCopy">复制 SQL</n-button>
        <n-button @click="showModal = false">关闭</n-button>
        <n-button type="success" :loading="executing" @click="handleExecute">确认执行</n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<script lang="ts" setup>
  import { ref, computed } from 'vue';
  import { useMessage } from 'naive-ui';
  import hljs from 'highlight.js/lib/core';
  import sql from 'highlight.js/lib/languages/sql';
  import 'highlight.js/styles/vs2015.css';

  hljs.registerLanguage('sql', sql);

  const props = defineProps<{
    show: boolean;
    ddl: string;
    loading?: boolean;
  }>();

  const emit = defineEmits(['update:show', 'execute']);

  const message = useMessage();
  const executing = ref(false);

  const showModal = computed({
    get: () => props.show,
    set: (val) => emit('update:show', val),
  });

  const ddlContent = computed(() => props.ddl || '');

  async function handleCopy() {
    try {
      await navigator.clipboard.writeText(ddlContent.value);
      message.success('已复制到剪贴板');
    } catch {
      message.error('复制失败');
    }
  }

  function handleExecute() {
    emit('execute');
  }
</script>

<style lang="less" scoped>
  .ddl-preview {
    background: #1e1e1e;
    border-radius: 8px;
    padding: 20px;
    max-height: 500px;
    overflow: auto;

    :deep(.n-code) {
      font-size: 14px;
      line-height: 1.6;
      color: #d4d4d4;

      .hljs {
        background: transparent;
        color: #d4d4d4;
      }
    }
  }
</style>

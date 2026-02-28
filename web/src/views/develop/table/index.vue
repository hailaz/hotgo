<template>
  <div>
    <n-card :bordered="false" class="proCard" title="数据表管理">
      <div class="n-layout-page-header">
        <n-space align="center" :size="16">
          <n-form-item label="数据库" :show-feedback="false" label-placement="left">
            <n-select
              v-model:value="dbName"
              :options="dbOptions"
              style="width: 200px"
              placeholder="请选择数据库"
              @update:value="handleDbChange"
            />
          </n-form-item>
          <n-form-item label="表名搜索" :show-feedback="false" label-placement="left">
            <n-input
              v-model:value="searchTable"
              placeholder="输入表名搜索"
              clearable
              style="width: 200px"
              @update:value="handleSearch"
            />
          </n-form-item>
        </n-space>
      </div>

      <n-space class="mt-4 mb-4" align="center">
        <n-button type="primary" @click="handleCreate">
          <template #icon>
            <n-icon><PlusOutlined /></n-icon>
          </template>
          新建数据表
        </n-button>
        <n-button @click="loadTableList">
          <template #icon>
            <n-icon><ReloadOutlined /></n-icon>
          </template>
          刷新
        </n-button>
      </n-space>

      <n-data-table
        :columns="columns"
        :data="tableList"
        :loading="loading"
        :row-key="(row) => row.tableName"
        :scroll-x="1000"
        striped
      />

      <!-- 表结构详情弹窗 -->
      <n-modal
        v-model:show="showViewModal"
        :mask-closable="true"
        preset="card"
        :title="'表结构详情 - ' + viewData.tableName"
        :style="{ width: '80%' }"
      >
        <n-tabs type="line">
          <n-tab-pane name="columns" tab="字段信息">
            <n-data-table
              :columns="viewColumnsCols"
              :data="viewData.columns || []"
              :row-key="(row) => row.name"
              size="small"
              :scroll-x="900"
              striped
            />
          </n-tab-pane>
          <n-tab-pane name="indexes" tab="索引信息">
            <n-data-table
              :columns="viewIndexesCols"
              :data="viewData.indexes || []"
              :row-key="(row) => row.name"
              size="small"
              striped
            />
          </n-tab-pane>
        </n-tabs>
      </n-modal>
    </n-card>
  </div>
</template>

<script lang="ts" setup>
  import { h, onMounted, ref } from 'vue';
  import { NTag, NButton, NSpace, NIcon, useDialog, useMessage } from 'naive-ui';
  import { DbSelect, TableList, TableView, TableDrop } from '@/api/develop/genTable';
  import { useRouter } from 'vue-router';
  import { PlusOutlined, ReloadOutlined } from '@vicons/antd';

  const router = useRouter();
  const message = useMessage();
  const dialog = useDialog();

  const dbName = ref('default');
  const searchTable = ref('');
  const dbOptions = ref<any[]>([]);
  const tableList = ref<any[]>([]);
  const loading = ref(false);

  // 表结构详情
  const showViewModal = ref(false);
  const viewData = ref<any>({});

  // 列表列定义
  const columns = [
    {
      title: '表名',
      key: 'tableName',
      width: 260,
      render(row) {
        return h('span', { style: 'font-family: monospace; font-weight: 500;' }, row.tableName);
      },
    },
    {
      title: '表注释',
      key: 'tableComment',
      width: 200,
      ellipsis: { tooltip: true },
    },
    {
      title: '存储引擎',
      key: 'engine',
      width: 100,
      render(row) {
        return row.engine
          ? h(NTag, { size: 'small', type: 'info', bordered: false }, { default: () => row.engine })
          : '-';
      },
    },
    {
      title: '数据行数',
      key: 'tableRows',
      width: 100,
      render(row) {
        return h('span', {}, row.tableRows?.toLocaleString() ?? '0');
      },
    },
    {
      title: '字符集',
      key: 'tableCollation',
      width: 160,
      ellipsis: { tooltip: true },
    },
    {
      title: '创建时间',
      key: 'createTime',
      width: 180,
    },
    {
      title: '操作',
      key: 'actions',
      width: 320,
      fixed: 'right' as const,
      render(row) {
        return h(NSpace, { size: 4 }, () => [
          h(
            NButton,
            { size: 'small', quaternary: true, type: 'default', onClick: () => handleView(row) },
            { default: () => '查看结构' }
          ),
          h(
            NButton,
            { size: 'small', quaternary: true, type: 'info', onClick: () => handleEdit(row) },
            { default: () => '编辑结构' }
          ),
          h(
            NButton,
            { size: 'small', quaternary: true, type: 'success', onClick: () => handleGenCode(row) },
            { default: () => '生成代码' }
          ),
          h(
            NButton,
            { size: 'small', quaternary: true, type: 'error', onClick: () => handleDrop(row) },
            { default: () => '删除' }
          ),
        ]);
      },
    },
  ];

  // 表结构详情 - 字段列配置
  const viewColumnsCols = [
    { title: '字段名', key: 'name', width: 160 },
    { title: '数据类型', key: 'dataType', width: 100 },
    { title: '完整类型', key: 'columnType', width: 160 },
    {
      title: '允许空',
      key: 'isNullable',
      width: 80,
      render(row) {
        return h(
          NTag,
          { size: 'small', type: row.isNullable ? 'success' : 'warning', bordered: false },
          { default: () => (row.isNullable ? '是' : '否') }
        );
      },
    },
    { title: '默认值', key: 'defaultValue', width: 120 },
    { title: '注释', key: 'comment', width: 200, ellipsis: { tooltip: true } },
    {
      title: '主键',
      key: 'isPrimaryKey',
      width: 70,
      render(row) {
        return row.isPrimaryKey
          ? h(NTag, { size: 'small', type: 'error', bordered: false }, { default: () => 'PK' })
          : '';
      },
    },
    {
      title: '自增',
      key: 'isAutoInc',
      width: 70,
      render(row) {
        return row.isAutoInc
          ? h(NTag, { size: 'small', type: 'info', bordered: false }, { default: () => 'AI' })
          : '';
      },
    },
  ];

  // 表结构详情 - 索引列配置
  const viewIndexesCols = [
    { title: '索引名', key: 'name', width: 200 },
    {
      title: '索引类型',
      key: 'type',
      width: 120,
      render(row) {
        const typeMap = { UNIQUE: 'warning', FULLTEXT: 'info', INDEX: 'default' };
        return h(
          NTag,
          { size: 'small', type: typeMap[row.type] || 'default', bordered: false },
          { default: () => row.type }
        );
      },
    },
    {
      title: '关联字段',
      key: 'columns',
      render(row) {
        return h(NSpace, { size: 4 }, () =>
          (row.columns || []).map((col) =>
            h(NTag, { size: 'small', bordered: false }, { default: () => col })
          )
        );
      },
    },
  ];

  // 加载数据库列表
  async function loadDbList() {
    try {
      const res = await DbSelect();
      dbOptions.value = res || [];
      if (dbOptions.value.length > 0 && !dbOptions.value.find((d) => d.value === dbName.value)) {
        dbName.value = dbOptions.value[0].value;
      }
    } catch (e) {
      console.error(e);
    }
  }

  // 加载表列表
  async function loadTableList() {
    loading.value = true;
    try {
      const res = await TableList({ dbName: dbName.value, tableName: searchTable.value });
      tableList.value = res.list || [];
    } catch (e) {
      console.error(e);
    } finally {
      loading.value = false;
    }
  }

  function handleDbChange() {
    loadTableList();
  }

  let searchTimer: any = null;
  function handleSearch() {
    clearTimeout(searchTimer);
    searchTimer = setTimeout(() => {
      loadTableList();
    }, 300);
  }

  // 新建数据表
  function handleCreate() {
    router.push({ path: '/develop/table/design', query: { dbName: dbName.value } });
  }

  // 查看表结构
  async function handleView(row) {
    try {
      const res = await TableView({ dbName: dbName.value, tableName: row.tableName });
      viewData.value = res;
      showViewModal.value = true;
    } catch (e) {
      console.error(e);
    }
  }

  // 编辑表结构
  function handleEdit(row) {
    router.push({
      path: '/develop/table/design',
      query: { dbName: dbName.value, tableName: row.tableName, mode: 'edit' },
    });
  }

  // 一键生成代码
  function handleGenCode(row) {
    router.push({
      name: 'develop_code',
      query: { dbName: dbName.value, tableName: row.tableName },
    });
  }

  // 删除表
  function handleDrop(row) {
    dialog.warning({
      title: '危险操作',
      content: `确定要删除表 "${row.tableName}" 吗？此操作不可恢复！`,
      positiveText: '确定删除',
      negativeText: '取消',
      onPositiveClick: async () => {
        try {
          await TableDrop({ dbName: dbName.value, tableName: row.tableName });
          message.success('删除成功');
          await loadTableList();
        } catch (e) {
          console.error(e);
        }
      },
    });
  }

  onMounted(async () => {
    await loadDbList();
    await loadTableList();
  });
</script>

<style lang="less" scoped>
  .n-layout-page-header {
    padding-bottom: 12px;
  }
</style>

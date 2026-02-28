<template>
  <div>
    <n-card :bordered="false" class="proCard">
      <template #header>
        <n-space align="center">
          <n-button text @click="handleBack">
            <template #icon>
              <n-icon><ArrowLeftOutlined /></n-icon>
            </template>
          </n-button>
          <span>{{ isEditMode ? '编辑表结构' : '新建数据表' }}</span>
          <n-tag v-if="isEditMode && formData.tableName" size="small" type="info">
            {{ formData.tableName }}
          </n-tag>
        </n-space>
      </template>

      <n-spin :show="pageLoading">
        <!-- 基本信息 -->
        <n-card title="基本信息" size="small" class="mb-4">
          <n-form ref="baseFormRef" :model="formData" label-placement="left" :label-width="80">
            <n-grid :cols="4" :x-gap="16">
              <n-form-item-gi label="数据库" path="dbName">
                <n-select
                  v-model:value="formData.dbName"
                  :options="dbOptions"
                  placeholder="请选择"
                  :disabled="isEditMode"
                />
              </n-form-item-gi>
              <n-form-item-gi label="表名" path="tableName">
                <n-input-group>
                  <n-input-group-label v-if="tablePrefix">{{ tablePrefix }}</n-input-group-label>
                  <n-input
                    v-model:value="formData.tableName"
                    placeholder="请输入表名"
                    :disabled="isEditMode"
                  />
                </n-input-group>
              </n-form-item-gi>
              <n-form-item-gi label="表注释" path="comment">
                <n-input v-model:value="formData.comment" placeholder="请输入表注释" />
              </n-form-item-gi>
              <n-form-item-gi label="存储引擎" path="engine">
                <n-select
                  v-model:value="formData.engine"
                  :options="engineOptions"
                  placeholder="请选择"
                />
              </n-form-item-gi>
            </n-grid>
          </n-form>
        </n-card>

        <!-- 字段设计 -->
        <n-card title="字段设计" size="small" class="mb-4">
          <template #header-extra>
            <n-space :size="8">
              <n-button size="small" type="info" quaternary @click="openPresetModal"
                >预设字段</n-button
              >
              <n-button size="small" type="primary" @click="addColumn">添加字段</n-button>
            </n-space>
          </template>

          <n-data-table
            :columns="columnDesignCols"
            :data="formData.columns"
            :row-key="(_, index) => index"
            size="small"
            :scroll-x="1200"
            :max-height="450"
          />
        </n-card>

        <!-- 索引设计 -->
        <n-card title="索引设计" size="small" class="mb-4">
          <template #header-extra>
            <n-button size="small" type="primary" @click="addIndex">添加索引</n-button>
          </template>

          <n-data-table
            :columns="indexDesignCols"
            :data="formData.indexes"
            :row-key="(_, index) => index"
            size="small"
          />
          <n-empty
            v-if="!formData.indexes.length"
            description="暂无索引"
            size="small"
            class="py-4"
          />
        </n-card>
      </n-spin>

      <!-- 底部操作栏 -->
      <div class="action-bar">
        <n-space>
          <n-button @click="handleBack">返回列表</n-button>
          <n-button type="info" @click="handlePreviewDDL" :loading="previewLoading">
            预览 SQL
          </n-button>
          <n-button type="success" @click="handleSubmit" :loading="submitLoading">
            {{ isEditMode ? '保存修改' : '执行建表' }}
          </n-button>
          <n-button
            v-if="!isEditMode"
            type="primary"
            @click="handleSubmitAndGen"
            :loading="submitLoading"
          >
            执行并生成代码
          </n-button>
        </n-space>
      </div>
    </n-card>

    <!-- DDL 预览弹窗 -->
    <PreviewDDL v-model:show="showPreview" :ddl="previewDDL" @execute="handleExecuteFromPreview" />

    <!-- 预设字段弹窗 -->
    <n-modal
      v-model:show="showPresetModal"
      preset="card"
      title="选择预设字段"
      :style="{ width: '780px' }"
      :mask-closable="true"
      :segmented="{ content: true, footer: true }"
    >
      <div class="preset-scroll">
        <template v-for="group in presetGroups" :key="group.label">
          <div class="preset-group-title">{{ group.label }}</div>
          <div class="preset-grid">
            <div
              v-for="item in getGroupPresets(group)"
              :key="item.key"
              class="preset-card"
              :class="{
                'preset-card-selected': presetSelected.includes(item.key),
                'preset-card-disabled': isPresetExist(item.key),
              }"
              @click="togglePreset(item.key)"
            >
              <div class="preset-card-header">
                <n-checkbox
                  :checked="presetSelected.includes(item.key)"
                  :disabled="isPresetExist(item.key)"
                  @update:checked="() => togglePreset(item.key)"
                  @click.stop
                />
                <span class="preset-card-title">{{ item.label }}</span>
                <n-tag :type="getTagType(item.tag)" size="tiny" :bordered="false">{{
                  item.tag
                }}</n-tag>
              </div>
              <div class="preset-card-desc">{{ item.description }}</div>
              <div class="preset-card-fields">
                <n-tag
                  v-for="f in item.fields.split(', ')"
                  :key="f"
                  size="tiny"
                  :bordered="false"
                  class="preset-field-tag"
                  >{{ f }}</n-tag
                >
              </div>
              <div v-if="isPresetExist(item.key)" class="preset-card-badge">
                <svg viewBox="0 0 16 16" width="12" height="12" fill="currentColor">
                  <path
                    d="M13.78 4.22a.75.75 0 0 1 0 1.06l-7.25 7.25a.75.75 0 0 1-1.06 0L2.22 9.28a.75.75 0 0 1 1.06-1.06L6 10.94l6.72-6.72a.75.75 0 0 1 1.06 0Z"
                  />
                </svg>
                <span>已添加</span>
              </div>
            </div>
          </div>
        </template>
      </div>
      <template #footer>
        <n-space justify="space-between" align="center">
          <span style="color: #999; font-size: 13px">
            已选 {{ presetSelected.length }} 项，共 {{ presetSelectedFieldCount }} 个字段
          </span>
          <n-space :size="8">
            <n-button @click="showPresetModal = false">取消</n-button>
            <n-button
              type="primary"
              :disabled="presetSelected.length === 0"
              @click="handlePresetBatchAdd"
            >
              添加选中字段
            </n-button>
          </n-space>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
  import { h, onMounted, ref, computed } from 'vue';
  import {
    NButton,
    NIcon,
    NInput,
    NInputNumber,
    NSelect,
    NCheckbox,
    NSpace,
    NTag,
    useMessage,
  } from 'naive-ui';
  import { useRouter, useRoute } from 'vue-router';
  import { ArrowLeftOutlined } from '@vicons/antd';
  import {
    DbSelect,
    TableView,
    TableCreate,
    TableEdit,
    PreviewDDL as PreviewDDLApi,
  } from '@/api/develop/genTable';
  import {
    dataTypeOptions,
    indexTypeOptions,
    engineOptions,
    integerTypes,
    decimalTypes,
    lengthTypes,
    newColumn,
    newIndex,
    presetColumns,
    presetMeta,
    presetGroups,
    getFieldWeight,
  } from './components/model';
  import type { ColumnItem, IndexItem } from './components/model';
  import PreviewDDL from './components/PreviewDDL.vue';

  const router = useRouter();
  const route = useRoute();
  const message = useMessage();

  const isEditMode = ref(false);
  const pageLoading = ref(false);
  const submitLoading = ref(false);
  const previewLoading = ref(false);
  const showPreview = ref(false);
  const previewDDL = ref('');
  const dbOptions = ref<any[]>([]);
  const tablePrefix = ref('');
  const baseFormRef = ref();

  const formData = ref<{
    dbName: string;
    tableName: string;
    comment: string;
    engine: string;
    columns: ColumnItem[];
    indexes: IndexItem[];
  }>({
    dbName: 'default',
    tableName: '',
    comment: '',
    engine: 'InnoDB',
    columns: [],
    indexes: [],
  });

  // 可选字段名列表（用于索引关联字段选择）
  const columnNameOptions = computed(() =>
    formData.value.columns.filter((c) => c.name).map((c) => ({ label: c.name, value: c.name }))
  );

  // 预设字段弹窗
  const showPresetModal = ref(false);
  const presetSelected = ref<string[]>([]);

  function openPresetModal() {
    presetSelected.value = [];
    showPresetModal.value = true;
  }

  // 获取分组下的预设列表
  function getGroupPresets(group: (typeof presetGroups)[number]) {
    return presetMeta.filter((m) => group.tags.includes(m.tag));
  }

  // tag 颜色映射
  function getTagType(tag: string): 'success' | 'info' | 'warning' | 'error' | 'default' {
    const map: Record<string, 'success' | 'info' | 'warning' | 'error' | 'default'> = {
      必备: 'success',
      常用: 'info',
      用户: 'warning',
      关联: 'default',
      结构: 'default',
      审计: 'default',
      内容: 'info',
      媒体: 'warning',
      业务: 'success',
      日志: 'error',
      统计: 'info',
      扩展: 'default',
    };
    return map[tag] || 'default';
  }

  // 判断预设是否已在表中
  function isPresetExist(key: string): boolean {
    const preset = presetColumns[key];
    if (!preset) return false;
    const existing = formData.value.columns.map((c) => c.name);
    return preset.every((p) => existing.includes(p.name));
  }

  // 切换预设选中状态
  function togglePreset(key: string) {
    if (isPresetExist(key)) return;
    const idx = presetSelected.value.indexOf(key);
    if (idx >= 0) {
      presetSelected.value.splice(idx, 1);
    } else {
      presetSelected.value.push(key);
    }
  }

  // 已选预设包含的字段总数
  const presetSelectedFieldCount = computed(() => {
    return presetSelected.value.reduce((sum, key) => {
      const preset = presetColumns[key];
      if (!preset) return sum;
      const existing = formData.value.columns.map((c) => c.name);
      return sum + preset.filter((p) => !existing.includes(p.name)).length;
    }, 0);
  });

  // 批量添加预设字段（添加后按权重自动排序）
  function handlePresetBatchAdd() {
    const existing = formData.value.columns.map((c) => c.name);
    let addedCount = 0;
    for (const key of presetSelected.value) {
      const preset = presetColumns[key];
      if (!preset) continue;
      const toAdd = preset.filter((p) => !existing.includes(p.name));
      for (const p of toAdd) {
        formData.value.columns.push({ ...p });
        existing.push(p.name);
        addedCount++;
      }
    }
    if (addedCount > 0) {
      // 按字段权重排序，保持规整顺序
      formData.value.columns.sort((a, b) => getFieldWeight(a.name) - getFieldWeight(b.name));
      message.success(`已添加 ${addedCount} 个字段`);
    } else {
      message.warning('所选预设字段均已存在');
    }
    presetSelected.value = [];
    showPresetModal.value = false;
  }

  // 字段设计列
  const columnDesignCols = [
    {
      title: '字段名',
      key: 'name',
      width: 140,
      render(row, index) {
        return h(NInput, {
          value: row.name,
          size: 'small',
          placeholder: '字段名',
          onUpdateValue: (v) => (formData.value.columns[index].name = v),
        });
      },
    },
    {
      title: '类型',
      key: 'dataType',
      width: 130,
      render(row, index) {
        return h(NSelect, {
          value: row.dataType,
          size: 'small',
          options: dataTypeOptions,
          filterable: true,
          placeholder: '类型',
          onUpdateValue: (v) => {
            formData.value.columns[index].dataType = v;
            // 根据类型自动调整默认长度
            if (v === 'varchar') formData.value.columns[index].length = 255;
            else if (v === 'bigint') formData.value.columns[index].length = 20;
            else if (v === 'int') formData.value.columns[index].length = 11;
            else if (v === 'tinyint') formData.value.columns[index].length = 4;
          },
        });
      },
    },
    {
      title: '长度',
      key: 'length',
      width: 80,
      render(row, index) {
        return lengthTypes.includes(row.dataType)
          ? h(NInputNumber, {
              value: row.length,
              size: 'small',
              min: 0,
              showButton: false,
              placeholder: '长度',
              onUpdateValue: (v) => (formData.value.columns[index].length = v || 0),
            })
          : h('span', { style: 'color: #999' }, '-');
      },
    },
    {
      title: '小数位',
      key: 'decimal',
      width: 80,
      render(row, index) {
        return decimalTypes.includes(row.dataType)
          ? h(NInputNumber, {
              value: row.decimal,
              size: 'small',
              min: 0,
              showButton: false,
              placeholder: '小数位',
              onUpdateValue: (v) => (formData.value.columns[index].decimal = v || 0),
            })
          : h('span', { style: 'color: #999' }, '-');
      },
    },
    {
      title: '非空',
      key: 'isNullable',
      width: 60,
      render(row, index) {
        return h(NCheckbox, {
          checked: !row.isNullable,
          onUpdateChecked: (v) => (formData.value.columns[index].isNullable = !v),
        });
      },
    },
    {
      title: '默认值',
      key: 'defaultValue',
      width: 120,
      render(row, index) {
        return h(NInput, {
          value: row.defaultValue,
          size: 'small',
          placeholder: '默认值',
          onUpdateValue: (v) => (formData.value.columns[index].defaultValue = v),
        });
      },
    },
    {
      title: '注释',
      key: 'comment',
      width: 160,
      render(row, index) {
        return h(NInput, {
          value: row.comment,
          size: 'small',
          placeholder: '注释',
          onUpdateValue: (v) => (formData.value.columns[index].comment = v),
        });
      },
    },
    {
      title: '主键',
      key: 'isPrimaryKey',
      width: 55,
      render(row, index) {
        return h(NCheckbox, {
          checked: row.isPrimaryKey,
          onUpdateChecked: (v) => (formData.value.columns[index].isPrimaryKey = v),
        });
      },
    },
    {
      title: '自增',
      key: 'isAutoInc',
      width: 55,
      render(row, index) {
        return h(NCheckbox, {
          checked: row.isAutoInc,
          onUpdateChecked: (v) => (formData.value.columns[index].isAutoInc = v),
        });
      },
    },
    {
      title: '无符号',
      key: 'isUnsigned',
      width: 65,
      render(row, index) {
        return integerTypes.includes(row.dataType)
          ? h(NCheckbox, {
              checked: row.isUnsigned,
              onUpdateChecked: (v) => (formData.value.columns[index].isUnsigned = v),
            })
          : h('span', { style: 'color: #999' }, '-');
      },
    },
    {
      title: '操作',
      key: 'actions',
      width: 120,
      fixed: 'right' as const,
      render(_, index) {
        return h(NSpace, { size: 2 }, () => [
          h(
            NButton,
            {
              size: 'tiny',
              quaternary: true,
              disabled: index === 0,
              onClick: () => moveColumn(index, -1),
            },
            { default: () => '↑' }
          ),
          h(
            NButton,
            {
              size: 'tiny',
              quaternary: true,
              disabled: index === formData.value.columns.length - 1,
              onClick: () => moveColumn(index, 1),
            },
            { default: () => '↓' }
          ),
          h(
            NButton,
            { size: 'tiny', quaternary: true, type: 'error', onClick: () => removeColumn(index) },
            { default: () => '删除' }
          ),
        ]);
      },
    },
  ];

  // 索引设计列
  const indexDesignCols = [
    {
      title: '索引名',
      key: 'name',
      width: 200,
      render(row, index) {
        return h(NInput, {
          value: row.name,
          size: 'small',
          placeholder: '留空自动生成',
          onUpdateValue: (v) => (formData.value.indexes[index].name = v),
        });
      },
    },
    {
      title: '索引类型',
      key: 'type',
      width: 200,
      render(row, index) {
        return h(NSelect, {
          value: row.type,
          size: 'small',
          options: indexTypeOptions,
          onUpdateValue: (v) => (formData.value.indexes[index].type = v),
        });
      },
    },
    {
      title: '关联字段',
      key: 'columns',
      render(row, index) {
        return h(NSelect, {
          value: row.columns,
          size: 'small',
          options: columnNameOptions.value,
          multiple: true,
          placeholder: '选择字段',
          onUpdateValue: (v) => (formData.value.indexes[index].columns = v),
        });
      },
    },
    {
      title: '操作',
      key: 'actions',
      width: 80,
      render(_, index) {
        return h(
          NButton,
          { size: 'small', quaternary: true, type: 'error', onClick: () => removeIndex(index) },
          { default: () => '删除' }
        );
      },
    },
  ];

  // 字段操作
  function addColumn() {
    formData.value.columns.push(newColumn());
  }

  function removeColumn(index: number) {
    formData.value.columns.splice(index, 1);
  }

  function moveColumn(index: number, direction: number) {
    const target = index + direction;
    if (target < 0 || target >= formData.value.columns.length) return;
    const temp = formData.value.columns[index];
    formData.value.columns[index] = formData.value.columns[target];
    formData.value.columns[target] = temp;
    formData.value.columns = [...formData.value.columns];
  }

  // 预设字段（单个添加）
  function handlePresetSelect(key: string) {
    const preset = presetColumns[key];
    if (preset) {
      const existing = formData.value.columns.map((c) => c.name);
      const toAdd = preset.filter((p) => !existing.includes(p.name));
      if (toAdd.length === 0) {
        message.warning('预设字段已存在');
        return;
      }
      formData.value.columns.push(...toAdd.map((p) => ({ ...p })));
      formData.value.columns.sort((a, b) => getFieldWeight(a.name) - getFieldWeight(b.name));
      message.success(`已添加 ${toAdd.length} 个字段`);
    }
  }

  // 索引操作
  function addIndex() {
    formData.value.indexes.push(newIndex());
  }

  function removeIndex(index: number) {
    formData.value.indexes.splice(index, 1);
  }

  // 校验表单
  function validateForm(): boolean {
    if (!formData.value.dbName) {
      message.error('请选择数据库');
      return false;
    }
    if (!formData.value.tableName) {
      message.error('请输入表名');
      return false;
    }
    if (!/^[a-zA-Z_][a-zA-Z0-9_]*$/.test(formData.value.tableName)) {
      message.error('表名格式不正确，只允许字母、数字和下划线，且以字母或下划线开头');
      return false;
    }
    if (formData.value.columns.length === 0) {
      message.error('请至少添加一个字段');
      return false;
    }
    for (let i = 0; i < formData.value.columns.length; i++) {
      const col = formData.value.columns[i];
      if (!col.name) {
        message.error(`第 ${i + 1} 个字段名不能为空`);
        return false;
      }
      if (!col.dataType) {
        message.error(`字段 ${col.name} 的数据类型不能为空`);
        return false;
      }
    }
    return true;
  }

  // 预览DDL
  async function handlePreviewDDL() {
    if (!validateForm()) return;
    previewLoading.value = true;
    try {
      const res = await PreviewDDLApi({
        ...formData.value,
        isEdit: isEditMode.value,
      });
      previewDDL.value = res.ddl;
      showPreview.value = true;
    } catch (e) {
      console.error(e);
    } finally {
      previewLoading.value = false;
    }
  }

  // 执行建表/保存修改
  async function handleSubmit() {
    if (!validateForm()) return;
    submitLoading.value = true;
    try {
      if (isEditMode.value) {
        await TableEdit(formData.value);
        message.success('修改成功');
      } else {
        await TableCreate(formData.value);
        message.success('建表成功');
      }
      handleBack();
    } catch (e) {
      console.error(e);
    } finally {
      submitLoading.value = false;
    }
  }

  // 执行并跳转生成代码
  async function handleSubmitAndGen() {
    if (!validateForm()) return;
    submitLoading.value = true;
    try {
      await TableCreate(formData.value);
      message.success('建表成功，正在跳转代码生成...');
      router.push({
        name: 'develop_code',
        query: { dbName: formData.value.dbName, tableName: formData.value.tableName },
      });
    } catch (e) {
      console.error(e);
    } finally {
      submitLoading.value = false;
    }
  }

  // 从预览弹窗执行
  async function handleExecuteFromPreview() {
    await handleSubmit();
    showPreview.value = false;
  }

  function handleBack() {
    router.push({ path: '/develop/table' });
  }

  // 加载数据库列表
  async function loadDbList() {
    try {
      const res = await DbSelect();
      dbOptions.value = res || [];
    } catch (e) {
      console.error(e);
    }
  }

  // 编辑模式 - 加载已有表结构
  async function loadTableView() {
    pageLoading.value = true;
    try {
      const res = await TableView({
        dbName: formData.value.dbName,
        tableName: formData.value.tableName,
      });
      formData.value.comment = res.tableComment || '';
      formData.value.engine = res.engine || 'InnoDB';
      formData.value.columns = (res.columns || []).map((col) => ({
        name: col.name,
        dataType: col.dataType,
        length: col.length || 0,
        decimal: col.decimal || 0,
        isNullable: col.isNullable,
        defaultValue: col.defaultValue || '',
        comment: col.comment || '',
        isPrimaryKey: col.isPrimaryKey,
        isAutoInc: col.isAutoInc,
        isUnsigned: col.isUnsigned,
      }));
      formData.value.indexes = (res.indexes || []).map((idx) => ({
        name: idx.name,
        type: idx.type,
        columns: idx.columns || [],
      }));
    } catch (e) {
      console.error(e);
    } finally {
      pageLoading.value = false;
    }
  }

  onMounted(async () => {
    await loadDbList();

    const query = route.query;
    if (query.dbName) {
      formData.value.dbName = query.dbName as string;
    }

    if (query.mode === 'edit' && query.tableName) {
      isEditMode.value = true;
      formData.value.tableName = query.tableName as string;
      await loadTableView();
    } else {
      // 新建模式：默认添加 id 主键
      handlePresetSelect('id主键');
    }
  });
</script>

<style lang="less" scoped>
  .mb-4 {
    margin-bottom: 16px;
  }

  .py-4 {
    padding: 16px 0;
  }

  .action-bar {
    position: sticky;
    bottom: 0;
    background: #fff;
    padding: 16px 0 0;
    border-top: 1px solid #efeff5;
    margin-top: 8px;
    z-index: 10;
  }

  .preset-scroll {
    max-height: 520px;
    overflow-y: auto;
    padding-right: 4px;
  }

  .preset-group-title {
    font-size: 13px;
    font-weight: 600;
    color: #333;
    margin: 12px 0 8px;
    padding-left: 2px;

    &:first-child {
      margin-top: 0;
    }
  }

  .preset-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 12px;
    padding: 4px 0;
  }

  .preset-card {
    position: relative;
    border: 1px solid #e8e8e8;
    border-radius: 8px;
    padding: 12px;
    cursor: pointer;
    transition: all 0.2s;
    background: #fafafa;

    &:hover {
      border-color: #36ad6a;
      box-shadow: 0 2px 8px rgb(54 173 106 / 10%);
    }

    &-selected {
      border-color: #36ad6a;
      background: #f0faf4;
      box-shadow: 0 0 0 1px #36ad6a;
    }

    &-disabled {
      cursor: not-allowed;
      background: #f6f6f6;
      border-color: #d9d9d9;
      border-style: dashed;

      &:hover {
        border-color: #d9d9d9;
        box-shadow: none;
      }

      .preset-card-title {
        color: #999;
      }

      .preset-card-desc {
        color: #bbb;
      }

      .preset-field-tag {
        opacity: 0.5;
      }
    }

    &-header {
      display: flex;
      align-items: center;
      gap: 6px;
      margin-bottom: 6px;
    }

    &-title {
      font-weight: 500;
      font-size: 13px;
      flex: 1;
    }

    &-desc {
      font-size: 12px;
      color: #666;
      margin-bottom: 8px;
      line-height: 1.4;
    }

    &-fields {
      display: flex;
      flex-wrap: wrap;
      gap: 4px;
    }

    &-badge {
      position: absolute;
      top: 8px;
      right: 8px;
      display: inline-flex;
      align-items: center;
      gap: 3px;
      padding: 2px 8px;
      border-radius: 10px;
      font-size: 11px;
      font-weight: 500;
      color: #36ad6a;
      background: #e8f8ef;
    }
  }

  .preset-field-tag {
    font-family: SFMono-Regular, Consolas, monospace;
    font-size: 11px !important;
  }
</style>

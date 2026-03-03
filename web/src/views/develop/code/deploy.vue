<template>
  <div>
    <n-spin :show="show" description="正在生成配置信息...">
      <n-card class="proCard">
        <n-tabs
          type="card"
          class="card-tabs"
          :default-value="value"
          animated
          tab-style="min-width: 80px;"
          style="margin: 0 -4px"
          pane-style="padding-left: 4px; padding-right: 4px; box-sizing: border-box;"
          :on-update:value="updateTabs"
          ref="tabsRef"
          @close="handleClose"
          @add="handleAdd"
        >
          <n-tab-pane v-for="panel in panels" :key="panel" :name="panel">
            <template v-if="panel === '基本信息'">
              <BaseInfo v-model:value="genInfo" :selectList="selectList" />
            </template>

            <template v-if="panel === '主表字段'">
              <EditMasterCell v-model:value="genInfo" :selectList="selectList" />
            </template>
          </n-tab-pane>

          <n-tab-pane
            v-for="panel in slavePanels"
            :key="panel"
            :name="panel"
            v-show="slavePanels.length > 0"
          >
            <EditSlaveCell
              v-model:value="genInfo"
              :uuid="slaveMap[panel]"
              :selectList="selectList"
            />
          </n-tab-pane>

          <template #suffix>
            <n-space>
              <n-button type="default" @click="handleBack">返回列表</n-button>
              <n-button type="primary" :loading="formBtnPreviewLoading" @click="preview"
                >预览代码</n-button
              >
              <n-button type="success" :loading="formBtnLoading" @click="submitBuild"
                >提交生成</n-button
              >
              <n-button type="info" dashed @click="submitSave">仅保存配置</n-button>
              <n-button type="error" ghost :loading="cleanBtnLoading" @click="handleClean"
                >清除代码</n-button
              >
            </n-space>
          </template>
        </n-tabs>

        <n-modal
          v-model:show="showModal"
          :block-scroll="false"
          :mask-closable="false"
          :show-icon="false"
          preset="card"
          title="预览代码"
          style="width: 95%"
        >
          <PreviewTab :previewModel="previewModel" />
          <template #action>
            <n-space justify="end">
              <n-button @click="() => (showModal = false)">关闭</n-button>
              <n-button type="info" :loading="formBtnLoading" @click="submitBuild"
                >提交生成</n-button
              >
            </n-space>
          </template>
        </n-modal>

        <n-modal
          v-model:show="showCleanModal"
          :block-scroll="false"
          :mask-closable="false"
          :show-icon="false"
          preset="card"
          title="清除生成的代码文件"
          style="width: 720px"
        >
          <n-alert type="warning" style="margin-bottom: 12px">
            请勾选需要删除的文件，确认后将从磁盘上删除这些文件，此操作不可恢复！
          </n-alert>
          <n-data-table
            :columns="cleanColumns"
            :data="cleanFileList"
            :row-key="(row) => row.path"
            v-model:checked-row-keys="cleanCheckedKeys"
            size="small"
            :max-height="400"
          />
          <n-checkbox v-model:checked="cleanMenu" style="margin-top: 12px">
            同时清除数据库中对应的菜单权限
          </n-checkbox>
          <template #action>
            <n-space justify="end">
              <n-button @click="() => (showCleanModal = false)">取消</n-button>
              <n-button
                type="error"
                :loading="cleanSubmitLoading"
                :disabled="cleanCheckedKeys.length === 0"
                @click="submitClean"
              >
                确认删除 ({{ cleanCheckedKeys.length }})
              </n-button>
            </n-space>
          </template>
        </n-modal>
      </n-card>
    </n-spin>
  </div>
</template>

<script lang="ts" setup>
  import { onMounted, ref, h, watch } from 'vue';
  import { useRouter } from 'vue-router';
  import { useDialog, useMessage, useNotification, NTag, DataTableColumns } from 'naive-ui';
  import BaseInfo from './components/BaseInfo.vue';
  import EditMasterCell from './components/EditMasterCell.vue';
  import EditSlaveCell from './components/EditSlaveCell.vue';
  import { Selects, View, Preview, Build, Edit, Clean } from '@/api/develop/code';
  import { selectListObj, newState, formatColumns } from '@/views/develop/code/components/model';
  import PreviewTab from '@/views/develop/code/components/PreviewTab.vue';
  import { isJsonString } from '@/utils/is';

  interface Props {
    genId?: number;
  }
  const props = withDefaults(defineProps<Props>(), { genId: 0 });
  const router = useRouter();
  const genId = Number(router.currentRoute.value.params.id ?? props.genId);
  const show = ref(false);
  const message = useMessage();
  const selectList = ref<any>(selectListObj);
  const genInfo = ref(newState(null));
  const tabsRef = ref();
  const value = ref('基本信息');
  const panels = ref(['基本信息', '主表字段']);
  const slaveMap = ref<any>([]);
  const slavePanels = ref<any>([]);
  const showModal = ref(false);
  const formBtnLoading = ref(false);
  const formBtnPreviewLoading = ref(false);
  const previewModel = ref<any>();
  const dialog = useDialog();
  const notification = useNotification();

  // 清除代码相关
  const showCleanModal = ref(false);
  const cleanBtnLoading = ref(false);
  const cleanSubmitLoading = ref(false);
  const cleanFileList = ref<any[]>([]);
  const cleanCheckedKeys = ref<string[]>([]);
  const cleanMenu = ref(false);

  const cleanColumns: DataTableColumns<any> = [
    { type: 'selection' },
    { title: '文件类型', key: 'key', width: 140 },
    { title: '文件路径', key: 'path', ellipsis: { tooltip: true } },
    {
      title: '状态',
      key: 'exists',
      width: 100,
      render(row) {
        return row.exists
          ? h(NTag, { type: 'success', size: 'small' }, { default: () => '已存在' })
          : h(NTag, { type: 'default', size: 'small' }, { default: () => '不存在' });
      },
    },
  ];

  onMounted(async () => {
    if (genId < 1 && props.genId < 1) {
      message.error('生成ID不正确，请检查！');
      return;
    }
    await getGenInfo();
    await loadSelect();
  });

  async function getGenInfo() {
    let tmp = await View({ id: genId });
    // 导入主表数据
    tmp.masterColumns = formatColumns(tmp.masterColumns);

    // 导入生成选项
    if (isJsonString(tmp.options)) {
      tmp.options = JSON.parse(tmp.options);
    }

    // 预设流程
    if (!tmp.options.presetStep) {
      tmp.options.presetStep = {
        formGridCols: 1,
      };
    }

    // 树表
    if (!tmp.options.tree) {
      tmp.options.tree = {
        titleColumn: null,
        styleType: 1,
      };
    }

    genInfo.value = tmp;
  }

  watch(
    genInfo,
    (newVal, _oldVal) => {
      if (newVal.genType >= 10 && newVal.genType < 20) {
        handleAdd('主表字段');
      } else {
        handleClose('主表字段');
      }

      if (newVal && newVal.options && newVal.options.join !== undefined) {
        slavePanels.value = [];
        for (let i = 0; i < newVal.options.join.length; i++) {
          if (newVal.options.join[i]?.alias) {
            handleSlaveAdd(
              '关联表[ ' + newVal.options.join[i].alias + ' ]',
              newVal.options.join[i]
            );
          }
        }
      }
    },
    {
      deep: true, // 是否深度监听
    }
  );

  function updateTabs(value: string | number) {
    console.log('value:' + value);
  }

  function handleAdd(name: string) {
    const nameIndex = panels.value.findIndex((panelName) => panelName === name);
    if (!~nameIndex) {
      panels.value.push(name);
    }
  }

  function handleSlaveAdd(name: string, join) {
    const nameIndex = slavePanels.value.findIndex((panelName) => panelName === name);
    if (!~nameIndex) {
      slavePanels.value.push(name);
      slaveMap.value[name] = join.uuid;
    }
  }

  function handleClose(name: string) {
    const nameIndex = panels.value.findIndex((panelName) => panelName === name);
    if (!~nameIndex) return;
    panels.value.splice(nameIndex, 1);
    if (name === value.value) {
      value.value = panels.value[Math.min(nameIndex, panels.value.length - 1)];
    }
  }

  const loadSelect = async () => {
    selectList.value = await Selects({});
  };

  function preview() {
    formBtnPreviewLoading.value = true;
    Preview(genInfo.value)
      .then((res) => {
        previewModel.value = res;
        showModal.value = true;
      })
      .finally(() => {
        formBtnPreviewLoading.value = false;
      });
  }

  function submitBuild() {
    dialog.info({
      title: '提示',
      content: '你确定要提交生成吗？',
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: () => {
        formBtnLoading.value = true;
        Build(genInfo.value)
          .then((_res) => {
            buildSuccessNotify();
          })
          .finally(() => {
            formBtnLoading.value = false;
          });
      },
    });
  }

  function submitSave() {
    dialog.info({
      title: '提示',
      content: '你确定要保存生成配置吗？',
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: () => {
        Edit(genInfo.value).then((_res) => {
          message.success('操作成功');
        });
      },
    });
  }

  function buildSuccessNotify() {
    let count = 10;
    const n = notification.success({
      title: '生成提交成功',
      content: `如果你使用的热编译，页面将在 ${count} 秒后自动刷新即可生效。否则请手动重启服务后刷新页面！`,
      duration: 10000,
      closable: false,
      onAfterEnter: () => {
        const minusCount = () => {
          count--;
          n.content = `如果你使用的热编译，页面将在 ${count} 秒后自动刷新即可生效。否则请手动重启服务后刷新页面！`;
          if (count > 0) {
            window.setTimeout(minusCount, 1000);
          }
        };
        window.setTimeout(minusCount, 1000);
      },
      onAfterLeave: () => {
        location.reload();
      },
    });
  }

  function handleClean() {
    cleanBtnLoading.value = true;
    Preview(genInfo.value)
      .then((res) => {
        if (!res || !res.views) {
          message.warning('没有获取到生成文件信息');
          return;
        }

        const fileKeyNameMap: Record<string, string> = {
          'api.go': 'Go API 定义',
          'input.go': 'Go 输入模型',
          'controller.go': 'Go 控制器',
          'logic.go': 'Go 业务逻辑',
          'router.go': 'Go 路由',
          'web.api.ts': '前端 API',
          'web.model.ts': '前端数据模型',
          'web.index.vue': '前端列表页',
          'web.edit.vue': '前端编辑页',
          'web.view.vue': '前端详情页',
          'source.sql': '菜单SQL',
          'dao.go': 'DAO 模型',
          'dao.internal.go': 'DAO Internal',
          'do.go': 'DO 模型',
          'entity.go': 'Entity 模型',
        };

        const files: any[] = [];
        const existPaths: string[] = [];
        for (const [key, file] of Object.entries(res.views) as [string, any][]) {
          if (!file.path) continue;
          // meth=3(skip) 表示文件已存在, meth=1(create) 表示不存在
          const exists = file.meth === 3 || file.meth === 2;
          files.push({
            key: fileKeyNameMap[key] || key,
            path: file.path,
            exists,
          });
          if (exists) {
            existPaths.push(file.path);
          }
        }

        cleanFileList.value = files;
        cleanCheckedKeys.value = existPaths;
        cleanMenu.value = false;
        showCleanModal.value = true;
      })
      .finally(() => {
        cleanBtnLoading.value = false;
      });
  }

  function submitClean() {
    if (cleanCheckedKeys.value.length === 0) {
      message.warning('请至少选择一个文件');
      return;
    }

    const confirmMsg = cleanMenu.value
      ? `即将删除 ${cleanCheckedKeys.value.length} 个文件并清除对应菜单权限，此操作不可恢复，确定继续吗？`
      : `即将删除 ${cleanCheckedKeys.value.length} 个文件，此操作不可恢复，确定继续吗？`;

    dialog.warning({
      title: '确认删除',
      content: confirmMsg,
      positiveText: '确定删除',
      negativeText: '取消',
      onPositiveClick: () => {
        cleanSubmitLoading.value = true;
        Clean({ id: genId, files: cleanCheckedKeys.value, cleanMenu: cleanMenu.value })
          .then((res) => {
            showCleanModal.value = false;
            const msgs: string[] = [];
            if (res.count > 0) {
              msgs.push(`成功删除 ${res.count} 个文件`);
            }
            if (res.menuCount > 0) {
              msgs.push(`清除 ${res.menuCount} 个菜单`);
            }
            if (res.failed && res.failed.length > 0) {
              msgs.push(`${res.failed.length} 个文件删除失败`);
              message.warning(msgs.join('，'));
            } else {
              message.success(msgs.join('，') || '操作完成');
            }
          })
          .finally(() => {
            cleanSubmitLoading.value = false;
          });
      },
    });
  }

  function handleBack() {
    dialog.info({
      title: '提示',
      content: '你确定要返回生成列表？系统不会主动保存更改',
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: () => {
        router.push({ name: 'develop_code' });
      },
    });
  }
</script>

<style scoped>
  .card-tabs .n-tabs-nav--bar-type {
    padding-left: 4px;
  }
</style>
<style lang="less" scoped>
  ::v-deep(.alert-margin) {
    margin-bottom: 20px;
  }

  ::v-deep(.tag-margin) {
    margin-bottom: 10px;
  }

  ::v-deep(.code-scrollbar) {
    height: calc(100vh - 300px);
    background: #282b2e;
    color: #e0e2e4;
    padding: 10px;
  }
</style>

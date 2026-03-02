<template>
  <div>
    <div class="n-layout-page-header" v-if="!isModal">
      <n-card :bordered="false" title="菜单管理">
        在这里可以管理编辑系统下的所有菜单导航和分配相应的菜单权限
      </n-card>
    </div>
    <n-grid cols="1 s:1 m:1 l:3 xl:3 2xl:3" responsive="screen" :x-gap="12">
      <n-gi span="1">
        <n-card :segmented="{ content: true }" :bordered="false" size="small" class="proCard">
          <template #header>
            <n-space>
              <n-button type="info" icon-placement="left" @click="openCreateDrawer">
                <template #icon>
                  <div class="flex items-center">
                    <n-icon size="14">
                      <PlusOutlined />
                    </n-icon>
                  </div>
                </template>
                添加菜单
              </n-button>
              <n-button
                type="info"
                icon-placement="left"
                @click="openChildCreateDrawer"
                :disabled="formParams.id == 0"
              >
                <template #icon>
                  <div class="flex items-center">
                    <n-icon size="14">
                      <PlusOutlined />
                    </n-icon>
                  </div>
                </template>
                添加子菜单
              </n-button>
              <n-button type="primary" icon-placement="left" @click="packHandle">
                全部{{ expandedKeys.length ? '收起' : '展开' }}
                <template #icon>
                  <div class="flex items-center">
                    <n-icon size="14">
                      <AlignLeftOutlined />
                    </n-icon>
                  </div>
                </template>
              </n-button>
              <n-button
                type="error"
                icon-placement="left"
                @click="handleBatchDelete"
                :disabled="checkedKeys.length === 0"
              >
                <template #icon>
                  <div class="flex items-center">
                    <n-icon size="14">
                      <DeleteOutlined />
                    </n-icon>
                  </div>
                </template>
                删除
              </n-button>
            </n-space>
          </template>
          <div class="w-full menu">
            <n-input type="text" v-model:value="pattern" placeholder="输入菜单名称或权限路径搜索">
              <template #suffix>
                <n-icon size="18" class="cursor-pointer">
                  <SearchOutlined />
                </n-icon>
              </template>
            </n-input>
            <div class="py-3 menu-list">
              <template v-if="loading">
                <div class="flex items-center justify-center py-4">
                  <n-spin size="medium" />
                </div>
              </template>
              <template v-else>
                <n-tree
                  block-line
                  cascade
                  checkable
                  :virtual-scroll="true"
                  :pattern="pattern"
                  :filter="filterTreeNode"
                  :data="treeOption"
                  :expandedKeys="expandedKeys"
                  :checkedKeys="checkedKeys"
                  style="max-height: 650px; overflow: hidden"
                  @update:selected-keys="selectedTree"
                  @update:expanded-keys="onExpandedKeys"
                  @update:checked-keys="onCheckedKeys"
                />
              </template>
            </div>
          </div>
        </n-card>
      </n-gi>
      <n-gi span="2">
        <n-card :segmented="{ content: true }" :bordered="false" size="small" class="proCard">
          <template #header>
            <n-space>
              <n-icon size="18">
                <FormOutlined />
              </n-icon>
              <span>编辑菜单{{ treeItemTitle ? `：${treeItemTitle}` : '' }}</span>
            </n-space>
          </template>

          <n-result
            v-show="formParams.id == 0"
            status="info"
            title="提示"
            description="从菜单列表中选择一项进行编辑"
          />
          <EditForm
            v-if="formParams.id > 0"
            v-model:formParams="formParams"
            v-model:treeOption="treeOption"
            @reload-table="loadTreeOption"
          />
        </n-card>
      </n-gi>
    </n-grid>
    <AddModal ref="addModalRef" v-model:treeOption="treeOption" @reload-table="loadTreeOption" />
  </div>
</template>
<script lang="ts" setup>
  import { computed, onMounted, ref, unref } from 'vue';
  import { useDialog, useMessage } from 'naive-ui';
  import {
    AlignLeftOutlined,
    DeleteOutlined,
    FormOutlined,
    PlusOutlined,
    SearchOutlined,
  } from '@vicons/antd';
  import { getMenuList, BatchDeleteMenu } from '@/api/system/menu';
  import { newState, State, loadOptions } from '@/views/permission/menu/model';
  import EditForm from '@/views/permission/menu/editForm.vue';
  import AddModal from '@/views/permission/menu/addModal.vue';

  const dialog = useDialog();
  const message = useMessage();

  const addModalRef = ref();
  const loading = ref(false);
  const treeOption = ref([]);
  const pattern = ref('');
  const expandedKeys = ref([]);
  const checkedKeys = ref<number[]>([]);
  const formParams = ref<State>(newState(null));
  const treeItemTitle = computed(() => {
    if (formParams.value.id > 0) {
      return formParams.value.label + ' #' + formParams.value.id;
    }
    return '';
  });
  const isModal = defineModel<boolean>('isModal', { default: false });

  function openCreateDrawer() {
    addModalRef.value.openModal(null);
  }

  function openChildCreateDrawer() {
    const state = newState(null);
    state.pid = formParams.value.id;
    state.type = formParams.value.type;
    addModalRef.value.openModal(state);
  }

  function selectedTree(keys: number[], option: any[]) {
    let item = null;
    if (keys.length) {
      item = option[0];
    }
    formParams.value = newState(item);
  }

  function packHandle() {
    if (expandedKeys.value.length) {
      expandedKeys.value = [];
    } else {
      expandedKeys.value = unref(treeOption).map((item: any) => item.key as string) as [];
    }
  }

  function onExpandedKeys(keys) {
    expandedKeys.value = keys;
  }

  function onCheckedKeys(keys: number[]) {
    checkedKeys.value = keys;
  }

  // 从勾选的节点中提取顶级父节点ID（过滤掉那些父节点也被勾选的子节点）
  function getTopLevelCheckedIds(nodes: any[], checkedSet: Set<number>): number[] {
    const result: number[] = [];
    for (const node of nodes) {
      if (checkedSet.has(node.key)) {
        // 当前节点被勾选，加入结果，不再遍历其子节点（后端会递归删除）
        result.push(node.key);
      } else if (node.children && node.children.length > 0) {
        // 当前节点未被勾选，继续在子节点中查找
        result.push(...getTopLevelCheckedIds(node.children, checkedSet));
      }
    }
    return result;
  }

  // 批量删除勾选的菜单
  function handleBatchDelete() {
    if (checkedKeys.value.length === 0) {
      message.warning('请先勾选要删除的菜单');
      return;
    }

    dialog.warning({
      title: '提示',
      content: `您确定要删除选中的 ${checkedKeys.value.length} 个菜单吗？勾选父菜单将同步删除其所有子菜单。`,
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: () => {
        const checkedSet = new Set(checkedKeys.value);
        const ids = getTopLevelCheckedIds(unref(treeOption), checkedSet);
        BatchDeleteMenu({ ids }).then(() => {
          message.success('删除成功');
          checkedKeys.value = [];
          formParams.value = newState(null);
          loadTreeOption();
        });
      },
    });
  }

  // 按名称和权限搜索
  function filterTreeNode(pattern: string, node: any) {
    if (!pattern) return true;
    const searchText = pattern.toLowerCase();

    const label = (node.label || node.title || '').toLowerCase();
    if (label.includes(searchText)) {
      return true;
    }

    const permissions = node.permissions || '';
    if (permissions) {
      const permissionsLower = permissions.toLowerCase();
      if (permissionsLower.includes(searchText)) {
        return true;
      }
    }
    return false;
  }

  // 加载菜单选项树
  function loadTreeOption() {
    const needLoading = treeOption.value.length == 0;
    if (needLoading) {
      loading.value = true;
    }

    getMenuList()
      .then((res) => {
        if (res.list && res.list.length > 0) {
          treeOption.value = res.list;
        } else {
          treeOption.value = [];
        }
      })
      .finally(() => {
        if (needLoading) {
          loading.value = false;
        }
      });
  }

  onMounted(() => {
    loadTreeOption();
    loadOptions();
  });
</script>

<template>
  <div>
    <n-drawer v-model:show="showModal" :width="dialogWidth">
      <n-drawer-content title="CURD列表详情" closable>
        <n-spin :show="loading" description="请稍候...">
          <n-descriptions label-placement="left" class="py-2" :column="1">
            <n-descriptions-item>
              <template #label>
                标题
              </template>
              {{ formValue.title }}
            </n-descriptions-item>
            <n-descriptions-item>
              <template #label>
                描述
              </template>
              <span v-html="formValue.description"></span>
            </n-descriptions-item>
            <n-descriptions-item>
              <template #label>
                内容
              </template>
              <span v-html="formValue.content"></span>
            </n-descriptions-item>
            <n-descriptions-item>
              <template #label>
                单图
              </template>
              <n-image style="margin-left: 10px; height: 100px; width: 100px" :src="formValue.image" />
            </n-descriptions-item>
            <n-descriptions-item>
              <template #label>
                附件
              </template>
<div class="upload-card" v-show="formValue.attachfile !== ''" @click="download(formValue.attachfile)">
                <div class="upload-card-item" style="height: 100px; width: 100px">
                  <div class="upload-card-item-info">
                    <div class="img-box">
                      <n-avatar :style="fileAvatarCSS">
                        {{ getFileExt(formValue.attachfile) }}
                      </n-avatar>
                    </div>
                  </div>
                </div>
              </div>
            </n-descriptions-item>
            <n-descriptions-item>
              <template #label>
                所在城市
              </template>
              {{ formValue.cityId }}
            </n-descriptions-item>
            <n-descriptions-item label="显示开关">
              <n-switch
            v-model:value="formValue.switch"
            :unchecked-value="2"
            :checked-value="1"
            :disabled="true"
            />
            </n-descriptions-item>
            <n-descriptions-item>
              <template #label>
                排序
              </template>
              {{ formValue.sort }}
            </n-descriptions-item>
          </n-descriptions>
        </n-spin>
      </n-drawer-content>
    </n-drawer>
  </div>
</template>
<script lang="ts" setup>
  import { useMessage } from 'naive-ui';
  import { View } from '@/api/curdDemo';
  import { State, newState } from './model';
  import { adaModalWidth } from '@/utils/hotgo';
  import { getFileExt } from '@/utils/urlUtils';
  import { useDictStore } from '@/store/modules/dict';

  const message = useMessage();
  const dict = useDictStore();
  const loading = ref(false);
  const showModal = ref(false);
  const formValue = ref(newState(null));
  const dialogWidth = computed(() => {
    return adaModalWidth(580);
  });
  const fileAvatarCSS = computed(() => {
    return {
      '--n-merged-size': `var(--n-avatar-size-override, 80px)`,
      '--n-font-size': `18px`,
    };
  });

  // 下载
  function download(url: string) {
    window.open(url);
  }

  // 打开模态框
  function openModal(state: State) {
    showModal.value = true;
    loading.value = true;
    View({ id: state.id })
      .then((res) => {
        formValue.value = res;
      })
      .finally(() => {
        loading.value = false;
      });
  }

  defineExpose({
    openModal,
  });
</script>

<style lang="less" scoped></style>
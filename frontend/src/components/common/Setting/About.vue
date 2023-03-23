<script setup lang='ts'>
import { computed, onMounted, ref } from 'vue'
import { NSpin } from 'naive-ui'
import { fetchChatConfig } from '@/api'
import pkg from '@/../package.json'
import { useAuthStore } from '@/store'

interface ConfigState {
  timeoutMs?: number
  reverseProxy?: string
  apiModel?: string
  socksProxy?: string
  httpsProxy?: string
  balance?: string
  accessToken?: string
}

const authStore = useAuthStore()

const loading = ref(false)

const config = ref<ConfigState>()

const isChatGPTAPI = computed<boolean>(() => !!authStore.isChatGPTAPI)

async function fetchConfig() {
  try {
    loading.value = true
    const { data } = await fetchChatConfig<ConfigState>()
    console.log(data)
    config.value = data
  }
  finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchConfig()
})
</script>

<template>
  <NSpin :show="loading">
    <div class="p-4 space-y-4">
      <h2 class="text-xl font-bold">
        Version - {{ pkg.version }}
      </h2>
      <div class="p-2 space-y-2 rounded-md bg-neutral-100 dark:bg-neutral-700">
        <p>
          此项目开源于
          https://github.com/xyhelper/xyhelper-desktop
          ，免费且基于 MIT 协议，没有任何形式的付费行为！
        </p>
        <p>
          如果您觉得这个软件有帮助，请考虑捐赠以支持我们的持续开发工作。
         <center></center> <img src="https://xyhelper.cn/donate.jpg" width="264" height="371" />
        </p>
        <p>
          项目界面基于  ChatGPT Web 开发，感谢作者的开源精神！
        </p>
      </div>
      <!-- <p>{{ $t("setting.api") }}：{{ config?.apiModel ?? '-' }}</p> -->
      <!-- <p v-if="isChatGPTAPI">
        {{ $t("setting.balance") }}：{{ config?.balance ?? '-' }}
        </p> -->
        <p v-if="!isChatGPTAPI">
          {{ $t("setting.reverseProxy") }}：{{ config?.reverseProxy ?? '-' }}
        </p>
        <p>AccessToken：{{ config?.accessToken ?? '-' }}</p>
        <!-- <p>{{ $t("setting.timeout") }}：{{ config?.timeoutMs ?? '-' }}</p> -->
        <!-- <p>{{ $t("setting.socks") }}：{{ config?.socksProxy ?? '-' }}</p> -->
        <!-- <p>{{ $t("setting.httpsProxy") }}：{{ config?.httpsProxy ?? '-' }}</p> -->
    </div>
  </NSpin>
</template>

<template>
  <div class="connection-status-inline" :class="statusClass">
    <span class="status-text">{{ statusText }}</span>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { useNetworkStatus } from '../composables/useNetworkStatus'

const { t } = useI18n()
const route = useRoute()
const { isOnline } = useNetworkStatus()

const statusText = computed(() => {
  return isOnline.value ? t('connection.online') : t('connection.offline')
})

// 페이지별로 상태 색상 결정
const statusClass = computed(() => {
  const isTransferPage = route.path === '/send-bitcoin'
  
  if (isTransferPage) {
    // 전송 페이지: 온라인이 좋음(초록), 오프라인이 나쁨(빨강)
    return {
      'online-good': isOnline.value,
      'offline-bad': !isOnline.value
    }
  } else {
    // 지갑 생성/확인 페이지: 오프라인이 좋음(초록), 온라인이 나쁨(빨강)
    return {
      'online-bad': isOnline.value,
      'offline-good': !isOnline.value
    }
  }
})
</script>

<style scoped>
.connection-status-inline {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

/* 좋은 상태 (초록색) */
.connection-status-inline.online-good,
.connection-status-inline.offline-good {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
  border: 1px solid rgba(34, 197, 94, 0.3);
}

/* 나쁜 상태 (빨간색) */
.connection-status-inline.online-bad,
.connection-status-inline.offline-bad {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.status-text {
  white-space: nowrap;
}
</style>
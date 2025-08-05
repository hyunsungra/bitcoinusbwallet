<template>
  <div class="screen">
    <div class="container">
      <div class="top-controls">
        <LanguageSelector />
      </div>
      
      <div class="wallet-header">
        <div class="bitcoin-logo">₿</div>
        <h1>Bitcoin USB Wallet</h1>
        <p class="subtitle">{{ $t('main.subtitle') }}</p>
        <p class="version-info">v0.0.1</p>
      </div>
      
      <div class="menu-buttons">
        <button class="menu-btn primary" @click="handleCreateWallet">
          <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 4v16m8-8H4"/>
          </svg>
          <span>{{ $t('main.create_wallet') }}</span>
        </button>
        
        <button class="menu-btn" @click="handleCheckWallet">
          <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
            <path d="M9 12l2 2 4-4"/>
          </svg>
          <span>{{ $t('main.check_wallet') }}</span>
        </button>
        
        <button class="menu-btn" @click="handleTransferBitcoin">
          <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14.5 10a4.5 4.5 0 0 0 0-9H12h-1.5m0 9H12h1.5m-1.5 0V1v9zm-1.5 0h1.5m-1.5 0H9a3 3 0 0 0-3 3v2a3 3 0 0 0 3 3h3.5M9 16l3-3m0 0l3 3m-3-3v9"/>
          </svg>
          <span>{{ $t('main.transfer_bitcoin') }}</span>
        </button>
        
      </div>
      
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Swal from 'sweetalert2'
import LanguageSelector from '../components/LanguageSelector.vue'
import { useNetworkStatus } from '../composables/useNetworkStatus'

const router = useRouter()
const { t } = useI18n()
const { isOnline } = useNetworkStatus()

// 지갑 생성 클릭 핸들러
const handleCreateWallet = async () => {
  if (isOnline.value) {
    // 온라인일 때 경고
    const result = await Swal.fire({
      icon: 'warning',
      title: t('security.online_warning_title'),
      html: t('security.online_create_warning'),
      showCancelButton: true,
      confirmButtonText: t('common.continue'),
      cancelButtonText: t('common.cancel'),
      confirmButtonColor: '#f7931a',
      cancelButtonColor: '#6b7280'
    })
    
    if (result.isConfirmed) {
      router.push('/create-wallet')
    }
  } else {
    // 오프라인일 때 바로 이동
    router.push('/create-wallet')
  }
}

// 지갑 확인 클릭 핸들러
const handleCheckWallet = async () => {
  if (isOnline.value) {
    // 온라인일 때 경고
    const result = await Swal.fire({
      icon: 'warning',
      title: t('security.online_warning_title'),
      html: t('security.online_check_warning'),
      showCancelButton: true,
      confirmButtonText: t('common.continue'),
      cancelButtonText: t('common.cancel'),
      confirmButtonColor: '#f7931a',
      cancelButtonColor: '#6b7280'
    })
    
    if (result.isConfirmed) {
      router.push('/check-wallet')
    }
  } else {
    // 오프라인일 때 바로 이동
    router.push('/check-wallet')
  }
}

// 비트코인 전송 클릭 핸들러
const handleTransferBitcoin = async () => {
  if (!isOnline.value) {
    // 오프라인일 때 경고 (전송은 온라인 필요)
    await Swal.fire({
      icon: 'error',
      title: t('security.offline_error_title'),
      text: t('security.offline_transfer_error'),
      confirmButtonText: t('common.ok'),
      confirmButtonColor: '#f7931a'
    })
  } else {
    // 온라인일 때 바로 이동
    router.push('/send-bitcoin')
  }
}

</script>

<style scoped>
.top-controls {
  position: absolute;
  top: 20px;
  right: 20px;
  z-index: 100;
}

.screen {
  position: relative;
  min-height: 100vh;
  background: linear-gradient(135deg, #1a1a2e 0%, #0f0f1e 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.container {
  text-align: center;
  max-width: 600px;
  padding: 40px;
}

.wallet-header {
  margin-bottom: 40px;
}

.bitcoin-logo {
  font-size: 80px;
  color: #f7931a;
  margin-bottom: 20px;
  text-shadow: 0 0 30px rgba(247, 147, 26, 0.5);
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.05); }
}

.wallet-header h1 {
  font-size: 36px;
  font-weight: 700;
  margin: 0 0 10px 0;
  line-height: 1.2;
  background: linear-gradient(135deg, #f7931a 0%, #ff9800 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.subtitle {
  color: rgba(255, 255, 255, 0.7);
  font-size: 18px;
  margin: 0 0 10px 0;
}

.version-info {
  color: rgba(255, 255, 255, 0.5);
  font-size: 14px;
  font-weight: 400;
  margin: 0;
}

.menu-buttons {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 40px;
}

.menu-btn {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 24px;
  border: none;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.1);
  color: white;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.menu-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.2);
}

.menu-btn.primary {
  background: linear-gradient(135deg, #f7931a 0%, #ff9800 100%);
  border: none;
  box-shadow: 0 4px 12px rgba(247, 147, 26, 0.3);
}

.menu-btn.primary:hover {
  background: linear-gradient(135deg, #ff9800 0%, #f7931a 100%);
  box-shadow: 0 6px 20px rgba(247, 147, 26, 0.4);
}

.menu-btn .icon {
  width: 20px;
  height: 20px;
}

.developer-info {
  margin-top: 20px;
  text-align: center;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}
</style>
<template>
  <div class="screen">
    <div class="container">
      <div class="wallet-details-header">
        <div class="page-header-row">
          <button class="back-btn" @click="$router.push('/check-wallet')">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M15 19l-7-7 7-7"/>
            </svg>
            <span>{{ $t('common.back') }}</span>
          </button>
          <div class="title-description">
            <h2>{{ $t('wallet.title') }}</h2>
            <span class="description-text">{{ $t('wallet.description') }}</span>
          </div>
          <NetworkStatus />
        </div>
      </div>
      
      <div class="wallet-details-layout">
        <div class="wallet-details-columns">
          <!-- 좌측: 니모닉, 패스프레이즈 -->
          <div class="left-column">
            <div class="detail-section">
              <div class="section-header">
                <h3>{{ $t('wallet.mnemonic_title') }}</h3>
                <button class="toggle-btn" @click="toggleMnemonic">
                  <span class="toggle-text">{{ showMnemonic ? $t('common.hide') : $t('common.show') }}</span>
                </button>
              </div>
              <div class="hidden-content" v-if="showMnemonic">
                <div class="mnemonic-words-display">
                  <span v-for="(word, index) in mnemonicWords" :key="index" class="mnemonic-word">
                    <span class="word-number">{{ index + 1 }}</span>
                    {{ word }}
                  </span>
                </div>
                <div class="warning-text">
                  <p><strong>{{ $t('wallet.mnemonic_warning') }}</strong></p>
                </div>
              </div>
            </div>

            <div class="detail-section">
              <div class="section-header">
                <h3>{{ $t('wallet.passphrase_title') }}</h3>
                <button class="toggle-btn" @click="togglePassphrase">
                  <span class="toggle-text">{{ showPassphrase ? $t('common.hide') : $t('common.show') }}</span>
                </button>
              </div>
              <div class="hidden-content" v-if="showPassphrase">
                <div class="passphrase-display">{{ passphrase || $t('wallet.no_passphrase') }}</div>
                <div class="warning-text">
                  <p><strong>{{ $t('wallet.passphrase_warning') }}</strong></p>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 우측: 공개키, 개인키 -->
          <div class="right-column">
            <div class="detail-section">
              <div class="section-header">
                <h3>{{ $t('wallet.public_address_title') }}</h3>
                <button class="copy-btn" @click="copyPublicAddress">{{ $t('common.copy') }}</button>
              </div>
              <div class="address-display">
                <div class="address-text">{{ publicAddress }}</div>
                <div class="qr-container">
                  <canvas id="public-address-qr"></canvas>
                </div>
              </div>
            </div>

            <div class="detail-section">
              <div class="section-header">
                <h3>{{ $t('wallet.private_key_title') }}</h3>
                <button class="toggle-btn" @click="togglePrivateKey">
                  <span class="toggle-text">{{ showPrivateKey ? $t('common.hide') : $t('common.show') }}</span>
                </button>
              </div>
              <div class="hidden-content" v-if="showPrivateKey">
                <div class="private-key-display">
                  <div class="key-text">{{ privateKey }}</div>
                  <div class="qr-container">
                    <canvas id="private-key-qr"></canvas>
                  </div>
                  <button class="copy-btn" @click="copyPrivateKey">{{ $t('wallet.copy_private_key') }}</button>
                </div>
                <div class="warning-text">
                  <p><strong>{{ $t('wallet.private_key_warning') }}</strong></p>
                </div>
              </div>
            </div>

          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Swal from 'sweetalert2'
import NetworkStatus from '../components/NetworkStatus.vue'
import QRCode from 'qrcode'


const router = useRouter()
const { t } = useI18n()

const showMnemonic = ref(false);
const showPassphrase = ref(false);
const showPrivateKey = ref(false);
const walletData = ref(null);

const mnemonicWords = computed(() => {
  if (!walletData.value?.mnemonic) return []
  return walletData.value.mnemonic.split(' ')
})

const passphrase = computed(() => walletData.value?.passphrase || '')
const publicAddress = computed(() => walletData.value?.address || '')
const privateKey = computed(() => walletData.value?.privateKeyWIF || '')

onMounted(async () => {
  // history.state에서 지갑 데이터 가져오기
  if (history.state?.walletData) {
    walletData.value = history.state.walletData
    
    // QR 코드 생성
    await nextTick()
    await generateQRCodes()
  } else {
    // 직접 접근한 경우 체크 페이지로 리다이렉트
    router.push('/check-wallet')
  }
})

// QR 코드 생성
const generateQRCodes = async () => {
  try {
    // 공개 주소 QR 코드
    if (publicAddress.value) {
      const addressCanvas = document.getElementById('public-address-qr')
      if (addressCanvas) {
        await QRCode.toCanvas(addressCanvas, publicAddress.value, {
          width: 100,
          margin: 1,
          color: {
            dark: '#000000',
            light: '#ffffff'
          }
        })
      }
    }
    
    // 개인키 QR 코드 (보임 상태에서만)
    if (showPrivateKey.value && privateKey.value) {
      const privateKeyCanvas = document.getElementById('private-key-qr')
      if (privateKeyCanvas) {
        await QRCode.toCanvas(privateKeyCanvas, privateKey.value, {
          width: 100,
          margin: 1,
          color: {
            dark: '#000000',
            light: '#ffffff'
          }
        })
      }
    }
  } catch (error) {
    console.error('QR 코드 생성 실패:', error)
  }
}

const toggleMnemonic = () => {
  showMnemonic.value = !showMnemonic.value;
};

const togglePassphrase = () => {
  showPassphrase.value = !showPassphrase.value;
};

const togglePrivateKey = async () => {
  showPrivateKey.value = !showPrivateKey.value;
  // 개인키 보임 상태에서 QR 코드 생성
  if (showPrivateKey.value) {
    await nextTick()
    await generateQRCodes()
  }
};



const copyPublicAddress = async () => {
  try {
    await navigator.clipboard.writeText(publicAddress.value);
    await Swal.fire({
      icon: 'success',
      title: t('clipboard.copy_complete'),
      text: t('clipboard.address_copied'),
      timer: 1500,
      showConfirmButton: false,
      toast: true,
      position: 'top-end'
    });
  } catch (error) {
    console.error('복사 실패:', error);
  }
};

const copyPrivateKey = async () => {
  const result = await Swal.fire({
    icon: 'warning',
    title: t('security.danger_title'),
    text: t('security.private_key_warning'),
    showCancelButton: true,
    confirmButtonText: t('common.copy'),
    cancelButtonText: t('common.cancel'),
    confirmButtonColor: '#ef4444',
    cancelButtonColor: '#6b7280'
  });
  
  if (result.isConfirmed) {
    try {
      await navigator.clipboard.writeText(privateKey.value);
      await Swal.fire({
        icon: 'success',
        title: t('clipboard.copy_complete'),
        text: t('clipboard.private_key_copied'),
        timer: 1500,
        showConfirmButton: false,
        toast: true,
        position: 'top-end'
      });
    } catch (error) {
      console.error('복사 실패:', error);
    }
  }
};
</script>

<style scoped>
.screen {
  min-height: 100vh;
  background: linear-gradient(135deg, #1a1a2e 0%, #0f0f1e 100%);
  color: white;
  padding: 20px;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
}

.wallet-details-header {
  margin-bottom: 0px;
}

.page-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 0;
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  color: white;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s ease;
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.back-btn svg {
  width: 16px;
  height: 16px;
}

.title-description {
  flex-grow: 1;
  margin-left: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-description h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #f7931a;
  white-space: nowrap;
}

.description-text {
  color: rgba(255, 255, 255, 0.7);
  font-size: 14px;
  white-space: nowrap;
}


.wallet-details-layout {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 30px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.wallet-details-columns {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 30px;
  align-items: start;
}

@media (max-width: 992px) {
  .wallet-details-columns {
    grid-template-columns: 1fr;
    gap: 30px;
  }
}

.left-column, .right-column {
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-height: fit-content;
}

.left-column {
  overflow-y: auto;
  max-height: 80vh;
}

.detail-section {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  padding: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.section-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 500;
}

.toggle-btn, .copy-btn {
  padding: 8px 16px;
  background: linear-gradient(135deg, #f7931a 0%, #ff9800 100%);
  border: none;
  border-radius: 6px;
  color: white;
  cursor: pointer;
  font-size: 12px;
  font-weight: 500;
  transition: all 0.2s ease;
}

.toggle-btn:hover, .copy-btn:hover {
  background: linear-gradient(135deg, #ff9800 0%, #f7931a 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(247, 147, 26, 0.3);
}

.hidden-content {
  margin-top: 16px;
}

.mnemonic-words-display {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
  margin-bottom: 16px;
}

@media (max-width: 768px) {
  .mnemonic-words-display {
    grid-template-columns: repeat(3, 1fr);
  }
}

.mnemonic-word {
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 6px;
  color: white;
  font-size: 13px;
  text-align: center;
}

.passphrase-display, .address-text, .key-text {
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  color: white;
  font-size: 14px;
  word-break: break-all;
  margin-bottom: 16px;
}

.qr-container {
  text-align: center;
  margin: 0px 0;
}

.qr-container canvas {
  border-radius: 8px;
  background: white;
  padding: 8px;
}

.balance-display {
  margin: 12px 0;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.balance-label {
  font-size: 14px;
  opacity: 0.8;
}

.balance-value {
  font-size: 16px;
  font-weight: 600;
  color: #10b981;
  font-family: monospace;
}

.balance-loading {
  font-size: 14px;
  opacity: 0.7;
  font-style: italic;
}


.warning-text {
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 8px;
  padding: 12px;
  margin-top: 16px;
}

.warning-text p {
  margin: 0;
  font-size: 13px;
  line-height: 1.4;
  color: #fecaca;
}

.warning-text p strong {
  color: #ef4444;
}

.private-key-display {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

@media (max-width: 768px) {
  .page-header-row {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }
  
  .title-description {
    margin-left: 0;
    flex-direction: column;
    gap: 4px;
    text-align: center;
  }
  
  .title-description h2,
  .description-text {
    white-space: normal;
  }
  
  .wallet-details-layout {
    padding: 20px;
  }
}
</style>
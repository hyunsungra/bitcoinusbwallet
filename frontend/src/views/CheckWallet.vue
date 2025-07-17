<template>
  <div class="screen">
    <div class="container">
      <div class="page-header-row">
        <button class="back-btn" @click="$router.push('/')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M15 19l-7-7 7-7"/>
          </svg>
          <span>{{ $t('common.back') }}</span>
        </button>
        <div class="title-description">
          <h2>{{ $t('check.title') }}</h2>
          <span class="description-text">{{ $t('check.description') }}</span>
        </div>
        <NetworkStatus />
      </div>
      
      <div class="wallet-form">
        <div class="form-group">
          <label for="wallet-file">{{ $t('check.file_label') }}</label>
          <div class="file-input-wrapper">
            <input type="file" id="wallet-file" accept=".wallet" style="display: none;" @change="handleFileSelect">
            <button type="button" class="file-select-btn" @click="selectFile">
              <svg class="file-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14,2 14,8 20,8"/>
                <line x1="16" y1="13" x2="8" y2="13"/>
                <line x1="16" y1="17" x2="8" y2="17"/>
                <polyline points="10,9 9,9 8,9"/>
              </svg>
              <span class="file-btn-text">{{ $t('check.choose_file') }}</span>
            </button>
            <span class="file-name-display">{{ fileName || $t('check.no_file_selected') }}</span>
          </div>
          <p class="file-help">{{ $t('check.file_help') }}</p>
        </div>
        
        <div class="form-group">
          <label for="wallet-unlock-password">{{ $t('check.password_label') }}</label>
          <div class="password-input-container">
            <input 
              :type="showPassword ? 'text' : 'password'" 
              id="wallet-unlock-password" 
              :placeholder="$t('check.password_placeholder')"
              v-model="password"
              @keyup.enter="unlockWallet"
            >
            <button type="button" class="password-toggle-btn" @click="showPassword = !showPassword">
              <svg class="eye-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                <circle cx="12" cy="12" r="3"/>
              </svg>
            </button>
          </div>
        </div>
        
        <div class="form-actions">
          <button class="action-btn primary large" @click="unlockWallet" :disabled="isUnlocking">
            <span v-if="isUnlocking">{{ t('wallet.opening') }}</span>
            <span v-else>{{ $t('check.unlock_button') }}</span>
          </button>
          <button class="action-btn secondary large" @click="$router.push('/')">{{ $t('common.cancel') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Swal from 'sweetalert2'
import NetworkStatus from '../components/NetworkStatus.vue'
import { useNetworkStatus } from '../composables/useNetworkStatus'

const router = useRouter()
const { t } = useI18n()
const { isOnline } = useNetworkStatus()

const fileName = ref('');
const filePath = ref('');
const password = ref('');
const showPassword = ref(false);
const isUnlocking = ref(false);

// 하이브리드 함수들
const CheckWallet = async (request) => {
  if (window.go && window.go.main && window.go.main.App) {
    return await window.go.main.App.CheckWallet(request);
  }
  // 백업 더미 데이터
  return {
    success: true,
    message: "성공",
    walletData: {
      version: "2.0",
      name: "",
      mnemonic: "",
      passphrase: "",
      address: "",
      privateKey: "",
      publicKey: "",
      createdAt: ""
    }
  }
}

const SelectWalletFile = async () => {
  if (window.go && window.go.main && window.go.main.App) {
    return await window.go.main.App.SelectWalletFile();
  }
  return ""
}

const selectFile = async () => {
  try {
    const selectedPath = await SelectWalletFile();
    if (selectedPath) {
      filePath.value = selectedPath;
      fileName.value = selectedPath.split(/[\\\/]/).pop(); // 파일명만 추출
    }
  } catch (error) {
    console.error('파일 선택 오류:', error);
    await Swal.fire({
      icon: 'error',
      title: t('alerts.error'),
      text: t('alerts.folder_select_error'),
      confirmButtonColor: '#f7931a'
    });
  }
};

const handleFileSelect = (event) => {
  const file = event.target.files[0];
  if (file) {
    fileName.value = file.name;
    filePath.value = file.path || file.name;
  }
};

const unlockWallet = async () => {
  if (!filePath.value) {
    await Swal.fire({
      icon: 'warning',
      title: t('alerts.warning'),
      text: t('wallet.select_file'),
      confirmButtonColor: '#f7931a'
    });
    return;
  }
  
  if (!password.value) {
    await Swal.fire({
      icon: 'warning',
      title: t('alerts.warning'),
      text: t('wallet.enter_password'),
      confirmButtonColor: '#f7931a'
    });
    return;
  }
  
  isUnlocking.value = true;
  
  try {
    const response = await CheckWallet({
      filePath: filePath.value,
      password: password.value
    });
    
    if (response && response.success) {
      // 성공적으로 지갑을 열었을 때 상세 페이지로 이동
      router.push({
        path: '/wallet-details',
        state: { walletData: response.walletData }
      });
    } else {
      await Swal.fire({
        icon: 'error',
        title: t('alerts.error'),
        text: response?.message || t('wallet.invalid_password'),
        confirmButtonColor: '#f7931a'
      });
    }
  } catch (error) {
    // console.error('지갑 열기 오류:', error);
    await Swal.fire({
      icon: 'error',
      title: t('alerts.error'),
      text: t('wallet.cannot_open'),
      confirmButtonColor: '#f7931a'
    });
  } finally {
    isUnlocking.value = false;
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
  max-width: 800px;
  margin: 0 auto;
}

.page-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
  padding: 20px 0;
}

/* 네트워크 경고 배너 */
.network-warning-banner {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 8px;
  margin-bottom: 20px;
  color: #ef4444;
  font-size: 14px;
  font-weight: 500;
}

.network-warning-banner .warning-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
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


.wallet-form {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 30px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  opacity: 0.9;
}

.file-input-wrapper {
  display: flex;
  gap: 12px;
  align-items: center;
}

.file-select-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: linear-gradient(135deg, #f7931a 0%, #ff9800 100%);
  border: none;
  border-radius: 8px;
  color: white;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.file-select-btn:hover {
  background: linear-gradient(135deg, #ff9800 0%, #f7931a 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(247, 147, 26, 0.3);
}

.file-icon {
  width: 16px;
  height: 16px;
}

.file-name-display {
  flex: 1;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  color: white;
  font-size: 14px;
}

.file-help {
  font-size: 12px;
  opacity: 0.7;
  margin: 8px 0 0 0;
}

.password-input-container {
  position: relative;
  display: flex;
}

.password-input-container input {
  flex: 1;
  padding: 12px 45px 12px 16px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  color: white;
  font-size: 14px;
  transition: all 0.2s ease;
}

.password-input-container input:focus {
  outline: none;
  border-color: rgba(255, 255, 255, 0.4);
  background: rgba(255, 255, 255, 0.15);
}

.password-input-container input::placeholder {
  color: rgba(255, 255, 255, 0.5);
}

/* 브라우저 기본 패스워드 아이콘 숨기기 */
.password-input-container input::-ms-reveal,
.password-input-container input::-ms-clear {
  display: none;
}

.password-input-container input::-webkit-credentials-auto-fill-button {
  display: none !important;
}

.password-input-container input::-webkit-contacts-auto-fill-button {
  display: none !important;
}

.password-toggle-btn {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.6);
  cursor: pointer;
  padding: 4px;
}

.password-toggle-btn:hover {
  color: white;
}

.eye-icon {
  width: 16px;
  height: 16px;
}

.form-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-top: 30px;
}

.action-btn {
  padding: 12px 32px;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  min-width: 120px;
}

.action-btn.large {
  padding: 14px 40px;
  font-size: 16px;
}

.action-btn.primary {
  background: linear-gradient(145deg, #10b981, #059669);
  color: white;
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
}

.action-btn.primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(16, 185, 129, 0.4);
}

.action-btn.secondary {
  background: rgba(255, 255, 255, 0.1);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.action-btn.secondary:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-1px);
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
  
  .wallet-form {
    padding: 20px;
  }
}
</style>
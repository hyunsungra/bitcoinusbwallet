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
          <h2>{{ $t('create.title') }}</h2>
          <span class="description-text">{{ $t('create.description') }}</span>
        </div>
        <NetworkStatus />
      </div>
      
      <div class="wallet-creation-layout">
        <div class="main-columns">
          <!-- 좌측: 니모닉 섹션 -->
          <div class="left-column">
            <div class="section-header">
              <h3>{{ $t('create.mnemonic_title') }}</h3>
              <div class="button-group">
                <button type="button" class="paste-btn-3d" @click="onPasteMnemonic" :disabled="isCreating">
                  <svg class="paste-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"></path>
                    <rect x="8" y="2" width="8" height="4" rx="1" ry="1"></rect>
                  </svg>
                  <span class="btn-text">{{ $t('create.paste_mnemonic') }}</span>
                </button>
                <button type="button" class="generate-btn-3d" @click="onGenerateMnemonic" :disabled="isCreating">
                  <svg class="refresh-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M23 4v6h-6M1 20v-6h6M20.49 9A9 9 0 0 0 5.64 5.64L1 10M3.51 15a9 9 0 0 0 14.85 3.36L23 14"/>
                  </svg>
                  <span class="btn-text">{{ $t('create.generate_new') }}</span>
                </button>
              </div>
            </div>
            <div class="mnemonic-inputs-grid">
              <input 
                v-for="(word, index) in 24" 
                :key="index"
                type="text" 
                class="mnemonic-input"
                :value="mnemonicWords[index] || ''"
                :placeholder="`${$t('create.word')} ${index + 1}`"
                readonly
                @click="onWordClick(index)"
                :data-index="index"
              >
            </div>
            <div class="mnemonic-info">
              <p>{{ $t('create.mnemonic_info') }}</p>
            </div>
            
            <div class="form-group">
              <div class="section-header">
                <h3>{{ $t('create.passphrase_label') }}</h3>
              </div>
              <input type="text" id="passphrase" v-model="passphrase" :placeholder="$t('create.passphrase_placeholder')" :disabled="isCreating">
              <p class="form-help">{{ $t('create.passphrase_help') }}</p>
            </div>

            <div class="warning-text">
              <p><strong>{{ $t('create.warning_title') }}</strong></p>
              <p v-html="$t('create.warning_content')"></p>
            </div>

          </div>
          
          <!-- 우측: 지갑 정보 입력 -->
          <div class="right-column">
            <div class="form-group">
              <div class="section-header">
                <h3>{{ $t('create.wallet_name_label') }}</h3>
              </div>
              <input type="text" id="wallet-name" v-model="walletName" :placeholder="$t('create.wallet_name_placeholder')" :disabled="isCreating">
            </div>
            
            <div class="form-group">
              <div class="section-header">
                <h3>{{ $t('create.password_label') }}</h3>
              </div>
              <input type="password" id="wallet-password" v-model="password" :placeholder="$t('create.password_placeholder')" :disabled="isCreating">
              <!-- <p class="form-help" v-html="$t('create.password_requirements')"></p> -->
              
              <!-- 비밀번호 강도 표시 -->
              <div v-if="password.length > 0" class="password-strength">
                <div class="strength-bar">
                  <div class="strength-fill" :class="'strength-' + passwordValidation.score"></div>
                </div>
                <div class="strength-text">{{ passwordValidation.strength }}</div>
                <div v-if="passwordValidation.errors.length > 0" class="password-errors">
                  <div v-for="error in passwordValidation.errors" :key="error" class="error-item">• {{ error }}</div>
                </div>
              </div>
            </div>
            
            <div class="form-group">
              <div class="section-header">
                <h3>{{ $t('create.password_confirm_label') }}</h3>
              </div>
              <input type="password" id="wallet-password-confirm" v-model="confirmPassword" :placeholder="$t('create.password_confirm_placeholder')" :disabled="isCreating">
              
              <!-- 비밀번호 일치 여부 표시 -->
              <div v-if="confirmPassword.length > 0" class="password-match">
                <div v-if="password === confirmPassword" class="match-success">✓ {{ $t('common.confirm') }}</div>
                <div v-else class="match-error">✗ {{ $t('alerts.validation_password_mismatch') }}</div>
              </div>
            </div>
            
            <div class="form-group">
              <div class="section-header">
                <h3>{{ $t('create.save_path_label') }}</h3>
              </div>
              <div class="save-path-input">
                <input 
                  type="text" 
                  id="save-path" 
                  v-model="savePath" 
                  :placeholder="$t('create.save_path_placeholder')"
                  readonly
                  :disabled="isCreating"
                >
                <button 
                  type="button" 
                  class="folder-select-btn" 
                  @click="selectSaveFolder"
                  :disabled="isCreating"
                >
                  📁 {{ $t('create.select_folder') }}
                </button>
              </div>
              <p class="form-help">{{ $t('create.save_path_help') }}</p>
            </div>
            
            <div class="bottom-actions">
              <button class="action-btn secondary large" @click="$router.push('/')" :disabled="isCreating">{{ $t('common.cancel') }}</button>
              <button class="action-btn primary large" @click="onCreateWallet" :disabled="isCreating">
                <span v-if="isCreating">{{ $t('alerts.info') }}...</span>
                <span v-else>{{ $t('create.create_button') }}</span>
              </button>
            </div>
            
          </div>
        </div>
        
        
      </div>
    </div>

    <!-- 니모닉 단어 선택 팝업 -->
    <div v-if="showWordModal" class="popup-overlay" @click="closeWordModal">
      <div class="popup-container" @click.stop>
        <div class="popup-header">
          <h3>{{ $t('create.select_word') }}</h3>
          <button class="popup-close" @click="closeWordModal">×</button>
        </div>
        <div class="popup-content">
          <div class="word-search">
            <input 
              type="text" 
              v-model="wordSearch" 
              :placeholder="$t('create.search_word')"
              id="word-search-input"
              ref="wordSearchInput"
            >
          </div>
          <div class="word-list" id="word-list">
            <div 
              v-for="word in filteredWords" 
              :key="word"
              class="word-list-item"
              :class="{ active: word === selectedWord }"
              @click="selectWord(word)"
            >
              {{ word }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Swal from 'sweetalert2'
import NetworkStatus from '../components/NetworkStatus.vue'
import { useNetworkStatus } from '../composables/useNetworkStatus'
// import { GenerateMnemonic, ValidatePassword, CreateWallet, GetBIP39WordList, SelectSaveDirectory } from '../../wailsjs/go/main/App'

const router = useRouter()
const { t } = useI18n()
const { isOnline } = useNetworkStatus()

// 임시 함수들 (Wails 바인딩 생성 후 자동으로 실제 함수 사용)
const GenerateMnemonic = async () => {
  if (window.go && window.go.main && window.go.main.App) {
    return await window.go.main.App.GenerateMnemonic();
  }
  // 백업 더미 데이터
  return "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art"
}

const ValidatePassword = async (password) => {
  if (window.go && window.go.main && window.go.main.App) {
    return await window.go.main.App.ValidatePassword(password);
  }
  // 백업 더미 검증
  return {
    isValid: password.length >= 8,
    errors: password.length < 8 ? [t('password.min_length')] : [],
    strength: password.length >= 8 ? t('password.strong') : t('password.weak'),
    score: password.length >= 8 ? 5 : 2
  }
}

const CreateWallet = async (request) => {
  if (window.go && window.go.main && window.go.main.App) {
    return await window.go.main.App.CreateWallet(request);
  }
  // 백업 더미 응답
  return {
    success: true,
    message: '지갑이 성공적으로 생성되었습니다',
    filePath: request.savePath ? `${request.savePath}/${request.name}.wallet` : `C:\\${request.name}.wallet`
  }
}

const GetBIP39WordList = async () => {
  if (window.go && window.go.main && window.go.main.App) {
    return await window.go.main.App.GetBIP39WordList();
  }
  // 백업 더미 단어 목록 (BIP39 일부)
  return [
    "abandon", "ability", "able", "about", "above", "absent", "absorb", "abstract", "absurd", "abuse",
    "access", "accident", "account", "accuse", "achieve", "acid", "acoustic", "acquire", "across", "act",
    "action", "actor", "actress", "actual", "adapt", "add", "addict", "address", "adjust", "admit",
    "adult", "advance", "advice", "aerobic", "affair", "afford", "afraid", "again", "against", "age"
  ]
}

const SelectSaveDirectory = async () => {
  if (window.go && window.go.main && window.go.main.App) {
    return await window.go.main.App.SelectSaveDirectory();
  }
  // 백업 기본 경로
  return "C:\\"
}


// 상태 관리
const mnemonicWords = ref([])
const walletName = ref('')
const passphrase = ref('')
const password = ref('')
const confirmPassword = ref('')
const passwordValidation = ref({
  isValid: false,
  errors: [],
  strength: '',
  score: 0
})
const isCreating = ref(false)
const savePath = ref('')

// 니모닉 단어 선택 모달 상태
const showWordModal = ref(false)
const selectedWordIndex = ref(-1)
const selectedWord = ref('')
const wordSearch = ref('')
const bip39Words = ref([])
const wordSearchInput = ref(null)

// 검색어로 필터링된 단어 목록
const filteredWords = computed(() => {
  if (!wordSearch.value) {
    return bip39Words.value
  }
  return bip39Words.value.filter(word => 
    word.toLowerCase().includes(wordSearch.value.toLowerCase())
  )
})

// 니모닉 붙여넣기 이벤트 리스너
const onPasteMnemonic = async () => {
  try {
    const result = await Swal.fire({
      title: t('create.paste_mnemonic_title'),
      html: `
        <div style="text-align: left; margin-bottom: 15px;">
          <p>${t('create.paste_mnemonic_description')}</p>
        </div>
        <textarea 
          id="mnemonic-textarea" 
          placeholder="${t('create.paste_mnemonic_placeholder')}" 
          style="width: 100%; height: 120px; padding: 10px; border: 1px solid #ddd; border-radius: 5px; font-size: 14px; resize: vertical; font-family: monospace;"
        ></textarea>
      `,
      showCancelButton: true,
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      confirmButtonColor: '#10b981',
      cancelButtonColor: '#6b7280',
      preConfirm: () => {
        const textarea = document.getElementById('mnemonic-textarea')
        const value = textarea.value.trim()
        
        if (!value) {
          Swal.showValidationMessage(t('create.paste_mnemonic_empty'))
          return false
        }
        
        // 띄어쓰기로 분리하여 단어 배열 생성
        const words = value.split(/\s+/).filter(word => word.length > 0)
        
        if (words.length !== 24) {
          Swal.showValidationMessage(t('create.paste_mnemonic_invalid_count', { count: words.length }))
          return false
        }
        
        return words
      }
    })
    
    if (result.isConfirmed && result.value) {
      mnemonicWords.value = result.value
      
      await Swal.fire({
        icon: 'success',
        title: t('create.paste_mnemonic_success'),
        text: t('create.paste_mnemonic_success_description'),
        confirmButtonColor: '#10b981',
        timer: 2000,
        showConfirmButton: false
      })
    }
  } catch (error) {
    console.error('Mnemonic paste error:', error)
    await Swal.fire({
      icon: 'error',
      title: t('alerts.error'),
      text: t('create.paste_mnemonic_error'),
      confirmButtonColor: '#10b981'
    })
  }
}

// 니모닉 생성 이벤트 리스너
const onGenerateMnemonic = async () => {
  try {
    const mnemonic = await GenerateMnemonic()
    if (mnemonic) {
      mnemonicWords.value = mnemonic.split(' ')
    } else {
      await Swal.fire({
        icon: 'error',
        title: t('alerts.error'),
        text: t('alerts.mnemonic_generate_failed'),
        confirmButtonColor: '#10b981'
      })
    }
  } catch (error) {
    console.error('Mnemonic generation error:', error)
    await Swal.fire({
      icon: 'error',
      title: t('alerts.error'),
      text: t('alerts.mnemonic_generate_error'),
      confirmButtonColor: '#10b981'
    })
  }
}

// 비밀번호 검증 이벤트 리스너 (실시간)
const onPasswordValidate = async () => {
  if (password.value.length === 0) {
    passwordValidation.value = {
      isValid: false,
      errors: [],
      strength: '',
      score: 0
    }
    return
  }
  
  try {
    const validation = await ValidatePassword(password.value)
    passwordValidation.value = validation
  } catch (error) {
    console.error('Password validation error:', error)
  }
}


// 지갑 생성 이벤트 리스너
const onCreateWallet = async () => {
  // 입력값 검증
  if (!(await validateInputs())) {
    return
  }
  
  isCreating.value = true
  
  try {
    const request = {
      name: walletName.value,
      password: password.value,
      mnemonic: mnemonicWords.value.join(' '),
      passphrase: passphrase.value,
      savePath: savePath.value
    }
    
    const response = await CreateWallet(request)
    
    if (response.success) {
      const result = await Swal.fire({
        icon: 'success',
        title: t('alerts.wallet_create_complete'),
        html: `${t('alerts.wallet_create_success')}<br><br><strong>${t('alerts.save_path')}</strong><br>${response.filePath}`,
        confirmButtonText: t('alerts.return_to_main'),
        confirmButtonColor: '#10b981',
        timer: 5000,
        timerProgressBar: true
      })
      router.push('/')
    } else {
      await Swal.fire({
        icon: 'error',
        title: t('alerts.wallet_create_failed'),
        text: response.message,
        confirmButtonColor: '#10b981'
      })
    }
  } catch (error) {
    // console.error('Wallet creation error:', error)
    await Swal.fire({
      icon: 'error',
      title: t('alerts.error'),
      text: t('alerts.wallet_create_error'),
      confirmButtonColor: '#10b981'
    })
  } finally {
    isCreating.value = false
  }
}

// 입력값 검증 함수
const validateInputs = async () => {
  if (mnemonicWords.value.length !== 24) {
    await Swal.fire({
      icon: 'warning',
      title: t('alerts.warning'),
      text: t('alerts.validation_mnemonic_required'),
      confirmButtonColor: '#10b981'
    })
    return false
  }
  
  // 니모닉 유효성 검사 제거 (부분 수정 허용)
  // coldwallet도 수정된 니모닉 사용 가능
  /*
  const mnemonicString = mnemonicWords.value.join(' ')
  if (!bip39.validateMnemonic(mnemonicString)) {
    await Swal.fire({
      icon: 'error',
      title: t('alerts.error'),
      text: '잘못된 니모닉입니다. 올바른 BIP39 단어를 사용해주세요.',
      confirmButtonColor: '#10b981'
    })
    return false
  }
  */
  
  if (!walletName.value.trim()) {
    await Swal.fire({
      icon: 'warning',
      title: t('alerts.warning'),
      text: t('alerts.validation_wallet_name_required'),
      confirmButtonColor: '#10b981'
    })
    return false
  }
  
  if (!passwordValidation.value.isValid) {
    await Swal.fire({
      icon: 'warning',
      title: t('alerts.warning'),
      text: t('alerts.validation_password_invalid'),
      confirmButtonColor: '#10b981'
    })
    return false
  }
  
  if (password.value !== confirmPassword.value) {
    await Swal.fire({
      icon: 'warning',
      title: t('alerts.warning'),
      text: t('alerts.validation_password_mismatch'),
      confirmButtonColor: '#10b981'
    })
    return false
  }
  
  return true
}


// 니모닉 단어 클릭 이벤트
const onWordClick = async (index) => {
  if (!bip39Words.value.length) {
    try {
      bip39Words.value = await GetBIP39WordList()
    } catch (error) {
      console.error('Failed to load BIP39 word list:', error)
      return
    }
  }
  
  selectedWordIndex.value = index
  selectedWord.value = mnemonicWords.value[index] || ''
  wordSearch.value = ''
  showWordModal.value = true
  
  // 모달이 열린 후 검색 입력창에 포커스
  setTimeout(() => {
    if (wordSearchInput.value) {
      wordSearchInput.value.focus()
    }
  }, 100)
}

// 단어 선택 (클릭 시 즉시 적용)
const selectWord = (word) => {
  selectedWord.value = word
  if (selectedWordIndex.value >= 0) {
    // 배열 길이가 부족하면 확장
    while (mnemonicWords.value.length <= selectedWordIndex.value) {
      mnemonicWords.value.push('')
    }
    mnemonicWords.value[selectedWordIndex.value] = word
  }
  closeWordModal()
}

// 모달 닫기
const closeWordModal = () => {
  showWordModal.value = false
  selectedWordIndex.value = -1
  selectedWord.value = ''
  wordSearch.value = ''
}

// 저장 폴더 선택
const selectSaveFolder = async () => {
  try {
    const selectedPath = await SelectSaveDirectory()
    if (selectedPath) {
      savePath.value = selectedPath
    }
  } catch (error) {
    console.error('Folder selection error:', error)
    // 사용자가 취소한 경우가 아닌 실제 오류만 알림 표시
    if (!error.message.includes('취소')) {
      await Swal.fire({
        icon: 'error',
        title: t('alerts.error'),
        text: t('alerts.folder_select_error'),
        confirmButtonColor: '#10b981'
      })
    }
  }
}

// 비밀번호 입력시 실시간 검증
watch(password, () => {
  onPasswordValidate()
})

// 페이지 로드시 니모닉 자동 생성
onGenerateMnemonic()
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

.page-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0px;
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


.wallet-creation-layout {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 30px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.main-columns {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 30px;
  margin-bottom: 30px;
}

@media (max-width: 1200px) {
  .main-columns {
    gap: 25px;
  }
}

@media (max-width: 992px) {
  .main-columns {
    grid-template-columns: 1fr;
    gap: 30px;
  }
}

.left-column, .right-column {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0px;
}

.section-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 500;
}

.generate-btn-3d {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: linear-gradient(135deg, #f7931a 0%, #ff9800 100%);
  border: none;
  border-radius: 8px;
  color: white;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s ease;
  box-shadow: 0 4px 12px rgba(247, 147, 26, 0.3);
}

.generate-btn-3d:hover:not(:disabled) {
  background: linear-gradient(135deg, #ff9800 0%, #f7931a 100%);
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(247, 147, 26, 0.4);
}

.generate-btn-3d:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.button-group {
  display: flex;
  gap: 8px;
}

.paste-btn-3d {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: linear-gradient(135deg, #10b981 0%, #16a34a 100%);
  border: none;
  border-radius: 8px;
  color: white;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s ease;
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
}

.paste-btn-3d:hover:not(:disabled) {
  background: linear-gradient(135deg, #16a34a 0%, #10b981 100%);
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(16, 185, 129, 0.4);
}

.paste-btn-3d:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.paste-icon {
  width: 16px;
  height: 16px;
}

.refresh-icon {
  width: 16px;
  height: 16px;
}

.mnemonic-inputs-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
  margin-bottom: 0px;
}

@media (max-width: 768px) {
  .mnemonic-inputs-grid {
    grid-template-columns: repeat(3, 1fr);
    gap: 6px;
  }
}

@media (max-width: 480px) {
  .mnemonic-inputs-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

.mnemonic-input {
  padding: 10px 12px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 6px;
  color: white;
  font-size: 13px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s ease;
  outline: none;
  max-width: 110px;
}

.mnemonic-input:hover {
  background: rgba(255, 255, 255, 0.15);
  border-color: rgba(255, 255, 255, 0.3);
}

.mnemonic-input:focus {
  border-color: rgba(255, 255, 255, 0.4);
}

.mnemonic-input::placeholder {
  color: rgba(255, 255, 255, 0.4);
  font-size: 11px;
}

.mnemonic-info {
  text-align: left;
  opacity: 0.8;
  font-size: 14px;
  margin-top: 0px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-size: 14px;
  font-weight: 500;
  opacity: 0.9;
}

.form-group input {
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  color: white;
  font-size: 14px;
  transition: all 0.2s ease;
}

.form-group input:focus {
  outline: none;
  border-color: rgba(255, 255, 255, 0.4);
  background: rgba(255, 255, 255, 0.15);
}

.form-group input:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.form-group input::placeholder {
  color: rgba(255, 255, 255, 0.5);
}

.form-help {
  font-size: 12px;
  opacity: 0.7;
  margin: 0;
  line-height: 1.4;
}

.password-strength {
  margin-top: 8px;
}

.strength-bar {
  width: 100%;
  height: 4px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 2px;
  overflow: hidden;
  margin-bottom: 4px;
}

.strength-fill {
  height: 100%;
  transition: all 0.3s ease;
}

.strength-0 { width: 0%; background: transparent; }
.strength-1 { width: 15%; background: #ef4444; }
.strength-2 { width: 30%; background: #f97316; }
.strength-3 { width: 45%; background: #eab308; }
.strength-4 { width: 60%; background: #84cc16; }
.strength-5 { width: 75%; background: #22c55e; }
.strength-6 { width: 90%; background: #10b981; }
.strength-7 { width: 100%; background: #059669; }

.strength-text {
  font-size: 12px;
  font-weight: 500;
  margin-bottom: 4px;
}

.password-errors {
  margin-top: 8px;
}

.error-item {
  font-size: 11px;
  color: #fca5a5;
  margin-bottom: 2px;
}

.password-match {
  margin-top: 8px;
  font-size: 12px;
}

.match-success {
  color: #bbf7d0;
  font-weight: 500;
}

.match-error {
  color: #fecaca;
  font-weight: 500;
}

.save-path-input {
  display: flex;
  gap: 12px;
  align-items: center;
}

.save-path-input input {
  flex: 1;
  background: rgba(255, 255, 255, 0.05);
  cursor: default;
}

.folder-select-btn {
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
  display: flex;
  align-items: center;
  gap: 6px;
}

.folder-select-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #ff9800 0%, #f7931a 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(247, 147, 26, 0.3);
}

.folder-select-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.warning-text {
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 8px;
  padding: 16px;
  margin-top: 16px;
}

.warning-text p {
  margin: 0 0 8px 0;
  font-size: 13px;
  line-height: 1.4;
  color: #fecaca;
}

.warning-text p:last-child {
  margin-bottom: 0;
}

.warning-text p strong {
  color: #ef4444;
}

.bottom-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  padding-top: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
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

.action-btn.primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(16, 185, 129, 0.4);
}

.action-btn.secondary {
  background: rgba(255, 255, 255, 0.1);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.action-btn.secondary:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-1px);
}

.action-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none !important;
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
  
  .bottom-actions {
    flex-direction: column;
    align-items: center;
  }
  
  .wallet-creation-layout {
    padding: 20px;
  }
}

/* 팝업 오버레이 */
.popup-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.popup-container {
  background: white;
  border-radius: 12px;
  width: 100%;
  max-width: 600px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
}

.popup-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #e5e7eb;
  background: #f8fafc;
}

.popup-header h3 {
  margin: 0;
  color: #1f2937;
  font-size: 18px;
  font-weight: 600;
}

.popup-close {
  background: none;
  border: none;
  font-size: 24px;
  color: #6b7280;
  cursor: pointer;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 0.2s ease;
}

.popup-close:hover {
  background: #e5e7eb;
  color: #374151;
}

.popup-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.word-search {
  padding: 16px 24px;
  border-bottom: 1px solid #e5e7eb;
}

.word-search input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s ease;
}

.word-search input:focus {
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.word-list {
  flex: 1;
  overflow-y: auto;
  padding: 16px 24px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
  gap: 6px;
  max-height: 400px;
}

.word-list-item {
  padding: 8px 12px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  color: #374151;
  cursor: pointer;
  font-size: 13px;
  text-align: center;
  transition: all 0.2s ease;
  user-select: none;
}

.word-list-item:hover {
  background: #e2e8f0;
  border-color: #cbd5e1;
}

.word-list-item.active {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

@media (max-width: 768px) {
  .popup-container {
    max-width: 95vw;
    max-height: 85vh;
  }
  
  .word-list {
    grid-template-columns: repeat(auto-fill, minmax(70px, 1fr));
    gap: 4px;
    padding: 12px 16px;
  }
  
  .popup-header,
  .word-search {
    padding-left: 16px;
    padding-right: 16px;
  }
}
</style>
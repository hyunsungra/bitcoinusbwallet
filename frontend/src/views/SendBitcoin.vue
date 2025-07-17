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
          <h2>{{ $t('send.title') }}</h2>
          <span class="description-text">{{ $t('send.description') }}</span>
        </div>
        <NetworkStatus />
      </div>
      
      <div class="wallet-details-layout">
        <div class="wallet-details-columns">
          <!-- 좌측 컬럼 - 지갑 로드 -->
          <div class="left-column">
            <div class="detail-section">
              <div class="section-header">
                <h3>{{ $t('send.wallet_load') }}</h3>
              </div>
              
              <!-- 지갑 파일 선택 -->
              <div class="form-group">
                <label for="transfer-wallet-file">{{ $t('send.wallet_file') }}</label>
                <div class="file-input-wrapper">
                  <input type="file" id="transfer-wallet-file" accept=".wallet" style="display: none;" @change="handleFileSelect">
                  <button type="button" class="file-select-btn" @click="selectFile">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                      <polyline points="14,2 14,8 20,8"/>
                      <line x1="16" y1="13" x2="8" y2="13"/>
                      <line x1="16" y1="17" x2="8" y2="17"/>
                      <polyline points="10,9 9,9 8,9"/>
                    </svg>
                    {{ $t('common.select') }}
                  </button>
                  <span class="file-name-display">{{ fileName || $t('send.no_file_selected') }}</span>
                </div>
              </div>
              
              <!-- 지갑 비밀번호 입력 -->
              <div class="form-group">
                <label for="transfer-wallet-password">{{ $t('send.wallet_password') }}</label>
                <div class="password-input-container">
                  <input 
                    :type="showPassword ? 'text' : 'password'"
                    id="transfer-wallet-password"
                    v-model="password"
                    :placeholder="$t('send.enter_password')"
                    @keyup.enter="loadWallet"
                    ref="passwordInput"
                  />
                  <button type="button" class="password-toggle-btn" @click="togglePasswordVisibility">
                    <svg v-if="showPassword" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                      <line x1="1" y1="1" x2="23" y2="23"/>
                    </svg>
                    <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                      <circle cx="12" cy="12" r="3"/>
                    </svg>
                  </button>
                </div>
              </div>
              
              <!-- 지갑 불러오기 버튼 -->
              <div class="form-actions">
                <button class="action-btn primary" @click="loadWallet" :disabled="!filePath || !password">
                  <svg v-if="false" class="loading-spinner" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                    <path d="M9 12l2 2 4-4"/>
                  </svg>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 11-7.778 7.778 5.5 5.5 0 017.777-7.777zm0 0L9.5 9.5m0 0V6m0 3.5H6"/>
                  </svg>
                  {{ $t('send.load_wallet') }}
                </button>
              </div>
              
              <!-- 지갑 정보 표시 -->
              <div class="wallet-info-section" v-if="walletData">
                <div class="info-group">
                  <label>{{ $t('send.public_address') }}</label>
                  <div class="address-text">{{ walletData.address }}</div>
                </div>
                <div class="info-group">
                  <label>{{ $t('send.balance') }}</label>
                  <div class="balance-row">
                    <div class="balance-amount">{{ formatBTC(balance) }} BTC</div>
                    <button class="action-btn secondary small" @click="checkBalance" :disabled="balanceLoading">
                      <svg v-if="balanceLoading" class="loading-spinner" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                        <path d="M9 12l2 2 4-4"/>
                      </svg>
                      <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                        <circle cx="12" cy="12" r="3"/>
                      </svg>
                      {{ balanceLoading ? $t('send.checking') : $t('send.check_balance') }}
                    </button>
                  </div>
                </div>
                <!-- 히스토리 보기 버튼 -->
                <div class="history-button-row">
                  <button class="action-btn secondary full-width" @click="viewHistory">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                    </svg>
                    {{ $t('send.view_history') }}
                  </button>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 우측 컬럼 - 전송 폼 -->
          <div class="right-column">
            <div class="detail-section" v-if="walletData">
              <div class="section-header">
                <h3>{{ $t('send.send_bitcoin') }}</h3>
              </div>
              
              <!-- 받는 주소 입력 -->
              <div class="form-group">
                <label for="receiver-address">{{ $t('send.recipient_address') }}</label>
                <input 
                  type="text" 
                  id="receiver-address" 
                  v-model="recipientAddress"
                  :placeholder="$t('send.enter_recipient_address')"
                  autocomplete="off"
                  spellcheck="false"
                />
                <p class="address-warning">{{ $t('send.address_warning') }}</p>
              </div>
              
              <!-- 네트워크 표시 -->
              <div class="form-group">
                <label>{{ $t('send.network') }}</label>
                <div class="network-display">Bitcoin (BTC)</div>
              </div>
              
              <!-- 수수료 선택 -->
              <div class="form-group">
                <label>{{ $t('send.fee_title') }}</label>
                <div class="fee-options">
                  <label class="radio-option">
                    <input type="radio" name="fee-speed" value="slow" v-model="feeSpeed" @change="updateFeeDisplay">
                    <span class="radio-text">{{ $t('send.slow') }}</span>
                  </label>
                  <label class="radio-option">
                    <input type="radio" name="fee-speed" value="normal" v-model="feeSpeed" @change="updateFeeDisplay">
                    <span class="radio-text">{{ $t('send.normal') }}</span>
                  </label>
                  <label class="radio-option">
                    <input type="radio" name="fee-speed" value="fast" v-model="feeSpeed" @change="updateFeeDisplay">
                    <span class="radio-text">{{ $t('send.fast') }}</span>
                  </label>
                  <label class="radio-option">
                    <input type="radio" name="fee-speed" value="fastest" v-model="feeSpeed" @change="updateFeeDisplay">
                    <span class="radio-text">{{ $t('send.fastest') }}</span>
                  </label>
                  <label class="radio-option">
                    <input type="radio" name="fee-speed" value="custom" v-model="feeSpeed" @change="updateFeeDisplay">
                    <span class="radio-text">{{ $t('send.custom') }}</span>
                  </label>
                </div>
                <div class="fee-display">
                  <div class="selected-fee">{{ $t('send.selected_fee') }}: {{ formatBTC(selectedFeeInBTC) }} BTC ({{ selectedFeeInSatoshi.toLocaleString() }} satoshi)</div>
                  <div class="custom-fee-input" v-if="feeSpeed === 'custom'">
                    <input 
                      type="number" 
                      v-model="customFee"
                      placeholder="0.00002"
                      step="0.000001" 
                      min="0.00002"
                      @input="updateCustomFee"
                    />
                    <span class="custom-fee-unit">BTC</span>
                  </div>
                </div>
              </div>
              
              <!-- 전송 금액 입력 -->
              <div class="form-group">
                <label for="amount-to-send">{{ $t('send.amount_to_send') }}</label>
                <div class="amount-input-wrapper">
                  <input 
                    type="number" 
                    id="amount-to-send"
                    v-model="amount"
                    :placeholder="$t('send.amount_placeholder')"
                    step="0.00000001" 
                    min="0.00002"
                  />
                  <div class="amount-buttons">
                    <button type="button" class="amount-btn max-btn" @click="setMaxAmount">Max</button>
                    <button type="button" class="amount-btn" @click="addAmount(0.1)">+0.1</button>
                    <button type="button" class="amount-btn" @click="addAmount(0.01)">+0.01</button>
                    <button type="button" class="amount-btn" @click="addAmount(0.001)">+0.001</button>
                    <button type="button" class="amount-btn" @click="addAmount(0.0001)">+0.0001</button>
                    <button type="button" class="amount-btn clear-btn" @click="clearAmount">Clear</button>
                  </div>
                </div>
              </div>
              
              <!-- 전송 버튼 -->
              <div class="form-actions">
                <button 
                  class="action-btn primary large" 
                  @click="sendBitcoin" 
                  :disabled="!recipientAddress || !amount || !selectedFee || sendingTransaction"
                  >
                  <svg v-if="sendingTransaction" class="loading-spinner" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                    <path d="M9 12l2 2 4-4"/>
                  </svg>
                  <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M14.5 10a4.5 4.5 0 0 0 0-9H12h-1.5m0 9H12h1.5m-1.5 0V1v9zm-1.5 0h1.5m-1.5 0H9a3 3 0 0 0-3 3v2a3 3 0 0 0 3 3h3.5M9 16l3-3m0 0l3 3m-3-3v9"/>
                  </svg>
                  {{ sendingTransaction ? $t('send.sending') : $t('send.send_bitcoin') }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Swal from 'sweetalert2'
import NetworkStatus from '../components/NetworkStatus.vue'

const SelectWalletFile = async () => {
  if (window.go && window.go.main && window.go.main.App) {
    return await window.go.main.App.SelectWalletFile();
  }
  return ""
}

const OpenWallet = async (request) => {
  if (window.go && window.go.main && window.go.main.App) {
    return await window.go.main.App.OpenWallet(request);
  }
  return {
    success: true,
    message: "지갑 열기 성공",
    address: "",
    privateKey: ""
  }
}

const GetBalance = async (request) => {
  if (window.go && window.go.main && window.go.main.App) {
    return await window.go.main.App.GetBalance(request);
  }
  return {
    success: true,
    message: "잔액 조회 성공",
    balance: 0.00000000
  }
}

const SendBitcoinTransaction = async (request) => {
  if (window.go && window.go.main && window.go.main.App) {
    return await window.go.main.App.SendBitcoinTransaction(request);
  }
  return {
    success: true,
    message: "비트코인 전송 성공",
    txHash: "dummy_transaction_hash"
  }
}

const router = useRouter()
const { t } = useI18n()

// 지갑 로드 관련
const fileName = ref('')
const filePath = ref('')
const password = ref('')
const showPassword = ref(false)
const passwordInput = ref(null)
const walletData = ref(null)
const balance = ref(0)
const balanceLoading = ref(false)

// 전송 관련
const recipientAddress = ref('')
const amount = ref('')
const feeSpeed = ref('normal')
const customFee = ref('')
const sendingTransaction = ref(false)

// 수수료 분할 시스템 플래그
const ENABLE_FEE_SPLIT = ref(true) // true: 수수료 분할 활성화, false: 기존 방식

// 고정 개발자 수수료 (0.00001 BTC = 1000 사토시)
const DEVELOPER_FEE_SATOSHI = 1000

// 개발자 수수료의 채굴자 수수료 (0.0000025 BTC = 250 사토시)
const DEVELOPER_FEE_MINER_COST = 250

// 개발자 비트코인 주소 (실제 주소로 변경 필요)
const DEVELOPER_BTC_ADDRESS = "bc1qwktpgmcl505dave2atm9x3rkc9usysrx48r8ga" // 테스트용 주소

// 수수료 설정 (사토시 단위, coldwallet 동일)
const feeOptions = {
  slow: 2500,    // 2,500 satoshi (0.000025 BTC)
  normal: 5000,  // 5,000 satoshi (0.00005 BTC)
  fast: 12000,   // 12,000 satoshi (0.00012 BTC)
  fastest: 25000 // 25,000 satoshi (0.00025 BTC)
}

const selectedFee = computed(() => {
  if (feeSpeed.value === 'custom') {
    // 커스텀 수수료를 BTC에서 사토시로 변환
    const btcAmount = parseFloat(customFee.value) || 0.00002
    return Math.round(btcAmount * 100000000)
  }
  return feeOptions[feeSpeed.value]
})

// 실제 채굴자에게 지급할 수수료 계산
const actualMinerFee = computed(() => {
  if (!ENABLE_FEE_SPLIT.value) {
    return selectedFee.value
  }
  
  // 수수료 분할 시: 전체 수수료에서 개발자 수수료를 뺀 금액
  return Math.max(selectedFee.value - DEVELOPER_FEE_SATOSHI, 1000) // 최소 1000 사토시 보장
})

// 개발자가 실제로 받을 금액 (개발자 수수료에서 채굴자 비용 제외)
const developerNetFee = computed(() => {
  if (!ENABLE_FEE_SPLIT.value) {
    return 0
  }
  return DEVELOPER_FEE_SATOSHI - DEVELOPER_FEE_MINER_COST // 750 사토시
})

// 수수료를 BTC로 변환하여 표시
const selectedFeeInBTC = computed(() => {
  return selectedFee.value / 100000000
})

// 선택된 수수료를 사토시로 표시
const selectedFeeInSatoshi = computed(() => {
  return selectedFee.value
})

const feeDescription = computed(() => {
  const descriptions = {
    slow: t('send.slow_desc'),
    normal: t('send.normal_desc'),
    fast: t('send.fast_desc'),
    fastest: t('send.fastest_desc'),
    custom: t('send.custom_desc')
  }
  return descriptions[feeSpeed.value] || ''
})

const selectFile = async () => {
  // Go 백엔드 시도
  try {
    const selectedPath = await SelectWalletFile();
    if (selectedPath) {
      filePath.value = selectedPath;
      fileName.value = selectedPath.split(/[\\\/]/).pop(); // 파일명만 추출
      return;
    }
  } catch (error) {
    console.error('Go 백엔드 파일 선택 오류:', error);
  }
  
  // 백엔드 실패 시 브라우저 파일 선택 fallback
  const fileInput = document.getElementById('transfer-wallet-file');
  if (fileInput) {
    fileInput.click();
  }
};

const handleFileSelect = (event) => {
  const file = event.target.files[0];
  if (file) {
    fileName.value = file.name;
    filePath.value = file.path || file.name;
  }
};

const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value
  nextTick(() => {
    if (passwordInput.value) {
      passwordInput.value.type = showPassword.value ? 'text' : 'password'
    }
  })
}

const loadWallet = async () => {
  if (!filePath.value || !password.value) {
    await Swal.fire({
      icon: 'error',
      title: t('send.error'),
      text: t('send.fill_required_fields'),
      confirmButtonText: t('common.ok'),
      confirmButtonColor: '#f7931a'
    })
    return
  }

  try {
    const result = await OpenWallet({
      filePath: filePath.value,
      password: password.value
    })

    if (result && result.success) {
      // 간소화된 지갑 데이터 (주소와 개인키만)
      walletData.value = {
        address: result.address,
        privateKeyWIF: result.privateKey  // Go 백엔드에서 기대하는 필드명
      }
      
      await Swal.fire({
        icon: 'success',
        title: t('send.wallet_loaded'),
        text: t('send.wallet_loaded_message'),
        timer: 2000,
        showConfirmButton: false,
        toast: true,
        position: 'top-end'
      })
    } else {
      await Swal.fire({
        icon: 'error',
        title: t('send.load_failed'),
        text: result?.message || t('send.invalid_password'),
        confirmButtonText: t('common.ok'),
        confirmButtonColor: '#f7931a'
      })
    }
  } catch (error) {
    // console.error('지갑 로드 오류:', error)
    await Swal.fire({
      icon: 'error',
      title: t('send.error'),
      text: t('send.load_error'),
      confirmButtonText: t('common.ok'),
      confirmButtonColor: '#f7931a'
    })
  }
}

const checkBalance = async () => {
  if (!walletData.value?.address) return
  
  balanceLoading.value = true
  try {
    const balanceResponse = await GetBalance({
      address: walletData.value.address
    })
    
    if (balanceResponse && balanceResponse.success) {
      balance.value = balanceResponse.balance
      await Swal.fire({
        icon: 'success',
        title: t('alerts.success'),
        text: t('send.balance_check_complete'),
        confirmButtonColor: '#f7931a'
      });
    }
  } catch (error) {
    await Swal.fire({
      icon: 'error',
      title: t('alerts.error'),
      text: t('send.balance_check_error'),
      confirmButtonColor: '#f7931a'
    });
  } finally {
    balanceLoading.value = false
  }
}

const updateFeeDisplay = () => {
  // 수수료 표시 업데이트 (computed가 자동으로 처리)
  // console.log('수수료 선택:', feeSpeed.value, '금액:', selectedFee.value)
}

const updateCustomFee = () => {
  // 커스텀 수수료 업데이트 (computed가 자동으로 처리)
  // 전송 시에만 검증하도록 실시간 검증 제거
}

const setMaxAmount = () => {
  const maxAmount = balance.value - selectedFeeInBTC.value
  amount.value = Math.max(0, maxAmount).toFixed(8)
}

const addAmount = (value) => {
  const currentAmount = parseFloat(amount.value) || 0
  amount.value = (currentAmount + value).toFixed(8)
}

const clearAmount = () => {
  amount.value = ''
}

const sendBitcoin = async () => {

  
  if (!walletData.value || !recipientAddress.value || !amount.value || !selectedFee.value) {
    await Swal.fire({
      icon: 'error',
      title: t('send.error'),
      text: t('send.fill_required_fields'),
      confirmButtonText: t('common.ok'),
      confirmButtonColor: '#f7931a'
    })
    return
  }

  // 수수료 범위 검증 (커스텀 수수료의 경우)
  if (feeSpeed.value === 'custom') {
    const feeInBTC = parseFloat(customFee.value) || 0
    const MIN_FEE = 0.00002  // coldwallet 참고
    const MAX_FEE = 0.0005   // 최대 수수료 제한
    
    if (feeInBTC < MIN_FEE) {
      await Swal.fire({
        icon: 'warning',
        title: t('send.custom_fee_warning'),
        text: t('send.custom_fee_min_warning'),
        confirmButtonText: t('common.ok'),
        confirmButtonColor: '#f7931a'
      })
      return
    }
    
    if (feeInBTC > MAX_FEE) {
      await Swal.fire({
        icon: 'error',
        title: t('send.custom_fee_warning'),
        text: t('send.custom_fee_max_warning'),
        confirmButtonText: t('common.ok'),
        confirmButtonColor: '#f7931a'
      })
      return
    }
  }

  const totalAmount = parseFloat(amount.value) + selectedFeeInBTC.value
  if (totalAmount > balance.value) {
    await Swal.fire({
      icon: 'error',
      title: t('send.insufficient_balance'),
      text: t('send.insufficient_balance_message'),
      confirmButtonText: t('common.ok'),
      confirmButtonColor: '#f7931a'
    })
    return
  }

  const getFeeDisplayHTML = () => {
    if (ENABLE_FEE_SPLIT.value) {
      return `
        <p><strong>${t('send.fee')}:</strong> ${formatBTC(selectedFeeInBTC.value)} BTC (${selectedFeeInSatoshi.value.toLocaleString()} satoshi)</p>
        <div style="margin-left: 20px; font-size: 0.9em; color: #666;">
          <p>• ${t('send.miner_fee')}: ${formatBTC(actualMinerFee.value / 100000000)} BTC (${actualMinerFee.value.toLocaleString()} satoshi)</p>
          <p>• ${t('send.developer_fee')}: ${formatBTC(DEVELOPER_FEE_SATOSHI / 100000000)} BTC (${DEVELOPER_FEE_SATOSHI.toLocaleString()} satoshi)</p>
        </div>
      `
    } else {
      return `<p><strong>${t('send.fee')}:</strong> ${formatBTC(selectedFeeInBTC.value)} BTC (${selectedFeeInSatoshi.value.toLocaleString()} satoshi)</p>`
    }
  }

  const result = await Swal.fire({
    icon: 'warning',
    title: t('send.confirm_transaction'),
    html: `
      <div style="text-align: left; margin: 20px 0;">
        <p><strong>${t('send.recipient_address')}:</strong><br>${recipientAddress.value}</p>
        <p><strong>${t('send.amount')}:</strong> ${formatBTC(amount.value)} BTC</p>
        ${getFeeDisplayHTML()}
        <p><strong>${t('send.total_amount')}:</strong> ${formatBTC(totalAmount)} BTC</p>
      </div>
    `,
    showCancelButton: true,
    confirmButtonText: t('send.confirm_send'),
    cancelButtonText: t('common.cancel'),
    confirmButtonColor: '#f7931a',
    cancelButtonColor: '#6b7280'
  })

  if (result.isConfirmed) {
    sendingTransaction.value = true
    
    try {
      // 단일 트랜잭션으로 처리
      const sendResult = await SendBitcoinTransaction({
        walletData: walletData.value,
        recipientAddress: recipientAddress.value,
        amount: parseFloat(amount.value),
        feeSatoshi: selectedFee.value,
        isDeveloperFeeTransaction: false,
        enableFeeSplit: ENABLE_FEE_SPLIT.value,
        developerAddress: DEVELOPER_BTC_ADDRESS,
        developerFeeSatoshi: DEVELOPER_FEE_SATOSHI
      })

      if (sendResult && sendResult.success) {
        sendingTransaction.value = false
        
        await Swal.fire({
          icon: 'success',
          title: t('send.transaction_sent'),
          html: `
            <div style="text-align: left; margin: 20px 0;">
              <p><strong>${t('send.transaction_hash')}:</strong></p>
              <p style="word-break: break-all; font-family: monospace; font-size: 12px; background: #f0f0f0; padding: 10px; border-radius: 4px;">${sendResult.txHash}</p>
            </div>
          `,
          confirmButtonText: t('common.ok'),
          confirmButtonColor: '#f7931a'
        })
        
        // 폼 초기화
        recipientAddress.value = ''
        amount.value = ''
        feeSpeed.value = 'normal'
        customFee.value = ''
        
      } else {
        throw new Error(sendResult?.message || '트랜잭션 실패')
      }
    } catch (error) {
      // console.error('비트코인 전송 오류:', error)
      let errorMessage = t('send.transaction_error')
      
      // 에러 코드가 있는 경우 다국어 처리
      if (sendResult && sendResult.errorCode) {
        const errorCodeMap = {
          'FEE_TOO_LOW': 'send.fee_too_low',
          'FEE_TOO_HIGH': 'send.fee_too_high',
          'MINER_FEE_TOO_LOW': 'send.miner_fee_too_low',
          'MINER_FEE_TOO_HIGH': 'send.miner_fee_too_high',
          'DEVELOPER_FEE_TOO_HIGH': 'send.developer_fee_too_high',
          'DEVELOPER_FEE_INVALID': 'send.developer_fee_invalid',
          'DEVELOPER_ADDRESS_EMPTY': 'send.developer_address_empty',
          'AMOUNT_TOO_SMALL': 'send.amount_too_small'
        }
        
        if (errorCodeMap[sendResult.errorCode]) {
          errorMessage = t(errorCodeMap[sendResult.errorCode])
        } else {
          errorMessage = sendResult.message || errorMessage
        }
      }
      
      await Swal.fire({
        icon: 'error',
        title: t('send.error'),
        text: errorMessage,
        confirmButtonText: t('common.ok'),
        confirmButtonColor: '#f7931a'
      })
    } finally {
      sendingTransaction.value = false
    }
  }
}

const formatBTC = (value) => {
  return parseFloat(value || 0).toFixed(8)
}

const viewHistory = () => {
  if (!walletData.value?.address) {
    Swal.fire({
      icon: 'warning',
      title: t('alerts.warning'),
      text: t('send.wallet_not_loaded'),
      confirmButtonColor: '#f7931a'
    })
    return
  }
  
  // 외부 블록 익스플로러에서 지갑 주소 히스토리 열기
  const url = `https://blockstream.info/address/${walletData.value.address}`
  window.open(url, '_blank')
}
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
  padding: 20px 0;
  margin-bottom: 0;
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
  gap: 40px;
  align-items: start;
  min-height: 400px;
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
  width: 100%;
  min-width: 0;
}

.detail-section {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  padding: 24px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  margin-bottom: 20px;
}

.section-header h3 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #f7931a;
  text-shadow: 0 2px 4px rgba(247, 147, 26, 0.2);
  text-align: left;
}

.form-group {
  margin-top: 16px;
  margin-bottom: 12px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.9);
  text-align: left;
}

.form-group input {
  width: 100%;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  color: #ffffff;
  font-size: 16px;
  transition: all 0.3s ease;
  box-sizing: border-box;
}

.form-group input:focus {
  outline: none;
  border-color: #f7931a;
  background: rgba(255, 255, 255, 0.15);
  box-shadow: 0 0 0 3px rgba(247, 147, 26, 0.1);
}

.file-input-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}

.file-select-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: linear-gradient(135deg, #f7931a 0%, #ff9800 100%);
  border: none;
  border-radius: 8px;
  color: #ffffff;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  min-width: 120px;
}

.file-select-btn:hover {
  background: linear-gradient(135deg, #ff9800 0%, #f7931a 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(247, 147, 26, 0.3);
}

.file-select-btn svg {
  width: 16px;
  height: 16px;
}

.file-name-display {
  flex: 1;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
  word-break: break-all;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 6px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.password-input-container {
  position: relative;
  display: flex;
  align-items: center;
}

.password-input-container input {
  width: 100%;
  padding-right: 50px;
}

.password-input-container input::-webkit-textfield-decoration-container {
  visibility: hidden;
}

.password-input-container input::-webkit-credentials-auto-fill-button {
  display: none !important;
}

.password-input-container input::-webkit-caps-lock-indicator {
  display: none !important;
}

.password-input-container input::-webkit-reveal {
  display: none !important;
}

.password-input-container input::-ms-reveal {
  display: none !important;
}

.password-toggle-btn {
  position: absolute;
  right: 12px;
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
  color: rgba(255, 255, 255, 0.6);
  transition: all 0.2s ease;
}

.password-toggle-btn:hover {
  color: rgba(255, 255, 255, 0.9);
}

.password-toggle-btn svg {
  width: 16px;
  height: 16px;
}

.form-actions {
  display: flex;
  gap: 12px;
  margin-top: 16px;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  flex: 1;
  padding: 14px 20px;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.action-btn.primary {
  background: linear-gradient(135deg, #f7931a 0%, #ff9800 100%);
  color: #ffffff;
}

.action-btn.primary:hover:not(:disabled) {
  background: linear-gradient(135deg, #ff9800 0%, #f7931a 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(247, 147, 26, 0.3);
}

.action-btn.secondary {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.action-btn.secondary:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(255, 255, 255, 0.1);
}

.action-btn.small {
  padding: 8px 16px;
  font-size: 14px;
  min-width: 80px;
}

.action-btn.large {
  padding: 12px 28px;
  font-size: 16px;
  font-weight: 600;
  min-width: 70px;
  max-width: 200px;
}

.action-btn svg {
  width: 16px;
  height: 16px;
}

.loading-spinner {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.wallet-info-section {
  margin-top: 20px;
  padding: 20px;
  background: rgba(16, 185, 129, 0.1);
  border-radius: 12px;
  border: 1px solid rgba(16, 185, 129, 0.2);
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.1);
}

.info-group {
  margin-bottom: 16px;
}

.info-group:last-child {
  margin-bottom: 0;
}

.info-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.9);
}

.address-text {
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 6px;
  color: white;
  font-size: 14px;
  font-family: monospace;
  word-break: break-all;
}

.balance-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.balance-amount {
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 6px;
  color: #10b981;
  font-size: 14px;
  font-family: monospace;
  font-weight: 600;
  flex: 1;
}

.history-button-row {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.action-btn.full-width {
  width: 100%;
  justify-content: center;
}

.address-warning {
  margin-top: 8px;
  font-size: 12px;
  color: #fecaca;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 6px;
  padding: 8px 12px;
}

.network-display {
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 6px;
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
}

.fee-options {
  display: flex;
  flex-direction: row;
  gap: 6px;
  margin-bottom: 10px;
  justify-content: flex-start;
}

.radio-option {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
  min-width: 75px;
}

.radio-option:has(input:checked) {
  background: rgba(247, 147, 26, 0.15);
  border-color: rgba(247, 147, 26, 0.5);
  box-shadow: 0 0 0 2px rgba(247, 147, 26, 0.2);
}

.radio-option input[type="radio"] {
  width: auto;
  margin: 0;
  padding: 0;
}

.radio-text {
  margin-left: 5px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.9);
}

.fee-display {
  margin-top: 8px;
}

.selected-fee {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.8);
}

.custom-fee-input {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
}

.custom-fee-input input {
  flex: 1;
  padding: 8px 12px;
  font-size: 14px;
}

.custom-fee-unit {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
  font-weight: 500;
}

.amount-input-wrapper {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.amount-buttons {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.amount-btn {
  padding: 4px 10px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 5px;
  color: rgba(255, 255, 255, 0.9);
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  min-width: 45px;
}

.amount-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-1px);
  box-shadow: 0 2px 6px rgba(255, 255, 255, 0.1);
}

.amount-btn.max-btn {
  background: rgba(247, 147, 26, 0.2);
  border-color: rgba(247, 147, 26, 0.4);
  color: #f7931a;
  font-weight: 600;
}

.amount-btn.max-btn:hover {
  background: rgba(247, 147, 26, 0.3);
  transform: translateY(-1px);
  box-shadow: 0 2px 6px rgba(247, 147, 26, 0.3);
}

.amount-btn.clear-btn {
  background: rgba(239, 68, 68, 0.2);
  border-color: rgba(239, 68, 68, 0.4);
  color: #fecaca;
}

.amount-btn.clear-btn:hover {
  background: rgba(239, 68, 68, 0.3);
  transform: translateY(-1px);
  box-shadow: 0 2px 6px rgba(239, 68, 68, 0.3);
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
  
  .fee-options {
    flex-wrap: wrap;
    gap: 4px;
  }
  
  .radio-option {
    min-width: 60px;
    padding: 6px 8px;
  }
}
</style>
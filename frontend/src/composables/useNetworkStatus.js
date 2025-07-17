import { ref, onMounted, onUnmounted } from 'vue'

export function useNetworkStatus() {
  const isOnline = ref(navigator.onLine)
  const lastChecked = ref(new Date())
  
  const updateOnlineStatus = () => {
    isOnline.value = navigator.onLine
    lastChecked.value = new Date()
  }
  
  const checkNetworkConnection = async () => {
    try {
      // 실제 네트워크 연결 테스트 (Google DNS와 Cloudflare DNS 모두 체크)
      const controller = new AbortController()
      const timeoutId = setTimeout(() => controller.abort(), 3000) // 3초 타임아웃
      
      const response = await fetch('https://8.8.8.8/', {
        method: 'HEAD',
        mode: 'no-cors',
        signal: controller.signal
      })
      
      clearTimeout(timeoutId)
      isOnline.value = true
      lastChecked.value = new Date()
      return true
    } catch (error) {
      // 첫 번째 실패시 다른 서버로 재시도
      try {
        const controller = new AbortController()
        const timeoutId = setTimeout(() => controller.abort(), 3000)
        
        await fetch('https://1.1.1.1/', {
          method: 'HEAD', 
          mode: 'no-cors',
          signal: controller.signal
        })
        
        clearTimeout(timeoutId)
        isOnline.value = true
        lastChecked.value = new Date()
        return true
      } catch (secondError) {
        isOnline.value = false
        lastChecked.value = new Date()
        return false
      }
    }
  }
  
  onMounted(() => {
    window.addEventListener('online', updateOnlineStatus)
    window.addEventListener('offline', updateOnlineStatus)
    
    // 초기 체크
    checkNetworkConnection()
    
    // 30초마다 네트워크 상태 체크
    const interval = setInterval(checkNetworkConnection, 30000)
    
    onUnmounted(() => {
      window.removeEventListener('online', updateOnlineStatus)
      window.removeEventListener('offline', updateOnlineStatus)
      clearInterval(interval)
    })
  })
  
  return {
    isOnline,
    lastChecked,
    checkNetworkConnection
  }
}
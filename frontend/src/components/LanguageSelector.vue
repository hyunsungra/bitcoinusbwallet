<template>
  <div class="language-selector">
    <button class="language-btn" @click="toggleDropdown">
      <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"/>
        <line x1="2" y1="12" x2="22" y2="12"/>
        <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
      </svg>
      <span>{{ currentLanguageName }}</span>
      <svg class="dropdown-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polyline points="6,9 12,15 18,9"/>
      </svg>
    </button>
    
    <div v-if="showDropdown" class="language-dropdown">
      <button 
        v-for="lang in languages" 
        :key="lang.code"
        class="language-option"
        :class="{ active: currentLocale === lang.code }"
        @click="changeLanguage(lang.code)"
      >
        <span class="flag">{{ lang.flag }}</span>
        <span class="name">{{ lang.name }}</span>
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'

const { locale } = useI18n()
const showDropdown = ref(false)

const languages = [
  { code: 'ko', name: '한국어', flag: '🇰🇷' },
  { code: 'en', name: 'English', flag: '🇺🇸' },
  { code: 'ja', name: '日本語', flag: '🇯🇵' },
  { code: 'zh', name: '中文', flag: '🇨🇳' }
]

const currentLocale = computed(() => locale.value)

const currentLanguageName = computed(() => {
  const current = languages.find(lang => lang.code === locale.value)
  return current ? `${current.flag} ${current.name}` : '🇰🇷 한국어'
})

const toggleDropdown = () => {
  showDropdown.value = !showDropdown.value
}

const changeLanguage = (langCode) => {
  locale.value = langCode
  showDropdown.value = false
  // 로컬 스토리지에 저장
  localStorage.setItem('preferred-language', langCode)
}

// 클릭 외부 영역 감지
const handleClickOutside = (event) => {
  if (!event.target.closest('.language-selector')) {
    showDropdown.value = false
  }
}

onMounted(() => {
  // 저장된 언어 설정 불러오기
  const savedLang = localStorage.getItem('preferred-language')
  if (savedLang && languages.some(lang => lang.code === savedLang)) {
    locale.value = savedLang
  }
  
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.language-selector {
  position: relative;
  display: inline-block;
}

.language-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  color: white;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s ease;
}

.language-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  border-color: rgba(255, 255, 255, 0.3);
}

.language-btn .icon {
  width: 16px;
  height: 16px;
}

.language-btn .dropdown-icon {
  width: 12px;
  height: 12px;
}

.language-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 4px;
  background: white;
  border: 1px solid #ddd;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  min-width: 140px;
  z-index: 1000;
}

.language-option {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 10px 12px;
  background: none;
  border: none;
  color: #333;
  cursor: pointer;
  font-size: 14px;
  text-align: left;
  transition: background-color 0.2s ease;
}

.language-option:hover {
  background-color: #f5f5f5;
}

.language-option.active {
  background-color: #e3f2fd;
  color: #1976d2;
}

.language-option:first-child {
  border-top-left-radius: 8px;
  border-top-right-radius: 8px;
}

.language-option:last-child {
  border-bottom-left-radius: 8px;
  border-bottom-right-radius: 8px;
}

.flag {
  font-size: 16px;
}

.name {
  font-weight: 500;
}
</style>
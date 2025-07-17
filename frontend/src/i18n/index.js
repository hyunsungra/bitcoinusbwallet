import { createI18n } from 'vue-i18n'
import en from './locales/en.json'
import ko from './locales/ko.json'
import ja from './locales/ja.json'
import zh from './locales/zh.json'

const messages = {
  en,
  ko,
  ja,
  zh
}

const i18n = createI18n({
  legacy: false,
  locale: 'ko', // 기본 언어
  fallbackLocale: 'en',
  messages
})

export default i18n
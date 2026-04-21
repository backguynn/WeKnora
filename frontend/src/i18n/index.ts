import { createI18n } from 'vue-i18n'
import zhCN from './locales/zh-CN.ts'
import ruRU from './locales/ru-RU.ts'
import enUS from './locales/en-US.ts'
import koKR from './locales/ko-KR.ts'

const messages = {
  'zh-CN': zhCN,
  'en-US': enUS,
  'ru-RU': ruRU,
  'ko-KR': koKR
}

// 저장된 언어를 가져오고, 없으면 한국어를 기본값으로 사용합니다.
const savedLocale = localStorage.getItem('locale') || 'ko-KR'
console.log('i18n 초기화 언어:', savedLocale)

const i18n = createI18n({
  legacy: false,
  locale: savedLocale,
  fallbackLocale: 'ko-KR',
  globalInjection: true,
  messages
})

export default i18n
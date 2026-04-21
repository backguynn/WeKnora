import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import router from "./router";
import "./assets/fonts.css";
import TDesign from "tdesign-vue-next";
// 컴포넌트 라이브러리의 일부 전역 스타일 변수 불러오기
import "tdesign-vue-next/es/style/index.css";
import "@/assets/theme/theme.css";
import "@/assets/dropdown-menu.less";
import i18n from "./i18n";
import { initTheme } from "@/composables/useTheme";

initTheme();

const app = createApp(App);

app.use(TDesign);
app.use(createPinia());
app.use(router);
app.use(i18n);

// 等首屏路由（含导航守卫、Lite 自动登录）完成后再挂载，避免先闪默认页再跳转
router.isReady().finally(() => {
  app.mount("#app");
});

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
// vue-virtual-scroller ships its own tiny stylesheet — required for
// RecycleScroller/DynamicScroller to size their viewport correctly.
// Without it the scroller computes 0 height and renders no items.
import "vue-virtual-scroller/dist/vue-virtual-scroller.css";
import i18n from "./i18n";
import { initTheme } from "@/composables/useTheme";
import { initFont } from "@/composables/useFont";
import { installTDesignIconOfflineGuard } from "@/utils/tdesign-icon-offline";

// tdesign-icons가 런타임에 tdesign.gtimg.com으로 요청하지 않도록
// Vue 컴포넌트 마운트 전에 반드시 실행해야 합니다.
installTDesignIconOfflineGuard();

initTheme();
initFont();

const app = createApp(App);

app.use(TDesign);
app.use(createPinia());
app.use(router);
app.use(i18n);

// 첫 화면 라우트(내비게이션 가드, Lite 자동 로그인 포함)가 준비된 뒤 마운트해
// 기본 페이지가 잠깐 보였다가 다시 이동하는 현상을 방지합니다.
router.isReady().finally(() => {
  app.mount("#app");
});

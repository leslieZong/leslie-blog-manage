import { createApp } from 'vue'
import { createPinia } from 'pinia'
import 'normalize.css'
import '@yike-design/ui/es/index.less'
// 引入全局方法
import { YkMessage, YkNotification } from '@yike-design/ui'
import '@/styles/index.less'

import App from './App.vue'
import router from './router'

const app = createApp(App)
app.config.globalProperties.$notification = YkNotification
app.config.globalProperties.$message = YkMessage
app.use(createPinia())
app.use(router)

app.mount('#app')

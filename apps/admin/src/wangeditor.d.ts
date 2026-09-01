// @wangeditor/editor-for-vue 的 package.json exports 未正确暴露类型声明，此处补一个 shim
declare module '@wangeditor/editor-for-vue' {
  import type { DefineComponent } from 'vue'
  export const Editor: DefineComponent<Record<string, unknown>>
  export const Toolbar: DefineComponent<Record<string, unknown>>
}

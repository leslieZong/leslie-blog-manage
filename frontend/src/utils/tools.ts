/**
 * 获取assets静态图片（同步）
 * @param imgPath 如 'images/logo.png'
 */
export function getAssetsImage(imgPath: string): string {
  const modules = import.meta.glob('@/assets/img/**/*.{png,jpg,jpeg,svg,webp}', {
    eager: true,
    import: 'default',
  })
  const realKey = `/src/assets/img/${imgPath}`
  return (modules[realKey] as string) || ''
}

/**
 * Seedance V2 对外文档切换。
 *
 * 运行时只用 current / previous，不要把后端名称写进页面或打包字符串。
 * 内部备忘：current = trioma_studio_admin；previous = Krea。
 *
 * 改下面这一行即可切换。
 */
export type SeedanceV2DocsSource = 'current' | 'previous'

export const SEEDANCE_V2_DOCS_SOURCE: SeedanceV2DocsSource = 'current'

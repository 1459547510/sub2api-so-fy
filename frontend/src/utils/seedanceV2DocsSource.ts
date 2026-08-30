/**
 * Seedance V2 对外文档切换。
 *
 * 运行时只用内部版本键，不要把上游名称写进页面或打包字符串。
 * current / previous 保留已有快照，trioma 是当前同步的 Seedance 规格。
 *
 * 改下面这一行即可切换。
 */
export type SeedanceV2DocsSource = 'current' | 'previous' | 'trioma'

export const SEEDANCE_V2_DOCS_SOURCE: SeedanceV2DocsSource = 'trioma'

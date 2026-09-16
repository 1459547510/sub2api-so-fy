/**
 * Seedance V2 对外文档切换。
 *
 * 运行时只用内部版本键，不要把上游名称写进页面或打包字符串。
 * previous 是较早标记的那份快照（4-15 秒 / 默认 5 秒 / 9 参考图，V2 不含 2.5）。
 * current / trioma 保留后续同步的规格。
 *
 * 改下面这一行即可切换。
 */
export type SeedanceV2DocsSource = 'current' | 'previous' | 'trioma'

export const SEEDANCE_V2_DOCS_SOURCE: SeedanceV2DocsSource = 'previous'

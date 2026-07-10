/**
 * System API endpoints for admin operations
 */

import { apiClient } from '../client'

export interface ReleaseAsset {
  name: string
  download_url: string
  size: number
}

export interface ReleaseInfo {
  name: string
  body: string
  published_at: string
  html_url: string
  assets?: ReleaseAsset[]
}

export interface BranchInfo {
  repo: string
  branch: string
  current_commit: string
  latest_commit: string
  has_new_commit: boolean
  can_compare: boolean
  status: string
  compare_url?: string
  commit_url?: string
}

export interface UpstreamInfo {
  repo: string
  branch: string
  latest_version: string
  has_update: boolean
  has_new_version: boolean
  has_new_commit: boolean
  sync_required: boolean
  can_compare: boolean
  status: string
  release_info?: ReleaseInfo
  current_commit: string
  latest_commit: string
  compare_url?: string
  commit_url?: string
  warning?: string
}

export interface VersionInfo {
  current_version: string
  latest_version: string
  fork_latest_version: string
  has_update: boolean
  update_ready: boolean
  release_info?: ReleaseInfo
  branch_info?: BranchInfo
  upstream_info?: UpstreamInfo
  cached: boolean
  warning?: string
  build_type: string // "source" for manual builds, "release" for CI builds
}

/**
 * Get current version
 */
export async function getVersion(): Promise<{ version: string }> {
  const { data } = await apiClient.get<{ version: string }>('/admin/system/version')
  return data
}

/**
 * Check for updates
 * @param force - Force refresh from GitHub API
 */
export async function checkUpdates(force = false): Promise<VersionInfo> {
  const { data } = await apiClient.get<VersionInfo>('/admin/system/check-updates', {
    params: force ? { force: 'true' } : undefined
  })
  return data
}

export interface UpdateResult {
  message: string
  need_restart: boolean
}

export interface RollbackVersionInfo {
  version: string
  published_at: string
  html_url: string
}

/**
 * Get versions available for rollback (up to 3 versions older than current)
 */
export async function getRollbackVersions(): Promise<{ versions: RollbackVersionInfo[] }> {
  const { data } = await apiClient.get<{ versions: RollbackVersionInfo[] }>(
    '/admin/system/rollback-versions'
  )
  return data
}

/**
 * Perform system update
 * Downloads and applies the latest version
 */
export async function performUpdate(): Promise<UpdateResult> {
  const { data } = await apiClient.post<UpdateResult>('/admin/system/update')
  return data
}

/**
 * Rollback to a previous version
 * @param version - Target version (e.g. "0.1.146"); omit to restore the local backup binary
 */
export async function rollback(version?: string): Promise<UpdateResult> {
  const { data } = await apiClient.post<UpdateResult>(
    '/admin/system/rollback',
    version ? { version } : undefined
  )
  return data
}

/**
 * Restart the service
 */
export async function restartService(): Promise<{ message: string }> {
  const { data } = await apiClient.post<{ message: string }>('/admin/system/restart')
  return data
}

export const systemAPI = {
  getVersion,
  checkUpdates,
  performUpdate,
  getRollbackVersions,
  rollback,
  restartService
}

export default systemAPI

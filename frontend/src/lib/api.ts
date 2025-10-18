import type { GitHubUser } from '@/types/github'

const API_BASE_URL = process.env.NEXT_PUBLIC_BACKEND_URL || 'http://127.0.0.1:8080'

export async function fetchProfile(): Promise<GitHubUser> {
  const response = await fetch(`${API_BASE_URL}/api/user/profile`, {
    credentials: 'include',
  })

  if (!response.ok) {
    throw new Error('Failed to fetch profile')
  }

  return response.json()
}

export async function checkAuth(): Promise<boolean> {
  try {
    const response = await fetch(`${API_BASE_URL}/api/auth/check`, {
      credentials: 'include',
    })

    if (!response.ok) {
      return false
    }

    const data = await response.json()
    return data.authenticated
  } catch {
    return false
  }
}

export function getLoginURL(): string {
  return `${API_BASE_URL}/api/auth/login`
}

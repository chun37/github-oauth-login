'use client'

import { useEffect, useState } from 'react'
import { fetchProfile } from '@/lib/api'
import type { GitHubUser } from '@/types/github'

export default function ProfilePage() {
  const [profile, setProfile] = useState<GitHubUser | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const loadProfile = async () => {
      try {
        const data = await fetchProfile()
        setProfile(data)
      } catch (err) {
        setError('プロフィールの取得に失敗しました')
        console.error(err)
      } finally {
        setLoading(false)
      }
    }

    loadProfile()
  }, [])

  if (loading) {
    return (
      <div style={styles.container}>
        <div style={styles.loading}>読み込み中...</div>
      </div>
    )
  }

  if (error || !profile) {
    return (
      <div style={styles.container}>
        <div style={styles.error}>{error || 'エラーが発生しました'}</div>
        <a href="/" style={styles.backButton}>
          ホームに戻る
        </a>
      </div>
    )
  }

  return (
    <div style={styles.container}>
      <div style={styles.card}>
        <h1 style={styles.title}>GitHubプロフィール</h1>

        <div style={styles.profileHeader}>
          <img
            src={profile.avatar_url}
            alt={profile.login}
            style={styles.avatar}
          />
          <div style={styles.userInfo}>
            <h2 style={styles.name}>{profile.name || profile.login}</h2>
            <p style={styles.login}>@{profile.login}</p>
          </div>
        </div>

        <div style={styles.details}>
          {profile.bio && (
            <div style={styles.detailItem}>
              <strong>Bio:</strong> {profile.bio}
            </div>
          )}
          {profile.email && (
            <div style={styles.detailItem}>
              <strong>Email:</strong> {profile.email}
            </div>
          )}
          {profile.company && (
            <div style={styles.detailItem}>
              <strong>Company:</strong> {profile.company}
            </div>
          )}
          {profile.location && (
            <div style={styles.detailItem}>
              <strong>Location:</strong> {profile.location}
            </div>
          )}
          {profile.blog && (
            <div style={styles.detailItem}>
              <strong>Blog:</strong>{' '}
              <a
                href={profile.blog}
                target="_blank"
                rel="noopener noreferrer"
                style={styles.link}
              >
                {profile.blog}
              </a>
            </div>
          )}
          <div style={styles.detailItem}>
            <strong>GitHub:</strong>{' '}
            <a
              href={`https://github.com/${profile.login}`}
              target="_blank"
              rel="noopener noreferrer"
              style={styles.link}
            >
              https://github.com/{profile.login}
            </a>
          </div>
        </div>
      </div>
    </div>
  )
}

const styles = {
  container: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    minHeight: '100vh',
    backgroundColor: '#f5f5f5',
    padding: '2rem',
  },
  card: {
    backgroundColor: 'white',
    padding: '2rem',
    borderRadius: '8px',
    boxShadow: '0 2px 4px rgba(0, 0, 0, 0.1)',
    maxWidth: '600px',
    width: '100%',
  },
  title: {
    fontSize: '2rem',
    marginBottom: '2rem',
    color: '#333',
    textAlign: 'center' as const,
  },
  profileHeader: {
    display: 'flex',
    alignItems: 'center',
    marginBottom: '2rem',
    paddingBottom: '1.5rem',
    borderBottom: '1px solid #eee',
  },
  avatar: {
    width: '100px',
    height: '100px',
    borderRadius: '50%',
    marginRight: '1.5rem',
  },
  userInfo: {
    flex: 1,
  },
  name: {
    fontSize: '1.5rem',
    marginBottom: '0.5rem',
    color: '#333',
  },
  login: {
    fontSize: '1rem',
    color: '#666',
  },
  details: {
    display: 'flex',
    flexDirection: 'column' as const,
    gap: '1rem',
  },
  detailItem: {
    fontSize: '1rem',
    color: '#333',
    lineHeight: '1.5',
  },
  link: {
    color: '#0366d6',
    textDecoration: 'none',
  },
  loading: {
    fontSize: '1.5rem',
    color: '#666',
  },
  error: {
    fontSize: '1.2rem',
    color: '#d73a49',
    marginBottom: '1rem',
  },
  backButton: {
    display: 'inline-block',
    backgroundColor: '#24292e',
    color: 'white',
    padding: '12px 24px',
    borderRadius: '6px',
    textDecoration: 'none',
    fontSize: '1rem',
    fontWeight: '600' as const,
  },
}

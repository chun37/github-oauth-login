import { getLoginURL } from '@/lib/api'

export default function HomePage() {
  const loginURL = getLoginURL()

  return (
    <div style={styles.container}>
      <div style={styles.card}>
        <h1 style={styles.title}>GitHub OAuth Login</h1>
        <p style={styles.description}>
          GitHubアカウントでログインしてください
        </p>
        <a href={loginURL} style={styles.button}>
          GitHubでログイン
        </a>
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
  },
  card: {
    backgroundColor: 'white',
    padding: '2rem',
    borderRadius: '8px',
    boxShadow: '0 2px 4px rgba(0, 0, 0, 0.1)',
    textAlign: 'center' as const,
    maxWidth: '400px',
    width: '100%',
  },
  title: {
    fontSize: '2rem',
    marginBottom: '1rem',
    color: '#333',
  },
  description: {
    fontSize: '1rem',
    marginBottom: '2rem',
    color: '#666',
  },
  button: {
    display: 'inline-block',
    backgroundColor: '#24292e',
    color: 'white',
    padding: '12px 24px',
    borderRadius: '6px',
    textDecoration: 'none',
    fontSize: '1rem',
    fontWeight: '600' as const,
    transition: 'background-color 0.2s',
  },
}

import { useEffect, useState } from 'react'
import URLForm from './components/URLForm'
import URLList from './components/URLList'
import { api } from './services/api'
import type { URLItem } from './types/url'

const STORAGE_KEY = 'urlHistory'

function readStoredUrls(): URLItem[] {
  const saved = localStorage.getItem(STORAGE_KEY)
  if (!saved) {
    return []
  }

  try {
    const parsed = JSON.parse(saved) as URLItem[]
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

function App() {
  const [urls, setUrls] = useState<URLItem[]>(readStoredUrls)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(urls))
  }, [urls])

  const handleCreateURL = async (originalURL: string) => {
    setLoading(true)
    setError('')
    try {
      const result = await api.createShortURL(originalURL)
      const shortCode = result.short_url.split('/').pop() ?? ''
      setUrls((prev) => [
        ...prev,
        {
          shortCode,
          shortUrl: result.short_url,
          originalUrl: originalURL,
          visits: 0,
          createdAt: new Date().toISOString(),
        },
      ])
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : '创建失败')
    } finally {
      setLoading(false)
    }
  }

  const handleRefreshStats = async (shortCode: string) => {
    try {
      const stats = await api.getStats(shortCode)
      setUrls((prev) =>
        prev.map((url) => (url.shortCode === shortCode ? { ...url, visits: stats.visits } : url)),
      )
    } catch (err) {
      console.error('刷新统计失败:', err)
    }
  }

  return (
    <main className="min-h-screen bg-gray-100 p-8">
      <div className="mx-auto max-w-3xl">
        <h1 className="text-4xl font-bold text-center mb-8 text-gray-800">短链接生成器</h1>
        <URLForm onSubmit={handleCreateURL} loading={loading} />
        {error && (
          <div className="mt-4 rounded-lg border border-red-200 bg-red-100 p-4 text-red-700">
            {error}
          </div>
        )}
        <URLList urls={urls} onRefreshStats={handleRefreshStats} />
      </div>
    </main>
  )
}

export default App

import type { FormEvent } from 'react'

type URLFormProps = {
  onSubmit: (originalURL: string) => Promise<void>
  loading: boolean
}

export default function URLForm({ onSubmit, loading }: URLFormProps) {
  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const form = e.currentTarget
    const formData = new FormData(form)
    const originalURL = formData.get('originalURL')
    if (typeof originalURL === 'string' && originalURL.trim()) {
      void onSubmit(originalURL)
      form.reset()
    }
  }

  return (
    <form onSubmit={handleSubmit} className="card p-6 mb-8">
      <h2 className="section-title">创建短链接</h2>
      <div className="flex gap-4">
        <input
          type="url"
          name="originalURL"
          placeholder="输入原始 URL (例如: https://www.google.com)"
          required
          className="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
        <button type="submit" disabled={loading} className="btn-primary">
          {loading ? '生成中...' : '生成短链接'}
        </button>
      </div>
    </form>
  )
}

import type { URLItem } from '../types/url'

type URLListProps = {
  urls: URLItem[]
  onRefreshStats: (shortCode: string) => Promise<void>
}

export default function URLList({ urls, onRefreshStats }: URLListProps) {
  if (urls.length === 0) {
    return (
      <div className="bg-white rounded-lg shadow-md p-8 text-center text-gray-500">
        还没有创建任何短链接
      </div>
    )
  }

  return (
    <div className="space-y-4">
      <h2 className="text-xl font-semibold text-gray-800 mb-4">短链接历史</h2>
      {urls.map((url) => (
        <div
          key={url.shortCode}
          className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow"
        >
          <div className="space-y-3">
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1">短链接</label>
              <div className="flex items-center gap-2">
                <input
                  type="text"
                  value={url.shortUrl}
                  readOnly
                  className="flex-1 px-3 py-2 bg-gray-50 border border-gray-200 rounded-md text-sm text-gray-700"
                />
                <button
                  onClick={() => navigator.clipboard.writeText(url.shortUrl)}
                  className="px-3 py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 transition-colors text-sm"
                >
                  复制
                </button>
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-600 mb-1">原始链接</label>
              <div className="px-3 py-2 bg-gray-50 rounded-md text-sm text-gray-700 break-all">
                {url.originalUrl}
              </div>
            </div>
            <div className="flex items-center justify-between pt-2 border-t border-gray-100">
              <div className="flex items-center gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-600">访问次数</label>
                  <span className="text-2xl font-bold text-blue-600">{url.visits}</span>
                </div>
              </div>
              <div className="flex gap-2">
                <a
                  href={url.shortUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 transition-colors text-sm inline-flex items-center justify-center"
                >
                  跳转
                </a>
                <button
                  onClick={() => void onRefreshStats(url.shortCode)}
                  className="px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600 transition-colors text-sm"
                >
                  刷新统计
                </button>
              </div>
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}

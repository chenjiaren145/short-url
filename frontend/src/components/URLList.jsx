export default function URLList({ urls, onRefreshStats }) {
  if (urls.length === 0) {
    return (
      <div className="card p-8 text-center text-gray-500">
        还没有创建任何短链接
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <h2 className="section-title">短链接历史</h2>
      {urls.map((url) => (
        <div
          key={url.shortCode}
          className="card p-6 hover:shadow-lg transition-shadow"
        >
          <div className="space-y-3">
            <div>
              <label className="field-label mb-1">短链接</label>
              <div className="flex items-center gap-2">
                <input
                  type="text"
                  value={url.shortUrl}
                  readOnly
                  className="flex-1 px-3 py-2 bg-gray-50 border border-gray-200 rounded-md text-sm text-gray-700"
                />
                <button
                  onClick={() => navigator.clipboard.writeText(url.shortUrl)}
                  className="btn-neutral"
                >
                  复制
                </button>
              </div>
            </div>

            <div>
              <label className="field-label mb-1">原始链接</label>
              <div className="px-3 py-2 bg-gray-50 rounded-md text-sm text-gray-700 break-all">
                {url.originalUrl}
              </div>
            </div>

            <div className="flex items-center justify-between pt-2 border-t border-gray-100">
              <div className="flex items-center gap-4">
                <div>
                  <label className="field-label">访问次数</label>
                  <span className="text-2xl font-bold text-blue-600">{url.visits}</span>
                </div>
              </div>
              <button
                onClick={() => onRefreshStats(url.shortCode)}
                className="btn-success"
              >
                刷新统计
              </button>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}

import { useState, useEffect } from 'react';
import { api } from './services/api';
import URLForm from './components/URLForm';
import URLList from './components/URLList';

function App() {
  const [urls, setUrls] = useState(() => {
    const saved = localStorage.getItem('urlHistory');
    return saved ? JSON.parse(saved) : [];
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    localStorage.setItem('urlHistory', JSON.stringify(urls));
  }, [urls]);

  const handleCreateURL = async (originalURL) => {
    setLoading(true);
    setError('');
    try {
      const result = await api.createShortURL(originalURL);
      const shortCode = result.short_url.split('/').pop();
      setUrls(prev => [...prev, {
        shortCode,
        shortUrl: result.short_url,
        originalUrl: originalURL,
        visits: 0,
        createdAt: new Date().toISOString()
      }]);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleRefreshStats = async (shortCode) => {
    try {
      const stats = await api.getStats(shortCode);
      setUrls(prev => prev.map(url =>
        url.shortCode === shortCode ? { ...url, visits: stats.visits } : url
      ));
    } catch (err) {
      console.error('刷新统计失败:', err);
    }
  };

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="max-w-3xl mx-auto">
        <h1 className="text-4xl font-bold text-center mb-8 text-gray-800">
          短链接生成器
        </h1>
        
        <URLForm onSubmit={handleCreateURL} loading={loading} />
        
        {error && (
          <div className="mt-4 p-4 bg-red-100 text-red-700 rounded-lg">
            {error}
          </div>
        )}
        
        <URLList urls={urls} onRefreshStats={handleRefreshStats} />
      </div>
    </div>
  );
}

export default App;

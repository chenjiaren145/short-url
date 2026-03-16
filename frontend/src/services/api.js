const API_BASE = '/api';

export const api = {
  async createShortURL(originalURL) {
    const response = await fetch(`${API_BASE}/shorten`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ original_url: originalURL })
    });
    if (!response.ok) throw new Error('创建失败');
    return response.json();
  },

  async getStats(shortCode) {
    const response = await fetch(`${API_BASE}/${shortCode}/stats`);
    if (!response.ok) throw new Error('获取统计失败');
    return response.json();
  }
};

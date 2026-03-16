import type { CreateShortURLResponse, URLStatsResponse } from '../types/url'

const API_BASE = '/api'

export const api = {
  async createShortURL(originalURL: string): Promise<CreateShortURLResponse> {
    const response = await fetch(`${API_BASE}/shorten`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ original_url: originalURL }),
    })
    if (!response.ok) {
      throw new Error('创建失败')
    }
    return response.json() as Promise<CreateShortURLResponse>
  },

  async getStats(shortCode: string): Promise<URLStatsResponse> {
    const response = await fetch(`${API_BASE}/${shortCode}/stats`)
    if (!response.ok) {
      throw new Error('获取统计失败')
    }
    return response.json() as Promise<URLStatsResponse>
  },
}

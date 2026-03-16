export type URLItem = {
  shortCode: string
  shortUrl: string
  originalUrl: string
  visits: number
  createdAt: string
}

export type CreateShortURLResponse = {
  short_url: string
}

export type URLStatsResponse = {
  visits: number
}

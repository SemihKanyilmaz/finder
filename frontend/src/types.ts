export type ContentType = 'video' | 'article' | ''

export interface Content {
  id: string
  title: string
  type: ContentType
  source: string
  publishedAt: string
  score: number
  views: number
  likes: number
  readingTime: number
  reactions: number
}

export interface SearchResult {
  items: Content[]
  totalCount: number
  totalPages: number
  page: number
  pageSize: number
}

export interface SearchParams {
  q: string
  type: ContentType
  sortBy: 'score' | 'freshness'
  page: number
  pageSize: number
}

import { SearchParams, SearchResult } from './types'

export async function search(params: SearchParams): Promise<SearchResult> {
  const query = new URLSearchParams({
    q: params.q,
    ...(params.type && { type: params.type }),
    sortBy: params.sortBy,
    page: String(params.page),
    pageSize: String(params.pageSize),
  })

  const res = await fetch(`/search?${query}`)
  if (!res.ok) throw new Error(`Search failed: ${res.status}`)
  return res.json()
}

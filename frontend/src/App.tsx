import { useState, useEffect, useCallback } from 'react'
import {
  Box,
  Container,
  Typography,
  Stack,
  Pagination,
  Alert,
  CircularProgress,
  Divider,
} from '@mui/material'
import { search } from './api'
import { SearchParams, SearchResult } from './types'
import SearchBar from './components/SearchBar'
import ContentCard from './components/ContentCard'

export default function App() {
  const [result, setResult] = useState<SearchResult | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [params, setParams] = useState<SearchParams | null>(null)

  const handleSearch = useCallback(async (p: Omit<SearchParams, 'page'>, page = 1) => {
    const fullParams: SearchParams = { ...p, page }
    setParams(fullParams)
    setLoading(true)
    setError(null)
    try {
      const data = await search(fullParams)
      setResult(data)
    } catch {
      setError('Arama sırasında bir hata oluştu.')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    handleSearch({ q: '', type: '', sortBy: 'score', pageSize: 10 })
  }, [handleSearch])

  const handlePageChange = (_: React.ChangeEvent<unknown>, page: number) => {
    if (!params) return
    handleSearch(params, page)
  }

  return (
    <Box sx={{ minHeight: '100vh', bgcolor: 'grey.50', py: 4 }}>
      <Container maxWidth="md">
        <Stack gap={3}>
          <Typography variant="h4" fontWeight={700} textAlign="center">
            🔍 Finder
          </Typography>

          <SearchBar onSearch={(p) => handleSearch(p, 1)} loading={loading} />

          {error && <Alert severity="error">{error}</Alert>}

          {loading && (
            <Box display="flex" justifyContent="center" py={6}>
              <CircularProgress />
            </Box>
          )}

          {!loading && result && (
            <>
              <Stack direction="row" justifyContent="space-between" alignItems="center">
                <Typography variant="body2" color="text.secondary">
                  {result.totalCount} sonuç bulundu
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Sayfa {result.page} / {result.totalPages}
                </Typography>
              </Stack>

              <Divider />

              {result.items.length === 0 ? (
                <Alert severity="info">Sonuç bulunamadı.</Alert>
              ) : (
                <Stack gap={1.5}>
                  {result.items.map((item) => (
                    <ContentCard key={`${item.id}-${item.source}`} content={item} />
                  ))}
                </Stack>
              )}

              {result.totalPages > 1 && (
                <Box display="flex" justifyContent="center">
                  <Pagination
                    count={result.totalPages}
                    page={result.page}
                    onChange={handlePageChange}
                    color="primary"
                  />
                </Box>
              )}
            </>
          )}
        </Stack>
      </Container>
    </Box>
  )
}

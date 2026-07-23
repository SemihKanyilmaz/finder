import { useState } from 'react'
import {
  Stack,
  TextField,
  ToggleButton,
  ToggleButtonGroup,
  Button,
  MenuItem,
  Select,
  FormControl,
  InputLabel,
  InputAdornment,
} from '@mui/material'
import SearchIcon from '@mui/icons-material/Search'
import { ContentType, SearchParams } from '../types'

interface Props {
  onSearch: (params: Omit<SearchParams, 'page'>) => void
  loading: boolean
}

export default function SearchBar({ onSearch, loading }: Props) {
  const [q, setQ] = useState('')
  const [type, setType] = useState<ContentType>('')
  const [sortBy, setSortBy] = useState<'score' | 'freshness'>('score')
  const [pageSize, setPageSize] = useState(10)

  const triggerSearch = (overrides?: Partial<{ q: string; type: ContentType; sortBy: 'score' | 'freshness'; pageSize: number }>) => {
    onSearch({
      q: overrides?.q ?? q,
      type: overrides?.type !== undefined ? overrides.type : type,
      sortBy: overrides?.sortBy ?? sortBy,
      pageSize: overrides?.pageSize ?? pageSize,
    })
  }

  const handleTypeChange = (_: React.MouseEvent, v: ContentType | null) => {
    const newType = v ?? ''
    setType(newType)
    triggerSearch({ type: newType })
  }

  const handleSortByChange = (value: 'score' | 'freshness') => {
    setSortBy(value)
    triggerSearch({ sortBy: value })
  }

  const handlePageSizeChange = (value: number) => {
    setPageSize(value)
    triggerSearch({ pageSize: value })
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') triggerSearch()
  }

  return (
    <Stack direction={{ xs: 'column', sm: 'row' }} gap={1.5} alignItems="center">
      <TextField
        fullWidth
        placeholder="Ara..."
        value={q}
        onChange={(e) => setQ(e.target.value)}
        onKeyDown={handleKeyDown}
        size="small"
        slotProps={{
          input: {
            startAdornment: (
              <InputAdornment position="start">
                <SearchIcon fontSize="small" />
              </InputAdornment>
            ),
          },
        }}
      />

      <ToggleButtonGroup
        value={type}
        exclusive
        onChange={handleTypeChange}
        size="small"
        sx={{ flexShrink: 0 }}
      >
        <ToggleButton value="">Tümü</ToggleButton>
        <ToggleButton value="video">Video</ToggleButton>
        <ToggleButton value="article">Makale</ToggleButton>
      </ToggleButtonGroup>

      <FormControl size="small" sx={{ minWidth: 140, flexShrink: 0 }}>
        <InputLabel>Sıralama</InputLabel>
        <Select
          value={sortBy}
          label="Sıralama"
          onChange={(e) => handleSortByChange(e.target.value as 'score' | 'freshness')}
        >
          <MenuItem value="score">Skor</MenuItem>
          <MenuItem value="freshness">Güncellik</MenuItem>
        </Select>
      </FormControl>

      <FormControl size="small" sx={{ minWidth: 90, flexShrink: 0 }}>
        <InputLabel>Limit</InputLabel>
        <Select
          value={pageSize}
          label="Limit"
          onChange={(e) => handlePageSizeChange(Number(e.target.value))}
        >
          <MenuItem value={5}>5</MenuItem>
          <MenuItem value={10}>10</MenuItem>
          <MenuItem value={20}>20</MenuItem>
          <MenuItem value={50}>50</MenuItem>
        </Select>
      </FormControl>

      <Button
        variant="contained"
        onClick={() => triggerSearch()}
        disabled={loading}
        sx={{ flexShrink: 0 }}
      >
        {loading ? 'Aranıyor...' : 'Ara'}
      </Button>
    </Stack>
  )
}

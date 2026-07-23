import {
  Card,
  CardContent,
  Chip,
  Typography,
  Stack,
  Tooltip,
} from '@mui/material'
import VideoLibraryIcon from '@mui/icons-material/VideoLibrary'
import ArticleIcon from '@mui/icons-material/Article'
import StarIcon from '@mui/icons-material/Star'
import { Content } from '../types'

interface Props {
  content: Content
}

export default function ContentCard({ content }: Props) {
  const isVideo = content.type === 'video'
  const date = new Date(content.publishedAt).toLocaleDateString('tr-TR', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })

  return (
    <Card variant="outlined" sx={{ borderRadius: 2 }}>
      <CardContent>
        <Stack direction="row" alignItems="flex-start" justifyContent="space-between" gap={1}>
          <Stack direction="row" alignItems="center" gap={1} flexShrink={1} minWidth={0}>
            {isVideo ? (
              <VideoLibraryIcon fontSize="small" color="error" sx={{ flexShrink: 0 }} />
            ) : (
              <ArticleIcon fontSize="small" color="primary" sx={{ flexShrink: 0 }} />
            )}
            <Typography variant="subtitle1" fontWeight={600} noWrap>
              {content.title}
            </Typography>
          </Stack>

          <Tooltip title="Score">
            <Chip
              icon={<StarIcon fontSize="small" />}
              label={content.score.toFixed(2)}
              size="small"
              color="warning"
              variant="outlined"
              sx={{ flexShrink: 0 }}
            />
          </Tooltip>
        </Stack>

        <Stack direction="row" gap={1} mt={1} flexWrap="wrap">
          <Chip label={content.type} size="small" variant="filled" />
          <Chip label={content.source} size="small" variant="outlined" />
          <Typography variant="caption" color="text.secondary" alignSelf="center">
            {date}
          </Typography>
        </Stack>

        <Stack direction="row" gap={2} mt={1}>
          {isVideo ? (
            <>
              <Typography variant="caption" color="text.secondary">
                👁 {content.views.toLocaleString()}
              </Typography>
              <Typography variant="caption" color="text.secondary">
                👍 {content.likes.toLocaleString()}
              </Typography>
            </>
          ) : (
            <>
              <Typography variant="caption" color="text.secondary">
                📖 {content.readingTime} dk
              </Typography>
              <Typography variant="caption" color="text.secondary">
                ❤️ {content.reactions}
              </Typography>
            </>
          )}
        </Stack>
      </CardContent>
    </Card>
  )
}

import React, { FC } from 'react'
import { useQueryRecipe } from '../hooks/useQueryRecipe'
import {
  Card,
  CardContent,
  Typography,
  CircularProgress,
  Alert,
  Box,
} from '@mui/material'
import { RestaurantMenu } from '@mui/icons-material'

export const RecipeSuggestions: FC = () => {
  const { data: recipe, isLoading, error } = useQueryRecipe()

  if (isLoading) {
    return (
      <Card>
        <CardContent sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
          <CircularProgress />
        </CardContent>
      </Card>
    )
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ mb: 2 }}>
        レシピの取得に失敗しました。しばらく経ってから再度お試しください。
      </Alert>
    )
  }

  if (!recipe) {
    return (
      <Alert severity="info" sx={{ mb: 2 }}>
        期限切れ間近の食材がありません。
      </Alert>
    )
  }

  return (
    <Card>
      <CardContent>
        <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
          <RestaurantMenu sx={{ mr: 1 }} />
          <Typography variant="h6" component="div">
            おすすめレシピ
          </Typography>
        </Box>
        <Typography
          variant="body1"
          component="div"
          sx={{
            whiteSpace: 'pre-line',
            '& h1, & h2, & h3': {
              fontSize: '1.2rem',
              fontWeight: 'bold',
              my: 2,
            },
            '& ul': {
              pl: 2,
              my: 1,
            },
            '& li': {
              my: 0.5,
            },
          }}
        >
          {recipe}
        </Typography>
      </CardContent>
    </Card>
  )
}

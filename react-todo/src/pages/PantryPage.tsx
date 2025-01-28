import React, { FC } from 'react'
import { useNavigate } from 'react-router-dom'
import { useQueryClient } from '@tanstack/react-query'
import { FoodItemForm } from '../components/FoodItemForm'
import { FoodItemList } from '../components/FoodList'
import { RecipeSuggestions } from '../components/RecipeSuggestions'
import { Container, Typography, Box, Button } from '@mui/material'
import { LogoutOutlined } from '@mui/icons-material'

export const PantryPage: FC = () => {
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const logout = async () => {
    try {
      await queryClient.resetQueries()
      navigate('/')
    } catch (err: any) {
      alert(err.message)
    }
  }

  return (
    <Container maxWidth="md">
      <Box sx={{ my: 4 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 4 }}>
          <Typography variant="h4" component="h1">
            スマートパントリー
          </Typography>
          <Button
            variant="outlined"
            color="inherit"
            onClick={logout}
            startIcon={<LogoutOutlined />}
          >
            ログアウト
          </Button>
        </Box>

        <FoodItemForm />

        <Box sx={{ mb: 4 }}>
          <Typography variant="h5" component="h2" sx={{ mb: 2 }}>
            食材一覧
          </Typography>
          <FoodItemList />
        </Box>

        <Box>
          <Typography variant="h5" component="h2" sx={{ mb: 2 }}>
            レシピ提案
          </Typography>
          <RecipeSuggestions />
        </Box>
      </Box>
    </Container>
  )
}

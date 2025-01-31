import React, { FC } from 'react'
import { useNavigate } from 'react-router-dom'
import { useQueryClient } from '@tanstack/react-query'
import { useQueryFoodItems } from '../hooks/useQueryFoodItems'
import { useQueryRecipe } from '../hooks/useQueryRecipe'
import { Typography, Box, Grid, Paper, CircularProgress } from '@mui/material'
import { Recipe } from '../types'

export const PantryPage: FC = () => {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { data: foodItems = [], isLoading: isFoodLoading } = useQueryFoodItems()
  const { data: recipes = [], isLoading: isRecipeLoading } = useQueryRecipe()

  // 日数計算用のヘルパー関数
  const getDaysUntilExpiry = (expiryDateStr: string): number => {
    const today = new Date()
    today.setHours(0, 0, 0, 0)
    const expiry = new Date(expiryDateStr)
    expiry.setHours(0, 0, 0, 0)
    return Math.ceil((expiry.getTime() - today.getTime()) / (1000 * 3600 * 24))
  }

  // 期限切れ食材を抽出
  const expiredItems = foodItems.filter((item) => {
    if (!item?.expiry_date) return false
    try {
      return getDaysUntilExpiry(item.expiry_date) < 0
    } catch (e) {
      console.error('Invalid date:', item.expiry_date)
      return false
    }
  })

  // 期限切れが近い食材（7日以内）を抽出
  const nearExpiryItems = foodItems.filter((item) => {
    if (!item?.expiry_date) return false
    try {
      const daysUntilExpiry = getDaysUntilExpiry(item.expiry_date)
      return daysUntilExpiry <= 7 && daysUntilExpiry >= 0
    } catch (e) {
      console.error('Invalid date:', item.expiry_date)
      return false
    }
  })

  // 残量が少ない食材を抽出
  const lowStockItems = foodItems.filter(
    (item) => item?.quantity && item.quantity <= 2
  )

  // 日付をフォーマットする関数
  const formatDate = (dateStr: string): string => {
    try {
      return new Date(dateStr).toLocaleDateString('ja-JP')
    } catch (e) {
      console.error('Invalid date:', dateStr)
      return '無効な日付'
    }
  }

  if (isFoodLoading || isRecipeLoading) {
    return (
      <Box
        display="flex"
        justifyContent="center"
        alignItems="center"
        minHeight="80vh"
      >
        <CircularProgress />
      </Box>
    )
  }

  return (
    <Box>
      {/* ページタイトル */}
      <Typography
        variant="h5"
        sx={{
          mb: 4,
          fontWeight: 600,
          color: '#1a1a1a',
          position: 'relative',
          '&::after': {
            content: '""',
            position: 'absolute',
            bottom: -8,
            left: 0,
            width: 40,
            height: 3,
            background: 'linear-gradient(45deg, #2196F3 30%, #21CBF3 90%)',
            borderRadius: '2px',
          },
        }}
      >
        パントリー管理
      </Typography>

      <Grid container spacing={3}>
        {/* 期限切れ食材 */}
        <Grid item xs={12} md={6}>
          <Paper
            elevation={0}
            sx={{
              p: 3,
              height: '100%',
              borderRadius: 2,
              background: 'linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%)',
              border: '1px solid rgba(0, 0, 0, 0.08)',
              position: 'relative',
              overflow: 'hidden',
              transition: 'transform 0.2s ease-in-out',
              '&:hover': {
                transform: 'translateY(-2px)',
              },
            }}
          >
            <Typography
              variant="h6"
              sx={{
                mb: 3,
                fontWeight: 600,
                color: 'error.main',
                display: 'flex',
                alignItems: 'center',
                gap: 1,
              }}
            >
              期限切れの食材
            </Typography>
            {expiredItems.length > 0 ? (
              <Box>
                {expiredItems.map((item) => (
                  <Box
                    key={item.id}
                    sx={{
                      p: 2,
                      mb: 2,
                      borderRadius: 1,
                      bgcolor: 'rgba(255, 255, 255, 0.8)',
                      border: '1px solid rgba(0, 0, 0, 0.04)',
                      '&:last-child': {
                        mb: 0,
                      },
                    }}
                  >
                    <Typography sx={{ fontWeight: 500 }}>
                      {item.title}
                    </Typography>
                    <Typography
                      variant="body2"
                      sx={{
                        color: 'text.secondary',
                        mt: 0.5,
                      }}
                    >
                      期限: {formatDate(item.expiry_date)}
                    </Typography>
                  </Box>
                ))}
              </Box>
            ) : (
              <Typography
                variant="body2"
                sx={{
                  color: 'text.secondary',
                  fontStyle: 'italic',
                }}
              >
                期限切れの食材はありません
              </Typography>
            )}
          </Paper>
        </Grid>

        {/* 期限切れ間近の食材 */}
        <Grid item xs={12} md={6}>
          <Paper
            elevation={0}
            sx={{
              p: 3,
              height: '100%',
              borderRadius: 2,
              background: 'linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%)',
              border: '1px solid rgba(0, 0, 0, 0.08)',
              position: 'relative',
              overflow: 'hidden',
              transition: 'transform 0.2s ease-in-out',
              '&:hover': {
                transform: 'translateY(-2px)',
              },
            }}
          >
            <Typography
              variant="h6"
              sx={{
                mb: 3,
                fontWeight: 600,
                color: 'warning.main',
                display: 'flex',
                alignItems: 'center',
                gap: 1,
              }}
            >
              期限切れ間近の食材
            </Typography>
            {nearExpiryItems.length > 0 ? (
              <Box>
                {nearExpiryItems.map((item) => (
                  <Box
                    key={item.id}
                    sx={{
                      p: 2,
                      mb: 2,
                      borderRadius: 1,
                      bgcolor: 'rgba(255, 255, 255, 0.8)',
                      border: '1px solid rgba(0, 0, 0, 0.04)',
                      '&:last-child': {
                        mb: 0,
                      },
                    }}
                  >
                    <Typography sx={{ fontWeight: 500 }}>
                      {item.title}
                    </Typography>
                    <Typography
                      variant="body2"
                      sx={{
                        color: 'text.secondary',
                        mt: 0.5,
                      }}
                    >
                      期限: {formatDate(item.expiry_date)}
                    </Typography>
                  </Box>
                ))}
              </Box>
            ) : (
              <Typography
                variant="body2"
                sx={{
                  color: 'text.secondary',
                  fontStyle: 'italic',
                }}
              >
                期限切れが近い食材はありません
              </Typography>
            )}
          </Paper>
        </Grid>

        {/* 在庫状況 */}
        <Grid item xs={12} md={6}>
          <Paper
            elevation={0}
            sx={{
              p: 3,
              height: '100%',
              borderRadius: 2,
              background: 'linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%)',
              border: '1px solid rgba(0, 0, 0, 0.08)',
              position: 'relative',
              overflow: 'hidden',
              transition: 'transform 0.2s ease-in-out',
              '&:hover': {
                transform: 'translateY(-2px)',
              },
            }}
          >
            <Typography
              variant="h6"
              sx={{
                mb: 3,
                fontWeight: 600,
                color: 'info.main',
                display: 'flex',
                alignItems: 'center',
                gap: 1,
              }}
            >
              在庫状況
            </Typography>
            {lowStockItems.length > 0 ? (
              <Box>
                {lowStockItems.map((item) => (
                  <Box
                    key={item.id}
                    sx={{
                      p: 2,
                      mb: 2,
                      borderRadius: 1,
                      bgcolor: 'rgba(255, 255, 255, 0.8)',
                      border: '1px solid rgba(0, 0, 0, 0.04)',
                      '&:last-child': {
                        mb: 0,
                      },
                    }}
                  >
                    <Typography sx={{ fontWeight: 500 }}>
                      {item.title}
                    </Typography>
                    <Typography
                      variant="body2"
                      sx={{
                        color: 'text.secondary',
                        mt: 0.5,
                      }}
                    >
                      残量: {item.quantity}個
                    </Typography>
                  </Box>
                ))}
              </Box>
            ) : (
              <Typography
                variant="body2"
                sx={{
                  color: 'text.secondary',
                  fontStyle: 'italic',
                }}
              >
                在庫が少ない食材はありません
              </Typography>
            )}
          </Paper>
        </Grid>

        {/* おすすめレシピ */}
        <Grid item xs={12}>
          <Paper
            elevation={0}
            sx={{
              p: 3,
              borderRadius: 2,
              background: 'linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%)',
              border: '1px solid rgba(0, 0, 0, 0.08)',
              position: 'relative',
              overflow: 'hidden',
              transition: 'transform 0.2s ease-in-out',
              '&:hover': {
                transform: 'translateY(-2px)',
              },
            }}
          >
            <Typography
              variant="h6"
              sx={{
                mb: 3,
                fontWeight: 600,
                color: 'success.main',
                display: 'flex',
                alignItems: 'center',
                gap: 1,
              }}
            >
              今日のおすすめレシピ
            </Typography>
            {recipes.length > 0 ? (
              <Box>
                {recipes.slice(0, 3).map((recipe: Recipe, index) => (
                  <Box
                    key={index}
                    sx={{
                      p: 2,
                      mb: 2,
                      borderRadius: 1,
                      bgcolor: 'rgba(255, 255, 255, 0.8)',
                      border: '1px solid rgba(0, 0, 0, 0.04)',
                      '&:last-child': {
                        mb: 0,
                      },
                    }}
                  >
                    <Typography sx={{ fontWeight: 500 }}>
                      {recipe.title}
                    </Typography>
                    <Typography
                      variant="body2"
                      sx={{
                        color: 'text.secondary',
                        mt: 0.5,
                      }}
                    >
                      使用食材: {recipe.ingredients?.join(', ') || '情報なし'}
                    </Typography>
                  </Box>
                ))}
              </Box>
            ) : (
              <Typography
                variant="body2"
                sx={{
                  color: 'text.secondary',
                  fontStyle: 'italic',
                }}
              >
                現在おすすめできるレシピはありません
              </Typography>
            )}
          </Paper>
        </Grid>
      </Grid>
    </Box>
  )
}

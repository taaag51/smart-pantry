import React, { FC } from 'react'
import { useNavigate } from 'react-router-dom'
import { useQueryClient } from '@tanstack/react-query'
import { useQueryFoodItems } from '../hooks/useQueryFoodItems'
import { useQueryRecipe } from '../hooks/useQueryRecipe'
import {
  Typography,
  Box,
  Button,
  Grid,
  Paper,
  CircularProgress,
} from '@mui/material'
import { Recipe } from '../types'

export const PantryPage: FC = () => {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { data: foodItems, isLoading: isFoodLoading } = useQueryFoodItems()
  const { data: recipes, isLoading: isRecipeLoading } = useQueryRecipe()

  const logout = async () => {
    try {
      await queryClient.resetQueries()
      navigate('/')
    } catch (err: any) {
      alert(err.message)
    }
  }

  // 日数計算用のヘルパー関数
  const getDaysUntilExpiry = (expiryDate: Date): number => {
    const today = new Date()
    today.setHours(0, 0, 0, 0)
    const expiry = new Date(expiryDate)
    expiry.setHours(0, 0, 0, 0)
    return Math.ceil((expiry.getTime() - today.getTime()) / (1000 * 3600 * 24))
  }

  // 期限切れ食材を抽出
  const expiredItems = foodItems?.filter((item) => {
    if (!item.expiry_date) return false
    return getDaysUntilExpiry(item.expiry_date) < 0
  })

  // 期限切れが近い食材（7日以内）を抽出
  const nearExpiryItems = foodItems?.filter((item) => {
    if (!item.expiry_date) return false
    const daysUntilExpiry = getDaysUntilExpiry(item.expiry_date)
    return daysUntilExpiry <= 7 && daysUntilExpiry >= 0
  })

  // 残量が少ない食材を抽出
  const lowStockItems = foodItems?.filter(
    (item) => item.quantity && item.quantity <= 2
  )

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
    <Box
      sx={{
        maxWidth: '1200px',
        mx: 'auto',
        px: { xs: 2, sm: 3 },
        py: { xs: 4, sm: 6 },
        bgcolor: '#FFFFFF',
      }}
    >
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          mb: { xs: 4, sm: 6 },
          pb: 3,
          borderBottom: '1px solid rgba(0, 0, 0, 0.08)',
        }}
      >
        <Typography
          variant="h4"
          component="h1"
          sx={{
            fontWeight: 400,
            letterSpacing: '-0.5px',
            color: 'rgba(0, 0, 0, 0.87)',
            position: 'relative',
            '&::after': {
              content: '""',
              position: 'absolute',
              bottom: -8,
              left: 0,
              width: '40px',
              height: '2px',
              backgroundColor: 'primary.main',
              borderRadius: '2px',
            },
          }}
        >
          パントリーダッシュボード
        </Typography>
        <Button
          variant="text"
          onClick={logout}
          sx={{
            color: 'rgba(0, 0, 0, 0.6)',
            textTransform: 'none',
            fontSize: '0.875rem',
            position: 'relative',
            '&:hover': {
              backgroundColor: 'transparent',
              '&::after': {
                width: '100%',
              },
            },
            '&::after': {
              content: '""',
              position: 'absolute',
              bottom: 6,
              left: 0,
              width: 0,
              height: '1px',
              backgroundColor: 'rgba(0, 0, 0, 0.6)',
              transition: 'width 0.2s ease-in-out',
            },
          }}
        >
          ログアウト
        </Button>
      </Box>

      <Grid container spacing={3}>
        {/* 期限切れ食材 */}
        <Grid item xs={12} md={6}>
          <Paper
            elevation={0}
            sx={{
              p: 3,
              borderRadius: 1,
              border: '1px solid rgba(0, 0, 0, 0.12)',
              position: 'relative',
              overflow: 'hidden',
              transition: 'all 0.2s ease-in-out',
              '&:hover': {
                borderColor: 'rgba(0, 0, 0, 0.24)',
                transform: 'translateY(-1px)',
                boxShadow: '0 2px 4px rgba(0,0,0,0.02)',
              },
              '&::before': {
                content: '""',
                position: 'absolute',
                top: 0,
                left: 0,
                width: '4px',
                height: '100%',
                backgroundColor: 'error.main',
              },
            }}
          >
            <Typography
              variant="h6"
              sx={{
                mb: 3,
                fontWeight: 500,
                color: 'rgba(0, 0, 0, 0.87)',
                letterSpacing: '0.5px',
                display: 'flex',
                alignItems: 'center',
                '&::before': {
                  content: '""',
                  width: '6px',
                  height: '6px',
                  borderRadius: '50%',
                  backgroundColor: 'error.main',
                  marginRight: '12px',
                },
              }}
            >
              期限切れの食材
            </Typography>
            {expiredItems && expiredItems.length > 0 ? (
              <Box>
                {expiredItems.map((item) => (
                  <Box
                    key={item.id}
                    sx={{
                      py: 2,
                      borderBottom: '1px solid rgba(0, 0, 0, 0.08)',
                      '&:last-child': {
                        borderBottom: 'none',
                      },
                    }}
                  >
                    <Typography
                      sx={{
                        fontWeight: 500,
                        color: 'rgba(0, 0, 0, 0.87)',
                      }}
                    >
                      {item.title}
                    </Typography>
                    <Typography
                      variant="body2"
                      sx={{
                        color: 'rgba(0, 0, 0, 0.6)',
                        mt: 0.5,
                      }}
                    >
                      期限:{' '}
                      {new Date(item.expiry_date).toLocaleDateString('ja-JP')}
                    </Typography>
                  </Box>
                ))}
              </Box>
            ) : (
              <Typography
                variant="body2"
                sx={{
                  color: 'rgba(0, 0, 0, 0.6)',
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
              borderRadius: 1,
              border: '1px solid rgba(0, 0, 0, 0.12)',
              position: 'relative',
              overflow: 'hidden',
              transition: 'all 0.2s ease-in-out',
              '&:hover': {
                borderColor: 'rgba(0, 0, 0, 0.24)',
                transform: 'translateY(-1px)',
                boxShadow: '0 2px 4px rgba(0,0,0,0.02)',
              },
              '&::before': {
                content: '""',
                position: 'absolute',
                top: 0,
                left: 0,
                width: '4px',
                height: '100%',
                backgroundColor: 'warning.main',
              },
            }}
          >
            <Typography
              variant="h6"
              sx={{
                mb: 3,
                fontWeight: 500,
                color: 'rgba(0, 0, 0, 0.87)',
                letterSpacing: '0.5px',
                display: 'flex',
                alignItems: 'center',
                '&::before': {
                  content: '""',
                  width: '6px',
                  height: '6px',
                  borderRadius: '50%',
                  backgroundColor: 'warning.main',
                  marginRight: '12px',
                },
              }}
            >
              期限切れ間近の食材
            </Typography>
            {nearExpiryItems && nearExpiryItems.length > 0 ? (
              <Box>
                {nearExpiryItems.map((item) => (
                  <Box
                    key={item.id}
                    sx={{
                      py: 2,
                      borderBottom: '1px solid rgba(0, 0, 0, 0.08)',
                      '&:last-child': {
                        borderBottom: 'none',
                      },
                    }}
                  >
                    <Typography
                      sx={{
                        fontWeight: 500,
                        color: 'rgba(0, 0, 0, 0.87)',
                      }}
                    >
                      {item.title}
                    </Typography>
                    <Typography
                      variant="body2"
                      sx={{
                        color: 'rgba(0, 0, 0, 0.6)',
                        mt: 0.5,
                      }}
                    >
                      期限:{' '}
                      {new Date(item.expiry_date).toLocaleDateString('ja-JP')}
                    </Typography>
                  </Box>
                ))}
              </Box>
            ) : (
              <Typography
                variant="body2"
                sx={{
                  color: 'rgba(0, 0, 0, 0.6)',
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
              borderRadius: 1,
              border: '1px solid rgba(0, 0, 0, 0.12)',
              position: 'relative',
              overflow: 'hidden',
              transition: 'all 0.2s ease-in-out',
              '&:hover': {
                borderColor: 'rgba(0, 0, 0, 0.24)',
                transform: 'translateY(-1px)',
                boxShadow: '0 2px 4px rgba(0,0,0,0.02)',
              },
              '&::before': {
                content: '""',
                position: 'absolute',
                top: 0,
                left: 0,
                width: '4px',
                height: '100%',
                backgroundColor: 'info.main',
              },
            }}
          >
            <Typography
              variant="h6"
              sx={{
                mb: 3,
                fontWeight: 500,
                color: 'rgba(0, 0, 0, 0.87)',
                letterSpacing: '0.5px',
                display: 'flex',
                alignItems: 'center',
                '&::before': {
                  content: '""',
                  width: '6px',
                  height: '6px',
                  borderRadius: '50%',
                  backgroundColor: 'info.main',
                  marginRight: '12px',
                },
              }}
            >
              在庫状況
            </Typography>
            {lowStockItems && lowStockItems.length > 0 ? (
              <Box>
                {lowStockItems.map((item) => (
                  <Box
                    key={item.id}
                    sx={{
                      py: 2,
                      borderBottom: '1px solid rgba(0, 0, 0, 0.08)',
                      '&:last-child': {
                        borderBottom: 'none',
                      },
                    }}
                  >
                    <Typography
                      sx={{
                        fontWeight: 500,
                        color: 'rgba(0, 0, 0, 0.87)',
                      }}
                    >
                      {item.title}
                    </Typography>
                    <Typography
                      variant="body2"
                      sx={{
                        color: 'rgba(0, 0, 0, 0.6)',
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
                  color: 'rgba(0, 0, 0, 0.6)',
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
              borderRadius: 1,
              border: '1px solid rgba(0, 0, 0, 0.12)',
              position: 'relative',
              overflow: 'hidden',
              transition: 'all 0.2s ease-in-out',
              '&:hover': {
                borderColor: 'rgba(0, 0, 0, 0.24)',
                transform: 'translateY(-1px)',
                boxShadow: '0 2px 4px rgba(0,0,0,0.02)',
              },
              '&::before': {
                content: '""',
                position: 'absolute',
                top: 0,
                left: 0,
                width: '4px',
                height: '100%',
                backgroundColor: 'success.main',
              },
            }}
          >
            <Typography
              variant="h6"
              sx={{
                mb: 3,
                fontWeight: 500,
                color: 'rgba(0, 0, 0, 0.87)',
                letterSpacing: '0.5px',
                display: 'flex',
                alignItems: 'center',
                '&::before': {
                  content: '""',
                  width: '6px',
                  height: '6px',
                  borderRadius: '50%',
                  backgroundColor: 'success.main',
                  marginRight: '12px',
                },
              }}
            >
              今日のおすすめレシピ
            </Typography>
            {recipes && recipes.length > 0 ? (
              <Box>
                {recipes.slice(0, 3).map((recipe: Recipe, index) => (
                  <Box
                    key={index}
                    sx={{
                      py: 2,
                      borderBottom: '1px solid rgba(0, 0, 0, 0.08)',
                      '&:last-child': {
                        borderBottom: 'none',
                      },
                    }}
                  >
                    <Typography
                      sx={{
                        fontWeight: 500,
                        color: 'rgba(0, 0, 0, 0.87)',
                      }}
                    >
                      {recipe.title}
                    </Typography>
                    <Typography
                      variant="body2"
                      sx={{
                        color: 'rgba(0, 0, 0, 0.6)',
                        mt: 0.5,
                      }}
                    >
                      使用食材: {recipe.ingredients?.join(', ') || '情報なし'}
                    </Typography>
                  </Box>
                ))}
                <Box sx={{ mt: 3, textAlign: 'right' }}>
                  <Button
                    variant="text"
                    onClick={() => navigate('/recipes')}
                    sx={{
                      color: 'rgba(0, 0, 0, 0.6)',
                      textTransform: 'none',
                      fontWeight: 500,
                      position: 'relative',
                      '&:hover': {
                        backgroundColor: 'transparent',
                        '&::after': {
                          width: '100%',
                        },
                      },
                      '&::after': {
                        content: '""',
                        position: 'absolute',
                        bottom: 6,
                        left: 0,
                        width: 0,
                        height: '1px',
                        backgroundColor: 'rgba(0, 0, 0, 0.6)',
                        transition: 'width 0.2s ease-in-out',
                      },
                    }}
                  >
                    レシピをもっと見る
                  </Button>
                </Box>
              </Box>
            ) : (
              <Typography
                variant="body2"
                sx={{
                  color: 'rgba(0, 0, 0, 0.6)',
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

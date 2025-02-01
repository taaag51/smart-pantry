import { FC } from 'react'
import {
  Box,
  Card,
  CardContent,
  Typography,
  List,
  ListItem,
  ListItemText,
  CircularProgress,
  Container,
  Button,
  Chip,
} from '@mui/material'
import { useQueryFoodItems } from '../hooks/useQueryFoodItems'
import { useQueryRecipe } from '../hooks/useQueryRecipe'
import { FoodItem, Recipe } from '../types'
import ReactMarkdown from 'react-markdown'

export const RecipeSuggestions: FC = () => {
  const { data: foodItems, isLoading: isFoodLoading } = useQueryFoodItems()
  const {
    data: recipes,
    isLoading: isRecipeLoading,
    refetch,
    isFetching,
  } = useQueryRecipe()

  if (isFoodLoading || isRecipeLoading || isFetching) {
    return (
      <Box
        display="flex"
        justifyContent="center"
        alignItems="center"
        minHeight="100vh"
      >
        <CircularProgress />
      </Box>
    )
  }

  return (
    <Container maxWidth="md" sx={{ py: 4 }}>
      <Typography variant="h4" component="h1" gutterBottom>
        レシピ提案
      </Typography>

      <Box mb={4}>
        <Typography variant="h6" gutterBottom>
          現在の食材
        </Typography>
        <Card>
          <CardContent>
            <List>
              {foodItems && foodItems.length > 0 ? (
                foodItems.map((item: FoodItem) => {
                  const expiryDate = new Date(item.expiry_date)
                  const daysUntilExpiry = Math.ceil(
                    (expiryDate.getTime() - new Date().getTime()) /
                      (1000 * 3600 * 24)
                  )
                  const isExpiringSoon =
                    daysUntilExpiry >= 0 && daysUntilExpiry <= 7

                  return (
                    <ListItem
                      key={item.id}
                      sx={{
                        backgroundColor: isExpiringSoon
                          ? 'rgba(255, 193, 7, 0.1)'
                          : 'transparent',
                        borderLeft: isExpiringSoon
                          ? '4px solid #ffc107'
                          : 'none',
                      }}
                    >
                      <ListItemText
                        primary={
                          <Box
                            component="span"
                            sx={{ display: 'flex', alignItems: 'center' }}
                          >
                            {item.title}
                            {isExpiringSoon && (
                              <Typography
                                component="span"
                                sx={{
                                  ml: 1,
                                  fontSize: '0.8rem',
                                  color: 'warning.main',
                                  bgcolor: 'warning.light',
                                  px: 1,
                                  py: 0.5,
                                  borderRadius: 1,
                                }}
                              >
                                期限まであと{daysUntilExpiry}日
                              </Typography>
                            )}
                          </Box>
                        }
                        secondary={`残量: ${
                          item.quantity
                        } (期限: ${expiryDate.toLocaleDateString()})`}
                      />
                    </ListItem>
                  )
                })
              ) : (
                <ListItem>
                  <ListItemText
                    primary="食材が登録されていません"
                    secondary="食材を追加してレシピの提案を受けることができます"
                  />
                </ListItem>
              )}
            </List>
          </CardContent>
        </Card>
      </Box>

      <Box>
        <Box
          display="flex"
          justifyContent="space-between"
          alignItems="center"
          mb={2}
        >
          <Typography variant="h6">おすすめレシピ</Typography>
          <Button
            variant="contained"
            onClick={() => refetch()}
            disabled={isFetching}
          >
            レシピを更新
          </Button>
        </Box>
        <Card>
          <CardContent>
            <Box sx={{ mb: 2 }}>
              <Typography variant="body2" color="warning.main" sx={{ mb: 2 }}>
                ※
                賞味期限が7日以内の食材を優先的に使用したレシピを提案しています
              </Typography>
            </Box>
            {recipes && recipes.length > 0 ? (
              recipes.map((recipe: Recipe, index) => (
                <Box
                  key={recipe.id}
                  sx={{
                    mb: index < recipes.length - 1 ? 4 : 0,
                    pb: index < recipes.length - 1 ? 2 : 0,
                    borderBottom:
                      index < recipes.length - 1
                        ? '1px solid rgba(0, 0, 0, 0.12)'
                        : 'none',
                  }}
                >
                  <Typography variant="h6" gutterBottom>
                    {recipe.title}
                  </Typography>
                  <Box sx={{ mb: 2 }}>
                    {recipe.ingredients.map((ingredient, idx) => (
                      <Chip
                        key={idx}
                        label={ingredient}
                        sx={{ mr: 1, mb: 1 }}
                        variant="outlined"
                      />
                    ))}
                  </Box>
                  <ReactMarkdown className="markdown-body">
                    {recipe.instructions}
                  </ReactMarkdown>
                </Box>
              ))
            ) : (
              <Typography variant="body1" color="text.secondary" align="center">
                現在おすすめのレシピはありません。
                <br />
                食材を追加するか、「レシピを更新」ボタンをクリックしてください。
              </Typography>
            )}
          </CardContent>
        </Card>
      </Box>
    </Container>
  )
}

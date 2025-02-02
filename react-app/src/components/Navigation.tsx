import { FC } from 'react'
import { useMutateAuth } from '../hooks/useMutateAuth' // 修正
import { AppBar, Toolbar, Button, Box } from '@mui/material'
import { useNavigate, useLocation } from 'react-router-dom'
import RestaurantMenuIcon from '@mui/icons-material/RestaurantMenu'
import KitchenIcon from '@mui/icons-material/Kitchen'
import MenuBookIcon from '@mui/icons-material/MenuBook'

export const Navigation: FC = () => {
  const navigate = useNavigate()
  const location = useLocation()

  const { logoutMutation } = useMutateAuth()

  return (
    <AppBar position="static">
      <Toolbar>
        <Box sx={{ flexGrow: 1, display: 'flex', gap: 2 }}>
          <Button
            color="inherit"
            startIcon={<KitchenIcon />}
            onClick={() => navigate('/pantry')}
            sx={{
              borderBottom:
                location.pathname === '/pantry' ? '2px solid white' : 'none',
            }}
          >
            パントリー
          </Button>
          <Button
            color="inherit"
            startIcon={<RestaurantMenuIcon />}
            onClick={() => navigate('/food')}
            sx={{
              borderBottom:
                location.pathname === '/food' ? '2px solid white' : 'none',
            }}
          >
            食材管理
          </Button>
          <Button
            color="inherit"
            startIcon={<MenuBookIcon />}
            onClick={() => navigate('/recipes')}
            sx={{
              borderBottom:
                location.pathname === '/recipes' ? '2px solid white' : 'none',
            }}
          >
            レシピ提案
          </Button>
          <Button
            color="inherit"
            onClick={async () => {
              console.log('ログアウトボタンがクリックされました')
              try {
                await logoutMutation.mutateAsync()
                console.log('ログアウト処理が完了しました')
              } catch (error) {
                console.error('ログアウト処理でエラーが発生しました:', error)
              }
            }}
          >
            ログアウト
          </Button>
        </Box>
      </Toolbar>
    </AppBar>
  )
}

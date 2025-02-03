import { FC } from 'react'
import { useMutateAuth } from '../hooks/useMutateAuth' // 修正
import { AppBar, Toolbar, Button, Box } from '@mui/material'
import { useNavigate, useLocation } from 'react-router-dom'
import RestaurantMenuIcon from '@mui/icons-material/RestaurantMenu'
import KitchenIcon from '@mui/icons-material/Kitchen'
import MenuBookIcon from '@mui/icons-material/MenuBook'
import LogoutIcon from '@mui/icons-material/Logout'

export const Navigation: FC = () => {
  const navigate = useNavigate()
  const location = useLocation()

  const { logoutMutation } = useMutateAuth()

  return (
    <AppBar position="static" sx={{ zIndex: 1100 }}>
      <Toolbar sx={{ position: 'relative' }}>
        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            width: '100%',
            position: 'relative',
          }}
        >
          <Box sx={{ display: 'flex', gap: 2 }}>
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
          </Box>
          <Box
            sx={{
              marginLeft: 'auto',
              display: 'flex',
              alignItems: 'center',
              paddingRight: 2,
              position: 'relative',
              zIndex: 1200,
            }}
          >
            <Button
              color="inherit"
              variant="outlined"
              startIcon={<LogoutIcon />}
              onClick={(e) => {
                e.preventDefault()
                e.stopPropagation()
                console.log('ログアウトボタンがクリックされました')
                logoutMutation.mutate()
              }}
              sx={{
                minWidth: '120px',
                height: '40px',
                border: '2px solid #fff',
                borderRadius: '4px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                gap: 1,
                cursor: 'pointer',
                padding: '4px 16px',
                position: 'relative',
                '&:hover': {
                  backgroundColor: 'rgba(255, 255, 255, 0.1)',
                },
                '&:active': {
                  backgroundColor: 'rgba(255, 255, 255, 0.2)',
                },
                '&::before': {
                  content: '""',
                  position: 'absolute',
                  top: -8,
                  right: -8,
                  bottom: -8,
                  left: -8,
                  cursor: 'pointer',
                },
              }}
            >
              ログアウト
            </Button>
          </Box>
        </Box>
      </Toolbar>
    </AppBar>
  )
}

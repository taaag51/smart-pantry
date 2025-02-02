import React, { FC, ReactNode } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import {
  Box,
  Drawer,
  AppBar,
  Toolbar,
  List,
  Typography,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Button,
} from '@mui/material'
import {
  // Dashboard as DashboardIcon,
  Kitchen as KitchenIcon,
  Restaurant as RestaurantIcon,
  Inventory as InventoryIcon,
  Logout as LogoutIcon,
} from '@mui/icons-material'

const drawerWidth = 240

interface Props {
  children: ReactNode
}

export const Layout: FC<Props> = ({ children }) => {
  const navigate = useNavigate()
  const location = useLocation()
  const isAuthenticated = true // TODO: 認証状態の管理を実装

  const handleLogout = () => {
    // TODO: 認証状態のクリアを実装
    navigate('/login')
  }

  const menuItems = [
    // {
    //   text: 'ダッシュボード',
    //   path: '/dashboard',
    //   icon: <DashboardIcon />,
    //   requiresAuth: true,
    // },
    {
      text: 'パントリー管理',
      path: '/pantry',
      icon: <KitchenIcon />,
      requiresAuth: true,
    },
    {
      text: '食材管理',
      path: '/food',
      icon: <InventoryIcon />,
      requiresAuth: true,
    },
    {
      text: 'レシピ提案',
      path: '/recipes',
      icon: <RestaurantIcon />,
      requiresAuth: true,
    },
  ]

  const handleNavigation = (path: string) => {
    if (isAuthenticated) {
      navigate(path)
    } else {
      navigate('/login')
    }
  }

  return (
    <Box sx={{ display: 'flex', minHeight: '100vh' }}>
      {/* サイドバー */}
      <Drawer
        variant="permanent"
        sx={{
          width: drawerWidth,
          flexShrink: 0,
          '& .MuiDrawer-paper': {
            width: drawerWidth,
            boxSizing: 'border-box',
            bgcolor: '#1a1a1a',
            color: 'white',
          },
        }}
      >
        {/* ロゴ */}
        <Box
          sx={{
            p: 2,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            height: 64,
            borderBottom: '1px solid rgba(255, 255, 255, 0.1)',
          }}
        >
          <Typography
            variant="h6"
            sx={{
              fontFamily: 'Outfit, sans-serif',
              fontWeight: 600,
              letterSpacing: '0.5px',
              color: 'white',
            }}
          >
            SMART PANTRY
          </Typography>
        </Box>

        {/* メニュー */}
        <List sx={{ mt: 2 }}>
          {menuItems.map((item) => (
            <ListItem key={item.text} disablePadding>
              <ListItemButton
                onClick={() => handleNavigation(item.path)}
                selected={location.pathname === item.path}
                sx={{
                  py: 1.5,
                  '&.Mui-selected': {
                    bgcolor: 'rgba(255, 255, 255, 0.08)',
                    '&:hover': {
                      bgcolor: 'rgba(255, 255, 255, 0.12)',
                    },
                  },
                  '&:hover': {
                    bgcolor: 'rgba(255, 255, 255, 0.04)',
                  },
                }}
              >
                <ListItemIcon sx={{ color: 'inherit', minWidth: 40 }}>
                  {item.icon}
                </ListItemIcon>
                <ListItemText
                  primary={item.text}
                  primaryTypographyProps={{
                    fontSize: '0.875rem',
                    fontWeight: location.pathname === item.path ? 500 : 400,
                  }}
                />
              </ListItemButton>
            </ListItem>
          ))}
        </List>
      </Drawer>

      {/* メインコンテンツ */}
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          bgcolor: '#f5f5f5',
          minHeight: '100vh',
        }}
      >
        {/* ヘッダー */}
        <AppBar
          position="fixed"
          sx={{
            width: `calc(100% - ${drawerWidth}px)`,
            ml: `${drawerWidth}px`,
            bgcolor: 'white',
            boxShadow: 'none',
            borderBottom: '1px solid rgba(0, 0, 0, 0.08)',
          }}
        >
          <Toolbar sx={{ justifyContent: 'flex-end' }}>
            <Button
              onClick={handleLogout}
              startIcon={<LogoutIcon />}
              sx={{
                color: 'text.secondary',
                textTransform: 'none',
                fontSize: '0.875rem',
                fontWeight: 500,
                '&:hover': {
                  bgcolor: 'rgba(0, 0, 0, 0.04)',
                },
              }}
            >
              ログアウト
            </Button>
          </Toolbar>
        </AppBar>

        {/* コンテンツ */}
        <Box sx={{ p: 3, mt: 8 }}>{children}</Box>
      </Box>
    </Box>
  )
}

import React, { FC } from 'react'
import { Box, Button, Typography, Paper } from '@mui/material'
import { FoodItem as FoodItemType } from '../types'
import { useMutateFoodItem } from '../hooks/useMutateFoodItem'

interface Props {
  foodItem: FoodItemType
}

export const FoodItem: FC<Props> = ({ foodItem }) => {
  const { deleteFoodItemMutation } = useMutateFoodItem()

  const getDaysUntilExpiry = (expiryDateStr: string): number => {
    const today = new Date()
    today.setHours(0, 0, 0, 0)
    const expiry = new Date(expiryDateStr)
    expiry.setHours(0, 0, 0, 0)
    return Math.ceil((expiry.getTime() - today.getTime()) / (1000 * 3600 * 24))
  }

  const isExpired = (expiryDateStr: string) => {
    return getDaysUntilExpiry(expiryDateStr) < 0
  }

  const isExpiringSoon = (expiryDateStr: string) => {
    const daysUntilExpiry = getDaysUntilExpiry(expiryDateStr)
    return daysUntilExpiry >= 0 && daysUntilExpiry <= 7
  }

  const getStatusColor = () => {
    if (isExpired(foodItem.expiry_date)) {
      return 'error.main'
    }
    if (isExpiringSoon(foodItem.expiry_date)) {
      return 'warning.main'
    }
    return 'info.main'
  }

  return (
    <Paper
      component="li"
      elevation={0}
      sx={{
        p: 2.5,
        mb: 2,
        borderRadius: 2,
        background: 'linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%)',
        border: '1px solid rgba(0, 0, 0, 0.08)',
        position: 'relative',
        overflow: 'hidden',
        transition: 'transform 0.2s ease-in-out',
        listStyle: 'none',
        '&:hover': {
          transform: 'translateY(-2px)',
        },
      }}
    >
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'flex-start',
        }}
      >
        <Box>
          <Typography
            sx={{
              fontSize: '1rem',
              fontWeight: 500,
              color: 'rgba(0, 0, 0, 0.87)',
              mb: 1,
              display: 'flex',
              alignItems: 'center',
              gap: 1,
              '&::before': {
                content: '""',
                width: 8,
                height: 8,
                borderRadius: '50%',
                bgcolor: getStatusColor(),
                flexShrink: 0,
              },
            }}
          >
            {foodItem.title}
          </Typography>
          <Box
            sx={{
              display: 'flex',
              alignItems: 'center',
              gap: 3,
              color: 'rgba(0, 0, 0, 0.6)',
            }}
          >
            <Typography
              variant="body2"
              sx={{
                fontSize: '0.875rem',
                display: 'flex',
                alignItems: 'center',
              }}
            >
              数量: {foodItem.quantity}
            </Typography>
            <Typography
              variant="body2"
              sx={{
                fontSize: '0.875rem',
                display: 'flex',
                alignItems: 'center',
              }}
            >
              期限: {new Date(foodItem.expiry_date).toLocaleDateString('ja-JP')}
            </Typography>
            <Typography
              variant="body2"
              sx={{
                fontSize: '0.875rem',
                color: getStatusColor(),
                fontWeight: 500,
                display: 'flex',
                alignItems: 'center',
              }}
            >
              残り{getDaysUntilExpiry(foodItem.expiry_date)}日
            </Typography>
          </Box>
        </Box>
        <Button
          onClick={() => deleteFoodItemMutation.mutate(foodItem.id)}
          variant="text"
          sx={{
            minWidth: 'auto',
            color: 'rgba(0, 0, 0, 0.6)',
            fontSize: '0.875rem',
            fontWeight: 500,
            textTransform: 'none',
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
          削除
        </Button>
      </Box>
    </Paper>
  )
}

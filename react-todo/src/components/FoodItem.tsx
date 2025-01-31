import React, { FC } from 'react'
import { Box, Button, Typography, Paper } from '@mui/material'
import { FoodItem as FoodItemType } from '../types'
import { useMutateFoodItem } from '../hooks/useMutateFoodItem'

interface Props {
  foodItem: FoodItemType
}

export const FoodItem: FC<Props> = ({ foodItem }) => {
  const { deleteFoodItemMutation } = useMutateFoodItem()

  const getDaysUntilExpiry = (expiryDate: Date): number => {
    const today = new Date()
    today.setHours(0, 0, 0, 0)
    const expiry = new Date(expiryDate)
    expiry.setHours(0, 0, 0, 0)
    return Math.ceil((expiry.getTime() - today.getTime()) / (1000 * 3600 * 24))
  }

  const isExpired = (expiryDate: Date) => {
    return getDaysUntilExpiry(expiryDate) < 0
  }

  const isExpiringSoon = (expiryDate: Date) => {
    const daysUntilExpiry = getDaysUntilExpiry(expiryDate)
    return daysUntilExpiry >= 0 && daysUntilExpiry <= 7
  }

  const getAccentColor = () => {
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
        borderRadius: 1,
        border: '1px solid rgba(0, 0, 0, 0.12)',
        position: 'relative',
        overflow: 'hidden',
        transition: 'all 0.2s ease-in-out',
        listStyle: 'none',
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
          backgroundColor: getAccentColor(),
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
              fontSize: '0.975rem',
              fontWeight: 500,
              color: 'rgba(0, 0, 0, 0.87)',
              mb: 1,
              display: 'flex',
              alignItems: 'center',
              '&::before': {
                content: '""',
                width: '4px',
                height: '4px',
                borderRadius: '50%',
                backgroundColor: getAccentColor(),
                marginRight: '8px',
              },
            }}
          >
            {foodItem.title}
          </Typography>
          <Box
            sx={{
              display: 'flex',
              alignItems: 'center',
              gap: 2.5,
              color: 'rgba(0, 0, 0, 0.6)',
              position: 'relative',
              '&::before': {
                content: '""',
                position: 'absolute',
                left: -12,
                top: '50%',
                width: '1px',
                height: '70%',
                transform: 'translateY(-50%)',
                backgroundColor: 'rgba(0, 0, 0, 0.08)',
              },
            }}
          >
            <Typography
              variant="body2"
              sx={{
                fontSize: '0.875rem',
              }}
            >
              数量: {foodItem.quantity}
            </Typography>
            <Typography
              variant="body2"
              sx={{
                fontSize: '0.875rem',
                position: 'relative',
                '&::before': {
                  content: '""',
                  position: 'absolute',
                  left: -16,
                  top: '50%',
                  width: '4px',
                  height: '4px',
                  borderRadius: '50%',
                  backgroundColor: 'rgba(0, 0, 0, 0.2)',
                  transform: 'translateY(-50%)',
                },
              }}
            >
              期限: {new Date(foodItem.expiry_date).toLocaleDateString('ja-JP')}
            </Typography>
            <Typography
              variant="body2"
              sx={{
                fontSize: '0.875rem',
                color: getAccentColor(),
                fontWeight: 500,
                position: 'relative',
                '&::before': {
                  content: '""',
                  position: 'absolute',
                  left: -16,
                  top: '50%',
                  width: '4px',
                  height: '4px',
                  borderRadius: '50%',
                  backgroundColor: 'rgba(0, 0, 0, 0.2)',
                  transform: 'translateY(-50%)',
                },
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

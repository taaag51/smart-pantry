import React, { FC, useState } from 'react'
import { useMutateFoodItem } from '../hooks/useMutateFoodItem'
import {
  TextField,
  Button,
  Box,
  FormControl,
  InputLabel,
  Input,
  InputAdornment,
} from '@mui/material'
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider'
import { DatePicker } from '@mui/x-date-pickers/DatePicker'
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns'
import jaLocale from 'date-fns/locale/ja'

export const FoodItemForm: FC = () => {
  const [title, setTitle] = useState('')
  const [quantity, setQuantity] = useState(1)
  const [expiryDate, setExpiryDate] = useState<Date | null>(new Date())
  const { createFoodItemMutation } = useMutateFoodItem()

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (title === '' || !expiryDate) return

    createFoodItemMutation.mutate({
      title,
      quantity,
      expiry_date: expiryDate,
    })
    setTitle('')
    setQuantity(1)
    setExpiryDate(new Date())
  }

  return (
    <Box component="form" onSubmit={handleSubmit} sx={{ mb: 4 }}>
      <FormControl fullWidth sx={{ mb: 2 }}>
        <TextField
          label="食材名"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
        />
      </FormControl>

      <FormControl fullWidth sx={{ mb: 2 }}>
        <InputLabel htmlFor="quantity">数量</InputLabel>
        <Input
          id="quantity"
          type="number"
          value={quantity}
          onChange={(e) => setQuantity(parseInt(e.target.value))}
          endAdornment={<InputAdornment position="end">個</InputAdornment>}
          inputProps={{ min: 1 }}
          required
        />
      </FormControl>

      <FormControl fullWidth sx={{ mb: 2 }}>
        <LocalizationProvider
          dateAdapter={AdapterDateFns}
          adapterLocale={jaLocale}
        >
          <DatePicker
            label="賞味期限"
            value={expiryDate}
            onChange={(newValue) => setExpiryDate(newValue)}
            format="yyyy/MM/dd"
          />
        </LocalizationProvider>
      </FormControl>

      <Button
        variant="contained"
        color="primary"
        type="submit"
        fullWidth
        sx={{ mt: 2 }}
      >
        追加
      </Button>
    </Box>
  )
}

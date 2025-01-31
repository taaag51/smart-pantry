import React, { FC } from 'react'
import { useFoodItemForm } from '../hooks/useFoodItemForm'
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

/**
 * 食材登録フォームコンポーネント
 *
 * このコンポーネントは以下の機能を提供します：
 * - 食材名の入力
 * - 数量の入力（1以上の整数）
 * - 賞味期限の入力（日付選択）
 * - フォームの送信処理
 *
 * フォームの状態管理とバリデーション処理は useFoodItemForm カスタムフックに
 * 分離されており、UIのみに集中しています。
 */
export const FoodItemForm: FC = () => {
  // フォームの状態と操作メソッドを取得
  const {
    title,
    quantity,
    expiryDate,
    handleTitleChange,
    handleQuantityChange,
    handleExpiryDateChange,
    handleSubmit,
  } = useFoodItemForm()

  return (
    <Box component="form" onSubmit={handleSubmit} sx={{ mb: 4 }}>
      {/* 食材名入力フィールド */}
      <FormControl fullWidth sx={{ mb: 2 }}>
        <TextField
          label="食材名"
          value={title}
          onChange={handleTitleChange}
          required
          inputProps={{
            'aria-label': '食材名入力',
          }}
        />
      </FormControl>

      {/* 数量入力フィールド */}
      <FormControl fullWidth sx={{ mb: 2 }}>
        <InputLabel htmlFor="quantity">数量</InputLabel>
        <Input
          id="quantity"
          type="number"
          value={quantity}
          onChange={handleQuantityChange}
          endAdornment={<InputAdornment position="end">個</InputAdornment>}
          inputProps={{
            min: 1,
            'aria-label': '数量入力',
          }}
          required
        />
      </FormControl>

      {/* 賞味期限入力フィールド */}
      <FormControl fullWidth sx={{ mb: 2 }}>
        <LocalizationProvider
          dateAdapter={AdapterDateFns}
          adapterLocale={jaLocale}
        >
          <DatePicker
            label="賞味期限"
            value={expiryDate}
            onChange={handleExpiryDateChange}
            format="yyyy/MM/dd"
            slotProps={{
              textField: {
                'aria-label': '賞味期限選択',
              },
            }}
          />
        </LocalizationProvider>
      </FormControl>

      {/* 送信ボタン */}
      <Button
        variant="contained"
        color="primary"
        type="submit"
        fullWidth
        sx={{ mt: 2 }}
        aria-label="食材を追加"
      >
        追加
      </Button>
    </Box>
  )
}

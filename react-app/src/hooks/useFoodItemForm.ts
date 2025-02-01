import { useState, ChangeEvent, FormEvent } from 'react'
import { useMutateFoodItem } from './useMutateFoodItem'

/**
 * 食材フォームのカスタムフック
 *
 * このフックは以下の機能を提供します：
 * - フォームの状態管理（食材名、数量、賞味期限）
 * - 入力値のバリデーション
 * - フォーム送信処理
 * - フォームのリセット
 *
 * @returns フォームの状態と操作メソッド
 */
export const useFoodItemForm = () => {
  // フォームの状態
  const [title, setTitle] = useState('')
  const [quantity, setQuantity] = useState(1)
  const [expiryDate, setExpiryDate] = useState<Date | null>(new Date())

  // 食材作成のミューテーションを取得
  const { createFoodItemMutation } = useMutateFoodItem()

  /**
   * 食材名の変更ハンドラー
   * @param e 変更イベント
   */
  const handleTitleChange = (e: ChangeEvent<HTMLInputElement>) => {
    setTitle(e.target.value)
  }

  /**
   * 数量の変更ハンドラー
   * 数値以外の入力や1未満の値は無視されます
   * @param e 変更イベント
   */
  const handleQuantityChange = (e: ChangeEvent<HTMLInputElement>) => {
    const value = parseInt(e.target.value)
    if (!isNaN(value) && value >= 1) {
      setQuantity(value)
    }
  }

  /**
   * 賞味期限の変更ハンドラー
   * @param date 選択された日付
   */
  const handleExpiryDateChange = (date: Date | null) => {
    setExpiryDate(date)
  }

  /**
   * フォーム送信ハンドラー
   * バリデーションを行い、問題なければ食材を作成します
   * @param e 送信イベント
   */
  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()

    // バリデーション
    if (title.trim() === '' || !expiryDate) {
      return
    }

    // 食材の作成
    createFoodItemMutation.mutate({
      title: title.trim(),
      quantity,
      expiry_date: expiryDate.toISOString(), // Date型をISO文字列に変換
    })

    // フォームのリセット
    resetForm()
  }

  /**
   * フォームの状態をリセット
   */
  const resetForm = () => {
    setTitle('')
    setQuantity(1)
    setExpiryDate(new Date())
  }

  return {
    // 状態
    title,
    quantity,
    expiryDate,
    // ハンドラー
    handleTitleChange,
    handleQuantityChange,
    handleExpiryDateChange,
    handleSubmit,
  }
}

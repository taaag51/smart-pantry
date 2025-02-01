import React, { FC, useState } from 'react'
import { useQueryFoodItems } from '../hooks/useQueryFoodItems'
import { useMutateFoodItem } from '../hooks/useMutateFoodItem'
import { FoodItem as FoodItemComponent } from './FoodItem'
import { FoodItem } from '../types'

export const FoodList: FC = () => {
  const [title, setTitle] = useState('')
  const [quantity, setQuantity] = useState(1)
  const [expiryDate, setExpiryDate] = useState('')
  const { data: foodItems } = useQueryFoodItems() as { data: FoodItem[] }
  const { createFoodItemMutation } = useMutateFoodItem()

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (title && quantity && expiryDate) {
      createFoodItemMutation.mutate({
        title,
        quantity,
        expiry_date: expiryDate, // Date型からstring型に変更
      })
      setTitle('')
      setQuantity(1)
      setExpiryDate('')
    }
  }

  return (
    <div className="p-6">
      <form onSubmit={handleSubmit} className="mb-6 space-y-4">
        <div>
          <label
            htmlFor="title"
            className="block text-sm font-medium text-gray-700"
          >
            食材名
          </label>
          <input
            type="text"
            id="title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            required
          />
        </div>
        <div>
          <label
            htmlFor="quantity"
            className="block text-sm font-medium text-gray-700"
          >
            数量
          </label>
          <input
            type="number"
            id="quantity"
            value={quantity}
            onChange={(e) => setQuantity(parseInt(e.target.value))}
            min="1"
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            required
          />
        </div>
        <div>
          <label
            htmlFor="expiryDate"
            className="block text-sm font-medium text-gray-700"
          >
            賞味期限
          </label>
          <input
            type="date"
            id="expiryDate"
            value={expiryDate}
            onChange={(e) => setExpiryDate(e.target.value)}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            required
          />
        </div>
        <button
          type="submit"
          className="w-full px-4 py-2 text-white bg-indigo-600 rounded hover:bg-indigo-700 focus:outline-none"
        >
          追加
        </button>
      </form>

      <div>
        <h2 className="text-xl font-bold mb-4">食材一覧</h2>
        <ul>
          {foodItems?.map((foodItem: FoodItem) => (
            <FoodItemComponent key={foodItem.id} foodItem={foodItem} />
          ))}
        </ul>
      </div>
    </div>
  )
}

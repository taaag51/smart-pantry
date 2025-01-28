import React, { FC } from 'react'
import { FoodItem as FoodItemType } from '../types'
import { useMutateFoodItem } from '../hooks/useMutateFoodItem'

interface Props {
  foodItem: FoodItemType
}

export const FoodItem: FC<Props> = ({ foodItem }) => {
  const { updateFoodItemMutation, deleteFoodItemMutation } = useMutateFoodItem()
  const isExpiringSoon = (expiryDate: Date) => {
    const today = new Date()
    const daysUntilExpiry = Math.ceil(
      (new Date(expiryDate).getTime() - today.getTime()) / (1000 * 3600 * 24)
    )
    return daysUntilExpiry <= 7 && daysUntilExpiry > 0
  }

  const isExpired = (expiryDate: Date) => {
    const today = new Date()
    return new Date(expiryDate) < today
  }

  return (
    <li className="my-3">
      <div className="flex justify-between items-center p-4 bg-white border rounded-lg shadow-lg">
        <div>
          <span
            className={`text-lg font-semibold ${
              isExpired(foodItem.expiryDate)
                ? 'text-red-600'
                : isExpiringSoon(foodItem.expiryDate)
                ? 'text-yellow-600'
                : 'text-gray-700'
            }`}
          >
            {foodItem.title}
          </span>
          <div className="mt-1 text-sm text-gray-500">
            <p>数量: {foodItem.quantity}</p>
            <p>
              賞味期限:{' '}
              {new Date(foodItem.expiryDate).toLocaleDateString('ja-JP', {
                year: 'numeric',
                month: 'long',
                day: 'numeric',
              })}
              {isExpiringSoon(foodItem.expiryDate) && (
                <span className="ml-2 text-yellow-600">
                  (期限切れまで
                  {Math.ceil(
                    (new Date(foodItem.expiryDate).getTime() -
                      new Date().getTime()) /
                      (1000 * 3600 * 24)
                  )}
                  日)
                </span>
              )}
            </p>
          </div>
        </div>
        <div className="flex space-x-2">
          <button
            onClick={() => deleteFoodItemMutation.mutate(foodItem.id)}
            className="px-3 py-2 text-white bg-red-500 rounded hover:bg-red-600 focus:outline-none"
          >
            削除
          </button>
        </div>
      </div>
    </li>
  )
}

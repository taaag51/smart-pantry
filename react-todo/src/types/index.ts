export type Task = {
  id: number
  title: string
  created_at: Date
  updated_at: Date
}

export type User = {
  email: string
  password: string
}

export type CsrfToken = {
  csrf_token: string
}

export type FoodItem = {
  id: number
  title: string
  quantity: number
  expiryDate: Date
}

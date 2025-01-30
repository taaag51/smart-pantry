export type Task = {
  id: number
  title: string
  created_at: Date
  updated_at: Date
}

export type Credential = {
  email: string
  password: string
}

export type User = {
  id: number
  email: string
  created_at: Date
  updated_at: Date
}

export type CsrfToken = {
  csrf_token: string
}

export type FoodItem = {
  id: number
  title: string
  quantity: number
  expiryDate: Date
  created_at?: Date
  updated_at?: Date
}

export type EditedTask = {
  id: number
  title: string
}

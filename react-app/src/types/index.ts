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
  expiry_date: string // Date型から文字列型に変更
  created_at?: string // 同様に文字列型に変更
  updated_at?: string // 同様に文字列型に変更
  user_id?: number
}

export type Recipe = {
  id: number
  title: string
  ingredients: string[]
  instructions: string
  created_at?: string // 同様に文字列型に変更
  updated_at?: string // 同様に文字列型に変更
}

export type EditedTask = {
  id: number
  title: string
}

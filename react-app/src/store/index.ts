import { create } from 'zustand'

type EditedTask = {
  id: number
  title: string
}

type AuthState = {
  isAuthenticated: boolean
  accessToken: string | null
  setAuth: (token: string | null) => void
  clearAuth: () => void
}

type TaskState = {
  editedTask: EditedTask
  updateEditedTask: (payload: EditedTask) => void
  resetEditedTask: () => void
}

type State = AuthState & TaskState

const useStore = create<State>((set) => ({
  // 認証状態
  isAuthenticated: false,
  accessToken: null,
  setAuth: (token) => {
    console.log('認証状態を更新:', { token, isAuthenticated: !!token })
    set({
      isAuthenticated: !!token,
      accessToken: token,
    })
  },
  clearAuth: () =>
    set({
      isAuthenticated: false,
      accessToken: null,
    }),

  // タスク状態
  editedTask: { id: 0, title: '' },
  updateEditedTask: (payload) =>
    set({
      editedTask: payload,
    }),
  resetEditedTask: () => set({ editedTask: { id: 0, title: '' } }),
}))

export default useStore

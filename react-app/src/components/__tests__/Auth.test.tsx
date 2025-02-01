import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { Auth } from '../Auth'
import { useMutateAuth } from '../../hooks/useMutateAuth'
import { useError } from '../../hooks/useError'

// モックを作成
jest.mock('../../hooks/useMutateAuth')
jest.mock('../../hooks/useError')

describe('Auth Component', () => {
  const loginMutation = {
    mutateAsync: jest.fn(),
    isPending: false,
  }

  beforeEach(() => {
    ;(useMutateAuth as jest.Mock).mockReturnValue({
      loginMutation,
      registerMutation: { mutateAsync: jest.fn(), isPending: false },
    })
    ;(useError as jest.Mock).mockReturnValue({
      errorMessage: '',
      setErrorMessage: jest.fn(),
    })
  })

  test('ログイン成功', async () => {
    loginMutation.mutateAsync.mockResolvedValueOnce({}) // 成功を模擬

    render(<Auth />)

    fireEvent.change(screen.getByLabelText(/メールアドレス/i), {
      target: { value: 'tani@example.com' },
    })
    fireEvent.change(screen.getByLabelText(/パスワード/i), {
      target: { value: 'password123' },
    })
    fireEvent.click(screen.getByRole('button', { name: /ログイン/i }))

    await waitFor(() => {
      expect(loginMutation.mutateAsync).toHaveBeenCalledWith({
        email: 'tani@example.com',
        password: 'password123',
      })
    })
  })
})

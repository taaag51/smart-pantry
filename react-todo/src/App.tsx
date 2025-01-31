import React, { FC } from 'react'
import {
  BrowserRouter,
  Routes,
  Route,
  Navigate,
  Outlet,
  useNavigate,
} from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Auth } from './components/Auth'
import { Layout } from './components/Layout'
import { PantryPage } from './pages/PantryPage'
import { FoodList } from './components/FoodList'
import { RecipeSuggestions } from './components/RecipeSuggestions'
import { useAuth } from './hooks/useAuth'
import { CircularProgress, Box } from '@mui/material'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
      refetchOnWindowFocus: false,
    },
  },
})

const AuthListener: FC = () => {
  const navigate = useNavigate()

  React.useEffect(() => {
    const handleUnauthorized = () => {
      navigate('/', { replace: true })
    }

    window.addEventListener('unauthorized', handleUnauthorized)
    return () => {
      window.removeEventListener('unauthorized', handleUnauthorized)
    }
  }, [navigate])

  return null
}

const LoadingSpinner: FC = () => (
  <Box
    sx={{
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      height: '100vh',
    }}
  >
    <CircularProgress />
  </Box>
)

const ProtectedRoute: FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated, isLoading } = useAuth()

  if (isLoading) {
    return <LoadingSpinner />
  }

  if (!isAuthenticated) {
    return <Navigate to="/" replace />
  }

  return <>{children}</>
}

const AppRoutes: FC = () => {
  const { isAuthenticated, isLoading } = useAuth()

  if (isLoading) {
    return <LoadingSpinner />
  }

  return (
    <Routes>
      {/* 認証ページ */}
      <Route
        path="/"
        element={isAuthenticated ? <Navigate to="/pantry" replace /> : <Auth />}
      />

      {/* 認証が必要なページ */}
      <Route
        element={
          <ProtectedRoute>
            <Layout>
              <Outlet />
            </Layout>
          </ProtectedRoute>
        }
      >
        {/* パントリー管理 */}
        <Route path="/pantry" element={<PantryPage />} />

        {/* 食材管理 */}
        <Route path="/food" element={<FoodList />} />

        {/* レシピ提案 */}
        <Route path="/recipes" element={<RecipeSuggestions />} />

        {/* 404ページ */}
        <Route path="*" element={<Navigate to="/pantry" replace />} />
      </Route>
    </Routes>
  )
}

const App: FC = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <AuthListener />
        <AppRoutes />
      </BrowserRouter>
    </QueryClientProvider>
  )
}

export default App

import { FC } from 'react'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { Auth } from './components/Auth'
import { PantryPage } from './pages/PantryPage'
import axios from 'axios'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import { ThemeProvider, createTheme } from '@mui/material'
import { LocalizationProvider } from '@mui/x-date-pickers'
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns'
import { ja } from 'date-fns/locale'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
      refetchOnWindowFocus: false,
    },
  },
})

const theme = createTheme({
  palette: {
    primary: {
      main: '#2196f3',
    },
    secondary: {
      main: '#f50057',
    },
  },
  typography: {
    fontFamily: [
      '-apple-system',
      'BlinkMacSystemFont',
      '"Segoe UI"',
      'Roboto',
      '"Helvetica Neue"',
      'Arial',
      'sans-serif',
    ].join(','),
  },
})

export const App: FC = () => {
  axios.defaults.withCredentials = true
  return (
    <ThemeProvider theme={theme}>
      <LocalizationProvider dateAdapter={AdapterDateFns} adapterLocale={ja}>
        <QueryClientProvider client={queryClient}>
          <BrowserRouter>
            <Routes>
              <Route path="/" element={<Auth />} />
              <Route path="/pantry" element={<PantryPage />} />
            </Routes>
          </BrowserRouter>
          <ReactQueryDevtools />
        </QueryClientProvider>
      </LocalizationProvider>
    </ThemeProvider>
  )
}

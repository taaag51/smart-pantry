import { useState, FormEvent } from 'react'
import { useMutateAuth } from '../hooks/useMutateAuth'
import {
  Container,
  Box,
  Typography,
  TextField,
  Button,
  Paper,
  IconButton,
} from '@mui/material'
import {
  KitchenOutlined,
  SwapHoriz as SwapHorizIcon,
} from '@mui/icons-material'

export const Auth = () => {
  const [email, setEmail] = useState('')
  const [pw, setPw] = useState('')
  const [isLogin, setIsLogin] = useState(true)
  const { loginMutation, registerMutation } = useMutateAuth()

  const submitAuthHandler = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (isLogin) {
      loginMutation.mutate({
        email: email,
        password: pw,
      })
    } else {
      await registerMutation
        .mutateAsync({
          email: email,
          password: pw,
        })
        .then(() =>
          loginMutation.mutate({
            email: email,
            password: pw,
          })
        )
    }
  }

  return (
    <Container component="main" maxWidth="xs">
      <Box
        sx={{
          marginTop: 8,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
        }}
      >
        <Paper
          elevation={3}
          sx={{
            padding: 4,
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            width: '100%',
          }}
        >
          <Box
            sx={{
              display: 'flex',
              flexDirection: { xs: 'column', sm: 'row' },
              alignItems: 'center',
              gap: { xs: 1, sm: 2 },
              mb: 3,
            }}
          >
            <KitchenOutlined
              sx={{
                fontSize: { xs: 32, sm: 40 },
                color: 'primary.main',
              }}
            />
            <Typography
              component="h1"
              variant="h4"
              fontWeight="bold"
              sx={{
                fontSize: { xs: '1.5rem', sm: '2rem' },
                letterSpacing: 0.5,
                color: 'primary.main',
                textAlign: 'center',
                lineHeight: 1.2,
              }}
            >
              Smart Pantry
            </Typography>
          </Box>

          <Typography component="h2" variant="h6" sx={{ mb: 3 }}>
            {isLogin ? 'ログイン' : '新規アカウント作成'}
          </Typography>

          <Box
            component="form"
            onSubmit={submitAuthHandler}
            sx={{ width: '100%' }}
          >
            <TextField
              margin="normal"
              required
              fullWidth
              id="email"
              label="メールアドレス"
              name="email"
              autoComplete="email"
              autoFocus
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
            <TextField
              margin="normal"
              required
              fullWidth
              name="password"
              label="パスワード"
              type="password"
              id="password"
              autoComplete="current-password"
              value={pw}
              onChange={(e) => setPw(e.target.value)}
            />
            <Button
              type="submit"
              fullWidth
              variant="contained"
              sx={{ mt: 3, mb: 2 }}
              disabled={!email || !pw}
            >
              {isLogin ? 'ログイン' : '登録'}
            </Button>
          </Box>

          <Box sx={{ mt: 1, display: 'flex', alignItems: 'center' }}>
            <Typography variant="body2" color="text.secondary">
              {isLogin
                ? 'アカウントをお持ちでない方は'
                : 'アカウントをお持ちの方は'}
            </Typography>
            <IconButton
              onClick={() => setIsLogin(!isLogin)}
              size="small"
              sx={{ ml: 1 }}
            >
              <SwapHorizIcon />
            </IconButton>
          </Box>
        </Paper>
      </Box>
    </Container>
  )
}

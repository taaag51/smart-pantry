import { FC, ReactNode } from 'react'
import { Box } from '@mui/material'
import { Navigation } from './Navigation'

type Props = {
  children: ReactNode
}

export const Layout: FC<Props> = ({ children }) => {
  return (
    <Box>
      <Navigation />
      <Box component="main" sx={{ p: 3 }}>
        {children}
      </Box>
    </Box>
  )
}

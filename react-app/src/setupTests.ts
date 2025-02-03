// jest-dom adds custom jest matchers for asserting on DOM nodes.
// allows you to do things like:
// expect(element).toHaveTextContent(/react/i)
// learn more: https://github.com/testing-library/jest-dom

// TextEncoder/TextDecoder polyfill for Node environment
if (
  typeof global.TextEncoder === 'undefined' ||
  typeof global.TextDecoder === 'undefined'
) {
  const util = require('util')
  global.TextEncoder = util.TextEncoder
  global.TextDecoder = util.TextDecoder
}

import '@testing-library/jest-dom'
import React from 'react'
import { server } from './mocks/server'

// Mock the environment variables
process.env.REACT_APP_API_URL = 'http://localhost:8080'

// Mock the IntersectionObserver
const mockIntersectionObserver = jest.fn()
mockIntersectionObserver.mockReturnValue({
  observe: () => null,
  unobserve: () => null,
  disconnect: () => null,
})
window.IntersectionObserver = mockIntersectionObserver

// Mock matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: jest.fn().mockImplementation((query) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: jest.fn(), // deprecated
    removeListener: jest.fn(), // deprecated
    addEventListener: jest.fn(),
    removeEventListener: jest.fn(),
    dispatchEvent: jest.fn(),
  })),
})

// Mock MUI DatePicker
jest.mock('@mui/x-date-pickers', () => {
  const React = require('react')
  return {
    DatePicker: (props: {
      label?: string
      value: Date | null
      onChange: (date: Date | null) => void
      slotProps?: {
        textField?: {
          'aria-label'?: string
        }
      }
    }) => {
      return React.createElement('input', {
        type: 'date',
        'aria-label': props.slotProps?.textField?.['aria-label'] || props.label,
        value: props.value ? props.value.toISOString().split('T')[0] : '',
        onChange: (e: React.ChangeEvent<HTMLInputElement>) =>
          props.onChange(e.target.value ? new Date(e.target.value) : null),
      })
    },
    LocalizationProvider: ({ children }: { children: React.ReactNode }) =>
      children,
  }
})

// Mock date-fns locale
jest.mock('date-fns/locale/ja', () => ({
  default: {
    code: 'ja',
    formatLong: {},
    formatRelative: {},
    formatDistance: {},
    localize: {},
    match: {},
    options: {},
  },
}))

// MSWのセットアップ
beforeAll(() => {
  // MSWサーバーを起動
  server.listen()
})

afterEach(() => {
  // 各テスト後にハンドラーをリセット
  server.resetHandlers()
})

afterAll(() => {
  // テスト終了後にサーバーをクリーンアップ
  server.close()
})

// エラー処理のグローバル設定
const originalError = console.error
beforeAll(() => {
  console.error = (...args) => {
    if (
      /Warning: ReactDOM.render is no longer supported in React 18./.test(
        args[0]
      )
    ) {
      return
    }
    originalError.call(console, ...args)
  }
})

afterAll(() => {
  console.error = originalError
})

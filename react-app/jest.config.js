module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'jsdom',
  transform: {
    '^.+\\.tsx?$': 'ts-jest',
  },
  moduleNameMapper: {
    '\\.(css|less|scss|sass)$': 'identity-obj-proxy',
  },
  transformIgnorePatterns: [
    '/node_modules/(?!axios)', // axiosをトランスパイル対象に追加
  ],
  setupFilesAfterEnv: ['<rootDir>/src/setupTests.ts'],
}

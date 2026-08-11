/// <reference types="vite/client" />

interface ImportMetaEnv {
  /** Empty string = same-origin (nginx proxy in Docker). */
  readonly VITE_API_BASE_URL?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import PublicRouter from './routes/PublicRouter'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <PublicRouter />
  </StrictMode>,
)

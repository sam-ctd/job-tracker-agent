import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom"
import { Layout } from "./layout/Layout"
import { KanbanPage } from "./pages/KanbanPage"
import { NuevaPage } from "./pages/NuevaPage"

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<Layout />}>
          <Route index element={<Navigate to="/kanban" replace />} />
          <Route path="kanban" element={<KanbanPage />} />
          <Route path="nueva" element={<NuevaPage />} />
          <Route path="*" element={<Navigate to="/kanban" replace />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

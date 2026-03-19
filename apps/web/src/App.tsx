import { useEffect } from "react"
import { BrowserRouter, Navigate, Route, Routes, useNavigate } from "react-router-dom"
import { AuthProvider, useAuth } from "@/contexts/AuthContext"
import { setOnUnauthorized } from "@/lib/api"
import { AppLayout } from "@/components/layout/AppLayout"
import { Login } from "@/pages/Login"
import { Signup } from "@/pages/Signup"
import { Dashboard } from "@/pages/Dashboard"
import { CodingLog } from "@/pages/CodingLog"
import { Ideas } from "@/pages/Ideas"
import { LeetCode } from "@/pages/LeetCode"
import { Opportunities } from "@/pages/Opportunities"
import { Finance } from "@/pages/Finance"
import { AIGenerator } from "@/pages/AIGenerator"

function AuthSetup() {
  const { logout } = useAuth()
  const navigate = useNavigate()
  useEffect(() => {
    setOnUnauthorized(() => {
      logout()
      navigate("/login", { replace: true })
    })
    return () => setOnUnauthorized(null)
  }, [logout, navigate])
  return null
}

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuth()
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }
  return (
    <>
      <AuthSetup />
      {children}
    </>
  )
}

function PublicOnlyRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuth()
  if (isAuthenticated) {
    return <Navigate to="/" replace />
  }
  return <>{children}</>
}

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route
            path="/login"
            element={
              <PublicOnlyRoute>
                <Login />
              </PublicOnlyRoute>
            }
          />
          <Route
            path="/signup"
            element={
              <PublicOnlyRoute>
                <Signup />
              </PublicOnlyRoute>
            }
          />
          <Route
            path="/"
            element={
              <ProtectedRoute>
                <AppLayout />
              </ProtectedRoute>
            }
          >
            <Route index element={<Dashboard />} />
            <Route path="coding-log" element={<CodingLog />} />
            <Route path="ideas" element={<Ideas />} />
            <Route path="leetcode" element={<LeetCode />} />
            <Route path="opportunities" element={<Opportunities />} />
            <Route path="finance" element={<Finance />} />
            <Route path="ai" element={<AIGenerator />} />
          </Route>
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  )
}

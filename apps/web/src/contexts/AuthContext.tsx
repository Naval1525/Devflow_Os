import React, { createContext, useCallback, useContext, useState } from "react"
import { API_URL } from "@/lib/utils"

type AuthContextType = {
  token: string | null
  login: (email: string, password: string) => Promise<void>
  signup: (email: string, password: string) => Promise<void>
  logout: () => void
  isAuthenticated: boolean
}

const AuthContext = createContext<AuthContextType | null>(null)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setToken] = useState<string | null>(() => localStorage.getItem("token"))

  const login = useCallback(async (email: string, password: string) => {
    const res = await fetch(`${API_URL}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    })
    if (!res.ok) {
      const data = await res.json().catch(() => ({}))
      const message = res.status === 403 ? "Access denied" : (data.error || "Login failed")
      throw new Error(message)
    }
    const data = await res.json()
    const t = data.token
    if (t) {
      localStorage.setItem("token", t)
      setToken(t)
    }
  }, [])

  const signup = useCallback(async (email: string, password: string) => {
    const res = await fetch(`${API_URL}/auth/signup`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    })
    if (!res.ok) {
      const data = await res.json().catch(() => ({}))
      throw new Error(data.error || "Signup failed")
    }
    const data = await res.json()
    const t = data.token
    if (t) {
      localStorage.setItem("token", t)
      setToken(t)
    }
  }, [])

  const logout = useCallback(() => {
    localStorage.removeItem("token")
    setToken(null)
  }, [])

  return (
    <AuthContext.Provider
      value={{
        token,
        login,
        signup,
        logout,
        isAuthenticated: !!token,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error("useAuth must be used within AuthProvider")
  return ctx
}

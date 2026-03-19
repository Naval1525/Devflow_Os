import { API_URL, getAuthHeaders } from "./utils"

type OnUnauthorized = () => void

let onUnauthorized: OnUnauthorized | null = null

export function setOnUnauthorized(fn: OnUnauthorized | null) {
  onUnauthorized = fn
}

export async function apiFetch(
  path: string,
  options: RequestInit = {}
): Promise<Response> {
  const res = await fetch(`${API_URL}${path}`, {
    ...options,
    headers: {
      ...getAuthHeaders(),
      ...options.headers,
    },
  })
  if (res.status === 401 && onUnauthorized) {
    onUnauthorized()
  }
  return res
}

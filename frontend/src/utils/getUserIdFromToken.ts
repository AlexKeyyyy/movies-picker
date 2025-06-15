export function getUserIdFromToken(): number | null {
    const token = localStorage.getItem('token')
    if (!token) return null
  
    try {
      const payload = JSON.parse(atob(token.split('.')[1]))
      return payload.user_id || null
    } catch {
      return null
    }
  }  
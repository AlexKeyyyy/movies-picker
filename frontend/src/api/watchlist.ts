import { getUserIdFromToken } from '../utils/getUserIdFromToken'
import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
})

api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export const getWatchlist = () => {
  const userId = getUserIdFromToken()
  return api.get(`/users/${userId}/watchlist`)
}

export const addToWatchlist = (movieId: number) => {
  const userId = getUserIdFromToken()
  return api.post(`/users/${userId}/watchlist`, { movie_id: movieId })
}

export const removeFromWatchlist = (movieId: number) => {
  const userId = getUserIdFromToken()
  return api.delete(`/users/${userId}/watchlist/${movieId}`)
}

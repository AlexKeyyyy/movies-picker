import api from './client'

export interface Movie {
  movie_id: number
  title: string
  year: number
  poster_url: string
  description: string
  ratingKinopoisk?: number
}

export const searchMovies = (query: string) =>
  api.get<Movie[]>(`/movies/search?q=${encodeURIComponent(query)}`)

export const getMovieDetails = (id: string | undefined) =>
  api.get<Movie>(`/movies/${id}`)

export const listMovies = (page = 1, size = 20) =>
    api.get<Movie[]>(`/movies?page=${page}&size=${size}`)

export const listPopularMovies = (limit: number) =>
    api.get<Movie[]>(`http://localhost:8080/movies/popular?limit=${limit}`)
  

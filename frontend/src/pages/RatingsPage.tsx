import { useEffect, useState } from 'react'
import { List, Card, Rate, Typography, message } from 'antd'
import axios from 'axios'

interface RatedMovie {
  movie_id: number
  rating: number
  rated_at: string
  // после дозагрузки
  title?: string
  year?: number
  poster_url?: string
}

export default function RatingsPage() {
  const [ratedMovies, setRatedMovies] = useState<RatedMovie[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchRatingsAndDetails = async () => {
      try {
        const token = localStorage.getItem('token')
        if (!token) {
          message.error('Пользователь не авторизован')
          return
        }
        const { data: ratings } = await axios.get('/api/users/me/ratings', {
          headers: { Authorization: `Bearer ${token}` }
        })

        // Параллельно запросить данные фильмов
        const movieDetailsPromises = ratings.map((r: RatedMovie) =>
          axios.get(`/api/movies/${r.movie_id}`).then(res => res.data)
        )
        const moviesDetails = await Promise.all(movieDetailsPromises)

        // Объединяем рейтинги и детали
        const merged = ratings.map((r: RatedMovie) => {
          const details = moviesDetails.find((m: any) => m.movie_id === r.movie_id)
          return { ...r, ...details }
        })

        setRatedMovies(merged)
      } catch {
        message.error('Ошибка загрузки рейтингов')
      } finally {
        setLoading(false)
      }
    }
    fetchRatingsAndDetails()
  }, [])

  if (loading) return <div>Загрузка...</div>

  return (
    <div style={{ padding: 24 }}>
      <Typography.Title level={2}>Мои рейтинги</Typography.Title>
      <List
        grid={{ gutter: 16, column: 4 }}
        dataSource={ratedMovies ?? []}
        renderItem={movie => (
          <List.Item key={movie.movie_id}>
            <Card
              cover={
                <div
                  style={{
                    height: 300,
                    background: '#f0f0f0',
                    display: 'flex',
                    justifyContent: 'center',
                    alignItems: 'center',
                  }}
                >
                  <img
                    alt={movie.title || 'фильм'}
                    src={movie.poster_url || 'https://via.placeholder.com/300x450?text=No+Image'}
                    style={{
                      maxHeight: '100%',
                      maxWidth: '100%',
                      objectFit: 'contain',
                    }}
                  />
                </div>
              }
              hoverable
            >
              <Card.Meta
                title={`${movie.title || 'Без названия'} (${movie.year || '—'})`}
                description={
                  <>
                    <Rate disabled value={movie.rating} />{' '}
                    <Typography.Text>({movie.rating})</Typography.Text>
                  </>
                }
              />
            </Card>
          </List.Item>
        )}
      />
    </div>
  )
}

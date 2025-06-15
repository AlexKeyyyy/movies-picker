import { useParams } from 'react-router-dom'
import { useEffect, useState } from 'react'
import { Button, Card, Spin, Typography, message, Rate } from 'antd'
import axios from 'axios'

interface Movie {
  movie_id: number
  title: string
  year: number
  poster_url: string
  description: string
}

interface ReviewItem {
  video_id: string
  video_url: string
  title: string
  channel_title: string
  thumbnail_url: string
}

export default function MovieDetailsPage() {
  const { id } = useParams()
  const [movie, setMovie] = useState<Movie | null>(null)
  const [reviews, setReviews] = useState<ReviewItem[]>([])
  const [loading, setLoading] = useState(true)
  const [inWatchlist, setInWatchlist] = useState(false)
  const [saving, setSaving] = useState(false)
  const [rating, setRating] = useState<number | null>(null)
  const [savingRating, setSavingRating] = useState(false)

  const token = localStorage.getItem('token')

  useEffect(() => {
    setLoading(true)
    axios.get(`/api/movies/${id}`)
      .then(res => setMovie(res.data))
      .catch(() => message.error('Не удалось загрузить фильм'))
      .finally(() => setLoading(false))
  }, [id])

  useEffect(() => {
    if (!token) return
    // Проверка в списке и загрузка рейтинга пользователя для этого фильма
    const fetchUserData = async () => {
      try {
        // watchlist
        const wlRes = await axios.get('/api/users/me/watchlist', {
          headers: { Authorization: `Bearer ${token}` }
        })
        const ids = wlRes.data.map((m: Movie) => m.movie_id)
        setInWatchlist(ids.includes(Number(id)))

        // рейтинг
        const ratingsRes = await axios.get('/api/users/me/ratings', {
          headers: { Authorization: `Bearer ${token}` }
        })
        const myRating = ratingsRes.data.find((r: { movie_id: number }) => r.movie_id === Number(id))
        if (myRating) setRating(myRating.rating)
        else setRating(null)
      } catch (err) {
        console.error('Ошибка получения пользовательских данных:', err)
      }
    }
    fetchUserData()
  }, [id, token])

  useEffect(() => {
    axios.get(`/api/movies/${id}/reviews`)
      .then(res => {
        console.log('Reviews:', res.data)
        setReviews(res.data)})
      .catch(() => message.error('Не удалось загрузить обзоры'))
  }, [id])

  const handleAddToWatchlist = async () => {
    if (!movie || !token) return
    try {
      setSaving(true)
      await axios.post(
        `/api/users/me/watchlist`,
        { movie_id: movie.movie_id },
        { headers: { Authorization: `Bearer ${token}` } }
      )
      setInWatchlist(true)
      message.success('Фильм добавлен в список')
    } catch {
      message.error('Ошибка при добавлении')
    } finally {
      setSaving(false)
    }
  }

  const handleRemoveFromWatchlist = async () => {
    if (!movie || !token) return
    try {
      await axios.delete(`/api/users/me/watchlist/${movie.movie_id}`, {
        headers: { Authorization: `Bearer ${token}` }
      })
      setInWatchlist(false)
      message.success('Фильм удалён из списка')
    } catch {
      message.error('Ошибка при удалении')
    }
  }

  const handleRatingChange = async (value: number) => {
    if (!movie || !token) return
    try {
      setSavingRating(true)
      await axios.post(
        `/api/users/me/ratings`,
        { movie_id: movie.movie_id, rating: value },
        { headers: { Authorization: `Bearer ${token}` } }
      )
      setRating(value)
      message.success('Рейтинг сохранён')
    } catch {
      message.error('Ошибка при сохранении рейтинга')
    } finally {
      setSavingRating(false)
    }
  }

  if (loading) return <Spin style={{ display: 'block', margin: '100px auto' }} />
  if (!movie) return null

  return (
    <Card
      style={{ maxWidth: 800, margin: '40px auto' }}
      cover={
        <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', background: '#f0f2f5' }}>
          <img
            alt={movie.title}
            src={movie.poster_url}
            style={{
              width: '100%',
              height: 'auto',
              maxHeight: 500,
              objectFit: 'contain',
              backgroundColor: '#f0f2f5'
            }}
          />
        </div>
      }      
    >
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Typography.Title level={3} style={{ margin: 0 }}>
          {movie.title} ({movie.year})
        </Typography.Title>
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          {inWatchlist ? (
            <Button danger onClick={handleRemoveFromWatchlist}>
              Убрать из списка
            </Button>
          ) : (
            <Button loading={saving} type="primary" onClick={handleAddToWatchlist}>
              К просмотру
            </Button>
          )}
          <Rate
            allowClear
            count={5}
            value={rating ?? 0}
            onChange={handleRatingChange}
            disabled={savingRating}
          />
        </div>
      </div>

      <Typography.Paragraph>{movie.description}</Typography.Paragraph>

      <Typography.Title level={4} style={{ marginTop: 32 }}>
        Обзоры с YouTube
      </Typography.Title>
      {reviews.length === 0 ? (
        <Typography.Text>Обзоры не найдены</Typography.Text>
      ) : (
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 16 }}>
          {reviews.map((r) => (
            <Card
              key={r.video_id}
              hoverable
              cover={<img alt={r.title} src={r.thumbnail_url} />}
              onClick={() => window.open(r.video_url, '_blank')}
            >
              <Typography.Text strong>{r.title}</Typography.Text>
              <br />
              <Typography.Text type="secondary">{r.channel_title}</Typography.Text>
            </Card>
          ))}
        </div>
      )}
    </Card>
  )
}

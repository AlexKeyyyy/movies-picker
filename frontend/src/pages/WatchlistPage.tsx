import { useEffect, useState } from 'react'
import { List, Card, message, Button } from 'antd'
import axios from 'axios'
import { useNavigate } from 'react-router-dom'

interface Movie {
  movie_id: number
  title: string
  year: number
  poster_url: string
  description: string
}

export default function WatchlistPage() {
  const [movies, setMovies] = useState<Movie[]>([])
  const navigate = useNavigate()

  const fetchWatchlist = async () => {
    try {
      const token = localStorage.getItem('token')
      const res = await axios.get('/api/users/me/watchlist', {
        headers: { Authorization: `Bearer ${token}` }
      })
      console.log('WatchList from server: ', res);
      setMovies(res.data)
    } catch (err) {
      console.error(err)
      message.error('Ошибка загрузки списка просмотра')
    }
  }

  const removeFromWatchlist = async (movieId: number) => {
    try {
      const token = localStorage.getItem('token')
      await axios.delete(`/api/users/me/watchlist/${movieId}`, {
        headers: { Authorization: `Bearer ${token}` }
      })
      setMovies(prev => prev.filter(m => m.movie_id !== movieId))
      message.success('Фильм удалён из списка')
    } catch (err) {
      console.error(err)
      message.error('Не удалось удалить фильм')
    }
  }

  useEffect(() => {
    fetchWatchlist()
  }, [])

  return (
    <div style={{ padding: 24 }}>
      <h2>Список к просмотру</h2>
      <List
        grid={{ gutter: 16, column: 4 }}
        dataSource={movies ?? []}
        renderItem={movie => (
          <List.Item>
            <Card
              cover={
                <div
                  style={{
                    height: 300,
                    background: '#f0f0f0',
                    display: 'flex',
                    justifyContent: 'center',
                    alignItems: 'center',
                    cursor: 'pointer',
                  }}
                  onClick={() => navigate(`/movies/${movie.movie_id}`)}
                >
                  <img
                    alt={movie.title}
                    src={movie.poster_url}
                    style={{
                      maxHeight: '100%',
                      maxWidth: '100%',
                      objectFit: 'contain',
                    }}
                  />
                </div>
              }              
              actions={[
                <Button danger onClick={() => removeFromWatchlist(movie.movie_id)}>Удалить</Button>
              ]}
              hoverable
            >
              <Card.Meta title={`${movie.title}`} description={movie.description} />
            </Card>
          </List.Item>
        )}
      />
    </div>
  )
}

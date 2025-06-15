import { useEffect, useState, useRef } from 'react'
import { List, Card, Typography, message, Input, Spin } from 'antd'
import { useNavigate } from 'react-router-dom'
import { listPopularMovies, searchMovies, Movie } from '../api/movies'  // добавим searchMovies
import { StarOutlined } from '@ant-design/icons'

const { Title } = Typography

export default function MoviesPage() {
  const [movies, setMovies] = useState<Movie[]>([])
  const [loading, setLoading] = useState(false)
  const [searchText, setSearchText] = useState('')
  const navigate = useNavigate()
  const debounceTimeout = useRef<NodeJS.Timeout | null>(null)

  // Изначально загружаем популярные фильмы
  const fetchPopularMovies = async () => {
    try {
      setLoading(true)
      const res = await listPopularMovies(20)
      const filtered = res.data.filter(movie => movie.title.trim() !== '')
      setMovies(filtered)
    } catch {
      message.error('Ошибка загрузки популярных фильмов')
    } finally {
      setLoading(false)
    }
  }

  // Поиск фильмов на сервере
  const fetchSearchMovies = async (query: string) => {
    if (!query.trim()) {
      // Если поисковый запрос пустой — показываем популярные фильмы
      fetchPopularMovies()
      return
    }

    try {
      setLoading(true)
      const res = await searchMovies(query)
      setMovies(res.data)
    } catch {
      message.error('Ошибка поиска фильмов')
    } finally {
      setLoading(false)
    }
  }

  // Запускаем начальную загрузку популярных фильмов
  useEffect(() => {
    fetchPopularMovies()
  }, [])

  // При изменении searchText запускаем поиск с дебаунсом
  useEffect(() => {
    if (debounceTimeout.current) clearTimeout(debounceTimeout.current)
    debounceTimeout.current = setTimeout(() => {
      fetchSearchMovies(searchText)
    }, 500) // 500ms debounce
  }, [searchText])

  return (
    <div style={{ padding: 24 }}>
      <Title level={2}>Популярные фильмы</Title>

      <Input
        placeholder="Поиск фильмов..."
        value={searchText}
        onChange={e => setSearchText(e.target.value)}
        style={{ marginBottom: 24, maxWidth: 400 }}
        allowClear
      />

      {loading ? (
        <Spin style={{ display: 'block', margin: '50px auto' }} />
      ) : (
        <List
          grid={{ gutter: 16, column: 4 }}
          dataSource={movies ?? []}
          locale={{ emptyText: 'Фильмы не найдены' }}
          renderItem={movie => (
            <List.Item>
              <Card
                hoverable
                onClick={() => navigate(`/movies/${movie.movie_id}`)}
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
              >
                <Card.Meta
                  title={`${movie.title} (${movie.year})`}
                  description={
                    movie.ratingKinopoisk !== undefined ? (
                      <span>
                        <StarOutlined style={{ color: '#fadb14', marginRight: 4 }} />
                        {movie.ratingKinopoisk.toFixed(1)}
                      </span>
                    ) : (
                      'Рейтинг отсутствует'
                    )
                  }
                />
              </Card>
            </List.Item>
          )}
        />
      )}
    </div>
  )
}

import { Layout } from 'antd'
import { Route, Routes } from 'react-router-dom'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import MoviesPage from './pages/MoviesPage'
import MovieDetailsPage from './pages/MovieDetailsPage'
import WatchlistPage from './pages/WatchlistPage'
import RatingsPage from './pages/RatingsPage'
import EditProfilePage from './pages/EditProfilePage'
import Navbar from './components/Navbar'
import PrivateRoute from './routes/PrivateRoute'

const { Header, Content } = Layout

export default function App() {
  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header>
        <Navbar />
      </Header>
      <Content style={{ padding: '24px' }}>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/movies/:id" element={<MovieDetailsPage />} />
          <Route
            path="/watchlist"
            element={
              <PrivateRoute>
                <WatchlistPage />
              </PrivateRoute>
            }
          />
          <Route
            path="/ratings"
            element={
              <PrivateRoute>
                <RatingsPage />
              </PrivateRoute>
            }
          />
          <Route
            path="/profile/edit"
            element={
            <PrivateRoute>
              <EditProfilePage />
            </PrivateRoute>
  }
/>
          <Route path="/" element={<MoviesPage />} />
        </Routes>
      </Content>
    </Layout>
  )
}

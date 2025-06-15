import { Menu } from 'antd'
import { useNavigate } from 'react-router-dom'
import {
  HomeOutlined,
  StarOutlined,
  LoginOutlined,
  UserAddOutlined,
  VideoCameraOutlined,
  LogoutOutlined,
  UserOutlined,
} from '@ant-design/icons'
import { useAuth } from '../context/AuthContext'

export default function Navbar() {
  const navigate = useNavigate()
  const { isAuthenticated, logout } = useAuth()

  const items = [
    { key: '/', icon: <HomeOutlined />, label: 'Главная' },
    ...(isAuthenticated
      ? [
          { key: '/watchlist', icon: <VideoCameraOutlined />, label: 'К просмотру' },
          { key: '/ratings', icon: <StarOutlined />, label: 'Мои оценки' },
          { key: '/profile/edit', icon: <UserOutlined />, label: 'Профиль' },
          { key: 'logout', icon: <LogoutOutlined />, label: 'Выйти' },
        ]
      : [
          { key: '/login', icon: <LoginOutlined />, label: 'Войти' },
          { key: '/register', icon: <UserAddOutlined />, label: 'Регистрация' },
        ]),
  ]

  const handleClick = ({ key }: { key: string }) => {
    if (key === 'logout') {
      logout()
      navigate('/')
    } else {
      navigate(key)
    }
  }

  return (
    <Menu
      mode="horizontal"
      theme="dark"
      onClick={handleClick}
      selectable={false}
      style={{ display: 'flex', justifyContent: 'center' }}
      items={items}
    />
  )
}

import { Button, Card, Form, Input, Typography, message } from 'antd'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import axios from 'axios'

export default function LoginPage() {
  const navigate = useNavigate()
  const { login } = useAuth()

  const onFinish = async (values: any) => {
    try {
      const res = await axios.post('/api/auth/login', values)
      login(res.data.access_token)
      localStorage.setItem('token', res.data.access_token)
      message.success('Успешный вход')
      navigate('/')
    } catch {
      message.error('Неверные учетные данные')
    }
  }

  return (
    <Card title="Вход" style={{ maxWidth: 400, margin: '100px auto' }}>
      <Form layout="vertical" onFinish={onFinish}>
        <Form.Item name="email" label="Email" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
        <Form.Item name="password" label="Пароль" rules={[{ required: true }]}>
          <Input.Password />
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit" block>
            Войти
          </Button>
        </Form.Item>
      </Form>
      <Typography.Text>
        Нет аккаунта? <a href="/register">Зарегистрируйтесь</a>
      </Typography.Text>
    </Card>
  )
}

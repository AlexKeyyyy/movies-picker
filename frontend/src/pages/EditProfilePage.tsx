import { useEffect, useState } from 'react'
import { Card, Form, Input, Button, message, Typography } from 'antd'
import axios from 'axios'

export default function EditProfilePage() {
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [initialEmail, setInitialEmail] = useState('')

  useEffect(() => {
    // Загрузим текущий профиль
    const token = localStorage.getItem('token')
    axios.get('/api/users/me', {
      headers: { Authorization: `Bearer ${token}` }
    }).then(res => {
      form.setFieldsValue({ email: res.data.email })
      setInitialEmail(res.data.email)
    }).catch(() => {
      message.error('Не удалось загрузить профиль')
    })
  }, [form])

  const onFinish = async (values: {
    email?: string
    password?: string
    current_password?: string
  }) => {
    const token = localStorage.getItem('token')
    setLoading(true)

    try {
      await axios.patch('/api/users/me', values, {
        headers: { Authorization: `Bearer ${token}` }
      })
      message.success('Профиль обновлён')
    } catch (err: any) {
      message.error(err?.response?.data || 'Ошибка при обновлении профиля')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card title="Редактировать профиль" style={{ maxWidth: 500, margin: '100px auto' }}>
      <Form form={form} layout="vertical" onFinish={onFinish}>
        <Form.Item
          name="email"
          label="Email"
          rules={[{ type: 'email', message: 'Некорректный email' }]}
        >
          <Input placeholder="Введите новый email (если нужно)" />
        </Form.Item>

        <Form.Item
  name="password"
  label="Новый пароль"
  rules={[
    {
      validator: (_, value) => {
        if (!value) {
          // Если поле пустое — валидация проходит, т.к. менять пароль необязательно
          return Promise.resolve()
        }

        const strongPasswordRegex = /^(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*()_+\[\]{};':"\\|,.<>\/?]).{8,}$/
        if (!strongPasswordRegex.test(value)) {
          return Promise.reject(
            new Error(
              'Пароль должен быть не короче 8 символов, содержать заглавную букву, цифру и спецсимвол'
            )
          )
        }
        return Promise.resolve()
      }
    }
  ]}
>
  <Input.Password placeholder="Оставьте пустым, если не меняется" />
</Form.Item>

        <Form.Item
          name="current_password"
          label="Текущий пароль"
          rules={[
            {
              required: true,
              message: 'Введите текущий пароль для подтверждения изменений',
            },
          ]}
        >
          <Input.Password placeholder="Обязательно для изменения email или пароля" />
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit" block loading={loading}>
            Сохранить изменения
          </Button>
        </Form.Item>
      </Form>

      <Typography.Text type="secondary">
        Если вы хотите сменить email или пароль — введите текущий пароль.
      </Typography.Text>
    </Card>
  )
}

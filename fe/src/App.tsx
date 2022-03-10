import React, { useEffect } from 'react';
import logo from './logo.svg';
import { Form, Input, Button } from 'antd';
import './App.css';


function App() {
  const [token ,setToken] = React.useState<string>('');
  const [username ,setUsername] = React.useState<string>('');

  useEffect(() => {
    async function refresh() {
      // You can await here
      const resp = await fetch('http://localhost:8080/refresh', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
      });
      console.log("status is " + resp.status);
      if (resp.status === 204) {
        console.log("exit")
        return
      }

      const d = await resp.json();
      if(resp.ok) {
        console.log('Login Success:', d);
        setToken(d.access_token);
        user(d.access_token);
      }
    }

    refresh();

  }, []);

  const user = async (token: string) => {
    const resp = await fetch('http://localhost:8080/user', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `${token}`,
      },
      credentials: 'include',
    });

    const d = await resp.json();
    if(resp.ok) {
      console.log('User Success:', d);
      setUsername(d.name);
    }
  }

  const onFinish = async (values: any) => {
    const resp = await fetch('http://localhost:8080/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({'name': values.username}),
      credentials: 'include',
    });
    const {data, errors} = await resp.json();
    if (resp.ok) {
      console.log('Login Success:', data);
      setToken(data.access_token);
      setUsername(values.username);

    } else {
      console.error(errors);
    }
  };

  const onFinishFailed = (errorInfo: any) => {
    console.log('Failed:', errorInfo);
  };

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Welcome { username }.
        </p>
        <Form
          name="basic"
          labelCol={{ span: 8 }}
          wrapperCol={{ span: 16 }}
          initialValues={{ remember: true }}
          onFinish={onFinish}
          onFinishFailed={onFinishFailed}
          autoComplete="off"
        >
          <Form.Item
            label="Username"
            name="username"
            rules={[{ required: true, message: 'Please input your username!' }]}
          >
            <Input />
          </Form.Item>
          
          <Form.Item wrapperCol={{ offset: 8, span: 16 }}>
            <Button type="primary" htmlType="submit">
              Login
            </Button>
          </Form.Item>
        </Form>
      </header>
    </div>
  );
}

export default App;

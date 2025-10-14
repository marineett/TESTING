import React from 'react';
import { Layout, Menu } from 'antd';
import { Link, useLocation } from 'react-router-dom';

const { Header } = Layout;

const Navbar: React.FC = () => {
  const location = useLocation();

  return (
    <Header style={{ position: 'fixed', zIndex: 1, width: '100%' }}>
      <div className="logo" />
      <Menu
        theme="dark"
        mode="horizontal"
        selectedKeys={[location.pathname]}
        items={[
          {
            key: '/',
            label: <Link to="/">Home</Link>,
          },
          {
            key: '/login',
            label: <Link to="/login">Login</Link>,
          },
          {
            key: '/register',
            label: <Link to="/register">Register</Link>,
          },
        ]}
      />
    </Header>
  );
};

export default Navbar; 
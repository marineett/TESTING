import React from 'react';
import { Routes, Route } from 'react-router-dom';
import { Layout } from 'antd';
import Navbar from './components/Navbar';
import HomePage from './pages/HomePage';
import Login from './pages/Login';
import RegistrationPage from './pages/RegistrationPage';
import ClientPage from './pages/ClientPage';
import RepetitorPage from './pages/RepetitorPage';
import ModeratorPage from './pages/ModeratorPage';
import AdminPage from './pages/AdminPage';
import ChatPage from './pages/ChatPage';

const { Content } = Layout;

const App: React.FC = () => {
  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Navbar />
      <Content style={{ padding: '24px', marginTop: 64 }}>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<RegistrationPage />} />
          <Route path="/client" element={<ClientPage />} />
          <Route path="/repetitor" element={<RepetitorPage />} />
          <Route path="/moderator" element={<ModeratorPage />} />
          <Route path="/admin" element={<AdminPage />} />
          <Route path="/chat" element={<ChatPage />} />
        </Routes>
      </Content>
    </Layout>
  );
};

export default App; 
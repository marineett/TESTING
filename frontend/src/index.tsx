import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import App from './App';
import 'antd/dist/reset.css'; // Updated import for antd v5
import './index.css'; // Добавлен импорт пользовательских стилей

// This line imports reportWebVitals, remove it if not used elsewhere
// import reportWebVitals from './reportWebVitals';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

root.render(
  <React.StrictMode>
    <BrowserRouter>
    <App />
    </BrowserRouter>
  </React.StrictMode>
);

// This line calls reportWebVitals, remove it
// reportWebVitals(); 
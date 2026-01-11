import ReactDOM from 'react-dom/client';
import App from './App';
import './index.css';
import { BrowserRouter } from 'react-router';
import { ConfigProvider } from 'antd';
import { StyleProvider } from '@ant-design/cssinjs';

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <StyleProvider layer>
    <ConfigProvider
      theme={{
        token: {
          fontSize: 18,
        },
      }}
    >
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </ConfigProvider>
  </StyleProvider>
);

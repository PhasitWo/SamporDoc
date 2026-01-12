import ReactDOM from 'react-dom/client';
import App from './App';
import './index.css';
import { BrowserRouter } from 'react-router';
import { StyleProvider } from '@ant-design/cssinjs';


ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <StyleProvider layer>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </StyleProvider>
);

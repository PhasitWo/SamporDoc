import { Routes, Route, Outlet, useNavigate } from 'react-router';
import { EventsOn } from '../wailsjs/runtime';
import { CreateReceipt } from './pages/CreateReceipt';
import { Setting } from './pages/Setting';

const Layout = () => (
  <div className="w-[100vw] h-[100vh] bg-gray-100 text-[18px]">
    <Outlet />
  </div>
);

function App() {
  const navigate = useNavigate();

  EventsOn('navigate', (route: string) => {
    navigate(route);
  });

  return (
    <Routes>
      <Route element={<Layout />}>
        <Route index element={<CreateReceipt />} />
        <Route path="/setting" element={<Setting />} />
      </Route>
    </Routes>
  );
}

export default App;

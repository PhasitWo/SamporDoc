import { Routes, Route, Outlet, useNavigate, useLocation } from 'react-router';
import { EventsOn } from '../wailsjs/runtime';
import CreateReceipt from './pages/CreateReceipt';
import Setting from './pages/Setting';
import { ConfigProvider, Layout, Menu, App } from 'antd';
import { useAppStore } from './store/useAppStore';
import { useEffect } from 'react';
import th from 'antd/es/date-picker/locale/th_TH';
import thTH from 'antd/es/locale/th_TH';
import dayjs from 'dayjs';
import 'dayjs/locale/th';
import buddhistEra from 'dayjs/plugin/buddhistEra';
import { ErrorBoundary } from 'react-error-boundary';
import ErrorFallback from './pages/ErrorFallback';
import { useShowBoundary } from './utils';
import CreateProcurement from './pages/CreateProcurement';
import Automove from './pages/Automove';

dayjs.locale('th');
dayjs.extend(buddhistEra);

const buddhistLocale: typeof th = {
  ...th,
  lang: {
    ...th.lang,
    fieldDateFormat: 'DD MMMM BBBB',
    fieldDateTimeFormat: 'DD-MM-BBBB HH:mm:ss',
    yearFormat: 'BBBB',
    cellYearFormat: 'BBBB',
  },
};

const globalBuddhistLocale: typeof thTH = {
  ...thTH,
  DatePicker: {
    ...thTH.DatePicker!,
    lang: buddhistLocale.lang,
  },
};

const primaryColorMap: Record<string, string | undefined> = {
  '/': '#00b96b',
  '/createProcurement': '#e68415',
  '/automove': '#f24141',
};

const AppLayout = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const init = useAppStore((s) => s.init);
  const { showBoundary } = useShowBoundary();

  useEffect(() => {
    init().catch(showBoundary);
  }, []);

  // Prevent Backspace/Delete from navigating back
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Backspace' || e.key === 'Delete') {
        const target = e.target as HTMLElement;
        const isEditable = target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable;
        if (!isEditable) {
          e.preventDefault();
        }
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, []);

  return (
    <>
      <Layout className="w-full h-full bg-white text-[14px]">
        <Layout.Header className="px-5">
          <Menu
            mode="horizontal"
            items={[
              { key: '/', label: 'สร้างใบเสร็จรับเงิน' },
              { key: '/createProcurement', label: 'สร้างจัดซื้อจัดจ้าง' },
              { key: '/automove', label: 'Auto Move' },
              { key: '/setting', label: 'ตั้งค่า' },
            ]}
            onClick={(info) => navigate(info.key)}
            style={{ flex: 1, minWidth: 0 }}
            selectedKeys={[location.pathname]}
            styles={{ itemContent: { paddingTop: '5px' } }}
          />
        </Layout.Header>
        <Layout.Content className="pt-5 pb-[100px] p-[50px]">
          <Outlet />
        </Layout.Content>
      </Layout>
    </>
  );
};

function MyApp() {
  const location = useLocation();
  const navigate = useNavigate();
  const init = useAppStore((s) => s.init);

  EventsOn('navigate', (route: string) => {
    navigate(route);
  });

  return (
    <ConfigProvider
      theme={{
        token: {
          ...(primaryColorMap[location.pathname] ? { colorPrimary: primaryColorMap[location.pathname] } : undefined),
          fontSize: 14,
        },
        components: {
          Layout: { headerBg: undefined },
        },
      }}
      locale={globalBuddhistLocale}
    >
      <ErrorBoundary
        FallbackComponent={ErrorFallback}
        onReset={() => {
          init();
          navigate('/');
        }}
      >
        <App>
          <Routes>
            <Route element={<AppLayout />}>
              <Route index element={<CreateReceipt />} />
              <Route path="/createProcurement" element={<CreateProcurement />} />
              <Route path="/automove" element={<Automove />} />
              <Route path="/setting" element={<Setting />} />
            </Route>
          </Routes>
        </App>
      </ErrorBoundary>
    </ConfigProvider>
  );
}

export default MyApp;

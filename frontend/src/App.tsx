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

const AppLayout = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const appStore = useAppStore();

  useEffect(() => {
    appStore.init();
  }, []);

  return (
    <>
      <Layout className="w-[100vw] h-[100vh] bg-white text-[18px] ">
        <Layout.Header className="px-5">
          <Menu
            mode="horizontal"
            defaultSelectedKeys={['createReceipt']}
            items={[
              { key: '/', label: 'สร้างใบเสร็จรับเงิน' },
              { key: '/setting', label: 'ตั้งค่า' },
            ]}
            onClick={(info) => navigate(info.key)}
            style={{ flex: 1, minWidth: 0 }}
            selectedKeys={[location.pathname]}
          />
        </Layout.Header>
        <Layout.Content className="pt-5">
          <Outlet />
        </Layout.Content>
      </Layout>
    </>
  );
};

function MyApp() {
  const navigate = useNavigate();
  const location = useLocation();

  EventsOn('navigate', (route: string) => {
    navigate(route);
  });

  return (
    <ConfigProvider
      theme={{
        token: {
          ...(location.pathname === '/' ? { colorPrimary: '#00b96b' } : undefined),
          fontSize: 18,
        },
        components: {
          Layout: { headerBg: undefined },
        },
      }}
      locale={globalBuddhistLocale}
    >
      <App>
        <Routes>
          <Route element={<AppLayout />}>
            <Route index element={<CreateReceipt />} />
            <Route path="/setting" element={<Setting />} />
          </Route>
        </Routes>
      </App>
    </ConfigProvider>
  );
}

export default MyApp;

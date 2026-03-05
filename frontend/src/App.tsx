import { Routes, Route, Outlet, useNavigate, useLocation } from 'react-router';
import CreateReceipt from './pages/CreateReceipt';
import Setting from './pages/Setting';
import { ConfigProvider, Layout, Menu, App, Button } from 'antd';
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
import { DatabaseOutlined, FileExcelOutlined, HomeOutlined, SettingOutlined, ThunderboltOutlined } from '@ant-design/icons';
import Home from './pages/Home';
import { primaryColorMap } from './constants';
import Success from './pages/Success';
import DBSetting from './pages/DBSetting';
import useApp from 'antd/es/app/useApp';

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
  const init = useAppStore((s) => s.init);
  const { showBoundary } = useShowBoundary();
  const isRevertedToDefaultCustomerDB = useAppStore((s) => s.isRevertedToDefaultCustomerDB);
  const { notification } = useApp();

  useEffect(() => {
    if (isRevertedToDefaultCustomerDB) {
      notification.warning({
        placement: 'bottomRight',
        pauseOnHover: true,
        duration: 10,
        actions: [
          <Button
            onClick={() => {
              notification.destroy();
              navigate('/DBSetting');
            }}
          >
            ไปที่ตั้งค่า Database
          </Button>,
        ],
        description:
          'คุณกำลังใช้ฐานข้อมูลลูกค้าแบบค่าตั้งต้น (DEFAULT) เนื่องจากไม่สามารถเชื่อมต่อกับฐานข้อมูลลูกค้าที่คุณเลือกได้',
      });
    }
  }, [isRevertedToDefaultCustomerDB]);

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
              { key: '/', itemIcon: <HomeOutlined /> },
              {
                key: '/createReceipt',
                icon: <FileExcelOutlined style={{ color: primaryColorMap['/createReceipt'] }} />,
                label: 'สร้างใบเสร็จรับเงิน',
              },
              {
                key: '/createProcurement',
                icon: <FileExcelOutlined style={{ color: primaryColorMap['/createProcurement'] }} />,
                label: 'สร้างจัดซื้อจัดจ้าง',
              },
              {
                key: '/automove',
                icon: <ThunderboltOutlined style={{ color: primaryColorMap['/automove'] }} />,
                label: 'Auto Move',
              },
              { key: '/setting', icon: <SettingOutlined style={{ color: primaryColorMap['/setting'] }} />, label: 'ตั้งค่า' },
              {
                key: '/DBSetting',
                icon: <DatabaseOutlined style={{ color: primaryColorMap['/DBSetting'] }} />,
                label: 'Database',
              },
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
              <Route index element={<Home />} />
              <Route path="/createReceipt" element={<CreateReceipt />} />
              <Route path="/createProcurement" element={<CreateProcurement />} />
              <Route path="/automove" element={<Automove />} />
              <Route path="/setting" element={<Setting />} />
              <Route path="/DBSetting" element={<DBSetting />} />
              <Route path="/success" element={<Success />} />
            </Route>
          </Routes>
        </App>
      </ErrorBoundary>
    </ConfigProvider>
  );
}

export default MyApp;

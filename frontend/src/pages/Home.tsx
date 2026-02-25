import { Card } from 'antd';
import { useNavigate } from 'react-router';
import { cn } from '../utils';
import { FileExcelOutlined, ThunderboltOutlined } from '@ant-design/icons';
import { primaryColorMap } from '../constants';

const items = [
  {
    key: '/createReceipt',
    label: 'สร้างใบเสร็จรับเงิน',
    icon: <FileExcelOutlined style={{ color: primaryColorMap['/createReceipt'] }} />,
  },
  {
    key: '/createProcurement',
    label: 'สร้างจัดซื้อจัดจ้าง',
    icon: <FileExcelOutlined style={{ color: primaryColorMap['/createProcurement'] }} />,
  },
  { key: '/automove', label: 'Auto Move', icon: <ThunderboltOutlined style={{ color: primaryColorMap['/automove'] }} /> },
];

export default function Home() {
  const navigate = useNavigate();
  return (
    <div className="flex flex-col items-center gap-6 mt-5">
      {items.map((item) => (
        <Card
          className={cn(
            'w-[400px] h-[120px] text-xl font-bold cursor-pointer',
            'flex items-center justify-center',
            `shadow-md hover:shadow-xl hover:scale-[1.06] transition-all duration-300`
          )}
          onClick={() => navigate(item.key)}
        >
          <span className="text-center break-words">
            <span className="mr-2">{item.icon}</span>
            {item.label}
          </span>
        </Card>
      ))}
    </div>
  );
}

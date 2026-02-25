import { Card } from 'antd';
import { useNavigate } from 'react-router';
import { cn } from '../utils';
import { primaryColorMap } from '../constants';

const items = [
  { key: '/createReceipt', label: 'สร้างใบเสร็จรับเงิน' },
  { key: '/createProcurement', label: 'สร้างจัดซื้อจัดจ้าง' },
  { key: '/automove', label: 'Auto Move' },
];

export default function Home() {
  const navigate = useNavigate();
  return (
    <div className="flex flex-row gap-10 mt-5">
      {items.map((item) => (
        <Card
          className={cn(
            'w-[180px] h-[110px] text-sm font-bold cursor-pointer',
            'flex items-center justify-center',
            `shadow-md hover:shadow-xl hover:scale-[1.20] transition-all duration-300`
          )}
          onClick={() => navigate(item.key)}
        >
          <span>{item.label}</span>
        </Card>
      ))}
    </div>
  );
}

import { Result, Button } from 'antd';
import { useNavigate, useSearchParams } from 'react-router';
import { CMDOpenFile } from '../../wailsjs/go/main/App';
import { useMemo } from 'react';

export default function Success() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const filePath = useMemo(() => searchParams.get('file'), [searchParams]);
  return (
    <Result
      status="success"
      title="สร้างเอกสารสำเร็จ!"
      subTitle={filePath ? `ไฟล์ถูกบันทึกที่: ${filePath}` : undefined}
      extra={[
        <Button type="primary" onClick={() => navigate('/')}>
          ไปหน้าหลัก
        </Button>,
        <Button
          onClick={() => {
            if (filePath) {
              CMDOpenFile(filePath);
            }
          }}
        >
          เปิดไฟล์
        </Button>,
      ]}
    />
  );
}

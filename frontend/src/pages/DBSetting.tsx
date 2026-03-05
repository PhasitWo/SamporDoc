import { Alert, App, AutoComplete, Button, Select } from 'antd';
import InputContainer from '../components/InputContainer';
import { useEffect, useState } from 'react';
import { OpenDBFileDialog, ResetupApp, SaveSetting } from '../../wailsjs/go/main/App';
import { useShowBoundary } from '../utils';
import { useAppStore } from '../store/useAppStore';
import FormContainer from '../components/FormContainer';

export default function DBSetting() {
  const [filePath, setFilePath] = useState('');
  const { showBoundary } = useShowBoundary();
  const { message } = App.useApp();
  const init = useAppStore((s) => s.init);
  const isRevertedToDefaultCustomerDB = useAppStore((s) => s.isRevertedToDefaultCustomerDB);
  const customerDBPath = useAppStore((s) => s.customerDBPath);

  useEffect(() => {
    setFilePath(customerDBPath);
  }, [customerDBPath, isRevertedToDefaultCustomerDB]);

  const handleSubmit = async () => {
    if (filePath.trim() === '') {
      return;
    }
    try {
      await SaveSetting({ CustomerDBPath: filePath });
      // reset app state
      await ResetupApp();
      await init();
      message.success('บันทึกสำเร็จ', 3);
    } catch (error) {
      showBoundary(error);
    }
  };

  return (
    <FormContainer>
      <InputContainer>
        <label>ไฟล์ฐานข้อมูลลูกค้า</label>
        <div className="flex gap-1">
          <Select
            className="w-full"
            allowClear
            options={[{ value: '{{DEFAULT}}', label: 'DEFAULT' }]}
            onChange={(value) => setFilePath(value)}
            value={filePath}
          />
          <Button
            type="default"
            onClick={async () => {
              const path = await OpenDBFileDialog();
              if (path !== '') {
                setFilePath(path);
              }
            }}
          >
            เลือก
          </Button>
        </div>
      </InputContainer>
      <Button
        className="mt-3 w-full"
        type="primary"
        onClick={handleSubmit}
        disabled={filePath.trim() === '' || filePath === customerDBPath}
      >
        บันทึก
      </Button>
      {isRevertedToDefaultCustomerDB && (
        <Alert
          className="mt-3"
          type="warning"
          showIcon
          description={'คุณกำลังใช้ฐานข้อมูลลูกค้าแบบค่าตั้งต้น (DEFAULT) เนื่องจากไม่สามารถเชื่อมต่อกับฐานข้อมูลลูกค้าที่คุณเลือกได้'}
        />
      )}
    </FormContainer>
  );
}

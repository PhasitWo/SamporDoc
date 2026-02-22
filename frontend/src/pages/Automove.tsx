import { Button, Input, message } from 'antd';
import FormContainer from '../components/FormContainer';
import InputContainer from '../components/InputContainer';
import { excel } from '../../wailsjs/go/models';
import { useRef, useState } from 'react';
import { AutoMoveBookOrder, GetBookOrderFromDataSourceFile, OpenExcelFileDialog } from '../../wailsjs/go/main/App';
import { getFileName } from '../utils';
import BookOrderTable from '../components/BookOrderTable';
import Asterisk from '../components/Asterisk';
import { WindowReload } from '../../wailsjs/runtime/runtime';

export default function Automove() {
  const [procurementFilePath, setProcurementFilePath] = useState<string>('');
  const [bookOrderData, setBookOrderData] = useState<{ filePath: string; data: excel.PublisherItem[] } | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const btnRef = useRef<HTMLButtonElement>(null);

  const handleChooseFile = async () => {
    const path = await OpenExcelFileDialog();
    if (path) {
      setProcurementFilePath(path);
    }
  };

  const handleImportBookOrder = async () => {
    const filePath = await OpenExcelFileDialog();
    if (filePath) {
      message.loading('กำลังนำเข้าไฟล์...', 2);
      try {
        const bookOrder = await GetBookOrderFromDataSourceFile(filePath);
        setTimeout(() => {
          message.destroy();
          setBookOrderData({ filePath: filePath, data: bookOrder });
          message.success('นำเข้าไฟล์สำเร็จ!', 3);
          setTimeout(() => {
            btnRef.current?.scrollIntoView({ behavior: 'smooth' });
          }, 100);
        }, 500);
      } catch (err) {
        message.error(`นำเข้าไฟล์ไม่สำเร็จ :${(err as Error).message}`, 10);
      }
    }
  };

  const handleSubmit = async () => {
    if (!procurementFilePath || !bookOrderData) {
      return;
    }
    try {
      setIsLoading(true);
      message.loading('กำลังเขียนไฟล์...', 2);
      await AutoMoveBookOrder(procurementFilePath, bookOrderData.filePath);
      setTimeout(() => {
        message.destroy();
        message.success('เขียนไฟล์สำเร็จ!', 3, WindowReload);
        setIsLoading(false);
      }, 500);
    } catch (err) {
      setIsLoading(false);
      message.error(`เขียนไฟล์ไม่สำเร็จ :${(err as Error).message}`, 10);
    }
  };

  return (
    <FormContainer>
      <InputContainer>
        <label>
          ไฟล์จัดซื้อจัดจ้าง
          <Asterisk />
        </label>
        <div className="flex gap-1">
          <Input readOnly onClick={handleChooseFile} value={procurementFilePath ?? ''} />
          <Button type="default" onClick={handleChooseFile}>
            เลือก
          </Button>
        </div>
      </InputContainer>
      <InputContainer>
        <div className="flex gap-3 items-baseline truncate mt-1">
          <Button
            className="w-fit mb-3"
            danger={Boolean(bookOrderData)}
            onClick={bookOrderData ? () => setBookOrderData(null) : handleImportBookOrder}
          >
            {bookOrderData ? 'ยกเลิก' : 'Import ไฟล์แยกสำนัก'}
          </Button>
          <span>{bookOrderData ? getFileName(bookOrderData.filePath) : ''}</span>
        </div>

        {bookOrderData && <BookOrderTable data={bookOrderData.data} />}
      </InputContainer>
      <InputContainer>
        <Button
          ref={btnRef}
          className="w-full"
          type="primary"
          disabled={!procurementFilePath || !bookOrderData || isLoading}
          onClick={handleSubmit}
        >
          Auto Move!
        </Button>
      </InputContainer>
    </FormContainer>
  );
}

import { Input, Button, AutoComplete, Select, DatePicker, Divider, InputNumber } from 'antd';
import ErrorAlertCard from '../components/ErrorAlertCard';
import { useNavigate } from 'react-router';
import InputContainer from '../components/InputContainer';
import Asterisk from '../components/Asterisk';
import { cn, getFileName, isValidWindowsFilename } from '../utils';
import FormContainer from '../components/FormContainer';
import HiddenDatePicker from '../components/HiddenDatePicker';
import BookOrderTable from '../components/BookOrderTable';
import {
  CustomerOptionType,
  ProcurementOutputTypeOptionType,
  ShopOptionType,
  useCreateProcurement,
} from '../hooks/useCreateProcurement';
import SyncOutlined from '@ant-design/icons/lib/icons/SyncOutlined';

export default function CreateProcurementPage() {
  const {
    data,
    setData,
    bookOrderData,
    setBookOrderData,
    isLoading,
    selectedShop,
    shopOptions,
    handleShopChange,
    selectedCustomer,
    customerOptions,
    handleCustomerChange,
    handleChooseDir,
    handleImportBookOrder,
    handleSubmit,
    readyToCreate,
    hiddenPickerRef,
    handleLoadNextNumber,
  } = useCreateProcurement();
  const navigate = useNavigate();

  return (
    <FormContainer>
      <HiddenDatePicker
        ref={hiddenPickerRef}
        value={data.deliveryDate}
        onChange={(date) => setData({ ...data, deliveryDate: date })}
      />
      <InputContainer>
        <label>
          ชื่อไฟล์
          <Asterisk />
        </label>
        <Input
          value={data.filename}
          onChange={(e) => {
            if (isValidWindowsFilename(e.target.value)) {
              setData({ ...data, filename: e.target.value });
            }
          }}
        />
      </InputContainer>
      <InputContainer>
        <label>
          บันทึกที่
          <Asterisk />
        </label>
        <div className="flex gap-1">
          <Input readOnly value={data.saveDir} onClick={handleChooseDir} />
          <Button type="default" onClick={handleChooseDir}>
            เลือก
          </Button>
        </div>
      </InputContainer>
      <InputContainer>
        <label>
          ร้าน
          <Asterisk />
        </label>
        <Select<string | undefined, ShopOptionType>
          allowClear
          options={shopOptions}
          onChange={handleShopChange}
          value={selectedShop?.slug}
        />
        <ErrorAlertCard
          messages={[
            selectedShop && !selectedShop.procurementFormPath && 'ขาดไฟล์ต้นแบบจัดซื้อจัดจ้าง',
            selectedShop && !selectedShop.procurementControlPath && 'ขาดไฟล์สมุดคุมใบส่งของ',
          ]}
          action={
            <Button danger ghost onClick={() => navigate(`/setting?shopSlug=${selectedShop?.slug}`)}>
              ไปที่ตั้งค่า
            </Button>
          }
        />
      </InputContainer>
      <div className="flex flex-row w-full gap-2">
        <InputContainer>
          <label>เลขที่ใบส่งของ</label>
          <div className="flex items-center relative">
            <SyncOutlined
              className={cn('absolute right-3 z-10 hover:scale-[1.2]', !selectedShop?.procurementControlPath && 'text-gray-200')}
              onClick={handleLoadNextNumber}
            />
            <InputNumber<string | number>
              value={data.deliveryNO}
              min='1'
              controls={false}
              onChange={(value) => setData({ ...data, deliveryNO: String(value) })}
              className="w-full"
            />
          </div>
        </InputContainer>
        <InputContainer>
          <label>ใบส่งของลงวันที่</label>
          <DatePicker
            value={data.deliveryDate}
            readOnly
            onClick={() => hiddenPickerRef.current?.nativeElement.click()}
            popupClassName="hidden"
          />
        </InputContainer>
      </div>
      <InputContainer>
        <label>
          ซื้อ
          <Asterisk />
        </label>
        <AutoComplete<string>
          allowClear
          options={[
            { value: 'วัสดุสำนักงาน' },
            { value: 'วัสดุการศึกษา' },
            { value: 'หนังสือเรียน' },
            { value: 'อื่นๆ โปรดระบุ', disabled: true },
          ]}
          onChange={(v) => setData({ ...data, buy: v })}
        />
      </InputContainer>
      <InputContainer>
        <label>โครงการ</label>
        <Input value={data.project} onChange={(e) => setData({ ...data, project: e.target.value })} />
      </InputContainer>
      <InputContainer>
        <label>จำนวนเงิน</label>
        <Input
          type="number"
          value={data.amount}
          onChange={(e) =>
            setData({
              ...data,
              amount: isNaN(e.target.valueAsNumber) ? 0 : e.target.valueAsNumber,
            })
          }
          min={0}
        />
      </InputContainer>
      <Divider />
      <InputContainer>
        <label>รูปแบบ</label>
        <Select<string, ProcurementOutputTypeOptionType>
          options={[
            { value: 'FULL', label: 'ฉบับเต็ม' },
            { value: 'ONLY_DELIVERY_NOTE', label: 'เฉพาะใบส่งของ' },
            { value: 'ONLY_QUOTATION', label: 'เฉพาะใบเสนอราคา' },
          ]}
          value={data.procurementOutputType}
          onChange={(_, option) => {
            if (option && !Array.isArray(option)) {
              setData({ ...data, procurementOutputType: option.value });
            }
          }}
        />
      </InputContainer>
      <Divider />
      <InputContainer>
        <div className="flex gap-3 items-baseline truncate">
          <Button
            className="w-fit mb-3"
            danger={Boolean(bookOrderData)}
            onClick={bookOrderData ? () => setBookOrderData(null) : handleImportBookOrder}
          >
            {bookOrderData ? 'ยกเลิก' : 'Import ไฟล์แยกสำนัก'}
          </Button>
          {bookOrderData && <span>{getFileName(bookOrderData.filePath)}</span>}
        </div>

        {bookOrderData && <BookOrderTable data={bookOrderData.data} />}
      </InputContainer>
      <Divider />
      <InputContainer>
        <label>
          ชื่อลูกค้า
          <Asterisk />
        </label>
        <AutoComplete<string, CustomerOptionType>
          value={data.customerName}
          options={customerOptions}
          onChange={handleCustomerChange}
          allowClear
          showSearch={{
            filterOption: (inputValue, option) => option?.meta.name.toLowerCase().includes(inputValue.toLowerCase()) ?? false,
          }}
        />
      </InputContainer>
      <InputContainer>
        <label>ที่อยู่</label>
        <Input
          value={data.address}
          onChange={(e) => setData({ ...data, address: e.target.value })}
          disabled={Boolean(selectedCustomer)}
        />
      </InputContainer>
      <InputContainer>
        <label>ประธานกรรมการ</label>
        <Input value={data.headCheckerName} onChange={(e) => setData({ ...data, headCheckerName: e.target.value })} />
      </InputContainer>
      <InputContainer>
        <label>กรรมการ1</label>
        <Input value={data.checker1Name} onChange={(e) => setData({ ...data, checker1Name: e.target.value })} />
      </InputContainer>
      <InputContainer>
        <label>กรรมการ2</label>
        <Input value={data.checker2Name} onChange={(e) => setData({ ...data, checker2Name: e.target.value })} />
      </InputContainer>
      <InputContainer>
        <label>เจ้าหน้าที่พัสดุ</label>
        <Input value={data.objectName} onChange={(e) => setData({ ...data, objectName: e.target.value })} />
      </InputContainer>
      <InputContainer>
        <label>หัวหน้าเจ้าหน้าที่</label>
        <Input value={data.headObjectName} onChange={(e) => setData({ ...data, headObjectName: e.target.value })} />
      </InputContainer>
      <InputContainer>
        <label>ผู้อำนวยการ</label>
        <Input value={data.bossName} onChange={(e) => setData({ ...data, bossName: e.target.value })} />
      </InputContainer>
      <Button className="mt-3 w-full" type="primary" disabled={!readyToCreate || isLoading} onClick={handleSubmit}>
        สร้างจัดซื้อจัดจ้าง
      </Button>
    </FormContainer>
  );
}

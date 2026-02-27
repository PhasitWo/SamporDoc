import { Input, Button, AutoComplete, Select, DatePicker, InputNumber } from 'antd';
import ErrorAlertCard from '../components/ErrorAlertCard';
import { useNavigate } from 'react-router';
import InputContainer from '../components/InputContainer';
import Asterisk from '../components/Asterisk';
import { cn, isValidWindowsFilename } from '../utils';
import FormContainer from '../components/FormContainer';
import HiddenDatePicker from '../components/HiddenDatePicker';
import { CustomerOptionType, ShopOptionType, useCreateReceipt } from '../hooks/useCreateReceipt';
import SyncOutlined from '@ant-design/icons/lib/icons/SyncOutlined';

export default function CreateReceiptPage() {
  const navigate = useNavigate();
  const {
    data,
    setData,
    isLoading,
    receiptType,
    receiptFormPath,
    receiptControlPath,
    setReceiptType,
    selectedShop,
    shopOptions,
    handleShopChange,
    selectedCustomer,
    customerOptions,
    handleCustomerChange,
    handleChooseDir,
    handleSubmit,
    readyToCreate,
    hiddenReceiptDateRef,
    hiddenDeliveryDateRef,
    handleLoadNextNumber,
  } = useCreateReceipt();

  return (
    <FormContainer>
      <HiddenDatePicker
        ref={hiddenReceiptDateRef}
        value={data.receiptDate}
        onChange={(date) => setData({ ...data, receiptDate: date })}
      />
      <HiddenDatePicker
        ref={hiddenDeliveryDateRef}
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
            selectedShop &&
              receiptFormPath === null &&
              `ขาดไฟล์ต้นแบบใบเสร็จรับเงิน ${receiptType === 'MAIN' ? '(เล่มหลัก)' : '(เล่มรอง)'} `,
            selectedShop &&
              receiptControlPath === null &&
              `ขาดไฟล์สมุดคุมใบเสร็จรับเงิน ${receiptType === 'MAIN' ? '(เล่มหลัก)' : '(เล่มรอง)'} `,
          ]}
          action={
            <Button danger ghost onClick={() => navigate(`/setting?shopSlug=${selectedShop?.slug}`)}>
              ไปที่ตั้งค่า
            </Button>
          }
          noTitle
        />
      </InputContainer>
      <div className="flex flex-row w-full gap-2">
        <InputContainer>
          <label>เล่มใบเสร็จ</label>
          <Select
            options={[
              { label: 'เล่มหลัก', value: 'MAIN' },
              { label: 'เล่มรอง', value: 'SEC' },
            ]}
            onChange={(value) => setReceiptType(value as 'MAIN' | 'SEC')}
            value={receiptType}
          />
        </InputContainer>
        <InputContainer>
          <label>เลขที่ใบเสร็จ</label>
          <div className="flex items-center relative">
            <SyncOutlined
              className={cn('absolute right-3 z-10 hover:scale-[1.2]', !selectedShop?.procurementControlPath && 'text-gray-200')}
              onClick={handleLoadNextNumber}
            />
            <InputNumber<string | number>
              value={data.receiptNO}
              min="1"
              controls={false}
              onChange={(value) => setData({ ...data, receiptNO: String(value) })}
              className="w-full"
            />
          </div>
        </InputContainer>
        <InputContainer>
          <label>ใบเสร็จลงวันที่</label>
          <DatePicker
            value={data.receiptDate}
            readOnly
            onClick={() => hiddenReceiptDateRef.current?.nativeElement.click()}
            popupClassName="hidden"
          />
        </InputContainer>
      </div>
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
        <label>
          รายละเอียดใบเสร็จ
          <Asterisk />
        </label>
        <AutoComplete<string>
          allowClear
          options={[
            { value: 'ค่าวัสดุสำนักงาน' },
            { value: 'ค่าวัสดุการศึกษา' },
            { value: 'ค่าหนังสือเรียน' },
            { value: 'อื่นๆ โปรดระบุ', disabled: true },
          ]}
          onChange={(v) => setData({ ...data, detail: v })}
        />
      </InputContainer>

      <div className="flex flex-row w-full gap-2">
        <InputContainer>
          <label>อ้างใบส่งของเลขที่</label>
          <Input onChange={(e) => setData({ ...data, deliveryNO: e.target.value })} />
        </InputContainer>
        <InputContainer>
          <label>ใบส่งของลงวันที่</label>
          <DatePicker
            value={data.deliveryDate}
            readOnly
            onClick={() => hiddenDeliveryDateRef.current?.nativeElement.click()}
            popupClassName="hidden"
          />
        </InputContainer>
      </div>
      <InputContainer>
        <label>
          จำนวนเงิน
          <Asterisk />
        </label>
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
      <Button className="mt-3 w-full" type="primary" disabled={!readyToCreate || isLoading} onClick={handleSubmit}>
        สร้างใบเสร็จ
      </Button>
    </FormContainer>
  );
}

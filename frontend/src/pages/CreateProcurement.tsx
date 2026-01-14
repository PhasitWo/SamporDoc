import { Input, Button, AutoComplete, Select, DatePicker, App, Divider, Radio } from 'antd';
import { useEffect, useMemo, useState } from 'react';
import { main, model } from '../../wailsjs/go/models';
import { GetNextControlNumber, OpenDirectoryDialog, CreateReceipt } from '../../wailsjs/go/main/App';
import type { DefaultOptionType } from 'antd/es/select';
import { Dayjs } from 'dayjs';
import ErrorAlertCard from '../components/ErrorAlertCard';
import { useNavigate } from 'react-router';
import InputContainer from '../components/InputContainer';
import { useAppStore } from '../store/useAppStore';
import Asterisk from '../components/Asterisk';
import { useShowBoundary } from '../utils';
import type { CheckboxGroupProps } from 'antd/es/checkbox';

interface ShopOptionType extends DefaultOptionType {
  meta: model.Shop;
}

interface CustomerOptionType extends DefaultOptionType {
  meta: model.Customer;
}

interface FormData {
  filename: string;
  saveDir: string;
  receiptNO: string;
  receiptDate: Dayjs | null;
  customerName: string;
  address: string;
  detail: string;
  deliveryNO: string;
  deliveryDate: Dayjs | null;
  amount: number;
}

export default function CreateProcurementPage() {
  const navigate = useNavigate();
  const { message } = App.useApp();
  const { showBoundary } = useShowBoundary();
  // form
  const [data, setData] = useState<FormData>({
    filename: '',
    saveDir: '',
    receiptNO: '',
    receiptDate: null,
    customerName: '',
    address: '',
    detail: '',
    deliveryNO: '',
    deliveryDate: null,
    amount: 0,
  });

  // shop
  const [selectedShop, setSelectedShop] = useState<model.Shop | null>(null);
  const shops = useAppStore((s) => s.shops);
  const shopOptions = useMemo<ShopOptionType[]>(
    () => shops.map<ShopOptionType>((s) => ({ value: s.slug, label: s.name, meta: s })),
    [shops]
  );

  useEffect(() => {
    (async () => {
      setData({ ...data, receiptNO: '' });
      if (selectedShop && selectedShop.receiptControlPath) {
        const nextNumber = await GetNextControlNumber(selectedShop.receiptControlPath);
        setData({ ...data, receiptNO: String(nextNumber) });
      }
    })();
  }, [selectedShop]);

  // customers
  const [selectedCustomer, setSelectedCustomer] = useState<model.Customer | null>(null);
  const customers = useAppStore((s) => s.customers);
  const customerOptions = useMemo<CustomerOptionType[]>(
    () => customers.map<CustomerOptionType>((c) => ({ value: c.ID, label: `${c.name} (ID: ${c.ID})`, meta: c })),
    [customers]
  );

  const handleCustomerChange = (value: string, option?: CustomerOptionType | CustomerOptionType[]) => {
    if (option && 'meta' in option && !Array.isArray(option)) {
      setSelectedCustomer(option.meta);
      setData({ ...data, customerName: option.meta.name, address: option.meta.address ?? '' });
    } else {
      setSelectedCustomer(null);
      setData({ ...data, customerName: value, address: '' });
    }
  };

  const readyToCreate = useMemo<boolean>(() => {
    if (selectedShop == null || selectedShop.receiptFormPath == undefined || selectedShop.receiptControlPath == undefined) {
      return false;
    }
    if (data.amount <= 0) {
      return false;
    }
    if (data.filename.trim() === '' || data.saveDir === '' || data.receiptNO.trim() === '' || data.customerName.trim() === '') {
      return false;
    }
    return true;
  }, [data, selectedShop]);

  const handleShopChange = (_: any, option?: ShopOptionType | ShopOptionType[]) => {
    if (option && !Array.isArray(option)) {
      setSelectedShop(option.meta);
    } else {
      setSelectedShop(null);
    }
  };

  const handleSubmit = async () => {
    try {
      if (
        !selectedShop ||
        !selectedShop.receiptFormPath ||
        !selectedShop.receiptControlPath ||
        data.filename.trim() === '' ||
        data.saveDir === '' ||
        data.receiptNO === '' ||
        data.amount <= 0 ||
        data.customerName.trim() === ''
      ) {
        return;
      }
      message.loading('สร้างไฟล์ใบเสร็จรับเงิน...');
      await CreateReceipt({
        TemplatePath: selectedShop.receiptFormPath,
        Filename: data.filename.trim(),
        OutputDir: data.saveDir,
        ReceiptNO: data.receiptNO,
        CustomerName: data.customerName.trim(),
        Amount: data.amount,
        ControlPath: selectedShop.receiptControlPath,
        Address: data.address.trim() || undefined,
        DeliveryNO: data.deliveryNO.trim() || undefined,
        Detail: data.detail.trim() || undefined,
        DeliveryDate: data.deliveryDate?.toISOString(),
        ReceiptDate: data.receiptDate?.toISOString(),
      });
      message.destroy();
      message.success('สร้างไฟล์ใบเสร็จรับเงินสำเร็จ!');
      // refetch
      useAppStore.getState().fetchCustomers();
    } catch (err: any) {
      showBoundary(err);
    }
  };

  return (
    <div className="mx-auto flex flex-col gap-3 items-center justify-center max-w-[500px]">
      <InputContainer>
        <label>
          ชื่อไฟล์
          <Asterisk />
        </label>
        <Input onChange={(e) => setData({ ...data, filename: e.target.value })} />
      </InputContainer>
      <InputContainer>
        <label>
          บันทึกที่
          <Asterisk />
        </label>
        <div className="flex gap-1">
          <Input readOnly value={data.saveDir} />
          <Button
            type="default"
            onClick={() =>
              OpenDirectoryDialog().then((path) => {
                if (path) {
                  setData({ ...data, saveDir: path });
                }
              })
            }
          >
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
            selectedShop && !selectedShop.receiptFormPath && 'ขาดไฟล์ต้นแบบ',
            selectedShop && selectedShop.receiptControlPath === null && 'ขาดไฟล์สมุดคุม',
          ]}
          action={
            <Button danger ghost onClick={() => navigate(`/setting?shopSlug=${selectedShop?.slug}`)}>
              ไปที่ตั้งค่า
            </Button>
          }
        />
      </InputContainer>
      <div className="flex flex-row w-[500px] gap-2">
        <InputContainer>
          <label>ใบส่งของเลขที่</label>
          <Input onChange={(e) => setData({ ...data, deliveryNO: e.target.value })} />
        </InputContainer>
        <InputContainer>
          <label>ใบส่งของลงวันที่</label>
          <DatePicker onChange={(date) => setData({ ...data, deliveryDate: date })} />
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
          onChange={(v) => setData({ ...data, detail: v })}
        />
      </InputContainer>
      <InputContainer>
        <label>โครงการ</label>
        <Input />
      </InputContainer>
      <InputContainer>
        <label>
          จำนวนเงิน
          <Asterisk />
        </label>
        <Input
          type="number"
          value={data.amount}
          onChange={(e) => setData({ ...data, amount: isNaN(e.target.valueAsNumber) ? 0 : e.target.valueAsNumber })}
          min={0}
        />
      </InputContainer>
      <Divider />
      <InputContainer>
        <label>จำนวนรายการ</label>
        <Radio.Group
          vertical
          options={[
            { value: 'lte', label: 'น้อยกว่าหรือเท่ากับ 11 รายการ' },
            { value: 'gt', label: 'มากกว่า 11 รายการ' },
            {
              value: 'custom',
              label: <Input placeholder="ระบุจำนวน" type="number" min={0} />,
            },
          ]}
          defaultValue="lte"
        />
      </InputContainer>
      <InputContainer>
        <label>รูปแบบ</label>
        <Select<string, DefaultOptionType>
          options={[
            { value: 'full', label: 'ฉบับเต็ม' },
            { value: 'only_delivery_note', label: 'เฉพาะใบส่งของ' },
            { value: 'only_quotation', label: 'เฉพาะใบเสนอราคา' },
          ]}
        />
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
        <Input />
      </InputContainer>
      <InputContainer>
        <label>กรรมการ1</label>
        <Input />
      </InputContainer>
      <InputContainer>
        <label>กรรมการ2</label>
        <Input />
      </InputContainer>
      <InputContainer>
        <label>เจ้าหน้าที่พัสดุ</label>
        <Input />
      </InputContainer>
      <InputContainer>
        <label>หัวหน้าเจ้าหน้าที่</label>
        <Input />
      </InputContainer>
      <InputContainer>
        <label>ผู้อำนวยการ</label>
        <Input />
      </InputContainer>

      <Button className="mt-3 w-full  mb-5" type="primary" disabled={!readyToCreate} onClick={handleSubmit}>
        สร้างจัดซื้อจัดจ้าง
      </Button>
    </div>
  );
}

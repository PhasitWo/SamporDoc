import { Input, Button, AutoComplete, Select, DatePicker, App } from 'antd';
import { useEffect, useMemo, useRef, useState } from 'react';
import { model } from '../../wailsjs/go/models';
import { GetNextControlNumber, OpenDirectoryDialog, CreateReceipt } from '../../wailsjs/go/main/App';
import type { DefaultOptionType } from 'antd/es/select';
import { Dayjs } from 'dayjs';
import ErrorAlertCard from '../components/ErrorAlertCard';
import { useNavigate } from 'react-router';
import InputContainer from '../components/InputContainer';
import { useAppStore } from '../store/useAppStore';
import Asterisk from '../components/Asterisk';
import { isValidWindowsFilename, useShowBoundary } from '../utils';
import { WindowReload } from '../../wailsjs/runtime/runtime';
import FormContainer from '../components/FormContainer';
import { PickerRef } from 'rc-picker';
import HiddenDatePicker from '../components/HiddenDatePicker';

interface ShopOptionType extends DefaultOptionType {
  meta: model.Shop;
}

interface CustomerOptionType extends DefaultOptionType {
  meta: model.Customer;
  value: number;
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

export default function CreateReceiptPage() {
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
  const [isLoading, setIsLoading] = useState(false);

  // shop
  const [selectedShop, setSelectedShop] = useState<model.Shop | null>(null);
  const shops = useAppStore((s) => s.shops);
  const shopOptions = useMemo<ShopOptionType[]>(
    () =>
      shops.map<ShopOptionType>((s) => ({
        value: s.slug,
        label: s.name,
        meta: s,
      })),
    [shops]
  );

  const [receiptType, setReceiptType] = useState<'MAIN' | 'SEC'>('MAIN');
  const [receiptFormPath, setReceiptFormPath] = useState<string | null>(null);
  const [receiptControlPath, setReceiptControlPath] = useState<string | null>(null);

  useEffect(() => {
    (async () => {
      setData({ ...data, receiptNO: '' });
      if (selectedShop) {
        setReceiptFormPath(
          receiptType === 'MAIN' ? (selectedShop.receiptMainFormPath ?? null) : (selectedShop.receiptSecFormPath ?? null)
        );
        const controlPath =
          receiptType === 'MAIN' ? (selectedShop.receiptMainControlPath ?? null) : (selectedShop.receiptSecControlPath ?? null);
        setReceiptControlPath(controlPath);
        if (controlPath) {
          const nextNumber = await GetNextControlNumber(controlPath);
          setData({ ...data, receiptNO: String(nextNumber) });
        }
      } else {
        setReceiptFormPath(null);
        setReceiptControlPath(null);
      }
    })();
  }, [selectedShop, receiptType]);

  // customers
  const [selectedCustomer, setSelectedCustomer] = useState<model.Customer | null>(null);
  const customers = useAppStore((s) => s.customers);
  const customerOptions = useMemo<CustomerOptionType[]>(
    () =>
      customers.map<CustomerOptionType>((c) => ({
        value: c.ID,
        label: `${c.name} (ID: ${c.ID})`,
        meta: c,
      })),
    [customers]
  );

  const handleCustomerChange = (value: string, option?: CustomerOptionType | CustomerOptionType[]) => {
    if (option && 'meta' in option && !Array.isArray(option)) {
      setSelectedCustomer(option.meta);
      setData({
        ...data,
        customerName: option.meta.name,
        address: option.meta.address ?? '',
      });
    } else {
      setSelectedCustomer(null);
      setData({ ...data, customerName: value, address: '' });
    }
  };

  const readyToCreate = useMemo<boolean>(() => {
    if (selectedShop == null || !receiptFormPath || !receiptControlPath) {
      return false;
    }
    if (data.amount <= 0) {
      return false;
    }
    if (
      data.filename.trim() === '' ||
      data.saveDir === '' ||
      data.receiptNO.trim() === '' ||
      data.customerName.trim() === '' ||
      data.detail.trim() === ''
    ) {
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
        !receiptFormPath ||
        !receiptControlPath ||
        data.filename.trim() === '' ||
        data.saveDir === '' ||
        data.receiptNO === '' ||
        data.amount <= 0 ||
        data.customerName.trim() === '' ||
        data.detail.trim() === ''
      ) {
        return;
      }

      setIsLoading(true);
      message.loading('สร้างไฟล์ใบเสร็จรับเงิน...');
      await CreateReceipt({
        TemplatePath: receiptFormPath,
        Filename: data.filename.trim(),
        OutputDir: data.saveDir,
        ReceiptNO: data.receiptNO,
        CustomerName: data.customerName.trim(),
        CustomerID: selectedCustomer?.ID,
        Amount: data.amount,
        ControlPath: receiptControlPath,
        Address: data.address.trim() || undefined,
        DeliveryNO: data.deliveryNO.trim() || undefined,
        Detail: data.detail.trim(),
        DeliveryDate: data.deliveryDate?.toISOString(),
        ReceiptDate: data.receiptDate?.toISOString(),
      });
      message.destroy();
      // refetch
      useAppStore.getState().fetchCustomers();
      message.success('สร้างไฟล์ใบเสร็จรับเงินสำเร็จ!', 3, WindowReload);
    } catch (err: any) {
      showBoundary(err);
    } finally {
      setIsLoading(false);
    }
  };

  const hiddenReceiptDateRef = useRef<PickerRef>(null);
  const hiddenDeliveryDateRef = useRef<PickerRef>(null);
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
            selectedShop &&
              !receiptFormPath &&
              `ขาดไฟล์ต้นแบบใบเสร็จรับเงิน ${receiptType === 'MAIN' ? '(เล่มหลัก)' : '(เล่มรอง)'} `,
            selectedShop &&
              !receiptControlPath &&
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
          <Input readOnly value={data.receiptNO} />
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
          options={[{ value: 'ค่าวัสดุสำนักงาน' }, { value: 'ค่าวัสดุการศึกษา' }, { value: 'อื่นๆ โปรดระบุ', disabled: true }]}
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

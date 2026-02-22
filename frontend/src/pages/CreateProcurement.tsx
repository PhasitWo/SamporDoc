import { Input, Button, AutoComplete, Select, DatePicker, App, Divider } from 'antd';
import type { PickerRef } from 'rc-picker';
import { useEffect, useMemo, useRef, useState } from 'react';
import { excel, model } from '../../wailsjs/go/models';
import { WindowReload } from '../../wailsjs/runtime/runtime';
import {
  GetNextControlNumber,
  OpenDirectoryDialog,
  CreateProcurement,
  OpenExcelFileDialog,
  GetBookOrderFromDataSourceFile,
} from '../../wailsjs/go/main/App';
import type { DefaultOptionType } from 'antd/es/select';
import { Dayjs } from 'dayjs';
import ErrorAlertCard from '../components/ErrorAlertCard';
import { useNavigate } from 'react-router';
import InputContainer from '../components/InputContainer';
import { useAppStore } from '../store/useAppStore';
import Asterisk from '../components/Asterisk';
import { getFileName, isValidWindowsFilename, useShowBoundary } from '../utils';
import FormContainer from '../components/FormContainer';
import HiddenDatePicker from '../components/HiddenDatePicker';
import BookOrderTable from '../components/BookOrderTable';

interface ShopOptionType extends DefaultOptionType {
  meta: model.Shop;
}

interface CustomerOptionType extends DefaultOptionType {
  meta: model.Customer;
  value: number;
}

type ProcurementOutputType = 'FULL' | 'ONLY_DELIVERY_NOTE' | 'ONLY_QUOTATION';

interface ProcurementOutputTypeOptionType extends DefaultOptionType {
  value: ProcurementOutputType;
}

interface FormData {
  filename: string;
  saveDir: string;
  deliveryNO: string;
  deliveryDate: Dayjs | null;
  buy: string;
  project: string;
  amount: number;
  procurementOutputType: ProcurementOutputType;
  customerName: string;
  address: string;
  headCheckerName: string;
  checker1Name: string;
  checker2Name: string;
  objectName: string;
  headObjectName: string;
  bossName: string;
}

// type QuantityType = 'LTE' | 'GT' | 'CUSTOM';

export default function CreateProcurementPage() {
  const navigate = useNavigate();
  const { message } = App.useApp();
  const { showBoundary } = useShowBoundary();
  // form
  const [data, setData] = useState<FormData>({
    filename: '',
    saveDir: '',
    deliveryNO: '',
    deliveryDate: null,
    customerName: '',
    address: '',
    buy: '',
    project: '',
    amount: 0,
    procurementOutputType: 'FULL',
    headCheckerName: '',
    checker1Name: '',
    checker2Name: '',
    objectName: '',
    headObjectName: '',
    bossName: '',
  });
  const [isLoading, setIsLoading] = useState(false);
  const [bookOrderData, setBookOrderData] = useState<{ filePath: string; data: excel.PublisherItem[] } | null>(null);
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

  useEffect(() => {
    (async () => {
      setData({ ...data, deliveryNO: '' });
      if (selectedShop && selectedShop.procurementControlPath) {
        const nextNumber = await GetNextControlNumber(selectedShop.procurementControlPath);
        setData({ ...data, deliveryNO: String(nextNumber) });
      }
    })();
  }, [selectedShop]);

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
        headCheckerName: option.meta.headCheckerName ?? '',
        checker1Name: option.meta.checker1Name ?? '',
        checker2Name: option.meta.checker2Name ?? '',
        objectName: option.meta.objectName ?? '',
        headObjectName: option.meta.headObjectName ?? '',
        bossName: option.meta.bossName ?? '',
      });
    } else {
      setSelectedCustomer(null);
      setData({
        ...data,
        customerName: value,
        address: '',
        headCheckerName: '',
        checker1Name: '',
        checker2Name: '',
        objectName: '',
        headObjectName: '',
        bossName: '',
      });
    }
  };

  const readyToCreate = useMemo<boolean>(() => {
    if (selectedShop == null || !selectedShop.procurementFormPath || !selectedShop.procurementControlPath) {
      return false;
    }
    if (data.amount <= 0) {
      return false;
    }
    if (
      data.filename.trim() === '' ||
      data.saveDir === '' ||
      data.deliveryNO.trim() === '' ||
      data.customerName.trim() === '' ||
      data.buy.trim() === ''
    ) {
      return false;
    }
    return true;
  }, [data, selectedShop]);

  const handleChooseDir = async () => {
    const path = await OpenDirectoryDialog();
    if (path) {
      setData({ ...data, saveDir: path });
    }
  };

  const handleShopChange = (_: any, option?: ShopOptionType | ShopOptionType[]) => {
    if (option && !Array.isArray(option)) {
      setSelectedShop(option.meta);
    } else {
      setSelectedShop(null);
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
          setBookOrderData({ filePath, data: bookOrder });
          message.success('นำเข้าไฟล์สำเร็จ!', 3);
        }, 500);
      } catch (err) {
        message.error(`นำเข้าไฟล์ไม่สำเร็จ :${(err as Error).message}`, 10);
      }
    }
  };

  const handleSubmit = async () => {
    try {
      if (
        !selectedShop ||
        !selectedShop.procurementFormPath ||
        !selectedShop.procurementControlPath ||
        data.filename.trim() === '' ||
        data.saveDir === '' ||
        data.deliveryNO === '' ||
        data.amount <= 0 ||
        data.customerName.trim() === '' ||
        data.buy === ''
      ) {
        return;
      }

      setIsLoading(true);
      message.loading('สร้างไฟล์จัดซื้อจัดจ้าง...');
      await CreateProcurement({
        TemplatePath: selectedShop.procurementFormPath,
        ControlPath: selectedShop.procurementControlPath,
        Filename: data.filename.trim(),
        OutputDir: data.saveDir,
        DeliveryNO: data.deliveryNO.trim(),
        DeliveryDate: data.deliveryDate?.toISOString(),
        Buy: data.buy.trim(),
        Project: data.project.trim() || undefined,
        Amount: data.amount,
        BookOrderPath: bookOrderData ? bookOrderData.filePath : undefined,
        ProcurementOutputType: data.procurementOutputType,
        CustomerName: data.customerName.trim(),
        CustomerID: selectedCustomer?.ID,
        Address: data.address.trim() || undefined,
        HeadCheckerName: data.headCheckerName.trim() || undefined,
        Checker1Name: data.checker1Name.trim() || undefined,
        Checker2Name: data.checker2Name.trim() || undefined,
        HeadObjectName: data.headObjectName.trim() || undefined,
        ObjectName: data.objectName.trim() || undefined,
        BossName: data.bossName.trim() || undefined,
      });
      message.destroy();
      useAppStore.getState().fetchCustomers();
      message.success('สร้างไฟล์จัดซื้อจัดจ้างสำเร็จ!', 3, WindowReload);
      // refetch
    } catch (err: any) {
      showBoundary(err);
    } finally {
      setIsLoading(false);
    }
  };

  const hiddenPickerRef = useRef<PickerRef>(null);
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
          <Button type="default" onClick={() => handleChooseDir}>
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
          <Input readOnly value={data.deliveryNO} />
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
      <Divider />
      {/* <InputContainer>
        <label>จำนวนรายการ</label>
        <Radio.Group
          vertical
          options={[
            { value: 'LTE', label: 'ไม่เกิน 11 รายการ' },
            { value: 'GT', label: 'มากกว่า 11 รายการ' },
            {
              value: 'CUSTOM',
              label:
                quantityType === 'CUSTOM' ? (
                  <Input
                    placeholder="ระบุจำนวน"
                    type="number"
                    min={0}
                    value={data.quantity}
                    onChange={(e) =>
                      setData({
                        ...data,
                        quantity: isNaN(e.target.valueAsNumber) ? 0 : e.target.valueAsNumber,
                      })
                    }
                  />
                ) : (
                  'ระบุจำนวน'
                ),
            },
          ]}
          defaultValue="LTE"
          value={quantityType}
          onChange={(e) => {
            setQuantityType(e.target.value);
          }}
        />
      </InputContainer> */}
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

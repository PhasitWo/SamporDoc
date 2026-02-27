import { Button, App, Divider } from 'antd';
import type { PickerRef } from 'rc-picker';
import { useEffect, useMemo, useRef, useState } from 'react';
import { excel, model } from '../../wailsjs/go/models';
import {
  GetNextControlNumber,
  OpenDirectoryDialog,
  CreateProcurement,
  OpenExcelFileDialog,
  GetBookOrderFromDataSourceFile,
  CMDOpenFile,
  GetControlData,
} from '../../wailsjs/go/main/App';
import type { DefaultOptionType } from 'antd/es/select';
import { Dayjs } from 'dayjs';
import { useNavigate } from 'react-router';
import { useAppStore } from '../store/useAppStore';
import { moneyFormat, useShowBoundary } from '../utils';

export interface ShopOptionType extends DefaultOptionType {
  meta: model.Shop;
}

export interface CustomerOptionType extends DefaultOptionType {
  meta: model.Customer;
  value: number;
}

export type ProcurementOutputType = 'FULL' | 'ONLY_DELIVERY_NOTE' | 'ONLY_QUOTATION';

export interface ProcurementOutputTypeOptionType extends DefaultOptionType {
  value: ProcurementOutputType;
}

export interface FormData {
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

export function useCreateProcurement() {
  const navigate = useNavigate();
  const { message, modal } = App.useApp();
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

  const handleLoadNextNumber = async () => {
    if (selectedShop?.procurementControlPath) {
      const nextNumber = await GetNextControlNumber(selectedShop.procurementControlPath);
      setData({ ...data, deliveryNO: String(nextNumber) });
      message.info('Load เลขที่ใบส่งของจากสมุดคุม');
    }
  };

  const readyToCreate = useMemo<boolean>(() => {
    if (selectedShop == null || !selectedShop.procurementFormPath || !selectedShop.procurementControlPath) {
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
        data.customerName.trim() === '' ||
        data.buy === ''
      ) {
        return;
      }
      setIsLoading(true);
      // get control data to check if deliveryNO is duplicated
      let overwriteDeliveryNO: number | undefined = undefined;
      const controlData = await GetControlData(selectedShop.procurementControlPath, Number(data.deliveryNO));
      if (controlData !== null) {
        const ok = await new Promise((res) => {
          const modalInstance = modal.warning({
            title: 'คุณต้องการเขียนทับรายการสมุดคุมนี้หรือไม่?',
            content: (
              <div>
                <div className="font-bold">รายการเดิม</div>
                <div>เลขที่ใบส่งของ: {controlData.NO}</div>
                <div>ลูกค้า: {controlData.CustomerName}</div>
                <div>รายละเอียด: {controlData.Detail}</div>
                <div>จำนวนเงิน: {moneyFormat(controlData.Amount)} บาท</div>
                <Divider />
                <div className="font-bold">รายการใหม่</div>
                <div>เลขที่ใบส่งของ: {data.deliveryNO}</div>
                <div>ลูกค้า: {data.customerName}</div>
                <div>รายละเอียด: {data.buy}</div>
                <div>จำนวนเงิน: {moneyFormat(data.amount)} บาท</div>
              </div>
            ),
            okCancel: true,
            onOk: () => res(true),
            onCancel: () => res(false),
          });
          setTimeout(() => {
            modalInstance.destroy();
            res(false);
          }, 30 * 1000); // auto cancel after 30 seconds
        });
        if (!ok) {
          return;
        }
        overwriteDeliveryNO = Number(data.deliveryNO);
      }
      message.loading('สร้างไฟล์จัดซื้อจัดจ้าง...');
      const outputPath = await CreateProcurement({
        TemplatePath: selectedShop.procurementFormPath,
        ControlPath: selectedShop.procurementControlPath,
        Filename: data.filename.trim(),
        OutputDir: data.saveDir,
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
        OverwriteDeliveryNO: overwriteDeliveryNO,
      });
      message.destroy();
      useAppStore.getState().fetchCustomers();
      navigate(`/success?file=${encodeURIComponent(outputPath)}`);
      // refetch
    } catch (err: any) {
      showBoundary(err);
    } finally {
      setIsLoading(false);
    }
  };

  const hiddenPickerRef = useRef<PickerRef>(null);
  return {
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
  };
}

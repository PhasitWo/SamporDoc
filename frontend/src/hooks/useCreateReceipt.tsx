import { Button, App, Divider } from 'antd';
import { useEffect, useMemo, useRef, useState } from 'react';
import { model } from '../../wailsjs/go/models';
import { GetNextControlNumber, OpenDirectoryDialog, CreateReceipt, CMDOpenFile, GetControlData } from '../../wailsjs/go/main/App';
import type { DefaultOptionType } from 'antd/es/select';
import { Dayjs } from 'dayjs';
import { useNavigate } from 'react-router';
import { useAppStore } from '../store/useAppStore';
import { moneyFormat, useShowBoundary } from '../utils';
import { PickerRef } from 'rc-picker';

export interface ShopOptionType extends DefaultOptionType {
  meta: model.Shop;
}

export interface CustomerOptionType extends DefaultOptionType {
  meta: model.Customer;
  value: number;
}

export interface FormData {
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

export function useCreateReceipt() {
  const navigate = useNavigate();
  const { message, modal } = App.useApp();
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
  const [receiptFormPath, setReceiptFormPath] = useState<string | null | undefined>(undefined);
  const [receiptControlPath, setReceiptControlPath] = useState<string | null | undefined>(undefined);

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
        setReceiptFormPath(undefined);
        setReceiptControlPath(undefined);
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

  const handleLoadNextNumber = async () => {
    if (receiptControlPath) {
      const nextNumber = await GetNextControlNumber(receiptControlPath);
      setData({ ...data, receiptNO: String(nextNumber) });
      message.info('Load เลขที่ใบเสร็จจากสมุดคุม');
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
      // get control data to check if ReceiptNO is duplicated
      let overwriteReceiptNO: number | undefined = undefined;
      const controlData = await GetControlData(receiptControlPath, Number(data.receiptNO));
      if (controlData !== null) {
        const ok = await new Promise((res) => {
          const modalInstance = modal.warning({
            title: 'คุณต้องการเขียนทับรายการสมุดคุมนี้หรือไม่?',
            content: (
              <div>
                <div className="font-bold">รายการเดิม</div>
                <div>เลขที่ใบเสร็จ: {controlData.NO}</div>
                <div>ลูกค้า: {controlData.CustomerName}</div>
                <div>รายละเอียด: {controlData.Detail}</div>
                <div>จำนวนเงิน: {moneyFormat(controlData.Amount)} บาท</div>
                <Divider />
                <div className="font-bold">รายการใหม่</div>
                <div>เลขที่ใบเสร็จ: {data.receiptNO}</div>
                <div>ลูกค้า: {data.customerName}</div>
                <div>รายละเอียด: {data.detail}</div>
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
        overwriteReceiptNO = Number(data.receiptNO);
      }

      message.loading('สร้างไฟล์ใบเสร็จรับเงิน...');
      const outputPath = await CreateReceipt({
        TemplatePath: receiptFormPath,
        Filename: data.filename.trim(),
        OutputDir: data.saveDir,
        CustomerName: data.customerName.trim(),
        CustomerID: selectedCustomer?.ID,
        Amount: data.amount,
        ControlPath: receiptControlPath,
        Address: data.address.trim() || undefined,
        DeliveryNO: data.deliveryNO.trim() || undefined,
        Detail: data.detail.trim(),
        DeliveryDate: data.deliveryDate?.toISOString(),
        ReceiptDate: data.receiptDate?.toISOString(),
        OverwriteReceiptNO: overwriteReceiptNO,
      });
      message.destroy();
      // refetch
      useAppStore.getState().fetchCustomers();
      navigate(`/success?file=${encodeURIComponent(outputPath)}`);

    } catch (err: any) {
      showBoundary(err);
    } finally {
      setIsLoading(false);
    }
  };

  const hiddenReceiptDateRef = useRef<PickerRef>(null);
  const hiddenDeliveryDateRef = useRef<PickerRef>(null);

  return {
    data,
    setData,
    receiptType,
    receiptFormPath,
    setReceiptFormPath,
    receiptControlPath,
    setReceiptControlPath,
    setReceiptType,
    isLoading,
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
  };
}

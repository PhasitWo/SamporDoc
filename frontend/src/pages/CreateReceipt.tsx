import { Input, Button, AutoComplete, Select, Form } from 'antd';
import { ComponentProps, useEffect, useMemo, useState } from 'react';
import { model } from '../../wailsjs/go/models';
import { GetAllShops } from '../../wailsjs/go/main/App';
import type { DefaultOptionType } from 'antd/es/select';

interface FormData {
  filename: string;
  saveDir?: string;
  receiptNO: number;
  receiptDate?: Date;
  customerName: string;
  address?: string;
  detail?: string;
  deliveryNO?: number;
  deliveryDate?: Date;
  amount: number;
}

export function CreateReceipt() {
  const [form] = Form.useForm();
  const [selectedShop, setSelectedShop] = useState<model.Shop | null>(null);
  const [shops, setShops] = useState<model.Shop[]>([]);
  const shopOptions = useMemo<DefaultOptionType[]>(
    () => shops.map<DefaultOptionType>((s) => ({ value: s.Slug, label: s.Name })),
    [shops]
  );

  useEffect(() => {
    (async () => {
      const shops = await GetAllShops();
      setShops(shops);
    })();
  }, []);

  const handleShopChange = (slug: string) => {
    const found = shops.find((v) => v.Slug === slug)!;
    setSelectedShop(found);
  };

  return (
    <Form form={form} className="mx-auto flex flex-col gap-3 items-center justify-center">
      <InputContainer>
        <label>ชื่อไฟล์</label>
        <Input />
      </InputContainer>
      <InputContainer>
        <label>บันทึกที่</label>
        <div className="flex gap-1">
          <Input readOnly />
          <Button type="primary">Browse</Button>
        </div>
      </InputContainer>
      <div className="flex flex-row w-[500px] gap-2">
        <InputContainer>
          <label>ร้าน</label>
          <Select<string> allowClear options={shopOptions} onChange={handleShopChange} value={selectedShop?.Slug} />
        </InputContainer>
        <InputContainer>
          <label>เลขที่ใบเสร็จ</label>
          <Input readOnly />
        </InputContainer>
        <InputContainer>
          <label>ใบเสร็จลงวันที่</label>
          <Input />
        </InputContainer>
      </div>
      <InputContainer>
        <label>ส่วนราชการ</label>
        <Input />
      </InputContainer>
      <InputContainer>
        <label>ที่อยู่</label>
        <Input />
      </InputContainer>
      <InputContainer>
        <label>รายละเอียดใบเสร็จ</label>
        <AutoComplete
          allowClear
          options={[{ value: 'ค่าวัสดุสำนักงาน' }, { value: 'ค่าวัสดุการศึกษา' }, { value: 'อื่นๆ โปรดระบุ', disabled: true }]}
        />
      </InputContainer>

      <div className="flex flex-row w-[500px] gap-2">
        <InputContainer>
          <label>อ้างใบส่งของเลขที่</label>
          <Input />
        </InputContainer>
        <InputContainer>
          <label>ใบส่งของลงวันที่</label>
          <Input />
        </InputContainer>
      </div>
      <InputContainer>
        <label>จำนวนเงิน</label>
        <Input />
      </InputContainer>
    </Form>
  );
}

function InputContainer({ children }: ComponentProps<'div'>) {
  return <div className={'flex flex-col w-[500px]'}>{children}</div>;
}

import { Input, Button, Select, App } from 'antd';
import { model } from '../../wailsjs/go/models';
import { UpdateShopBySlug, OpenExcelFileDialog } from '../../wailsjs/go/main/App';
import InputContainer from '../components/InputContainer';
import { useAppStore } from '../store/useAppStore';
import { DefaultOptionType } from 'antd/es/select';
import { useEffect, useMemo, useRef, useState } from 'react';
import { isDeepEqual } from '../utils';

interface ShopOptionType extends DefaultOptionType {
  meta: model.Shop;
}

export default function Setting() {
  const [selectedShop, setSelectedShop] = useState<model.Shop | null>(null);
  const shops = useAppStore((s) => s.shops);
  const shopOptions = useMemo<ShopOptionType[]>(
    () => shops.map<ShopOptionType>((s) => ({ value: s.slug, label: s.name, meta: s })),
    [shops]
  );

  const handleShopChange = (_: any, option?: ShopOptionType | ShopOptionType[]) => {
    if (option && !Array.isArray(option)) {
      setSelectedShop(option.meta);
    } else {
      setSelectedShop(null);
    }
  };

  return (
    <div className="mx-auto justify-center max-w-[500px]">
      <Select<string | undefined, ShopOptionType>
        allowClear
        options={shopOptions}
        onChange={handleShopChange}
        value={selectedShop?.slug}
        className="min-w-[50%] self-start"
        placeholder="เลือกร้าน"
      />
      {selectedShop && <SingleShopSetting data={selectedShop} />}
    </div>
  );
}

function SingleShopSetting({ data }: { data: model.Shop }) {
  const [shop, setShop] = useState<model.Shop>(data);
  const shopRef = useRef<model.Shop>();
  const [isDirty, setIsDirty] = useState(false);
  const { message } = App.useApp();

  useEffect(() => {
    setShop(data);
    shopRef.current = data;
  }, [data]);

  useEffect(() => {
    setIsDirty(!isDeepEqual(shop, shopRef.current));
  }, [shop, shopRef]);

  const handleSave = async () => {
    message.loading('บันทึกการตั้งค่า...');
    const result = await UpdateShopBySlug(shop);
    setShop(result);
    shopRef.current = result;
    message.destroy();
    message.success('บันทึกสำเร็จ', 3);
    useAppStore.getState().fetchShops();
  };

  return (
    <div className="w-[500px] mt-5 flex flex-col gap-3">
      <div className="font-bold text-2xl mb-3">{shop.name}</div>
      <InputContainer>
        <label>ไฟล์ต้นแบบใบเสร็จรับเงิน</label>
        <div className="flex gap-1">
          <Input readOnly value={shop.receiptFormPath ?? ''} />
          <Button
            type="default"
            onClick={async () => {
              const path = await OpenExcelFileDialog();
              if (path !== '') {
                setShop({ ...shop, receiptFormPath: path });
              }
            }}
          >
            เลือก
          </Button>
        </div>
      </InputContainer>
      <InputContainer>
        <label>ไฟล์สมุดคุมใบเสร็จรับเงิน</label>
        <div className="flex gap-1">
          <Input readOnly value={shop.receiptControlPath ?? ''} />
          <Button
            type="default"
            onClick={async () => {
              OpenExcelFileDialog().then((path) => {
                if (path !== '') {
                  setShop({ ...shop, receiptControlPath: path });
                }
              });
            }}
          >
            เลือก
          </Button>
        </div>
      </InputContainer>
      <Button className="mt-3 w-full" type="primary" onClick={handleSave} disabled={!isDirty}>
        บันทึก
      </Button>
    </div>
  );
}

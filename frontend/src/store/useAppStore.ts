import { create } from 'zustand';
import { model, setting } from '../../wailsjs/go/models';
import { GetAllCustomers, GetAllShops, GetUseRemoteCustomerDB, GetSetting } from '../../wailsjs/go/main/App';
import useApp from 'antd/es/app/useApp';

interface State {
  init: () => Promise<void>;
  useRemoteCustomerDB: boolean;
  isRevertedToDefaultCustomerDB: boolean;
  customerDBPath: string;
  fetchCustomerDBState: () => Promise<void>;
  shops: model.Shop[];
  fetchShops: () => Promise<void>;
  customers: model.Customer[];
  fetchCustomers: () => Promise<void>;
}

export const useAppStore = create<State>((set, get) => ({
  init: async () => {
    get().fetchShops();
    get().fetchCustomers();
    get().fetchCustomerDBState();
  },
  useRemoteCustomerDB: false,
  isRevertedToDefaultCustomerDB: false,
  customerDBPath: '{{DEFAULT}}',
  fetchCustomerDBState: async () => {
    const useRemoteCustomerDB = await GetUseRemoteCustomerDB();
    const { CustomerDBPath: customerDBPath } = await GetSetting();
    set({
      useRemoteCustomerDB,
      isRevertedToDefaultCustomerDB: customerDBPath !== '{{DEFAULT}}' && !useRemoteCustomerDB,
      customerDBPath,
    });
  },
  shops: [],
  fetchShops: async () => {
    const shops = await GetAllShops();
    set({ shops });
  },
  customers: [],
  fetchCustomers: async () => {
    const customers = await GetAllCustomers();
    set({ customers });
  },
}));

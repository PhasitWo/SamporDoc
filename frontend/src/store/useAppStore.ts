import { create } from 'zustand';
import { model } from '../../wailsjs/go/models';
import { GetAllCustomers, GetAllShops } from '../../wailsjs/go/main/App';

interface State {
  init: () => Promise<void>;
  shops: model.Shop[];
  fetchShops: () => Promise<void>;
  customers: model.Customer[];
  fetchCustomers: () => Promise<void>;
}

export const useAppStore = create<State>((set, get) => ({
  init: async () => {
    get().fetchShops();
    get().fetchCustomers();
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

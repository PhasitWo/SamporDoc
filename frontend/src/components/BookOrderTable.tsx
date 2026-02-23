import { ColumnsType } from 'antd/es/table';
import { excel } from '../../wailsjs/go/models';
import { Table } from 'antd';
import { useMemo } from 'react';
import { moneyFormat } from '../utils';

const columns: ColumnsType<excel.PublisherItem> = [
  {
    title: 'สำนักพิมพ์',
    dataIndex: 'Name',
    key: 'name',
  },
  {
    title: 'จำนวนเงิน (บาท)',
    dataIndex: 'TotalAmount',
    key: 'totalAmount',
    align: 'right',
    render: (value: number) => moneyFormat(value),
  },
];

export default function BookOrderTable({ data }: { data: excel.PublisherItem[] }) {
  const grandTotal = useMemo(() => data.reduce((acc, item) => acc + item.TotalAmount, 0), [data]);

  return (
    <Table<excel.PublisherItem>
      columns={columns}
      dataSource={data}
      footer={() => <span className='font-bold'>{`จำนวนเงินรวม: ${moneyFormat(grandTotal)} บาท`}</span>}
      pagination={false}
    />
  );
}

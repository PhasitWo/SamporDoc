import { DatePicker, DatePickerProps } from 'antd';
import { Dayjs } from 'dayjs';
import { forwardRef, useState } from 'react';
import type { PickerRef } from 'rc-picker';

const HiddenDatePicker = forwardRef<PickerRef, DatePickerProps<Dayjs, false>>(({ onOpenChange, ...props }, ref) => {
  const [isOpen, setIsOpen] = useState(false);

  const handleOpenChange = (open: boolean) => {
    setIsOpen(open);
    onOpenChange?.(open);
  };

  return (
    <>
      {isOpen && <div className="fixed inset-0 bg-gray-400/25 z-[1050]" />}
      <DatePicker
        ref={ref}
        className="z-[1000] invisible fixed top-[15%] left-1/2 -translate-x-1/2 -translate-y-1/2"
        popupClassName="z-[1051]"
        onOpenChange={handleOpenChange}
        builtinPlacements={{
          bottomLeft: {
            points: ['tc', 'bc'],
            offset: [0, 4],
            overflow: { adjustX: true, adjustY: true },
          },
        }}
        {...props}
      />
    </>
  );
});

HiddenDatePicker.displayName = 'HiddenDatePicker';

export default HiddenDatePicker;

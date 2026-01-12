import { ComponentProps } from 'react';

export default function InputContainer({ children }: ComponentProps<'div'>) {
  return <div className={'flex flex-col w-full gap-1'}>{children}</div>;
}

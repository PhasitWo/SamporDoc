import { ComponentProps } from 'react';
import { cn } from '../utils';

export default function FormContainer({ children, className }: ComponentProps<'div'>) {
  return (
    <div
      className={cn(
        'mx-auto flex flex-col gap-3 items-center justify-center max-w-[750px] overflow-y-auto',
        className
      )}
    >
      {children}
    </div>
  );
}

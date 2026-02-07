import { ComponentProps } from 'react';
import { cn } from '../utils';

export default function FormContainer({ children, className }: ComponentProps<'div'>) {
  return (
    <div
      className={cn(
        'mx-auto flex flex-col gap-3 items-center justify-center max-w-[500px] overflow-y-scroll',
        className
      )}
    >
      {children}
    </div>
  );
}

import { ComponentProps } from 'react';
import { cn } from '../utils';

export default function Asterisk({ className, ...rest }: ComponentProps<'span'>) {
  return (
    <span className={cn('text-red-500', className)} {...rest}>
      *
    </span>
  );
}

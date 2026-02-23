import { clsx, type ClassValue } from 'clsx';
import { useErrorBoundary } from 'react-error-boundary';
import { twMerge } from 'tailwind-merge';

export function isDeepEqual(obj1: any, obj2: any) {
  if (obj1 === obj2) return true;

  if (typeof obj1 !== 'object' || obj1 === null || typeof obj2 !== 'object' || obj2 === null) {
    return false;
  }

  const keys1 = Object.keys(obj1);
  const keys2 = Object.keys(obj2);

  if (keys1.length !== keys2.length) return false;

  for (let key of keys1) {
    if (!keys2.includes(key) || !isDeepEqual(obj1[key], obj2[key])) {
      return false;
    }
  }

  return true;
}

export function cn(...args: ClassValue[]): string {
  return twMerge(clsx(args));
}

export function useShowBoundary() {
  const { showBoundary: defaultShowBoundary } = useErrorBoundary();
  return {
    showBoundary: (error: any) => {
      if (error instanceof Error) {
        defaultShowBoundary(error);
      } else {
        defaultShowBoundary(new Error(JSON.stringify(error)));
      }
    },
  };
}

export const isValidWindowsFilename = (filename: string): boolean => {
  // 1. Check for reserved characters: < > : " / \ | ? *
  // 2. Check for control characters (0-31)
  const reservedChars = /[<>:"/\\|?*\x00-\x1f]/;

  // 3. Check for reserved names (CON, PRN, AUX, NUL, COM1-9, LPT1-9)
  // These are case-insensitive and cannot be the full filename or the base name before an extension
  const reservedNames = /^(con|prn|aux|nul|com[1-9]|lpt[1-9])(\..*)?$/i;

  // 4. Check for trailing spaces or periods (not allowed in Windows)
  const trailingDotsOrSpaces = /[.]$/;

  if (filename.length > 255) return false;
  if (reservedChars.test(filename)) return false;
  if (reservedNames.test(filename)) return false;
  if (trailingDotsOrSpaces.test(filename)) return false;

  return true;
};

const formatter = new Intl.NumberFormat('en-US', {
  style: 'decimal',
  minimumFractionDigits: 2,
  maximumFractionDigits: 2,
});
export const moneyFormat = (value: number): string => {
  return formatter.format(value);
};

export const getFileName = (path: string): string => {
  // This regex looks for the last / or \ and takes everything after it
  return path.split(/[\\/]/).pop() || '';
};
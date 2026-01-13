import { FallbackProps } from 'react-error-boundary';
import ErrorAlertCard from '../components/ErrorAlertCard';
import { Button } from 'antd';

export default function ErrorFallback({ error, resetErrorBoundary }: FallbackProps) {
  return (
    <div className="w-full p-10">
      <ErrorAlertCard
        messages={[error.message]}
        action={
          <Button danger ghost onClick={resetErrorBoundary}>
            ลองอีกครั้ง
          </Button>
        }
      />
    </div>
  );
}

import { Alert, AlertProps } from 'antd';
import { useMemo } from 'react';

function ErrorAlertCard({
  messages,
  title = 'เกิดข้อผิดพลาด',
  ...rest
}: Omit<AlertProps, 'description' | 'type'> & { messages?: Array<string | boolean | undefined | null> }) {
  const node = useMemo<React.ReactNode | undefined | null>(() => {
    const arr = messages?.filter((v) => typeof v === 'string');
    if (arr?.length == 0) {
      return null
    }
    return (
      <div>
        {arr?.map((v, i) => (
          <span key={i}>
            {v}
            <br />
          </span>
        ))}
      </div>
    );
  }, [messages]);
  if (!node) {
    return null;
  }
  return <Alert type="error" showIcon description={node} title={title} {...rest} />;
}
export default ErrorAlertCard;

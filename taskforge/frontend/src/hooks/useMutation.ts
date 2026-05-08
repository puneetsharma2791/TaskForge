import { useState } from 'react';

export function useMutation<TData, TVars>(
  mutationFn: (vars: TVars) => Promise<TData>,
  options?: { onSuccess?: (data: TData) => void; onError?: (error: string) => void }
) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [data, setData] = useState<TData | null>(null);

  const mutate = async (vars: TVars) => {
    setLoading(true);
    try {
      const result = await mutationFn(vars);
      setData(result);
      options?.onSuccess?.(result);
      return result;
    } catch (e) {
      const msg = e instanceof Error ? e.message : 'Unknown error';
      setError(msg);
      if (options?.onError) {
        options.onError(msg);
      }
      // FLAW: returns undefined on error, caller may not check
    } finally {
      setLoading(false);
    }
  };

  return { mutate, loading, error, data };
}

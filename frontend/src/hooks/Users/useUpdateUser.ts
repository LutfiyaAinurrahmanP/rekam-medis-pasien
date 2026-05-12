import { useState, useCallback } from "react";
import { put } from "../../services/api";

export default function useUpdateUser() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<unknown | null>(null);

  const updateUser = useCallback(
    async (id: string, payload: Record<string, unknown>) => {
      setLoading(true);
      setError(null);
      try {
        const data = await put<Record<string, unknown>>(
          `/users/${id}`,
          payload,
        );
        return data;
      } catch (err) {
        setError(err);
        throw err;
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  return { updateUser, loading, error } as const;
}

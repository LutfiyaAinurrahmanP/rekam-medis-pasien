import { useState, useCallback } from "react";
import { apiCall } from "../../services/api";

export interface DeletedUser {
  id: number;
  username: string;
  email: string;
  phone: string;
  role: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
}

interface DeletedUserListResponse {
  data: DeletedUser[];
  meta: {
    page: number;
    page_size: number;
    total_items: number;
    total_pages: number;
  };
}

interface UseDeletedUsersParams {
  page?: number;
  page_size?: number;
  search?: string;
  role?: string;
  is_active?: boolean;
  sort_by?: string;
  sort_dir?: string;
}

export const useDeletedUsers = () => {
  const [deletedUsers, setDeletedUsers] = useState<DeletedUser[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [meta, setMeta] = useState({
    page: 1,
    page_size: 10,
    total_items: 0,
    total_pages: 1,
  });

  const fetchDeletedUsers = useCallback(
    async (params: UseDeletedUsersParams = {}) => {
      setLoading(true);
      setError(null);

      try {
        const queryParams = new URLSearchParams();

        if (params.page) queryParams.append("page", String(params.page));
        if (params.page_size)
          queryParams.append("page_size", String(params.page_size));
        if (params.search) queryParams.append("search", params.search);
        if (params.role) queryParams.append("role", params.role);
        if (params.is_active !== undefined)
          queryParams.append("is_active", String(params.is_active));
        if (params.sort_by) queryParams.append("sort_by", params.sort_by);
        if (params.sort_dir) queryParams.append("sort_dir", params.sort_dir);

        const endpoint = `/users/deleted?${queryParams.toString()}`;
        const response = await apiCall<DeletedUserListResponse>(endpoint);

        if (response) {
          setDeletedUsers(response.data);
          setMeta(response.meta);
        }
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : "Failed to fetch deleted users";
        setError(errorMessage);
        console.error("Error fetching deleted users:", err);
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  return {
    deletedUsers,
    loading,
    error,
    meta,
    fetchDeletedUsers,
  };
};

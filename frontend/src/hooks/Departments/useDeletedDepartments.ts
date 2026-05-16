import { useState, useCallback } from "react";
import { apiCall } from "../../services/api";

export interface DeletedDepartment {
  id: number;
  name: string;
  code?: string;
  description?: string;
  floor_location?: string;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
}

interface DeletedDepartmentListResponse {
  data: DeletedDepartment[];
  meta: {
    page: number;
    page_size: number;
    total_items: number;
    total_pages: number;
  };
}

interface UseDeletedDepartmentsParams {
  page?: number;
  page_size?: number;
  search?: string;
  sort_by?: string;
  sort_dir?: string;
}

export const useDeletedDepartments = () => {
  const [deletedDepartments, setDeletedDepartments] = useState<
    DeletedDepartment[]
  >([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [meta, setMeta] = useState({
    page: 1,
    page_size: 10,
    total_items: 0,
    total_pages: 1,
  });

  const fetchDeletedDepartments = useCallback(
    async (params: UseDeletedDepartmentsParams = {}) => {
      setLoading(true);
      setError(null);

      try {
        const queryParams = new URLSearchParams();
        if (params.page) queryParams.append("page", String(params.page));
        if (params.page_size)
          queryParams.append("page_size", String(params.page_size));
        if (params.search) queryParams.append("search", params.search);
        if (params.sort_by) queryParams.append("sort_by", params.sort_by);
        if (params.sort_dir) queryParams.append("sort_dir", params.sort_dir);

        const endpoint = `/departments/deleted?${queryParams.toString()}`;
        const response = await apiCall<DeletedDepartmentListResponse>(endpoint);

        if (response) {
          setDeletedDepartments(response.data);
          setMeta(response.meta);
        }
      } catch (err) {
        const errorMessage =
          err instanceof Error
            ? err.message
            : "Failed to fetch deleted departments";
        setError(errorMessage);
        console.error("Error fetching deleted departments:", err);
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  return { deletedDepartments, loading, error, meta, fetchDeletedDepartments };
};

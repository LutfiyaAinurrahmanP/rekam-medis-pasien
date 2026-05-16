import { useState, useCallback } from "react";
import { apiCall } from "../../services/api";

export interface Department {
  id: number;
  name: string;
  code?: string;
  description?: string;
  floor_location?: string;
  created_at: string;
  updated_at: string;
}

interface DepartmentListResponse {
  data: Department[];
  meta: {
    page: number;
    page_size: number;
    total_items: number;
    total_pages: number;
  };
}

interface UseDepartmentsParams {
  page?: number;
  page_size?: number;
  search?: string;
  sort_by?: string;
  sort_dir?: string;
}

export const useDepartments = () => {
  const [departments, setDepartments] = useState<Department[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [meta, setMeta] = useState({
    page: 1,
    page_size: 10,
    total_items: 0,
    total_pages: 1,
  });

  const fetchDepartments = useCallback(
    async (params: UseDepartmentsParams = {}) => {
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

        const endpoint = `/departments?${queryParams.toString()}`;
        const response = await apiCall<DepartmentListResponse>(endpoint);

        if (response) {
          setDepartments(response.data);
          setMeta(response.meta);
        }
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : "Failed to fetch departments";
        setError(errorMessage);
        console.error("Error fetching departments:", err);
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  return {
    departments,
    loading,
    error,
    meta,
    fetchDepartments,
  };
};

import { useState, useCallback } from "react";
import { apiCall } from "../../services/api";
import type {
  DeletedPatient,
  DeletedPatientListResponse,
  PatientQueryParams,
} from "./patientTypes";

export type { DeletedPatient } from "./patientTypes";

export const useDeletedPatients = () => {
  const [deletedPatients, setDeletedPatients] = useState<DeletedPatient[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [meta, setMeta] = useState({
    page: 1,
    page_size: 10,
    total_items: 0,
    total_pages: 1,
  });

  const fetchDeletedPatients = useCallback(
    async (params: PatientQueryParams = {}) => {
      setLoading(true);
      setError(null);

      try {
        const queryParams = new URLSearchParams();

        if (params.page) queryParams.append("page", String(params.page));
        if (params.page_size)
          queryParams.append("page_size", String(params.page_size));
        if (params.search) queryParams.append("search", params.search);
        if (params.gender) queryParams.append("gender", params.gender);
        if (params.blood_type)
          queryParams.append("blood_type", params.blood_type);
        if (params.insurance_provider)
          queryParams.append("insurance_provider", params.insurance_provider);
        if (params.min_age !== undefined)
          queryParams.append("min_age", String(params.min_age));
        if (params.max_age !== undefined)
          queryParams.append("max_age", String(params.max_age));
        if (params.sort_by) queryParams.append("sort_by", params.sort_by);
        if (params.sort_dir) queryParams.append("sort_dir", params.sort_dir);

        const queryString = queryParams.toString();
        const endpoint = queryString
          ? `/patients/deleted?${queryString}`
          : "/patients/deleted";
        const response = await apiCall<DeletedPatientListResponse>(endpoint);

        if (response) {
          setDeletedPatients(response.data);
          setMeta(response.meta);
        }
      } catch (err) {
        const errorMessage =
          err instanceof Error
            ? err.message
            : "Failed to fetch deleted patients";
        setError(errorMessage);
        console.error("Error fetching deleted patients:", err);
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  return {
    deletedPatients,
    loading,
    error,
    meta,
    fetchDeletedPatients,
  };
};

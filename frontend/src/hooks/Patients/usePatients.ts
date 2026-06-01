import { useState, useCallback } from "react";
import { apiCall } from "../../services/api";
import type {
  Patient,
  PatientListResponse,
  PatientQueryParams,
} from "./patientTypes";

export type { Patient } from "./patientTypes";

export const usePatients = () => {
  const [patients, setPatients] = useState<Patient[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [meta, setMeta] = useState({
    page: 1,
    page_size: 10,
    total_items: 0,
    total_pages: 1,
  });

  const fetchPatients = useCallback(async (params: PatientQueryParams = {}) => {
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
      const endpoint = queryString ? `/patients?${queryString}` : "/patients";
      const response = await apiCall<PatientListResponse>(endpoint);

      if (response) {
        setPatients(response.data);
        setMeta(response.meta);
      }
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : "Failed to fetch patients";
      setError(errorMessage);
      console.error("Error fetching patients:", err);
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    patients,
    loading,
    error,
    meta,
    fetchPatients,
  };
};

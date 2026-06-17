export interface PatientPaginationMeta {
  page: number;
  page_size: number;
  total_items: number;
  total_pages: number;
}

export interface Patient {
  id: number;
  user_id?: number | null;
  patient_code: string;
  full_name: string;
  date_of_birth: string;
  age: number;
  gender: string;
  blood_type: string;
  phone: string;
  email: string;
  address: string;
  emergency_contact_name: string;
  emergency_contact_phone: string;
  insurance_number: string;
  insurance_provider: string;
  allergies: string;
  created_at: string;
  updated_at: string;
}

export interface DeletedPatient extends Patient {
  deleted_at: string | null;
}

export interface PatientListResponse {
  data: Patient[];
  meta: PatientPaginationMeta;
}

export interface DeletedPatientListResponse {
  data: DeletedPatient[];
  meta: PatientPaginationMeta;
}

export interface PatientQueryParams {
  page?: number;
  page_size?: number;
  search?: string;
  gender?: string;
  blood_type?: string;
  insurance_provider?: string;
  min_age?: number;
  max_age?: number;
  sort_by?: string;
  sort_dir?: string;
}

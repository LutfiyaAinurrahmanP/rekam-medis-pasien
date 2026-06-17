import { useCallback, useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router";
import { get, put } from "../../../services/api";
import { useAuth } from "../../../context/AuthContext";
import EditFormModal, {
  type EditFormFieldDef,
} from "../../../components/modals/EditFormModal";
import SuccessModal from "../../../components/ui/notification/SuccessModal";
import { getRolePatientsPath } from "../../../pages/Roles/shared/role-routing";
import type { Patient } from "../../../hooks/Patients/usePatients";

interface EditPatientModalProps {
  isOpen: boolean;
  id?: string;
  onClose?: () => void;
  onSuccess?: () => void;
}

export default function EditPatientModal({
  isOpen,
  id,
  onClose,
  onSuccess,
}: EditPatientModalProps) {
  const navigate = useNavigate();
  const location = useLocation();
  const { user } = useAuth();
  const role = location.pathname.split("/")[1] || undefined;

  const [loadingFetch, setLoadingFetch] = useState(true);
  const [errorsList, setErrorsList] = useState<string[]>([]);
  const [showSuccessModal, setShowSuccessModal] = useState(false);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);
  const [formData, setFormData] = useState<Record<string, string>>({
    patient_code: "",
    full_name: "",
    date_of_birth: "",
    gender: "",
    blood_type: "",
    phone: "",
    email: "",
    address: "",
    emergency_contact_name: "",
    emergency_contact_phone: "",
    insurance_number: "",
    insurance_provider: "",
    allergies: "",
  });

  const close = useCallback(() => {
    if (onClose) {
      onClose();
      return;
    }

    navigate(getRolePatientsPath(role || user?.role));
  }, [navigate, onClose, role, user?.role]);

  const handleSuccessClose = useCallback(() => {
    setShowSuccessModal(false);
    if (onSuccess) {
      onSuccess();
      return;
    }
    close();
  }, [close, onSuccess]);

  useEffect(() => {
    if (!showSuccessModal) {
      return undefined;
    }

    const timer = window.setTimeout(() => {
      handleSuccessClose();
    }, 3000);

    return () => window.clearTimeout(timer);
  }, [handleSuccessClose, showSuccessModal]);

  useEffect(() => {
    let mounted = true;

    async function fetchPatient() {
      if (!id) {
        setLoadingFetch(false);
        return;
      }

      setLoadingFetch(true);
      try {
        const data = await get<Patient>(`/patients/${id}`);
        if (!mounted) return;

        setFormData({
          patient_code: String(data.patient_code ?? ""),
          full_name: String(data.full_name ?? ""),
          date_of_birth: String(data.date_of_birth ?? ""),
          gender: String(data.gender ?? ""),
          blood_type: String(data.blood_type ?? ""),
          phone: String(data.phone ?? ""),
          email: String(data.email ?? ""),
          address: String(data.address ?? ""),
          emergency_contact_name: String(data.emergency_contact_name ?? ""),
          emergency_contact_phone: String(data.emergency_contact_phone ?? ""),
          insurance_number: String(data.insurance_number ?? ""),
          insurance_provider: String(data.insurance_provider ?? ""),
          allergies: String(data.allergies ?? ""),
        });
      } catch (err) {
        console.error("Failed to fetch patient for edit:", err);
      } finally {
        if (mounted) setLoadingFetch(false);
      }
    }

    if (isOpen) fetchPatient();

    return () => {
      mounted = false;
    };
  }, [id, isOpen]);

  const mapValidationTagToMessage = (field: string, tag: string) => {
    const normalizedField = field.charAt(0).toUpperCase() + field.slice(1);
    const normalizedTag = tag.toLowerCase().replace(/_/g, " ");

    if (normalizedTag === "required") {
      return `${normalizedField}: is required`;
    }

    if (normalizedTag === "email") {
      return `${normalizedField}: invalid email format`;
    }

    if (normalizedTag === "uuid") {
      return `${normalizedField}: must be a valid UUID`;
    }

    if (normalizedTag === "min" || normalizedTag === "min length") {
      return `${normalizedField}: value is too short`;
    }

    if (normalizedTag === "max" || normalizedTag === "max length") {
      return `${normalizedField}: value is too long`;
    }

    if (normalizedTag === "oneof") {
      return `${normalizedField}: invalid value`;
    }

    return `${normalizedField}: invalid value`;
  };

  const getErrorList = (unknownErr: unknown): string[] => {
    const raw = unknownErr as Record<string, unknown> | undefined;
    const list: string[] = [];

    if (raw && raw.errors && typeof raw.errors === "object") {
      const errorsObj = raw.errors as Record<string, unknown>;
      for (const [field, msg] of Object.entries(errorsObj)) {
        const label = field.charAt(0).toUpperCase() + field.slice(1);
        list.push(`${label}: ${String(msg)}`);
      }

      return list;
    }

    const rawError = raw && (raw.error ?? raw.message ?? "");

    if (typeof rawError === "string" && rawError.trim().length > 0) {
      const parts = rawError
        .split(/\r?\n|;|\|/)
        .map((part: string) => part.trim())
        .filter(Boolean);

      for (const part of parts) {
        const validationMatch = part.match(
          /Field validation for '([^']+)' failed on the '([^']+)' tag/i,
        );

        if (validationMatch) {
          const field = validationMatch[1];
          const tag = validationMatch[2];
          list.push(mapValidationTagToMessage(field, tag));
          continue;
        }

        const quoted = part.match(/'([^']+)'/);
        if (quoted) {
          const field = quoted[1];
          if (/required|cannot be empty|missing/i.test(part)) {
            list.push(`${field}: is required`);
            continue;
          }
          if (/invalid/i.test(part)) {
            list.push(`${field}: invalid value`);
            continue;
          }
        }

        const existsMatch = part.match(
          /(patient code|full name|email|phone) already exists/i,
        );
        if (existsMatch) {
          const field = existsMatch[1];
          list.push(
            `${field.charAt(0).toUpperCase() + field.slice(1)}: ${existsMatch[0]}`,
          );
          continue;
        }

        list.push(part);
      }

      if (list.length > 0) {
        return list;
      }
    }

    return ["Failed to update patient. Please try again."];
  };

  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement
    >,
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setErrorsList([]);

    if (!id) {
      setErrorsList(["Missing patient id"]);
      return;
    }

    try {
      const payload: Record<string, unknown> = {};

      for (const key of [
        "full_name",
        "date_of_birth",
        "gender",
        "blood_type",
        "phone",
        "email",
        "address",
        "emergency_contact_name",
        "emergency_contact_phone",
        "insurance_number",
        "insurance_provider",
        "allergies",
      ] as const) {
        const value = formData[key].trim();
        if (value) {
          payload[key] = value;
        }
      }

      const patient = await put<Patient>(`/patients/${id}`, payload);
      setSuccessMsg(`Patient "${patient.full_name}" updated successfully.`);
      setShowSuccessModal(true);
    } catch (err) {
      console.error("Update patient failed:", err);
      setErrorsList(getErrorList(err));
    }
  };

  const fields: EditFormFieldDef[] = [
    {
      name: "patient_code",
      label: "Patient Code",
      type: "text",
      disabled: true,
    },
    { name: "full_name", label: "Full Name", type: "text", required: true },
    {
      name: "date_of_birth",
      label: "Date of Birth",
      type: "date",
      required: true,
    },
    {
      name: "gender",
      label: "Gender",
      type: "select",
      required: true,
      options: [
        { value: "male", label: "Male" },
        { value: "female", label: "Female" },
        { value: "other", label: "Other" },
      ],
    },
    {
      name: "blood_type",
      label: "Blood Type",
      type: "text",
      placeholder: "O+, A+, B+, AB+, etc.",
    },
    { name: "phone", label: "Phone", type: "tel" },
    { name: "email", label: "Email", type: "email" },
    {
      name: "address",
      label: "Address",
      type: "textarea",
      rows: 3,
    },
    {
      name: "emergency_contact_name",
      label: "Emergency Contact Name",
      type: "text",
    },
    {
      name: "emergency_contact_phone",
      label: "Emergency Contact Phone",
      type: "tel",
    },
    {
      name: "insurance_number",
      label: "Insurance Number",
      type: "text",
    },
    {
      name: "insurance_provider",
      label: "Insurance Provider",
      type: "text",
    },
    {
      name: "allergies",
      label: "Allergies",
      type: "textarea",
      rows: 3,
    },
  ];

  return (
    <>
      <EditFormModal
        isOpen={isOpen}
        formData={formData}
        onChange={handleChange}
        onSubmit={handleSubmit}
        onClose={close}
        errorsList={errorsList}
        loading={loadingFetch}
        fields={fields}
        title="Edit Patient"
        description="Update patient details below"
        submitLabel="Update"
        cancelLabel="Cancel"
      />

      <SuccessModal
        isOpen={showSuccessModal}
        title="Patient Updated"
        message={successMsg ?? "Patient updated successfully."}
        buttonLabel="Continue"
        onButtonClick={handleSuccessClose}
      />
    </>
  );
}

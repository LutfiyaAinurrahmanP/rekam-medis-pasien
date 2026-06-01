import { useCallback, useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { post } from "../../../services/api";
import { useAuth } from "../../../context/AuthContext";
import CreateFormModal, {
  type CreateFormFieldDef,
} from "../../../components/modals/CreateFormModal";
import SuccessModal from "../../../components/ui/notification/SuccessModal";
import { getRolePatientsPath } from "../../../pages/Roles/shared/role-routing";
import type { Patient } from "../../../hooks/Patients/usePatients";

interface CreatePatientModalProps {
  isOpen: boolean;
  onClose?: () => void;
  onSuccess?: () => void;
}

export default function CreatePatientModal({
  isOpen,
  onClose,
  onSuccess,
}: CreatePatientModalProps) {
  const navigate = useNavigate();
  const { user } = useAuth();

  const [loading, setLoading] = useState(false);
  const [errorsList, setErrorsList] = useState<string[]>([]);
  const [showSuccessModal, setShowSuccessModal] = useState(false);
  const [successData, setSuccessData] = useState<{
    full_name: string;
    patient_code: string;
  } | null>(null);
  const [formData, setFormData] = useState<Record<string, string>>({
    user_id: "",
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

    navigate(getRolePatientsPath(user?.role));
  }, [navigate, onClose, user?.role]);

  const handleSuccessModalClose = useCallback(() => {
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
      handleSuccessModalClose();
    }, 3000);

    return () => window.clearTimeout(timer);
  }, [handleSuccessModalClose, showSuccessModal]);

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

    const topMsg = raw && raw.message ? String(raw.message).toLowerCase() : "";
    if (
      topMsg.includes("duplicate data") ||
      topMsg.includes("already exists")
    ) {
      const rawErrStr = raw && raw.error ? String(raw.error).toLowerCase() : "";
      if (rawErrStr.includes("patient code")) {
        list.push("Patient code: Patient code already exists");
      }
      if (rawErrStr.includes("email")) list.push("Email: Email already exists");
      if (rawErrStr.includes("phone")) list.push("Phone: Phone already exists");
      if (list.length === 0) list.push("Duplicate data: conflict detected");
      return list;
    }

    return ["Failed to create patient. Please try again."];
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
    setLoading(true);

    try {
      const payload: Record<string, unknown> = {
        full_name: formData.full_name.trim(),
        date_of_birth: formData.date_of_birth,
        gender: formData.gender,
      };

      const userId = formData.user_id.trim();
      if (userId) {
        payload.user_id = Number(userId);
      }

      const patientCode = formData.patient_code.trim();
      if (patientCode) {
        payload.patient_code = patientCode;
      }

      for (const key of [
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

      const patient = await post<Patient>("/patients", payload);
      setSuccessData({
        full_name: patient.full_name,
        patient_code: patient.patient_code,
      });
      setShowSuccessModal(true);
      setFormData({
        user_id: "",
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
    } catch (unknownErr) {
      console.error("Create patient error:", unknownErr);
      setErrorsList(getErrorList(unknownErr));
    } finally {
      setLoading(false);
    }
  };

  const fields: CreateFormFieldDef[] = [
    {
      name: "user_id",
      label: "User ID",
      type: "text",
      placeholder: "Optional user account ID",
    },
    {
      name: "patient_code",
      label: "Patient Code",
      type: "text",
      placeholder: "Leave blank to auto-generate",
    },
    {
      name: "full_name",
      label: "Full Name",
      type: "text",
      placeholder: "Enter full name",
      required: true,
    },
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
    {
      name: "phone",
      label: "Phone",
      type: "tel",
      placeholder: "Enter phone number",
    },
    {
      name: "email",
      label: "Email",
      type: "email",
      placeholder: "Enter email address",
    },
    {
      name: "address",
      label: "Address",
      type: "textarea",
      placeholder: "Enter address",
      rows: 3,
    },
    {
      name: "emergency_contact_name",
      label: "Emergency Contact Name",
      type: "text",
      placeholder: "Enter emergency contact name",
    },
    {
      name: "emergency_contact_phone",
      label: "Emergency Contact Phone",
      type: "tel",
      placeholder: "Enter emergency contact phone",
    },
    {
      name: "insurance_number",
      label: "Insurance Number",
      type: "text",
      placeholder: "Enter insurance number",
    },
    {
      name: "insurance_provider",
      label: "Insurance Provider",
      type: "text",
      placeholder: "Enter insurance provider",
    },
    {
      name: "allergies",
      label: "Allergies",
      type: "textarea",
      placeholder: "Enter allergies or medical notes",
      rows: 3,
    },
  ];

  return (
    <>
      <CreateFormModal
        isOpen={isOpen}
        formData={formData}
        onChange={handleChange}
        onSubmit={handleSubmit}
        onClose={close}
        errorsList={errorsList}
        loading={loading}
        fields={fields}
        title="Create Patient"
        description="Enter the patient details below"
        submitLabel="Create"
        cancelLabel="Cancel"
      />

      <SuccessModal
        isOpen={showSuccessModal}
        title="Patient Created Successfully"
        message={`Patient "${successData?.full_name ?? ""}" (${successData?.patient_code ?? ""}) has been created successfully.`}
        buttonLabel="Continue"
        onButtonClick={handleSuccessModalClose}
      />
    </>
  );
}

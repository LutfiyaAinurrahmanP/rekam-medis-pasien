import { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router";
import { get } from "../../../services/api";
import { useAuth } from "../../../context/AuthContext";
import DetailModal from "../../../components/modals/DetailModal";
import { formatDateIndonesian } from "../../../utils";
import { getRolePatientsPath } from "../../../pages/Roles/shared/role-routing";
import type { Patient } from "../../../hooks/Patients/usePatients";

interface ShowPatientModalProps {
  isOpen: boolean;
  id?: string;
  onClose?: () => void;
}

export default function ShowPatientModal({
  isOpen,
  id,
  onClose,
}: ShowPatientModalProps) {
  const navigate = useNavigate();
  const location = useLocation();
  const { user } = useAuth();
  const role = location.pathname.split("/")[1] || undefined;

  const [loading, setLoading] = useState(true);
  const [patient, setPatient] = useState<Patient | null>(null);

  useEffect(() => {
    let mounted = true;

    async function fetchPatient() {
      if (!id) {
        setPatient(null);
        setLoading(false);
        return;
      }

      setLoading(true);
      try {
        const data = await get<Patient>(`/patients/${id}`);
        if (!mounted) return;
        setPatient(data);
      } catch (err) {
        console.error("Failed to fetch patient:", err);
        if (!mounted) return;
        setPatient(null);
      } finally {
        if (!mounted) return;
        setLoading(false);
      }
    }

    if (isOpen) fetchPatient();

    return () => {
      mounted = false;
    };
  }, [id, isOpen]);

  const close = () => {
    if (onClose) {
      onClose();
      return;
    }

    navigate(getRolePatientsPath(role || user?.role));
  };

  return (
    <DetailModal
      isOpen={isOpen}
      loading={loading}
      data={patient as Record<string, unknown> | null}
      fields={[
        { key: "patient_code", label: "Patient Code" },
        { key: "user_id", label: "User ID" },
        { key: "full_name", label: "Full Name" },
        { key: "age", label: "Age" },
        {
          key: "date_of_birth",
          label: "Date of Birth",
          render: (v) =>
            formatDateIndonesian((v as string | null | undefined) ?? null),
        },
        {
          key: "gender",
          label: "Gender",
          render: (v) =>
            String(v ?? "-").replace(/^./, (char) => char.toUpperCase()),
        },
        { key: "blood_type", label: "Blood Type" },
        { key: "phone", label: "Phone" },
        { key: "email", label: "Email" },
        { key: "address", label: "Address" },
        { key: "emergency_contact_name", label: "Emergency Contact Name" },
        { key: "emergency_contact_phone", label: "Emergency Contact Phone" },
        { key: "insurance_number", label: "Insurance Number" },
        { key: "insurance_provider", label: "Insurance Provider" },
        { key: "allergies", label: "Allergies" },
        {
          key: "created_at",
          label: "Created At",
          render: (v) =>
            formatDateIndonesian((v as string | null | undefined) ?? null),
        },
        {
          key: "updated_at",
          label: "Updated At",
          render: (v) =>
            formatDateIndonesian((v as string | null | undefined) ?? null),
        },
      ]}
      onClose={close}
    />
  );
}

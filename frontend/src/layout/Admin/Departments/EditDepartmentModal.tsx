import { useCallback, useEffect, useState } from "react";
import { get, put } from "../../../services/api";
import CreateFormModal, {
  type CreateFormFieldDef,
} from "../../../components/modals/CreateFormModal";
import SuccessModal from "../../../components/ui/notification/SuccessModal";
import { useNavigate } from "react-router";

interface Props {
  isOpen: boolean;
  id?: string | undefined;
  onClose?: () => void;
  onSuccess?: () => void;
}

const fields: CreateFormFieldDef[] = [
  {
    name: "name",
    label: "Name",
    type: "text",
    placeholder: "Enter department name",
    required: true,
  },
  {
    name: "code",
    label: "Code",
    type: "text",
    placeholder: "Enter department code",
  },
  {
    name: "description",
    label: "Description",
    type: "text",
    placeholder: "Enter short description",
  },
  {
    name: "floor_location",
    label: "Floor/Location",
    type: "text",
    placeholder: "Enter floor or location",
  },
];

export default function EditDepartmentModal({
  isOpen,
  id,
  onClose,
  onSuccess,
}: Props) {
  const navigate = useNavigate();
  const [formData, setFormData] = useState<Record<string, string>>({
    name: "",
    code: "",
    description: "",
    floor_location: "",
  });
  const [loading, setLoading] = useState(false);
  const [errorsList, setErrorsList] = useState<string[]>([]);
  const [showSuccessModal, setShowSuccessModal] = useState(false);

  useEffect(() => {
    const load = async () => {
      if (!id) return;
      setLoading(true);
      try {
        const data = await get<{
          id: number;
          name: string;
          code?: string;
          description?: string;
          floor_location?: string;
        }>(`/departments/${id}`);
        setFormData({
          name: data.name,
          code: data.code ?? "",
          description: data.description ?? "",
          floor_location: data.floor_location ?? "",
        });
      } catch (err) {
        console.error("Failed to load department:", err);
      } finally {
        setLoading(false);
      }
    };

    if (isOpen) load();
  }, [id, isOpen]);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>,
  ) => {
    const { name, value } = e.target;
    setFormData((s) => ({ ...s, [name]: value }));
  };

  const handleClose = useCallback(() => {
    if (onClose) return onClose();
    navigate("/admin/departments");
  }, [navigate, onClose]);

  const handleSuccessClose = useCallback(() => {
    setShowSuccessModal(false);
    if (onSuccess) {
      onSuccess();
      return;
    }

    navigate("/admin/departments");
  }, [navigate, onSuccess]);

  useEffect(() => {
    if (!showSuccessModal) return undefined;

    const timer = window.setTimeout(() => {
      handleSuccessClose();
    }, 2000);

    return () => window.clearTimeout(timer);
  }, [handleSuccessClose, showSuccessModal]);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!id) return;
    setErrorsList([]);
    setLoading(true);
    try {
      await put(`/departments/${id}`, {
        name: formData.name,
        code: formData.code || undefined,
        description: formData.description || undefined,
        floor_location: formData.floor_location || undefined,
      });
      setShowSuccessModal(true);
    } catch (err) {
      console.error("Update department error:", err);
      setErrorsList(["Failed to update department."]);
    } finally {
      setLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <>
      <CreateFormModal
        isOpen={isOpen}
        formData={formData}
        onChange={handleChange}
        onSubmit={handleSubmit}
        onClose={handleClose}
        errorsList={errorsList}
        loading={loading}
        fields={fields}
        title="Edit Department"
        description="Update department details below"
        submitLabel="Save"
      />

      <SuccessModal
        isOpen={showSuccessModal}
        title="Department Updated"
        message={`Department has been updated successfully.`}
        buttonLabel="Continue"
        onButtonClick={handleSuccessClose}
      />
    </>
  );
}

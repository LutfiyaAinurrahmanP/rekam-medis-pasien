import { useCallback, useState, useEffect } from "react";
import { post } from "../../../services/api";
import CreateFormModal, {
  type CreateFormFieldDef,
} from "../../../components/modals/CreateFormModal";
import SuccessModal from "../../../components/ui/notification/SuccessModal";
import { useNavigate } from "react-router";

interface Props {
  isOpen: boolean;
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

export default function CreateDepartmentModal({
  isOpen,
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
    if (showSuccessModal) {
      const t = window.setTimeout(() => {
        setShowSuccessModal(false);
        if (onSuccess) onSuccess();
      }, 2000);
      return () => window.clearTimeout(t);
    }
    return undefined;
  }, [showSuccessModal, onSuccess]);

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

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setErrorsList([]);
    setLoading(true);

    try {
      await post("/departments", {
        name: formData.name,
        code: formData.code || undefined,
        description: formData.description || undefined,
        floor_location: formData.floor_location || undefined,
      });
      setShowSuccessModal(true);
      setFormData({ name: "", code: "", description: "", floor_location: "" });
    } catch (err) {
      console.error("Create department error:", err);
      setErrorsList(["Failed to create department."]);
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
        title="Create Department"
        description="Enter department details below"
        submitLabel="Create"
      />

      <SuccessModal
        isOpen={showSuccessModal}
        title="Department Created"
        message={`Department has been created successfully.`}
        buttonLabel="Continue"
        onButtonClick={() => {
          setShowSuccessModal(false);
          if (onSuccess) onSuccess();
        }}
      />
    </>
  );
}

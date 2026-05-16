import { useCallback, useEffect, useState } from "react";
import { get, put } from "../../../services/api";
import CreateFormModal, {
  type CreateFormFieldDef,
} from "../../../components/modals/CreateFormModal";
import SuccessModal from "../../../components/ui/notification/SuccessModal";
import { useNavigate } from "react-router";
import { useAuth } from "../../../context/AuthContext";
import { getRoleDepartmentsPath } from "../../../pages/Roles/shared/role-routing";

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
  const { user } = useAuth();
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
    navigate(getRoleDepartmentsPath(user?.role));
  }, [navigate, onClose, user?.role]);

  const mapValidationTagToMessage = (field: string, tag: string) => {
    const normalizedField = field.charAt(0).toUpperCase() + field.slice(1);
    const normalizedTag = tag.toLowerCase().replace(/_/g, " ");

    if (normalizedTag === "required") {
      return `${normalizedField}: is required`;
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

  const getErrorList = (unknownErr: unknown) => {
    const list: string[] = [];
    const raw = unknownErr as {
      message?: unknown;
      error?: unknown;
      errors?: unknown;
    };

    if (raw && raw.errors && typeof raw.errors === "object") {
      for (const [field, value] of Object.entries(
        raw.errors as Record<string, unknown>,
      )) {
        const label = field.charAt(0).toUpperCase() + field.slice(1);
        list.push(`${label}: ${String(value)}`);
      }

      if (list.length > 0) {
        return list;
      }
    }

    const rawError =
      typeof raw?.error === "string"
        ? raw.error
        : typeof raw?.message === "string"
          ? raw.message
          : "";

    if (!rawError) {
      return list;
    }

    const parts = rawError
      .split(/\r?\n|;|\|/)
      .map((part) => part.trim())
      .filter(Boolean);

    for (const part of parts) {
      const validationMatch = part.match(
        /Field validation for '([^']+)' failed on the '([^']+)' tag/i,
      );

      if (validationMatch) {
        list.push(
          mapValidationTagToMessage(validationMatch[1], validationMatch[2]),
        );
        continue;
      }

      const duplicateFieldMatch = part.match(
        /\b(name|code|description|floor[_\s-]?location)\b.*(duplicate|already exists|unique)/i,
      );

      if (duplicateFieldMatch) {
        const field = duplicateFieldMatch[1]
          .replace(/[_\s-]+/g, " ")
          .replace(/\b\w/g, (char) => char.toUpperCase());
        const fieldLabel =
          field === "Floor Location" ? "Floor/Location" : field;
        list.push(`${fieldLabel}: already exists`);
        continue;
      }

      list.push(part);
    }

    if (list.length > 0) {
      return list;
    }

    const topMsg = rawError.toLowerCase();
    if (
      topMsg.includes("duplicate data") ||
      topMsg.includes("already exists") ||
      topMsg.includes("unique constraint")
    ) {
      if (topMsg.includes("name")) list.push("Name: already exists");
      if (topMsg.includes("code")) list.push("Code: already exists");
      if (topMsg.includes("description")) {
        list.push("Description: already exists");
      }
      if (topMsg.includes("floor") || topMsg.includes("location")) {
        list.push("Floor/Location: already exists");
      }

      if (list.length === 0) {
        list.push("Duplicate data: conflict detected");
      }
    }

    return list;
  };

  const handleSuccessClose = useCallback(() => {
    setShowSuccessModal(false);
    if (onSuccess) {
      onSuccess();
      return;
    }

    navigate(getRoleDepartmentsPath(user?.role));
  }, [navigate, onSuccess, user?.role]);

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
      const list = getErrorList(err);
      setErrorsList(
        list.length > 0
          ? list
          : ["Failed to update department. Please try again."],
      );
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

import { useEffect, useState, useCallback } from "react";
import { post } from "../../../services/api";
import type { User } from "../../../hooks/Users/useUsers";
import SuccessModal from "../../../components/ui/notification/SuccessModal";
import CreateFormModal from "../../../components/modals/CreateFormModal";
import type { CreateFormFieldDef } from "../../../components/modals/CreateFormModal";
import { useNavigate } from "react-router";
import { getRoleUsersPath } from "../../../pages/Roles/shared/role-routing";

interface CreateUserModalProps {
  isOpen: boolean;
  role?: string;
  forcedRole?: string;
  onClose?: () => void; // optional callback for inline usage
  onSuccess?: () => void; // optional callback for inline usage
}

type CreateUserFormState = {
  username: string;
  email: string;
  phone: string;
  password: string;
  confirmPassword: string;
  role: string;
};

// Form field configuration
const baseFormFields: CreateFormFieldDef[] = [
  {
    name: "username",
    label: "Username",
    type: "text",
    placeholder: "Enter username",
    required: true,
  },
  {
    name: "email",
    label: "Email",
    type: "email",
    placeholder: "Enter email",
    required: true,
  },
  {
    name: "phone",
    label: "Phone",
    type: "tel",
    placeholder: "Enter phone number",
    required: true,
  },
  {
    name: "role",
    label: "Role",
    type: "select",
    required: true,
    options: [
      { value: "patient", label: "Patient" },
      { value: "doctor", label: "Doctor" },
      { value: "receptionist", label: "Receptionist" },
      { value: "admin", label: "Admin" },
      { value: "super_admin", label: "Super Admin" },
    ],
  },
  {
    name: "password",
    label: "Password",
    type: "password",
    placeholder: "Enter password",
    required: true,
  },
  {
    name: "confirmPassword",
    label: "Confirm Password",
    type: "password",
    placeholder: "Confirm password",
    required: true,
  },
];

export default function CreateUserModal({
  isOpen,
  role,
  forcedRole,
  onClose,
  onSuccess,
}: CreateUserModalProps) {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [errorsList, setErrorsList] = useState<string[]>([]);
  const [showSuccessModal, setShowSuccessModal] = useState(false);
  const [successData, setSuccessData] = useState<{ username: string } | null>(
    null,
  );
  const [formData, setFormData] = useState<CreateUserFormState>({
    username: "",
    email: "",
    phone: "",
    password: "",
    confirmPassword: "",
    role: forcedRole || role || "",
  });

  const formFields =
    forcedRole && forcedRole.length > 0
      ? baseFormFields.filter((field) => field.name !== "role")
      : baseFormFields;

  const close = useCallback(() => {
    if (onClose) {
      onClose();
    } else {
      // Navigate back to users list if no callback provided
      navigate(getRoleUsersPath(role));
    }
  }, [navigate, role, onClose]);

  const handleSuccessModalClose = useCallback(() => {
    setShowSuccessModal(false);
    if (onSuccess) {
      onSuccess();
    }
    close();
  }, [close, onSuccess]);

  // Auto-close success modal after 3 seconds
  useEffect(() => {
    if (showSuccessModal) {
      const timer = setTimeout(() => {
        handleSuccessModalClose();
      }, 3000);

      return () => clearTimeout(timer);
    }
  }, [showSuccessModal, handleSuccessModalClose]);

  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement
    >,
  ) => {
    const { name, value } = e.target;

    if (forcedRole && name === "role") {
      return;
    }

    setFormData((previous) => ({
      ...previous,
      [name]: value,
    }));
  };

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

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setErrorsList([]);

    if (formData.password !== formData.confirmPassword) {
      setErrorsList(["Password and confirm password do not match."]);
      return;
    }

    setLoading(true);

    try {
      await post<User>("/users", {
        username: formData.username,
        email: formData.email,
        phone: formData.phone,
        password: formData.password,
        role: forcedRole || formData.role || undefined,
      });

      // Store username for success modal
      setSuccessData({ username: formData.username });
      setShowSuccessModal(true);

      // Reset form
      setFormData({
        username: "",
        email: "",
        phone: "",
        password: "",
        confirmPassword: "",
        role: forcedRole || role || "",
      });
    } catch (unknownErr) {
      console.error("Create user error:", unknownErr);

      const raw = unknownErr as Record<string, unknown> | undefined;
      const list: string[] = [];

      // Prefer structured `errors` map if provided
      if (raw && raw.errors && typeof raw.errors === "object") {
        const errorsObj = raw.errors as Record<string, unknown>;
        for (const [field, msg] of Object.entries(errorsObj)) {
          const label = field.charAt(0).toUpperCase() + field.slice(1);
          list.push(`${label}: ${String(msg)}`);
        }

        setErrorsList(list);
        return;
      }

      // Backend may return a string in `error` with joined validation messages.
      // Try to split and extract field-level info.
      const rawError = raw && (raw.error ?? raw.message ?? "");

      if (typeof rawError === "string" && rawError.trim().length > 0) {
        // Split on newlines or semicolons
        const parts = rawError
          .split(/\r?\n|;|\|/)
          .map((p: string) => p.trim())
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

          // Try to extract a quoted field name or CamelCase field
          const quoted = part.match(/'([^']+)'/);
          if (quoted) {
            const field = quoted[1];
            // Make a friendly message
            if (/required|cannot be empty|missing/i.test(part)) {
              list.push(`${field}: is required`);
              continue;
            }
            if (/invalid/i.test(part)) {
              list.push(`${field}: invalid value`);
              continue;
            }
          }

          // Try to find patterns like "Username already exists"
          const existsMatch = part.match(
            /(username|email|phone) already exists/i,
          );
          if (existsMatch) {
            const field = existsMatch[1];
            list.push(
              `${field.charAt(0).toUpperCase() + field.slice(1)}: ${existsMatch[0]}`,
            );
            continue;
          }

          // Fallback: push the raw part (English assumed)
          list.push(part);
        }

        if (list.length > 0) {
          setErrorsList(list);
          return;
        }
      }

      // Finally, map common top-level messages
      const topMsg =
        raw && raw.message ? String(raw.message).toLowerCase() : "";
      if (
        topMsg.includes("duplicate data") ||
        topMsg.includes("already exists")
      ) {
        const rawErrStr =
          raw && raw.error ? String(raw.error).toLowerCase() : "";
        if (rawErrStr.includes("username"))
          list.push("Username: Username already exists");
        if (rawErrStr.includes("email"))
          list.push("Email: Email already exists");
        if (rawErrStr.includes("phone"))
          list.push("Phone: Phone already exists");
        if (list.length === 0) list.push("Duplicate data: conflict detected");
        setErrorsList(list);
        return;
      }

      // Generic fallback
      setErrorsList(["Failed to create user. Please try again."]);
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
        onClose={close}
        errorsList={errorsList}
        loading={loading}
        fields={formFields}
        title="Create User"
        description="Enter the user details below"
        submitLabel="Create"
        cancelLabel="Cancel"
      />

      {/* Success Modal - Auto closes after 3 seconds */}
      <SuccessModal
        isOpen={showSuccessModal}
        title="User Created Successfully"
        message={`User "${successData?.username}" has been created successfully.`}
        buttonLabel="Continue"
        onButtonClick={handleSuccessModalClose}
      />
    </>
  );
}

import { useEffect, useState, useCallback } from "react";
import { get } from "../../../services/api";
import type { User } from "../../../hooks/Users/useUsers";
import type { EditFormFieldDef } from "../../../components/modals/EditFormModal";
import { useNavigate } from "react-router";
import EditFormModal from "../../../components/modals/EditFormModal";
import SuccessModal from "../../../components/ui/notification/SuccessModal";
import useUpdateUser from "../../../hooks/Users/useUpdateUser";
import { getRoleUsersPath } from "../../../pages/Roles/shared/role-routing";

interface EditUserModalProps {
  isOpen: boolean;
  id?: string;
  role?: string;
  onClose?: () => void;
  onSuccess?: () => void;
}

export default function EditUserModal({
  isOpen,
  id,
  role,
  onClose,
  onSuccess,
}: EditUserModalProps) {
  const navigate = useNavigate();
  const [loadingFetch, setLoadingFetch] = useState(true);
  const [errorsList, setErrorsList] = useState<string[]>([]);
  const [formData, setFormData] = useState<Record<string, string>>({
    username: "",
    email: "",
    phone: "",
    role: role || "",
    is_active: "true",
  });
  const [showSuccessModal, setShowSuccessModal] = useState(false);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  const { updateUser, loading: updating } = useUpdateUser();

  const close = useCallback(() => {
    if (onClose) onClose();
    else navigate(getRoleUsersPath(role));
  }, [navigate, onClose, role]);

  const handleSuccessClose = useCallback(() => {
    setShowSuccessModal(false);
    if (onSuccess) onSuccess();
    close();
  }, [close, onSuccess]);

  // Auto-close success modal after 3 seconds (same behavior as CreateUserModal)
  useEffect(() => {
    if (showSuccessModal) {
      const timer = setTimeout(() => {
        handleSuccessClose();
      }, 3000);

      return () => clearTimeout(timer);
    }
    return undefined;
  }, [showSuccessModal, handleSuccessClose]);

  useEffect(() => {
    let mounted = true;
    async function fetchUser() {
      if (!id) {
        setLoadingFetch(false);
        return;
      }
      setLoadingFetch(true);
      try {
        const data = await get<User>(`/users/${id}`);
        if (mounted) {
          setFormData({
            username: String(data.username ?? ""),
            email: String(data.email ?? ""),
            phone: String(data.phone ?? ""),
            role: String(data.role ?? role ?? ""),
            is_active: String(Boolean(data.is_active)),
          });
        }
      } catch (err) {
        console.error("Failed to fetch user for edit:", err);
      } finally {
        if (mounted) setLoadingFetch(false);
      }
    }

    if (isOpen) fetchUser();
    return () => {
      mounted = false;
    };
  }, [id, isOpen, role]);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>,
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setErrorsList([]);
    if (!id) {
      setErrorsList(["Missing user id"]);
      return;
    }
    try {
      await updateUser(id, {
        username: formData.username,
        email: formData.email,
        phone: formData.phone,
        role: formData.role || undefined,
        is_active: formData.is_active === "true",
      });

      setSuccessMsg(`User "${formData.username}" updated successfully.`);
      setShowSuccessModal(true);
    } catch (err) {
      console.error("Update failed:", err);
      setErrorsList(["Failed to update user. Please try again."]);
    }
  };

  const fields: EditFormFieldDef[] = [
    { name: "username", label: "Username", type: "text", required: true },
    { name: "email", label: "Email", type: "email", required: true },
    { name: "phone", label: "Phone", type: "tel", required: true },
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
      name: "is_active",
      label: "Status",
      type: "select",
      options: [
        { value: "true", label: "Active" },
        { value: "false", label: "Inactive" },
      ],
    },
  ];

  if (!isOpen) return null;

  return (
    <>
      <EditFormModal
        isOpen={isOpen}
        formData={formData}
        onChange={handleChange}
        onSubmit={handleSubmit}
        onClose={close}
        errorsList={errorsList}
        loading={updating || loadingFetch}
        fields={fields}
        title="Edit User"
        description="Update user details below"
        submitLabel="Update"
        cancelLabel="Cancel"
      />

      <SuccessModal
        isOpen={showSuccessModal}
        title="User Updated"
        message={successMsg ?? "User updated successfully."}
        buttonLabel="Continue"
        onButtonClick={handleSuccessClose}
      />
    </>
  );
}

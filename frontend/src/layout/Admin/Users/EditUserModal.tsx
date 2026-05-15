import { useEffect, useState, useCallback } from "react";
import { get, patch } from "../../../services/api";
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

  // Reset password (admin) - new password only
  const [showReset, setShowReset] = useState(false);
  const [resetForm, setResetForm] = useState({
    newPassword: "",
    confirmPassword: "",
  });
  const [resetError, setResetError] = useState<string | null>(null);
  const [resetSaving, setResetSaving] = useState(false);
  const [resetSuccessOpen, setResetSuccessOpen] = useState(false);

  const handleResetChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setResetForm((p) => ({ ...p, [name]: value }));
  };

  const handleResetSubmit = async () => {
    setResetError(null);
    if (!id) return setResetError("Missing user id");
    if (!resetForm.newPassword || resetForm.newPassword.length < 8) {
      return setResetError("Password must be at least 8 characters");
    }
    if (resetForm.newPassword !== resetForm.confirmPassword) {
      return setResetError("Passwords do not match");
    }

    setResetSaving(true);
    try {
      await patch(`/users/${id}/reset-password`, {
        new_password: resetForm.newPassword,
      });
      setResetForm({ newPassword: "", confirmPassword: "" });
      setShowReset(false);
      setResetSuccessOpen(true);
    } catch (err) {
      console.error("Reset password failed:", err);
      setResetError(
        err instanceof Error ? err.message : "Failed to reset password",
      );
    } finally {
      setResetSaving(false);
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
      >
        {/* render default form fields so we can append reset-password UI */}
        <>
          {errorsList.length > 0 ? (
            <div className="mb-6 rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-700 dark:border-red-800 dark:bg-red-900/30 dark:text-red-400">
              <ul className="list-inside list-disc space-y-1">
                {errorsList.map((errMsg, idx) => (
                  <li key={idx}>{errMsg}</li>
                ))}
              </ul>
            </div>
          ) : null}

          <div className="grid grid-cols-1 gap-x-6 gap-y-5 sm:grid-cols-2">
            {fields.map((field) => (
              <div key={field.name} className="col-span-1">
                <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                  {field.label}
                  {field.required && <span className="text-red-500"> *</span>}
                </label>
                {field.type === "select" ? (
                  <select
                    name={field.name}
                    value={formData[field.name] || ""}
                    onChange={handleChange}
                    className="w-full rounded-lg border border-gray-300 px-4 py-2 text-sm focus:border-transparent focus:ring-2 focus:ring-brand-500 dark:border-gray-700 dark:bg-gray-900 dark:text-white dark:focus:ring-brand-500"
                    required={field.required}
                    disabled={updating || loadingFetch}
                  >
                    <option value="">Select {field.label.toLowerCase()}</option>
                    {field.options?.map((opt) => (
                      <option key={opt.value} value={opt.value}>
                        {opt.label}
                      </option>
                    ))}
                  </select>
                ) : (
                  <input
                    name={field.name}
                    value={formData[field.name] || ""}
                    onChange={handleChange}
                    type={field.type || "text"}
                    className="w-full rounded-lg border border-gray-300 px-4 py-2 text-sm focus:border-transparent focus:ring-2 focus:ring-brand-500 dark:border-gray-700 dark:bg-gray-900 dark:text-white dark:focus:ring-brand-500"
                    placeholder={field.placeholder}
                    required={field.required}
                    disabled={updating || loadingFetch}
                  />
                )}
              </div>
            ))}
          </div>

          {/* Reset password section (admin) */}
          <div className="mt-6 border-t pt-6">
            <h5 className="mb-3 text-sm font-medium text-gray-800 dark:text-white/90">
              Reset Password
            </h5>
            <p className="mb-4 text-xs text-gray-500">
              Set a new password for this user (no old password required).
            </p>

            {!showReset ? (
              <button
                type="button"
                onClick={() => setShowReset(true)}
                className="inline-flex items-center rounded-lg border border-gray-300 px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-700 dark:text-gray-300"
              >
                Reset Password
              </button>
            ) : (
              <div className="grid grid-cols-1 gap-y-3">
                <div>
                  <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                    New Password
                  </label>
                  <input
                    name="newPassword"
                    type="password"
                    value={resetForm.newPassword}
                    onChange={handleResetChange}
                    className="w-full rounded-lg border border-gray-300 px-4 py-2 text-sm focus:border-transparent focus:ring-2 focus:ring-brand-500 dark:border-gray-700 dark:bg-gray-900 dark:text-white dark:focus:ring-brand-500"
                  />
                </div>
                <div>
                  <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                    Confirm New Password
                  </label>
                  <input
                    name="confirmPassword"
                    type="password"
                    value={resetForm.confirmPassword}
                    onChange={handleResetChange}
                    className="w-full rounded-lg border border-gray-300 px-4 py-2 text-sm focus:border-transparent focus:ring-2 focus:ring-brand-500 dark:border-gray-700 dark:bg-gray-900 dark:text-white dark:focus:ring-brand-500"
                  />
                </div>
                {resetError ? (
                  <p className="text-sm text-red-600">{resetError}</p>
                ) : null}
                <div className="flex items-center gap-3">
                  <button
                    type="button"
                    onClick={() => setShowReset(false)}
                    className="inline-flex items-center rounded-lg border border-gray-300 px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
                  >
                    Cancel
                  </button>
                  <button
                    type="button"
                    onClick={handleResetSubmit}
                    disabled={resetSaving}
                    className="inline-flex items-center rounded-lg bg-brand-500 px-3 py-2 text-sm font-medium text-white hover:bg-brand-600 disabled:opacity-50"
                  >
                    {resetSaving ? "Processing..." : "Confirm Reset"}
                  </button>
                </div>
              </div>
            )}
          </div>

          <div className="mt-6 flex items-center justify-end gap-3">
            <button
              type="button"
              onClick={close}
              disabled={updating || loadingFetch}
              className="inline-flex items-center justify-center rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50 dark:border-gray-700 dark:text-gray-300 dark:hover:bg-white/[0.05]"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={updating || loadingFetch}
              className="inline-flex items-center justify-center rounded-lg bg-brand-500 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-brand-600 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {updating ? `Update...` : `Update`}
            </button>
          </div>
        </>
      </EditFormModal>

      <SuccessModal
        isOpen={showSuccessModal}
        title="User Updated"
        message={successMsg ?? "User updated successfully."}
        buttonLabel="Continue"
        onButtonClick={handleSuccessClose}
      />

      <SuccessModal
        isOpen={resetSuccessOpen}
        title="Password Reset"
        message="User password has been reset."
        buttonLabel="OK"
        onButtonClick={() => setResetSuccessOpen(false)}
      />
    </>
  );
}

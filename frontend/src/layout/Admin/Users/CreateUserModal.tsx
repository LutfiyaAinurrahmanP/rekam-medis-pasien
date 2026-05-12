import { useEffect, useState, useCallback } from "react";
import { post } from "../../../services/api";
import type { User } from "../../../hooks/Users/useUsers";
import SuccessModal from "../../../components/ui/notification/SuccessModal";

interface CreateUserModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
  role?: string;
}

type CreateUserFormState = {
  username: string;
  email: string;
  phone: string;
  password: string;
  confirmPassword: string;
  role: string;
};

export default function CreateUserModal({
  isOpen,
  onClose,
  onSuccess,
  role,
}: CreateUserModalProps) {
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
    role: role || "",
  });

  const handleSuccessModalClose = useCallback(() => {
    setShowSuccessModal(false);
    onClose();
    onSuccess();
  }, [onClose, onSuccess]);

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
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>,
  ) => {
    const { name, value } = e.target;
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
        role: formData.role || undefined,
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
        role: role || "",
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

  const handleClose = () => {
    if (!loading) {
      onClose();
    }
  };

  if (!isOpen) return null;

  return (
    <>
      {/* Backdrop dengan blur dan overlay gelap */}
      <div
        className="fixed inset-0 z-40 bg-black/50 backdrop-blur-sm transition-opacity"
        onClick={handleClose}
      />

      {/* Modal */}
      <div className="fixed left-1/2 top-1/2 z-50 w-full max-w-md -translate-x-1/2 -translate-y-1/2">
        <div className="rounded-lg border border-gray-200 bg-white shadow-xl dark:border-gray-700 dark:bg-gray-800">
          {/* Header */}
          <div className="border-b border-gray-200 px-6 py-4 dark:border-gray-700">
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                  Create User
                </h2>
                <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                  Enter the user details below
                </p>
              </div>
              <button
                onClick={handleClose}
                className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                disabled={loading}
              >
                <svg
                  className="h-6 w-6"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
            </div>
          </div>

          {/* Body */}
          <form onSubmit={handleSubmit} className="space-y-6 p-6">
            {errorsList.length > 0 ? (
              <div className="rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-700 dark:border-red-800 dark:bg-red-900/30 dark:text-red-400">
                <ul className="list-inside list-disc space-y-1">
                  {errorsList.map((errMsg, idx) => (
                    <li key={idx}>{errMsg}</li>
                  ))}
                </ul>
              </div>
            ) : null}

            <div className="space-y-4">
              <div>
                <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                  Username
                </label>
                <input
                  name="username"
                  value={formData.username}
                  onChange={handleChange}
                  type="text"
                  className="w-full rounded-lg border border-gray-300 px-4 py-2 text-sm focus:border-transparent focus:ring-2 focus:ring-brand-500 dark:border-gray-700 dark:bg-gray-900 dark:text-white dark:focus:ring-brand-500"
                  placeholder="Enter username"
                  required
                  disabled={loading}
                />
              </div>

              <div>
                <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                  Email
                </label>
                <input
                  name="email"
                  value={formData.email}
                  onChange={handleChange}
                  type="email"
                  className="w-full rounded-lg border border-gray-300 px-4 py-2 text-sm focus:border-transparent focus:ring-2 focus:ring-brand-500 dark:border-gray-700 dark:bg-gray-900 dark:text-white dark:focus:ring-brand-500"
                  placeholder="Enter email"
                  required
                  disabled={loading}
                />
              </div>

              <div>
                <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                  Phone
                </label>
                <input
                  name="phone"
                  value={formData.phone}
                  onChange={handleChange}
                  type="tel"
                  className="w-full rounded-lg border border-gray-300 px-4 py-2 text-sm focus:border-transparent focus:ring-2 focus:ring-brand-500 dark:border-gray-700 dark:bg-gray-900 dark:text-white dark:focus:ring-brand-500"
                  placeholder="Enter phone number"
                  required
                  disabled={loading}
                />
              </div>

              <div>
                <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                  Role
                </label>
                <select
                  name="role"
                  value={formData.role}
                  onChange={handleChange}
                  className="w-full rounded-lg border border-gray-300 px-4 py-2 text-sm focus:border-transparent focus:ring-2 focus:ring-brand-500 dark:border-gray-700 dark:bg-gray-900 dark:text-white dark:focus:ring-brand-500"
                  required
                  disabled={loading}
                >
                  <option value="">Select role</option>
                  <option value="patient">Patient</option>
                  <option value="doctor">Doctor</option>
                  <option value="receptionist">Receptionist</option>
                  <option value="admin">Admin</option>
                  <option value="super_admin">Super Admin</option>
                </select>
              </div>

              <div>
                <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                  Password
                </label>
                <input
                  name="password"
                  value={formData.password}
                  onChange={handleChange}
                  type="password"
                  className="w-full rounded-lg border border-gray-300 px-4 py-2 text-sm focus:border-transparent focus:ring-2 focus:ring-brand-500 dark:border-gray-700 dark:bg-gray-900 dark:text-white dark:focus:ring-brand-500"
                  placeholder="Enter password"
                  required
                  disabled={loading}
                />
              </div>

              <div>
                <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                  Confirm Password
                </label>
                <input
                  name="confirmPassword"
                  value={formData.confirmPassword}
                  onChange={handleChange}
                  type="password"
                  className="w-full rounded-lg border border-gray-300 px-4 py-2 text-sm focus:border-transparent focus:ring-2 focus:ring-brand-500 dark:border-gray-700 dark:bg-gray-900 dark:text-white dark:focus:ring-brand-500"
                  placeholder="Confirm password"
                  required
                  disabled={loading}
                />
              </div>
            </div>

            {/* Footer */}
            <div className="border-t border-gray-200 pt-6 dark:border-gray-700">
              <div className="flex gap-3">
                <button
                  type="button"
                  onClick={handleClose}
                  disabled={loading}
                  className="flex-1 rounded-lg border border-gray-300 px-4 py-2 font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 dark:border-gray-700 dark:text-gray-300 dark:hover:bg-white/[0.05]"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={loading}
                  className="flex-1 rounded-lg bg-brand-500 px-4 py-2 font-medium text-white hover:bg-brand-600 disabled:opacity-50"
                >
                  {loading ? "Creating..." : "Create"}
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>

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

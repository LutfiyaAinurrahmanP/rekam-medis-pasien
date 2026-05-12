import { useEffect, useState } from "react";
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
  const [error, setError] = useState<string | null>(null);
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

  // Auto-close success modal after 3 seconds
  useEffect(() => {
    if (showSuccessModal) {
      const timer = setTimeout(() => {
        handleSuccessModalClose();
      }, 3000);

      return () => clearTimeout(timer);
    }
  }, [showSuccessModal]);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>,
  ) => {
    const { name, value } = e.target;
    setFormData((previous) => ({
      ...previous,
      [name]: value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError(null);

    if (formData.password !== formData.confirmPassword) {
      setError("Password dan konfirmasi password tidak sama.");
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

      const err = unknownErr as {
        message?: string;
        errors?: Record<string, string>;
      };

      if (err && err.errors && typeof err.errors === "object") {
        const errorsObj = err.errors as Record<string, string>;
        const joined = Object.values(errorsObj).join("\n");
        setError(`${err.message || "Validation error"}: ${joined}`);
        return;
      }

      if (err && err.message && typeof err.message === "string") {
        setError(err.message);
        return;
      }

      setError("Gagal membuat user. Silakan coba lagi.");
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    if (!loading) {
      onClose();
    }
  };

  const handleSuccessModalClose = () => {
    setShowSuccessModal(false);
    onClose();
    onSuccess();
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
            {error ? (
              <div className="rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-700 dark:border-red-800 dark:bg-red-900/30 dark:text-red-400">
                {error}
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

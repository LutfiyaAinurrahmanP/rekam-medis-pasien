import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { useAuth } from "../../context/AuthContext";
import authService from "../../services/auth";
import { Modal } from "../ui/modal";
import Button from "../ui/button/Button";
import Input from "../form/input/InputField";
import Label from "../form/Label";

export default function UserAddressCard() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [actionType, setActionType] = useState<"deactivate" | "delete" | null>(
    null,
  );
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [form, setForm] = useState({
    password: "",
    reason: "",
  });

  useEffect(() => {
    if (!actionType) {
      setError(null);
      setForm({ password: "", reason: "" });
    }
  }, [actionType]);

  const closeModal = () => {
    setActionType(null);
  };

  const handleConfirm = async () => {
    if (!actionType) return;

    setIsSaving(true);
    setError(null);

    try {
      if (actionType === "deactivate") {
        await authService.deactivateMe({
          password: form.password,
          reason: form.reason.trim() || undefined,
        });
      } else {
        await authService.deleteMe({
          password: form.password,
          reason: form.reason.trim() || undefined,
        });
      }

      logout();
      navigate("/auth/login");
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to update account";
      setError(message);
    } finally {
      setIsSaving(false);
    }
  };

  const title =
    actionType === "delete" ? "Delete Account" : "Deactivate Account";
  const description =
    actionType === "delete"
      ? "This will permanently remove your account."
      : "This will deactivate your account until it is reactivated.";

  return (
    <>
      <div className="p-5 border border-gray-200 rounded-2xl dark:border-gray-800 lg:p-6">
        <div className="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between">
          <div>
            <h4 className="text-lg font-semibold text-gray-800 dark:text-white/90 lg:mb-6">
              Account Actions
            </h4>

            <div className="grid grid-cols-1 gap-4 lg:grid-cols-2 lg:gap-7 2xl:gap-x-32">
              <div>
                <p className="mb-2 text-xs leading-normal text-gray-500 dark:text-gray-400">
                  Username
                </p>
                <p className="text-sm font-medium text-gray-800 dark:text-white/90">
                  {user?.username || "-"}
                </p>
              </div>

              <div>
                <p className="mb-2 text-xs leading-normal text-gray-500 dark:text-gray-400">
                  Role
                </p>
                <p className="text-sm font-medium text-gray-800 dark:text-white/90">
                  {user?.role || "-"}
                </p>
              </div>

              <div>
                <p className="mb-2 text-xs leading-normal text-gray-500 dark:text-gray-400">
                  Status
                </p>
                <p className="text-sm font-medium text-gray-800 dark:text-white/90">
                  {user?.is_active ? "Active" : "Inactive"}
                </p>
              </div>

              <div>
                <p className="mb-2 text-xs leading-normal text-gray-500 dark:text-gray-400">
                  Email
                </p>
                <p className="text-sm font-medium text-gray-800 dark:text-white/90">
                  {user?.email || "-"}
                </p>
              </div>
            </div>
          </div>

          <div className="flex flex-col gap-3 w-full lg:w-auto">
            <button
              onClick={() => setActionType("deactivate")}
              className="flex w-full items-center justify-center gap-2 rounded-full border border-warning-300 bg-warning-50 px-4 py-3 text-sm font-medium text-warning-700 shadow-theme-xs hover:bg-warning-100 dark:border-warning-700 dark:bg-warning-500/10 dark:text-warning-300 lg:inline-flex lg:w-auto"
            >
              Deactivate Account
            </button>
            <button
              onClick={() => setActionType("delete")}
              className="flex w-full items-center justify-center gap-2 rounded-full border border-red-300 bg-red-50 px-4 py-3 text-sm font-medium text-red-700 shadow-theme-xs hover:bg-red-100 dark:border-red-700 dark:bg-red-500/10 dark:text-red-300 lg:inline-flex lg:w-auto"
            >
              Delete Account
            </button>
          </div>
        </div>
      </div>

      <Modal
        isOpen={Boolean(actionType)}
        onClose={closeModal}
        className="max-w-[700px] m-4"
      >
        <div className="relative w-full p-4 overflow-y-auto bg-white no-scrollbar rounded-3xl dark:bg-gray-900 lg:p-11">
          <div className="px-2 pr-14">
            <h4 className="mb-2 text-2xl font-semibold text-gray-800 dark:text-white/90">
              {title}
            </h4>
            <p className="mb-6 text-sm text-gray-500 dark:text-gray-400 lg:mb-7">
              {description}
            </p>
          </div>

          <div className="flex flex-col gap-7 px-2">
            <div>
              <h5 className="mb-5 text-lg font-medium text-gray-800 dark:text-white/90 lg:mb-6">
                Confirmation
              </h5>

              <div className="grid grid-cols-1 gap-x-6 gap-y-5 lg:grid-cols-2">
                <div className="col-span-2 lg:col-span-1">
                  <Label>Password</Label>
                  <Input
                    type="password"
                    value={form.password}
                    onChange={(e) =>
                      setForm({ ...form, password: e.target.value })
                    }
                  />
                </div>

                <div className="col-span-2">
                  <Label>Reason</Label>
                  <Input
                    type="text"
                    value={form.reason}
                    onChange={(e) =>
                      setForm({ ...form, reason: e.target.value })
                    }
                    placeholder="Optional reason"
                  />
                </div>
              </div>
            </div>
          </div>

          {error ? (
            <p className="px-2 mt-5 text-sm text-red-600">{error}</p>
          ) : null}

          <div className="flex items-center gap-3 px-2 mt-6 lg:justify-end">
            <Button size="sm" variant="outline" onClick={closeModal}>
              Close
            </Button>
            <Button size="sm" onClick={handleConfirm} disabled={isSaving}>
              {isSaving ? "Processing..." : title}
            </Button>
          </div>
        </div>
      </Modal>
    </>
  );
}

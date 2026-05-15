import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { useAuth } from "../../context/AuthContext";
import authService from "../../services/auth";
import { Modal } from "../ui/modal";
import Button from "../ui/button/Button";
import Input from "../form/input/InputField";
import Label from "../form/Label";
import SuccessModal from "../ui/notification/SuccessModal";

export default function UserInfoCard() {
  const { user, setUser, logout } = useAuth();
  const navigate = useNavigate();

  // Edit profile state
  const [editOpen, setEditOpen] = useState(false);
  const [editSaving, setEditSaving] = useState(false);
  const [editError, setEditError] = useState<string | null>(null);
  const [editForm, setEditForm] = useState({
    username: "",
    email: "",
    phone: "",
  });

  // Change password state
  const [pwOpen, setPwOpen] = useState(false);
  const [pwSaving, setPwSaving] = useState(false);
  const [pwError, setPwError] = useState<string | null>(null);
  const [pwForm, setPwForm] = useState({
    oldPassword: "",
    newPassword: "",
    confirmPassword: "",
  });

  // Deactivate / Delete state
  const [actionType, setActionType] = useState<"deactivate" | "delete" | null>(
    null,
  );
  const [actionSaving, setActionSaving] = useState(false);
  const [actionError, setActionError] = useState<string | null>(null);
  const [actionForm, setActionForm] = useState({ password: "" });

  // Shared success modal
  const [successOpen, setSuccessOpen] = useState(false);
  const [successTitle, setSuccessTitle] = useState("");
  const [successMessage, setSuccessMessage] = useState("");
  const [successCallback, setSuccessCallback] = useState<(() => void) | null>(
    null,
  );

  useEffect(() => {
    if (user) {
      setEditForm({
        username: user.username ?? "",
        email: user.email ?? "",
        phone: user.phone ?? "",
      });
    }
  }, [user]);

  useEffect(() => {
    if (!pwOpen) {
      setPwError(null);
      setPwForm({ oldPassword: "", newPassword: "", confirmPassword: "" });
    }
  }, [pwOpen]);

  useEffect(() => {
    if (!actionType) {
      setActionError(null);
      setActionForm({ password: "" });
    }
  }, [actionType]);

  useEffect(() => {
    if (!successOpen) return;
    const timer = window.setTimeout(() => {
      handleSuccessClose();
    }, 3000);

    return () => window.clearTimeout(timer);
  }, [successOpen]);

  const showSuccess = (title: string, message: string, cb?: () => void) => {
    setSuccessTitle(title);
    setSuccessMessage(message);
    setSuccessCallback(() => cb ?? null);
    setSuccessOpen(true);
  };

  const handleSuccessClose = () => {
    setSuccessOpen(false);
    if (successCallback) successCallback();
    setSuccessCallback(null);
  };

  const handleSaveProfile = async () => {
    setEditSaving(true);
    setEditError(null);
    try {
      const updated = await authService.updateMe({
        username: editForm.username.trim() || undefined,
        email: editForm.email.trim() || undefined,
        phone: editForm.phone.trim() || undefined,
      });
      setUser(updated);
      showSuccess("Profile Updated", "Your profile has been updated.", () =>
        setEditOpen(false),
      );
    } catch (err) {
      setEditError(
        err instanceof Error ? err.message : "Failed to update profile",
      );
    } finally {
      setEditSaving(false);
    }
  };

  const handleChangePassword = async () => {
    if (pwForm.newPassword !== pwForm.confirmPassword) {
      setPwError("New password and confirmation do not match.");
      return;
    }
    setPwSaving(true);
    setPwError(null);
    try {
      await authService.changeMyPassword({
        old_password: pwForm.oldPassword,
        new_password: pwForm.newPassword,
      });
      showSuccess("Password Changed", "Your password has been updated.", () =>
        setPwOpen(false),
      );
    } catch (err) {
      setPwError(
        err instanceof Error ? err.message : "Failed to change password",
      );
    } finally {
      setPwSaving(false);
    }
  };

  const handleConfirmAction = async () => {
    if (!actionType) return;
    setActionSaving(true);
    setActionError(null);
    try {
      if (actionType === "deactivate") {
        await authService.deactivateMe({ password: actionForm.password });
      } else {
        await authService.deleteMe({ password: actionForm.password });
      }

      const title =
        actionType === "delete" ? "Account Deleted" : "Account Deactivated";
      const message =
        actionType === "delete"
          ? "Your account has been permanently deleted."
          : "Your account has been deactivated.";

      showSuccess(title, message, () => {
        logout();
        navigate("/auth/login");
      });
    } catch (err) {
      setActionError(
        err instanceof Error ? err.message : "Failed to update account",
      );
    } finally {
      setActionSaving(false);
    }
  };

  return (
    <div className="p-5 border border-gray-200 rounded-2xl dark:border-gray-800 lg:p-6">
      <div className="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between">
        <div>
          <h4 className="text-lg font-semibold text-gray-800 dark:text-white/90 lg:mb-6">
            Personal Information
          </h4>

          <div className="grid grid-cols-1 gap-4 lg:gap-7 2xl:gap-x-32">
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
                Email address
              </p>
              <p className="text-sm font-medium text-gray-800 dark:text-white/90">
                {user?.email || "-"}
              </p>
            </div>

            <div>
              <p className="mb-2 text-xs leading-normal text-gray-500 dark:text-gray-400">
                Phone
              </p>
              <p className="text-sm font-medium text-gray-800 dark:text-white/90">
                {user?.phone || "-"}
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
          </div>
        </div>

        <div className="flex flex-col gap-3 w-full lg:w-auto">
          <button
            onClick={() => setEditOpen(true)}
            className="flex w-full items-center justify-center gap-2 rounded-full border border-gray-300 bg-white px-4 py-3 text-sm font-medium text-gray-700 shadow-theme-xs hover:bg-gray-50 hover:text-gray-800 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-white/[0.03] dark:hover:text-gray-200 lg:inline-flex lg:w-auto"
          >
            Edit Profile
          </button>
          <button
            onClick={() => setPwOpen(true)}
            className="flex w-full items-center justify-center gap-2 rounded-full border border-gray-300 bg-white px-4 py-3 text-sm font-medium text-gray-700 shadow-theme-xs hover:bg-gray-50 hover:text-gray-800 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-white/[0.03] dark:hover:text-gray-200 lg:inline-flex lg:w-auto"
          >
            Edit Password
          </button>
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

      {/* Edit Profile Modal */}
      <Modal
        isOpen={editOpen}
        onClose={() => setEditOpen(false)}
        className="max-w-[700px] m-4"
      >
        <div className="no-scrollbar relative w-full max-w-[700px] overflow-y-auto rounded-3xl bg-white p-4 dark:bg-gray-900 lg:p-11">
          <div className="px-2 pr-14">
            <h4 className="mb-2 text-2xl font-semibold text-gray-800 dark:text-white/90">
              Edit Profile
            </h4>
            <p className="mb-6 text-sm text-gray-500 dark:text-gray-400 lg:mb-7">
              Update your account information.
            </p>
          </div>

          <div className="flex flex-col gap-7 px-2">
            <div>
              <h5 className="mb-5 text-lg font-medium text-gray-800 dark:text-white/90 lg:mb-6">
                Account Information
              </h5>

              <div className="grid grid-cols-1 gap-y-5">
                <div>
                  <Label>Username</Label>
                  <Input
                    type="text"
                    value={editForm.username}
                    onChange={(e) =>
                      setEditForm({ ...editForm, username: e.target.value })
                    }
                  />
                </div>

                <div>
                  <Label>Email Address</Label>
                  <Input
                    type="email"
                    value={editForm.email}
                    onChange={(e) =>
                      setEditForm({ ...editForm, email: e.target.value })
                    }
                  />
                </div>

                <div>
                  <Label>Phone</Label>
                  <Input
                    type="text"
                    value={editForm.phone}
                    onChange={(e) =>
                      setEditForm({ ...editForm, phone: e.target.value })
                    }
                  />
                </div>
              </div>
            </div>
          </div>

          {editError ? (
            <p className="px-2 mt-5 text-sm text-red-600">{editError}</p>
          ) : null}

          <div className="flex items-center gap-3 px-2 mt-6 lg:justify-end">
            <Button
              size="sm"
              variant="outline"
              onClick={() => setEditOpen(false)}
            >
              Close
            </Button>
            <Button size="sm" onClick={handleSaveProfile} disabled={editSaving}>
              {editSaving ? "Saving..." : "Save Changes"}
            </Button>
          </div>
        </div>
      </Modal>

      {/* Change Password Modal */}
      <Modal
        isOpen={pwOpen}
        onClose={() => setPwOpen(false)}
        className="max-w-[700px] m-4"
      >
        <div className="relative w-full p-4 overflow-y-auto bg-white no-scrollbar rounded-3xl dark:bg-gray-900 lg:p-11">
          <div className="px-2 pr-14">
            <h4 className="mb-2 text-2xl font-semibold text-gray-800 dark:text-white/90">
              Change Password
            </h4>
            <p className="mb-6 text-sm text-gray-500 dark:text-gray-400 lg:mb-7">
              Update your account password.
            </p>
          </div>

          <div className="flex flex-col gap-7 px-2">
            <div>
              <h5 className="mb-5 text-lg font-medium text-gray-800 dark:text-white/90 lg:mb-6">
                Security Details
              </h5>

              <div className="grid grid-cols-1 gap-y-5">
                <div>
                  <Label>Old Password</Label>
                  <Input
                    type="password"
                    value={pwForm.oldPassword}
                    onChange={(e) =>
                      setPwForm({ ...pwForm, oldPassword: e.target.value })
                    }
                  />
                </div>

                <div>
                  <Label>New Password</Label>
                  <Input
                    type="password"
                    value={pwForm.newPassword}
                    onChange={(e) =>
                      setPwForm({ ...pwForm, newPassword: e.target.value })
                    }
                  />
                </div>

                <div>
                  <Label>Confirm New Password</Label>
                  <Input
                    type="password"
                    value={pwForm.confirmPassword}
                    onChange={(e) =>
                      setPwForm({ ...pwForm, confirmPassword: e.target.value })
                    }
                  />
                </div>
              </div>
            </div>
          </div>

          {pwError ? (
            <p className="px-2 mt-5 text-sm text-red-600">{pwError}</p>
          ) : null}

          <div className="flex items-center gap-3 px-2 mt-6 lg:justify-end">
            <Button
              size="sm"
              variant="outline"
              onClick={() => setPwOpen(false)}
            >
              Close
            </Button>
            <Button
              size="sm"
              onClick={handleChangePassword}
              disabled={pwSaving}
            >
              {pwSaving ? "Saving..." : "Change Password"}
            </Button>
          </div>
        </div>
      </Modal>

      {/* Deactivate / Delete Modal */}
      <Modal
        isOpen={!!actionType}
        onClose={() => setActionType(null)}
        className="max-w-[600px] m-4"
      >
        <div className="no-scrollbar relative w-full max-w-[600px] overflow-y-auto rounded-3xl bg-white p-6 dark:bg-gray-900">
          <div>
            <h4 className="mb-2 text-2xl font-semibold text-gray-800 dark:text-white/90">
              {actionType === "delete"
                ? "Delete Account"
                : "Deactivate Account"}
            </h4>
            <p className="mb-6 text-sm text-gray-500 dark:text-gray-400">
              {actionType === "delete"
                ? "This will permanently delete your account."
                : "This will deactivate your account. You can reactivate by contacting support."}
            </p>
          </div>

          <div className="grid grid-cols-1 gap-y-5">
            <div>
              <Label>Password</Label>
              <Input
                type="password"
                value={actionForm.password}
                onChange={(e) =>
                  setActionForm({ ...actionForm, password: e.target.value })
                }
              />
            </div>
          </div>

          {actionError ? (
            <p className="px-2 mt-5 text-sm text-red-600">{actionError}</p>
          ) : null}

          <div className="flex items-center gap-3 px-2 mt-6 lg:justify-end">
            <Button
              size="sm"
              variant="outline"
              onClick={() => setActionType(null)}
            >
              Cancel
            </Button>
            <Button
              size="sm"
              onClick={handleConfirmAction}
              disabled={actionSaving}
            >
              {actionSaving
                ? "Processing..."
                : actionType === "delete"
                  ? "Delete"
                  : "Deactivate"}
            </Button>
          </div>
        </div>
      </Modal>

      <SuccessModal
        isOpen={successOpen}
        title={successTitle}
        message={successMessage}
        buttonLabel="Close"
        onButtonClick={handleSuccessClose}
      />
    </div>
  );
}

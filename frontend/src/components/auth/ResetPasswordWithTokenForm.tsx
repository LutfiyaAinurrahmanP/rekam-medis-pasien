import { useState, useEffect } from "react";
import { Link, useSearchParams } from "react-router";
import Label from "../form/Label";
import Input from "../form/input/InputField";

export default function ResetPasswordWithTokenForm() {
  const [searchParams] = useSearchParams();
  const [resetToken, setResetToken] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");
  const [isSuccess, setIsSuccess] = useState(false);

  useEffect(() => {
    const token = searchParams.get("token");
    if (token) {
      setResetToken(token);
    }
  }, [searchParams]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!resetToken || !newPassword) {
      setMessage("Please fill in all fields");
      setIsSuccess(false);
      return;
    }

    if (newPassword !== confirmPassword) {
      setMessage("Passwords do not match");
      setIsSuccess(false);
      return;
    }

    if (newPassword.length < 8) {
      setMessage("Password must be at least 8 characters long");
      setIsSuccess(false);
      return;
    }

    setLoading(true);
    setMessage("");

    try {
      const response = await fetch(
        "http://localhost:8080/api/v1/auth/reset-password",
        {
          method: "PATCH",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            reset_token: resetToken,
            new_password: newPassword,
          }),
        },
      );

      const data = await response.json();

      if (data.success) {
        setMessage(data.message);
        setIsSuccess(true);
      } else {
        setMessage(data.message || "Failed to reset password");
        setIsSuccess(false);
      }
    } catch (error) {
      setMessage("Network error. Please try again.");
      setIsSuccess(false);
    } finally {
      setLoading(false);
    }
  };

  if (isSuccess) {
    return (
      <div className="flex flex-col flex-1 w-full overflow-y-auto no-scrollbar">
        <div className="flex flex-col justify-center flex-1 w-full max-w-md mx-auto">
          <div>
            <Link to="/" className="flex justify-center mb-4">
              <img
                width={80}
                height={18}
                src="/images/logo/app-logo.png"
                alt="Logo"
              />
            </Link>
            <div className="mb-5 sm:mb-8">
              <h1 className="mb-2 font-semibold text-gray-800 text-title-sm dark:text-white/90 sm:text-title-md">
                Password Reset Successful
              </h1>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                Your password has been successfully reset. You can now log in
                with your new password.
              </p>
            </div>
            <div className="bg-green-100 text-green-800 p-3 rounded-lg mb-5">
              {message}
            </div>
            <div className="space-y-4">
              <Link
                to="/auth/login"
                className="flex items-center justify-center w-full px-4 py-3 text-sm font-medium text-white transition rounded-lg bg-brand-500 shadow-theme-xs hover:bg-brand-600"
              >
                Go to Login
              </Link>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col flex-1 w-full overflow-y-auto no-scrollbar">
      <div className="flex flex-col justify-center flex-1 w-full max-w-md mx-auto">
        <div>
          <Link to="/" className="flex justify-center mb-4">
            <img
              width={80}
              height={18}
              src="/images/logo/app-logo.png"
              alt="Logo"
            />
          </Link>
          <div className="mb-5 sm:mb-8">
            <h1 className="mb-2 font-semibold text-gray-800 text-title-sm dark:text-white/90 sm:text-title-md">
              Reset Your Password
            </h1>
            <p className="text-sm text-gray-500 dark:text-gray-400">
              Enter your new password below.
            </p>
          </div>
          <form onSubmit={handleSubmit}>
            <div className="space-y-5">
              <div>
                <Label>
                  New Password<span className="text-error-500">*</span>
                </Label>
                <Input
                  type="password"
                  id="newPassword"
                  name="newPassword"
                  placeholder="Enter your new password"
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                />
              </div>

              <div>
                <Label>
                  Confirm New Password<span className="text-error-500">*</span>
                </Label>
                <Input
                  type="password"
                  id="confirmPassword"
                  name="confirmPassword"
                  placeholder="Confirm your new password"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                />
              </div>

              {message && (
                <div
                  className={`p-3 rounded-lg ${isSuccess ? "bg-green-100 text-green-800" : "bg-red-100 text-red-800"}`}
                >
                  {message}
                </div>
              )}

              <div>
                <button
                  type="submit"
                  disabled={loading}
                  className="flex items-center justify-center w-full px-4 py-3 text-sm font-medium text-white transition rounded-lg bg-brand-500 shadow-theme-xs hover:bg-brand-600 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {loading ? "Resetting..." : "Reset Password"}
                </button>
              </div>
            </div>
          </form>
          <div className="mt-5">
            <p className="text-sm font-normal text-center text-gray-700 dark:text-gray-400 sm:text-start">
              Remember your password?{" "}
              <Link
                to="/auth/login"
                className="text-brand-500 hover:text-brand-600 dark:text-brand-400"
              >
                Back to Login
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

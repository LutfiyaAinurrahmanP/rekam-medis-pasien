import { useState, useEffect } from "react";
import { Link, useSearchParams } from "react-router";
import Label from "../form/Label";
import Input from "../form/input/InputField";

export default function VerifyResetCodeForm() {
  const [searchParams] = useSearchParams();
  const [email, setEmail] = useState("");
  const [resetCode, setResetCode] = useState("");
  const [resetToken, setResetToken] = useState("");
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");
  const [isSuccess, setIsSuccess] = useState(false);
  const [resendLoading, setResendLoading] = useState(false);
  const [resendMessage, setResendMessage] = useState("");

  useEffect(() => {
    const emailParam = searchParams.get("email");
    if (emailParam) {
      setEmail(decodeURIComponent(emailParam));
    }
  }, [searchParams]);

  const handleResendCode = async () => {
    if (!email) {
      setResendMessage("Email is required to resend code");
      return;
    }

    setResendLoading(true);
    setResendMessage("");

    try {
      const response = await fetch(
        "http://localhost:8080/api/v1/auth/resend-reset-code",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ email }),
        },
      );

      const data = await response.json();

      if (data.success) {
        setResendMessage("A new reset code has been sent to your email");
      } else {
        setResendMessage(data.message || "Failed to resend code");
      }
    } catch (error) {
      setResendMessage("Network error. Please try again.");
    } finally {
      setResendLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email || !resetCode) {
      setMessage("Please enter both email and reset code");
      setIsSuccess(false);
      return;
    }

    setLoading(true);
    setMessage("");

    try {
      const response = await fetch(
        "http://localhost:8080/api/v1/auth/verify-reset-code",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ email, reset_code: resetCode }),
        },
      );

      const data = await response.json();

      if (data.success) {
        setMessage(data.message);
        setIsSuccess(true);
        setResetToken(data.data.reset_token);
      } else {
        setMessage(data.message || "Invalid reset code");
        setIsSuccess(false);
      }
    } catch (error) {
      setMessage("Network error. Please try again.");
      setIsSuccess(false);
    } finally {
      setLoading(false);
    }
  };

  if (isSuccess && resetToken) {
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
                Reset Code Verified
              </h1>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                Your reset code has been verified. You can now reset your
                password.
              </p>
            </div>
            <div className="bg-green-100 text-green-800 p-3 rounded-lg mb-5">
              {message}
            </div>
            <div className="space-y-4">
              <Link
                to={`/auth/reset-password-token?token=${resetToken}`}
                className="flex items-center justify-center w-full px-4 py-3 text-sm font-medium text-white transition rounded-lg bg-brand-500 shadow-theme-xs hover:bg-brand-600"
              >
                Continue to Reset Password
              </Link>
              <Link
                to="/auth/reset-password"
                className="flex items-center justify-center w-full px-4 py-3 text-sm font-medium text-gray-700 transition rounded-lg border border-gray-300 hover:bg-gray-50"
              >
                Back to Forgot Password
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
              Verify Reset Code
            </h1>
            <p className="text-sm text-gray-500 dark:text-gray-400">
              Enter the reset code sent to your email and your email address.
            </p>
          </div>
          <form onSubmit={handleSubmit}>
            <div className="space-y-5">
              {resendMessage && (
                <div className="p-3 rounded-lg bg-blue-100 text-blue-800 text-sm">
                  {resendMessage}
                </div>
              )}

              {message && !isSuccess && (
                <div className="p-3 rounded-lg bg-red-100 text-red-800">
                  {message}
                </div>
              )}

              {!email && (
                <div>
                  <Label>
                    Email<span className="text-error-500">*</span>
                  </Label>
                  <Input
                    type="email"
                    id="email"
                    name="email"
                    placeholder="Enter your email address"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                  />
                </div>
              )}

              {email && (
                <div>
                  <Label>Email</Label>
                  <div className="p-3 bg-gray-100 rounded-lg text-gray-700">
                    {email}
                  </div>
                </div>
              )}

              <div>
                <Label>
                  Reset Code<span className="text-error-500">*</span>
                </Label>
                <Input
                  type="text"
                  id="resetCode"
                  name="resetCode"
                  placeholder="Enter the 6-digit reset code"
                  value={resetCode}
                  onChange={(e) => setResetCode(e.target.value)}
                />
              </div>

              <div>
                <button
                  type="submit"
                  disabled={loading}
                  className="flex items-center justify-center w-full px-4 py-3 text-sm font-medium text-white transition rounded-lg bg-brand-500 shadow-theme-xs hover:bg-brand-600 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {loading ? "Verifying..." : "Verify Code"}
                </button>
              </div>
            </div>
          </form>
          <div className="mt-5">
            <p className="text-sm font-normal text-center text-gray-700 dark:text-gray-400 sm:text-start">
              Didn't receive the code?{" "}
              <button
                type="button"
                onClick={handleResendCode}
                disabled={resendLoading}
                className="text-brand-500 hover:text-brand-600 dark:text-brand-400 disabled:opacity-50 disabled:cursor-not-allowed text-sm font-medium"
              >
                {resendLoading ? "Sending..." : "Resend Code"}
              </button>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

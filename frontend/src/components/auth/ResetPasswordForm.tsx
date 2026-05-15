import { useState } from "react";
import { Link, useNavigate } from "react-router";
import Label from "../form/Label";
import Input from "../form/input/InputField";

export default function ResetPasswordForm() {
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");
  const [isSuccess, setIsSuccess] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email) {
      setMessage("Please enter your email address");
      setIsSuccess(false);
      return;
    }

    setLoading(true);
    setMessage("");

    try {
      const response = await fetch(
        "http://localhost:8080/api/v1/auth/forgot-password",
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
        setMessage(data.message);
        setIsSuccess(true);
        // Redirect to verify reset code page after 2 seconds with email param
        setTimeout(() => {
          navigate(
            `/auth/verify-reset-code?email=${encodeURIComponent(email)}`,
          );
        }, 2000);
      } else {
        setMessage(data.message || "An error occurred");
        setIsSuccess(false);
      }
    } catch (error) {
      setMessage("Network error. Please try again.");
      setIsSuccess(false);
    } finally {
      setLoading(false);
    }
  };

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
              Forgot Your Password?
            </h1>
            <p className="text-sm text-gray-500 dark:text-gray-400">
              Enter the email address linked to your account, and we will send
              you a link to reset your password.
            </p>
          </div>
          <form onSubmit={handleSubmit}>
            <div className="space-y-5">
              {message && (
                <div
                  className={`p-3 rounded-lg ${isSuccess ? "bg-green-100 text-green-800" : "bg-red-100 text-red-800"}`}
                >
                  {message}
                  {isSuccess && (
                    <div className="mt-2 text-sm">
                      Redirecting to verification page...
                    </div>
                  )}
                </div>
              )}

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

              <div>
                <button
                  type="submit"
                  disabled={loading}
                  className="flex items-center justify-center w-full px-4 py-3 text-sm font-medium text-white transition rounded-lg bg-brand-500 shadow-theme-xs hover:bg-brand-600 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {loading ? "Sending..." : "Send Reset Link"}
                </button>
              </div>
            </div>
          </form>
          <div className="mt-5">
            <p className="text-sm font-normal text-center text-gray-700 dark:text-gray-400 sm:text-start">
              Wait, I remember my password...{" "}
              <Link
                to="/auth/login"
                className="text-brand-500 hover:text-brand-600 dark:text-brand-400"
              >
                Click here
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

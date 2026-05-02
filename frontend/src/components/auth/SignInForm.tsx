import { useEffect, useState, type FormEvent } from "react";
import { Link, useNavigate } from "react-router";
import { EyeCloseIcon, EyeIcon } from "../../icons";
import Label from "../form/Label";
import Input from "../form/input/InputField";
import Checkbox from "../form/input/Checkbox";
import SuccessModal from "../ui/notification/SuccessModal";
import authService from "../../services/auth";

export default function SignInForm() {
  const navigate = useNavigate();
  const [showPassword, setShowPassword] = useState(false);
  const [isChecked, setIsChecked] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showSuccessModal, setShowSuccessModal] = useState(false);

  const [formData, setFormData] = useState({
    usernameOrEmail: "",
    password: "",
  });

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
    if (error) setError(null);
  };

  const handleSuccessModalClick = () => {
    setShowSuccessModal(false);
    navigate("/dashboard");
  };

  useEffect(() => {
    if (!showSuccessModal) return;

    const timer = setTimeout(() => {
      navigate("/dashboard");
    }, 3000);

    return () => clearTimeout(timer);
  }, [showSuccessModal, navigate]);

  const validateForm = (): boolean => {
    if (!formData.usernameOrEmail.trim()) {
      setError("Username atau email tidak boleh kosong");
      return false;
    }

    if (!formData.password.trim()) {
      setError("Password tidak boleh kosong");
      return false;
    }

    return true;
  };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    setError(null);

    if (!validateForm()) {
      return;
    }

    setIsLoading(true);

    try {
      const response = await authService.login({
        username_or_email: formData.usernameOrEmail,
        password: formData.password,
      });

      authService.setToken(response.token);
      setShowSuccessModal(true);
    } catch (err: unknown) {
      console.error("Login error:", err);

      if (typeof err === "object" && err !== null) {
        const errorObj = err as {
          status?: number;
          message?: string;
          error?: unknown;
          errors?: Record<string, string>;
        };

        const backendDetail =
          typeof errorObj.error === "string"
            ? errorObj.error
            : errorObj.error && typeof errorObj.error === "object"
              ? Object.values(errorObj.error as Record<string, unknown>)
                  .filter((value): value is string => typeof value === "string")
                  .join(" ")
              : "";

        if (backendDetail) {
          setError(backendDetail);
        } else if (errorObj.message) {
          setError(errorObj.message);
        } else {
          setError("Terjadi kesalahan saat masuk. Silahkan coba lagi.");
        }
      } else {
        setError("Terjadi kesalahan saat masuk. Silahkan coba lagi.");
      }
    } finally {
      setIsLoading(false);
    }
  };
  return (
    <>
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
                Masuk
              </h1>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                Harap masukkan username atau email dan password untuk masuk!
              </p>
            </div>
            <div>
              <div className="grid grid-cols-1 gap-3 sm:gap-5">
                <button className="inline-flex items-center justify-center gap-3 py-3 text-sm font-normal text-gray-700 transition-colors bg-gray-100 rounded-lg px-7 hover:bg-gray-200 hover:text-gray-800 dark:bg-white/5 dark:text-white/90 dark:hover:bg-white/10">
                  <svg
                    width="20"
                    height="20"
                    viewBox="0 0 20 20"
                    fill="none"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      d="M18.7511 10.1944C18.7511 9.47495 18.6915 8.94995 18.5626 8.40552H10.1797V11.6527H15.1003C15.0011 12.4597 14.4654 13.675 13.2749 14.4916L13.2582 14.6003L15.9087 16.6126L16.0924 16.6305C17.7788 15.1041 18.7511 12.8583 18.7511 10.1944Z"
                      fill="#4285F4"
                    />
                    <path
                      d="M10.1788 18.75C12.5895 18.75 14.6133 17.9722 16.0915 16.6305L13.274 14.4916C12.5201 15.0068 11.5081 15.3666 10.1788 15.3666C7.81773 15.3666 5.81379 13.8402 5.09944 11.7305L4.99473 11.7392L2.23868 13.8295L2.20264 13.9277C3.67087 16.786 6.68674 18.75 10.1788 18.75Z"
                      fill="#34A853"
                    />
                    <path
                      d="M5.10014 11.7305C4.91165 11.186 4.80257 10.6027 4.80257 9.99992C4.80257 9.3971 4.91165 8.81379 5.09022 8.26935L5.08523 8.1534L2.29464 6.02954L2.20333 6.0721C1.5982 7.25823 1.25098 8.5902 1.25098 9.99992C1.25098 11.4096 1.5982 12.7415 2.20333 13.9277L5.10014 11.7305Z"
                      fill="#FBBC05"
                    />
                    <path
                      d="M10.1789 4.63331C11.8554 4.63331 12.9864 5.34303 13.6312 5.93612L16.1511 3.525C14.6035 2.11528 12.5895 1.25 10.1789 1.25C6.68676 1.25 3.67088 3.21387 2.20264 6.07218L5.08953 8.26943C5.81381 6.15972 7.81776 4.63331 10.1789 4.63331Z"
                      fill="#EB4335"
                    />
                  </svg>
                  Masuk dengan Google
                </button>
              </div>
              <div className="relative py-3 sm:py-5">
                <div className="absolute inset-0 flex items-center">
                  <div className="w-full border-t border-gray-200 dark:border-gray-800"></div>
                </div>
                <div className="relative flex justify-center text-sm">
                  <span className="p-2 text-gray-400 bg-white dark:bg-gray-900 sm:px-5 sm:py-2">
                    Atau
                  </span>
                </div>
              </div>
              <form onSubmit={handleSubmit}>
                <div className="space-y-5">
                  {error && (
                    <div className="p-3 text-sm text-red-700 bg-red-100 border border-red-200 rounded-lg dark:bg-red-900/30 dark:text-red-400 dark:border-red-800">
                      {error}
                    </div>
                  )}
                  <div>
                    <Label>
                      Username atau Email
                      <span className="text-error-500">*</span>
                    </Label>
                    <Input
                      type="text"
                      name="usernameOrEmail"
                      value={formData.usernameOrEmail}
                      onChange={handleInputChange}
                      placeholder="Masukkan username atau email anda"
                      disabled={isLoading}
                    />
                  </div>
                  <div>
                    <Label>
                      Password<span className="text-error-500">*</span>
                    </Label>
                    <div className="relative">
                      <Input
                        type={showPassword ? "text" : "password"}
                        name="password"
                        value={formData.password}
                        onChange={handleInputChange}
                        placeholder="Masukkan password anda"
                        disabled={isLoading}
                      />
                      <span
                        onClick={() =>
                          !isLoading && setShowPassword(!showPassword)
                        }
                        className={`absolute z-30 -translate-y-1/2 cursor-pointer right-4 top-1/2 ${
                          isLoading ? "cursor-not-allowed opacity-50" : ""
                        }`}
                      >
                        {showPassword ? (
                          <EyeIcon className="fill-gray-500 dark:fill-gray-400 size-5" />
                        ) : (
                          <EyeCloseIcon className="fill-gray-500 dark:fill-gray-400 size-5" />
                        )}
                      </span>
                    </div>
                  </div>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <Checkbox
                        checked={isChecked}
                        onChange={setIsChecked}
                        disabled={isLoading}
                      />
                      <span className="block font-normal text-gray-700 text-theme-sm dark:text-gray-400">
                        Ingat saya
                      </span>
                    </div>
                    <Link
                      to="/auth/reset-password"
                      className="text-sm text-brand-500 hover:text-brand-600 dark:text-brand-400"
                    >
                      Lupa password?
                    </Link>
                  </div>
                  <div>
                    <button
                      type="submit"
                      disabled={isLoading}
                      className={`flex items-center justify-center w-full px-4 py-3 text-sm font-medium text-white transition rounded-lg shadow-theme-xs ${
                        isLoading
                          ? "bg-gray-400 cursor-not-allowed"
                          : "bg-brand-500 hover:bg-brand-600"
                      }`}
                    >
                      {isLoading ? "Memproses..." : "Masuk"}
                    </button>
                  </div>
                </div>
              </form>

              <div className="mt-5">
                <p className="text-sm font-normal text-center text-gray-700 dark:text-gray-400 sm:text-start">
                  Belum memiliki akun? {""}
                  <Link
                    to="/auth/register"
                    className="text-brand-500 hover:text-brand-600 dark:text-brand-400"
                  >
                    Daftar
                  </Link>
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <SuccessModal
        isOpen={showSuccessModal}
        title="Login berhasil!"
        message="Anda akan segera diarahkan ke dashboard atau dapat menekan tombol dibawah ini"
        buttonLabel="Lanjut ke dashboard"
        onButtonClick={handleSuccessModalClick}
      />
    </>
  );
}

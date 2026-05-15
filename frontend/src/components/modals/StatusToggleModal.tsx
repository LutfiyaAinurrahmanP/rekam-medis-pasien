import { useEffect, useState } from "react";
import { Modal } from "../ui/modal";

interface StatusToggleModalProps {
  isOpen: boolean;
  actionLabel: string;
  confirmLabel: string;
  onClose: () => void;
  onConfirm: () => Promise<void> | void;
}

export default function StatusToggleModal({
  isOpen,
  actionLabel,
  confirmLabel,
  onClose,
  onConfirm,
}: StatusToggleModalProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!isOpen) {
      setLoading(false);
      setError(null);
    }
  }, [isOpen]);

  if (!isOpen) return null;

  const handleClose = () => {
    setError(null);
    onClose();
  };

  const handleConfirm = async () => {
    setError(null);
    setLoading(true);

    try {
      await onConfirm();
      onClose();
    } catch (err: unknown) {
      console.error("Status toggle action failed:", err);
      let errorMessage = "An error occurred while updating the user status";

      if (err instanceof Error) {
        errorMessage = err.message;
      } else if (typeof err === "object" && err !== null && "message" in err) {
        errorMessage = String((err as Record<string, unknown>).message);
      }

      setError(errorMessage);
      setLoading(false);
    }
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={handleClose}
      className="relative w-full max-w-[540px] rounded-3xl bg-white p-6 overflow-hidden lg:p-10 dark:bg-gray-900"
    >
      <div className="text-center">
        <div className="relative z-1 mb-7 flex items-center justify-center">
          <svg
            className="fill-brand-50 dark:fill-brand-500/15"
            width="90"
            height="90"
            viewBox="0 0 90 90"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path d="M34.364 6.85053C38.6205 -2.28351 51.3795 -2.28351 55.636 6.85053C58.0129 11.951 63.5594 14.6722 68.9556 13.3853C78.6192 11.0807 86.5743 21.2433 82.2185 30.3287C79.7862 35.402 81.1561 41.5165 85.5082 45.0122C93.3019 51.2725 90.4628 63.9451 80.7747 66.1403C75.3648 67.3661 71.5265 72.2695 71.5572 77.9156C71.6123 88.0265 60.1169 93.6664 52.3918 87.3184C48.0781 83.7737 41.9219 83.7737 37.6082 87.3184C29.8831 93.6664 18.3877 88.0266 18.4428 77.9156C18.4735 72.2695 14.6352 67.3661 9.22531 66.1403C-0.462787 63.9451 -3.30193 51.2725 4.49185 45.0122C8.84391 41.5165 10.2138 35.402 7.78151 30.3287C3.42572 21.2433 11.3808 11.0807 21.0444 13.3853C26.4406 14.6722 31.9871 11.951 34.364 6.85053Z" />
          </svg>

          <span className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2">
            <svg
              className="fill-brand-600 dark:fill-brand-500"
              width="38"
              height="38"
              viewBox="0 0 38 38"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                fillRule="evenodd"
                clipRule="evenodd"
                d="M19 2.9375C10.1294 2.9375 2.9375 10.1294 2.9375 19C2.9375 27.8706 10.1294 35.0625 19 35.0625C27.8706 35.0625 35.0625 27.8706 35.0625 19C35.0625 10.1294 27.8706 2.9375 19 2.9375ZM19 5.9375C26.2147 5.9375 32.0625 11.7853 32.0625 19C32.0625 26.2147 26.2147 32.0625 19 32.0625C11.7853 32.0625 5.9375 26.2147 5.9375 19C5.9375 11.7853 11.7853 5.9375 19 5.9375ZM17.5 10.75C17.5 9.92157 18.1716 9.25 19 9.25C19.8284 9.25 20.5 9.92157 20.5 10.75V20.25C20.5 21.0784 19.8284 21.75 19 21.75C18.1716 21.75 17.5 21.0784 17.5 20.25V10.75ZM19 28.75C18.1716 28.75 17.5 28.0784 17.5 27.25C17.5 26.4216 18.1716 25.75 19 25.75C19.8284 25.75 20.5 26.4216 20.5 27.25C20.5 28.0784 19.8284 28.75 19 28.75Z"
              />
            </svg>
          </span>
        </div>

        <h4 className="mb-2 text-2xl font-semibold text-gray-800 dark:text-white/90 sm:text-title-sm">
          {actionLabel}
        </h4>
        <p className="text-sm leading-6 text-gray-500 dark:text-gray-400 mb-6">
          Confirm this action to update the user status immediately.
        </p>

        {error ? <p className="mb-4 text-sm text-red-600">{error}</p> : null}

        <div className="flex items-center justify-center w-full gap-3 mt-2">
          <button
            type="button"
            onClick={handleClose}
            className="flex justify-center w-full px-4 py-3 text-sm font-medium text-gray-700 rounded-lg bg-gray-100 shadow-theme-xs hover:bg-gray-200 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700 sm:w-auto"
            disabled={loading}
          >
            Cancel
          </button>
          <button
            type="button"
            onClick={handleConfirm}
            className="flex justify-center w-full px-4 py-3 text-sm font-medium text-white rounded-lg bg-brand-500 shadow-theme-xs hover:bg-brand-600 disabled:opacity-50 sm:w-auto"
            disabled={loading}
          >
            {loading ? "Updating..." : confirmLabel}
          </button>
        </div>
      </div>
    </Modal>
  );
}

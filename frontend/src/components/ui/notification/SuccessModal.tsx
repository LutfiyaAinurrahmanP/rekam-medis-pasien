import { CheckCircleIcon } from "../../../icons";

interface SuccessModalProps {
  title: string;
  message: string;
  buttonLabel: string;
  onButtonClick: () => void;
  isOpen: boolean;
}

const SuccessModal: React.FC<SuccessModalProps> = ({
  title,
  message,
  buttonLabel,
  onButtonClick,
  isOpen,
}) => {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="w-full max-w-md rounded-2xl bg-white p-6 shadow-2xl dark:bg-[#1E2634] sm:p-8">
        {/* Icon */}
        <div className="flex justify-center mb-6">
          <div className="flex items-center justify-center w-16 h-16 rounded-full bg-success-50">
            <CheckCircleIcon className="w-8 h-8 text-success-500" />
          </div>
        </div>

        {/* Content */}
        <div className="text-center">
          <h3 className="mb-3 text-lg font-semibold text-gray-800 dark:text-white/90 sm:text-xl">
            {title}
          </h3>
          <p className="mb-6 text-sm text-gray-600 dark:text-gray-400 sm:text-base">
            {message}
          </p>
        </div>

        {/* Button */}
        <button
          type="button"
          onClick={onButtonClick}
          className="w-full px-4 py-3 text-sm font-medium text-white transition rounded-lg bg-brand-500 hover:bg-brand-600 sm:text-base"
        >
          {buttonLabel}
        </button>
      </div>
    </div>
  );
};

export default SuccessModal;

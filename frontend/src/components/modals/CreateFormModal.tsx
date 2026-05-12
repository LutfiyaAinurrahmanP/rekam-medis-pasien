import React from "react";

export interface CreateFormFieldDef {
  name: string;
  label: string;
  type?: "text" | "email" | "tel" | "password" | "select";
  placeholder?: string;
  required?: boolean;
  options?: { value: string; label: string }[];
}

interface CreateFormModalProps {
  isOpen: boolean;
  formData: Record<string, string>;
  onChange: (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>,
  ) => void;
  onSubmit: (e: React.FormEvent<HTMLFormElement>) => void;
  onClose: () => void;
  errorsList: string[];
  loading: boolean;
  fields: CreateFormFieldDef[];
  title?: string;
  description?: string;
  submitLabel?: string;
  cancelLabel?: string;
  disabled?: boolean;
}

export default function CreateFormModal({
  isOpen,
  formData,
  onChange,
  onSubmit,
  onClose,
  errorsList,
  loading,
  fields,
  title = "Create",
  description = "Enter the details below",
  submitLabel = "Create",
  cancelLabel = "Cancel",
  disabled = false,
}: CreateFormModalProps) {
  if (!isOpen) return null;
  return (
    <>
      <div
        className="fixed inset-0 z-40 bg-black/50 backdrop-blur-sm transition-opacity"
        onClick={onClose}
      />

      {/* Modal */}
      <div className="fixed left-1/2 top-1/2 z-50 w-full max-w-md -translate-x-1/2 -translate-y-1/2">
        <div className="rounded-lg border border-gray-200 bg-white shadow-xl dark:border-gray-700 dark:bg-gray-800">
          {/* Header */}
          <div className="border-b border-gray-200 px-6 py-4 dark:border-gray-700">
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                  {title}
                </h2>
                <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                  {description}
                </p>
              </div>
              <button
                onClick={onClose}
                className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                disabled={loading || disabled}
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
          <form onSubmit={onSubmit} className="space-y-6 p-6">
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
              {fields.map((field) => (
                <div key={field.name}>
                  <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                    {field.label}
                    {field.required && <span className="text-red-500"> *</span>}
                  </label>
                  {field.type === "select" ? (
                    <select
                      name={field.name}
                      value={formData[field.name] || ""}
                      onChange={onChange}
                      className="w-full rounded-lg border border-gray-300 px-4 py-2 text-sm focus:border-transparent focus:ring-2 focus:ring-brand-500 dark:border-gray-700 dark:bg-gray-900 dark:text-white dark:focus:ring-brand-500"
                      required={field.required}
                      disabled={loading || disabled}
                    >
                      <option value="">
                        Select {field.label.toLowerCase()}
                      </option>
                      {field.options?.map((opt) => (
                        <option key={opt.value} value={opt.value}>
                          {opt.label}
                        </option>
                      ))}
                    </select>
                  ) : (
                    <input
                      name={field.name}
                      value={formData[field.name] || ""}
                      onChange={onChange}
                      type={field.type || "text"}
                      className="w-full rounded-lg border border-gray-300 px-4 py-2 text-sm focus:border-transparent focus:ring-2 focus:ring-brand-500 dark:border-gray-700 dark:bg-gray-900 dark:text-white dark:focus:ring-brand-500"
                      placeholder={field.placeholder}
                      required={field.required}
                      disabled={loading || disabled}
                    />
                  )}
                </div>
              ))}
            </div>

            {/* Footer */}
            <div className="border-t border-gray-200 pt-6 dark:border-gray-700">
              <div className="flex gap-3">
                <button
                  type="button"
                  onClick={onClose}
                  disabled={loading || disabled}
                  className="flex-1 rounded-lg border border-gray-300 px-4 py-2 font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 dark:border-gray-700 dark:text-gray-300 dark:hover:bg-white/[0.05]"
                >
                  {cancelLabel}
                </button>
                <button
                  type="submit"
                  disabled={loading || disabled}
                  className="flex-1 rounded-lg bg-brand-500 px-4 py-2 font-medium text-white hover:bg-brand-600 disabled:opacity-50"
                >
                  {loading ? `${submitLabel}...` : submitLabel}
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </>
  );
}

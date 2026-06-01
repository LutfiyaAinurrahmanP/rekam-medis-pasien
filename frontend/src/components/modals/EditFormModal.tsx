import React from "react";
import { Modal } from "../ui/modal";

export interface EditFormFieldDef {
  name: string;
  label: string;
  type?: "text" | "email" | "tel" | "password" | "date" | "select" | "textarea";
  placeholder?: string;
  required?: boolean;
  disabled?: boolean;
  rows?: number;
  options?: { value: string; label: string }[];
}

interface EditFormModalProps {
  isOpen: boolean;
  formData: Record<string, string>;
  onChange: (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement
    >,
  ) => void;
  onSubmit: (e: React.FormEvent<HTMLFormElement>) => void;
  onClose: () => void;
  errorsList: string[];
  loading: boolean;
  fields: EditFormFieldDef[];
  title?: string;
  description?: string;
  submitLabel?: string;
  cancelLabel?: string;
  disabled?: boolean;
  children?: React.ReactNode;
}

export default function EditFormModal({
  isOpen,
  formData,
  onChange,
  onSubmit,
  onClose,
  errorsList,
  loading,
  fields,
  title = "Edit",
  description = "Update the details below",
  submitLabel = "Update",
  cancelLabel = "Cancel",
  disabled = false,
  children,
}: EditFormModalProps) {
  if (!isOpen) return null;

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      className="max-w-[584px] p-5 lg:p-10"
      showCloseButton={false}
    >
      <form onSubmit={onSubmit}>
        <h4 className="mb-6 text-lg font-medium text-gray-800 dark:text-white/90">
          {title}
        </h4>

        {description ? (
          <p className="-mt-4 mb-6 text-sm text-gray-500 dark:text-gray-400">
            {description}
          </p>
        ) : null}

        {children ? (
          children
        ) : (
          <>
            {errorsList.length > 0 ? (
              <div className="mb-6 rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-700 dark:border-red-800 dark:bg-red-900/30 dark:text-red-400">
                <ul className="list-inside list-disc space-y-1">
                  {errorsList.map((errMsg, idx) => (
                    <li key={idx}>{errMsg}</li>
                  ))}
                </ul>
              </div>
            ) : null}

            <div className="grid grid-cols-1 gap-x-6 gap-y-5 sm:grid-cols-2">
              {fields.map((field) => (
                <div key={field.name} className="col-span-1">
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
                      disabled={loading || disabled || field.disabled}
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
                  ) : field.type === "textarea" ? (
                    <textarea
                      name={field.name}
                      value={formData[field.name] || ""}
                      onChange={onChange}
                      rows={field.rows ?? 3}
                      className="w-full rounded-lg border border-gray-300 px-4 py-2 text-sm focus:border-transparent focus:ring-2 focus:ring-brand-500 dark:border-gray-700 dark:bg-gray-900 dark:text-white dark:focus:ring-brand-500"
                      placeholder={field.placeholder}
                      required={field.required}
                      disabled={loading || disabled || field.disabled}
                    />
                  ) : (
                    <input
                      name={field.name}
                      value={formData[field.name] || ""}
                      onChange={onChange}
                      type={field.type || "text"}
                      className="w-full rounded-lg border border-gray-300 px-4 py-2 text-sm focus:border-transparent focus:ring-2 focus:ring-brand-500 dark:border-gray-700 dark:bg-gray-900 dark:text-white dark:focus:ring-brand-500"
                      placeholder={field.placeholder}
                      required={field.required}
                      disabled={loading || disabled || field.disabled}
                    />
                  )}
                </div>
              ))}
            </div>

            <div className="mt-6 flex items-center justify-end gap-3">
              <button
                type="button"
                onClick={onClose}
                disabled={loading || disabled}
                className="inline-flex items-center justify-center rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50 dark:border-gray-700 dark:text-gray-300 dark:hover:bg-white/[0.05]"
              >
                {cancelLabel}
              </button>
              <button
                type="submit"
                disabled={loading || disabled}
                className="inline-flex items-center justify-center rounded-lg bg-brand-500 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-brand-600 disabled:cursor-not-allowed disabled:opacity-50"
              >
                {loading ? `${submitLabel}...` : submitLabel}
              </button>
            </div>
          </>
        )}
      </form>
    </Modal>
  );
}

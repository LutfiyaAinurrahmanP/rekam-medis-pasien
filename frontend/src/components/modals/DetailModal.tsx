import React from "react";
import { Modal } from "../ui/modal";

interface FieldDef {
  key: string;
  label?: string;
  type?: "text" | "button";
  render?: (
    value: unknown,
    record?: Record<string, unknown>,
  ) => React.ReactNode;
}

interface DetailModalProps {
  isOpen: boolean;
  loading: boolean;
  data: Record<string, unknown> | null;
  fields: FieldDef[];
  onClose: () => void;
  children?: React.ReactNode;
}

export default function DetailModal({
  isOpen,
  loading,
  data,
  fields,
  onClose,
  children,
}: DetailModalProps) {
  if (!isOpen) return null;

  const renderValue = (key: string, render?: FieldDef["render"]) => {
    const value = data ? data[key] : undefined;
    if (render) return render(value, data ?? undefined);
    if (typeof value === "boolean") return value ? "Active" : "Inactive";
    if (value === null || value === undefined || value === "") return "-";
    // Format ISO date-like strings
    if (typeof value === "string" && /\d{4}-\d{2}-\d{2}T?/.test(value)) {
      try {
        const d = new Date(String(value));
        if (!isNaN(d.getTime())) return d.toLocaleString();
      } catch {
        // ignore
      }
    }
    return String(value);
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      className="max-w-[584px] p-5 lg:p-10"
      showCloseButton={false}
    >
      <div>
        <h4 className="mb-6 text-lg font-medium text-gray-800 dark:text-white/90">
          Details
        </h4>
        <p className="-mt-4 mb-6 text-sm text-gray-500 dark:text-gray-400">
          View information below
        </p>

        {children ? (
          children
        ) : loading ? (
          <div className="flex items-center justify-center py-8">
            <p className="text-gray-500">Loading...</p>
          </div>
        ) : !data ? (
          <div className="flex items-center justify-center py-8">
            <p className="text-red-500">Not found</p>
          </div>
        ) : (
          <>
            <div className="grid grid-cols-1 gap-x-6 gap-y-5 sm:grid-cols-2">
              {fields.map((f) => (
                <div key={f.key} className="col-span-1">
                  <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                    {f.label ?? f.key.charAt(0).toUpperCase() + f.key.slice(1)}
                  </label>
                  {f.type === "button" ? (
                    <div>{renderValue(f.key, f.render)}</div>
                  ) : (
                    <p className="rounded-lg border border-gray-300 px-4 py-2 text-sm text-gray-800 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90">
                      {renderValue(f.key, f.render)}
                    </p>
                  )}
                </div>
              ))}
            </div>

            <div className="mt-6 flex items-center justify-end gap-3">
              <button
                type="button"
                onClick={onClose}
                className="inline-flex items-center justify-center rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-gray-700 dark:text-gray-300 dark:hover:bg-white/[0.05]"
              >
                Close
              </button>
            </div>
          </>
        )}
      </div>
    </Modal>
  );
}

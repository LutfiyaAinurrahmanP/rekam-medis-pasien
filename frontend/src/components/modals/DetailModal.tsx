import React from "react";

interface FieldDef {
  key: string;
  label?: string;
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
}

export default function DetailModal({
  isOpen,
  loading,
  data,
  fields,
  onClose,
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
    <>
      <div
        className="fixed inset-0 z-40 bg-black/50 backdrop-blur-sm transition-opacity"
        onClick={onClose}
      />

      <div className="fixed left-1/2 top-1/2 z-50 w-full max-w-2xl -translate-x-1/2 -translate-y-1/2">
        <div className="rounded-lg border border-gray-200 bg-white shadow-xl dark:border-gray-700 dark:bg-gray-800">
          <div className="border-b border-gray-200 px-6 py-4 dark:border-gray-700">
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                  Details
                </h2>
                <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                  View information
                </p>
              </div>
              <button
                onClick={onClose}
                className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
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

          <div className="p-6">
            {loading ? (
              <div className="flex items-center justify-center p-8">
                <p className="text-gray-500">Loading...</p>
              </div>
            ) : !data ? (
              <div className="flex items-center justify-center p-8">
                <p className="text-red-500">Not found</p>
              </div>
            ) : (
              <div className="space-y-4">
                <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
                  {fields.map((f) => (
                    <div key={f.key}>
                      <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                        {f.label ??
                          f.key.charAt(0).toUpperCase() + f.key.slice(1)}
                      </label>
                      <p className="rounded-lg bg-gray-50 px-4 py-2 dark:bg-gray-900 dark:text-white">
                        {renderValue(f.key, f.render)}
                      </p>
                    </div>
                  ))}
                </div>

                <div className="flex gap-3 pt-8 border-t border-gray-200 dark:border-gray-700">
                  <button
                    onClick={onClose}
                    className="rounded-lg border border-gray-300 px-6 py-2 font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-700 dark:text-gray-300 dark:hover:bg-white/[0.05]"
                  >
                    Close
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </>
  );
}

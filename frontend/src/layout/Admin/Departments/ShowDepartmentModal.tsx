import { useEffect, useState } from "react";
import { get } from "../../../services/api";
import { Modal } from "../../../components/ui/modal";

interface Props {
  isOpen: boolean;
  id?: string | undefined;
  onClose?: () => void;
}

export default function ShowDepartmentModal({ isOpen, id, onClose }: Props) {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<{
    id?: number;
    name?: string;
    code?: string;
    description?: string;
    floor_location?: string;
    created_at?: string;
    updated_at?: string;
  } | null>(null);

  useEffect(() => {
    const load = async () => {
      if (!id) return;
      setLoading(true);
      try {
        const resp = await get<{
          id: number;
          name: string;
          code?: string;
          description?: string;
          floor_location?: string;
          created_at: string;
          updated_at?: string;
        }>(`/departments/${id}`);
        setData(resp);
      } catch (err) {
        console.error("Failed to fetch department:", err);
      } finally {
        setLoading(false);
      }
    };

    if (isOpen) load();
  }, [id, isOpen]);

  if (!isOpen) return null;

  const handleClose = () => {
    onClose?.();
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={handleClose}
      className="max-w-[584px] p-5 lg:p-10"
    >
      {loading ? (
        <div className="py-8 text-center text-sm text-gray-500 dark:text-gray-400">
          Loading...
        </div>
      ) : (
        <div>
          <h4 className="mb-6 text-lg font-medium text-gray-800 dark:text-white/90">
            Department Details
          </h4>
          <p className="-mt-4 mb-6 text-sm text-gray-500 dark:text-gray-400">
            View information below
          </p>

          <div className="grid grid-cols-1 gap-x-6 gap-y-5 sm:grid-cols-2">
            <div className="col-span-1">
              <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                Name
              </label>
              <input
                value={data?.name ?? "-"}
                readOnly
                className="w-full rounded-lg border border-gray-300 bg-gray-50 px-4 py-2 text-sm text-gray-800 outline-none dark:border-gray-700 dark:bg-gray-900 dark:text-white/90"
              />
            </div>

            <div className="col-span-1">
              <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                Code
              </label>
              <input
                value={data?.code ?? "-"}
                readOnly
                className="w-full rounded-lg border border-gray-300 bg-gray-50 px-4 py-2 text-sm text-gray-800 outline-none dark:border-gray-700 dark:bg-gray-900 dark:text-white/90"
              />
            </div>

            <div className="col-span-1 sm:col-span-2">
              <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                Description
              </label>
              <textarea
                value={data?.description ?? "-"}
                readOnly
                rows={4}
                className="w-full rounded-lg border border-gray-300 bg-gray-50 px-4 py-2 text-sm text-gray-800 outline-none dark:border-gray-700 dark:bg-gray-900 dark:text-white/90"
              />
            </div>

            <div className="col-span-1">
              <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                Floor / Location
              </label>
              <input
                value={data?.floor_location ?? "-"}
                readOnly
                className="w-full rounded-lg border border-gray-300 bg-gray-50 px-4 py-2 text-sm text-gray-800 outline-none dark:border-gray-700 dark:bg-gray-900 dark:text-white/90"
              />
            </div>

            <div className="col-span-1">
              <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                Created At
              </label>
              <input
                value={data?.created_at ?? "-"}
                readOnly
                className="w-full rounded-lg border border-gray-300 bg-gray-50 px-4 py-2 text-sm text-gray-800 outline-none dark:border-gray-700 dark:bg-gray-900 dark:text-white/90"
              />
            </div>

            <div className="col-span-1">
              <label className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-400">
                Last Updated
              </label>
              <input
                value={data?.updated_at ?? "-"}
                readOnly
                className="w-full rounded-lg border border-gray-300 bg-gray-50 px-4 py-2 text-sm text-gray-800 outline-none dark:border-gray-700 dark:bg-gray-900 dark:text-white/90"
              />
            </div>
          </div>

          <div className="mt-6 flex items-center justify-end gap-3">
            <button
              type="button"
              onClick={handleClose}
              className="inline-flex items-center justify-center rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-gray-700 dark:text-gray-300 dark:hover:bg-white/[0.05]"
            >
              Close
            </button>
          </div>
        </div>
      )}
    </Modal>
  );
}

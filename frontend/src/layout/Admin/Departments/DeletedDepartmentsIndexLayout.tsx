import { useCallback, useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router";
import PageBreadcrumb from "../../../components/common/PageBreadCrumb";
import TableHeaderComponent from "../../../components/tables/TableHeader";
import Pagination from "../../../components/tables/DataTables/TableThree/Pagination";
import {
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
} from "../../../components/ui/table";
import Button from "../../../components/ui/button/Button";
import DeleteModal from "../../../components/modals/DeleteModal";
import SuccessModal from "../../../components/ui/notification/SuccessModal";
import { patch, del } from "../../../services/api";
import {
  useDeletedDepartments,
  type DeletedDepartment,
} from "../../../hooks/Departments/useDeletedDepartments";

interface Props {
  onBackToDepartments?: () => void;
}

const formatDateTime = (dateString: string | null) => {
  if (!dateString) return "-";
  const date = new Date(dateString);
  if (Number.isNaN(date.getTime())) return "-";
  return date.toLocaleString("en-GB", {
    day: "numeric",
    month: "long",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};

export default function DeletedDepartmentsIndexLayout({
  onBackToDepartments,
}: Props = {}) {
  const navigate = useNavigate();
  const { deletedDepartments, loading, error, meta, fetchDeletedDepartments } =
    useDeletedDepartments();
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [search, setSearch] = useState("");
  const [selectedDeleteId, setSelectedDeleteId] = useState<string | undefined>(
    undefined,
  );
  const [selectedDeleteName, setSelectedDeleteName] = useState<
    string | undefined
  >(undefined);
  const [showSuccessModal, setShowSuccessModal] = useState(false);
  const [successMessage, setSuccessMessage] = useState("");
  const debounceRef = useRef<number | null>(null);

  const triggerFetch = useCallback(
    (page: number, pageSize: number, searchValue: string) => {
      fetchDeletedDepartments({
        page,
        page_size: pageSize,
        search: searchValue.trim() || undefined,
      });
    },
    [fetchDeletedDepartments],
  );

  useEffect(() => {
    triggerFetch(currentPage, rowsPerPage, search);
  }, [currentPage, rowsPerPage, search, triggerFetch]);

  useEffect(() => {
    return () => {
      if (debounceRef.current) window.clearTimeout(debounceRef.current);
    };
  }, []);

  useEffect(() => {
    if (meta.total_pages > 0 && currentPage > meta.total_pages)
      setCurrentPage(meta.total_pages);
  }, [currentPage, meta.total_pages]);

  useEffect(() => {
    if (showSuccessModal) {
      const t = window.setTimeout(() => {
        setShowSuccessModal(false);
        setSuccessMessage("");
      }, 3000);
      return () => window.clearTimeout(t);
    }
    return undefined;
  }, [showSuccessModal]);

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setSearch(value);
    setCurrentPage(1);
    if (debounceRef.current) window.clearTimeout(debounceRef.current);
    debounceRef.current = window.setTimeout(() => {
      triggerFetch(1, rowsPerPage, value);
    }, 150);
  };

  const handlePageChange = (page: number) => {
    if (page >= 1 && page <= meta.total_pages) setCurrentPage(page);
  };

  const handleRowsPerPageChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const next = Number.parseInt(e.target.value, 10);
    setRowsPerPage(next);
    setCurrentPage(1);
  };

  const handleRestore = async (d: DeletedDepartment) => {
    try {
      await patch(`/departments/${d.id}/restore`, {});
      setSuccessMessage(
        `Department "${d.name}" has been successfully restored.`,
      );
      setShowSuccessModal(true);
      triggerFetch(currentPage, rowsPerPage, search);
    } catch (err) {
      console.error("Failed to restore department:", err);
    }
  };

  const handleHardDelete = async () => {
    if (!selectedDeleteId) return;
    try {
      await del(`/departments/${selectedDeleteId}/hard-delete`);
      setSuccessMessage(
        `Department "${selectedDeleteName ?? ""}" has been permanently deleted.`,
      );
      setShowSuccessModal(true);
      triggerFetch(currentPage, rowsPerPage, search);
    } catch (err) {
      console.error("Failed to hard delete department:", err);
      throw err;
    }
  };

  return (
    <div className="space-y-4">
      <PageBreadcrumb pageTitle="Deleted Departments" />

      <div className="overflow-hidden rounded-xl bg-white dark:bg-white/[0.03]">
        {error ? (
          <div className="border-b border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700 dark:border-red-800 dark:bg-red-900/20 dark:text-red-300">
            {error}
          </div>
        ) : null}

        <TableHeaderComponent
          search={search}
          rowsPerPage={rowsPerPage}
          onSearchChange={handleSearchChange}
          onRowsPerPageChange={handleRowsPerPageChange}
          showDeleteButton={false}
          searchPlaceholder="Search deleted departments..."
          actionButtons={
            <Button
              variant="outline"
              size="sm"
              onClick={
                onBackToDepartments ?? (() => navigate("/admin/departments"))
              }
            >
              Back to Departments
            </Button>
          }
        />

        <div className="max-w-full overflow-x-auto custom-scrollbar">
          <Table>
            <TableHeader>
              <TableRow>
                <TableCell
                  isHeader
                  className="px-4 py-3 border border-gray-100 dark:border-white/[0.05]"
                >
                  <p className="font-medium text-gray-700 text-theme-xs dark:text-gray-400">
                    Name
                  </p>
                </TableCell>
                <TableCell
                  isHeader
                  className="px-4 py-3 border border-gray-100 dark:border-white/[0.05]"
                >
                  <p className="font-medium text-gray-700 text-theme-xs dark:text-gray-400">
                    Code
                  </p>
                </TableCell>
                <TableCell
                  isHeader
                  className="px-4 py-3 border border-gray-100 dark:border-white/[0.05]"
                >
                  <p className="font-medium text-gray-700 text-theme-xs dark:text-gray-400">
                    Floor/Location
                  </p>
                </TableCell>
                <TableCell
                  isHeader
                  className="px-4 py-3 border border-gray-100 dark:border-white/[0.05]"
                >
                  <p className="font-medium text-gray-700 text-theme-xs dark:text-gray-400">
                    Deleted At
                  </p>
                </TableCell>
                <TableCell
                  isHeader
                  className="px-4 py-3 border border-gray-100 dark:border-white/[0.05]"
                >
                  <p className="font-medium text-gray-700 text-theme-xs dark:text-gray-400">
                    Action
                  </p>
                </TableCell>
              </TableRow>
            </TableHeader>
            <TableBody>
              {loading && deletedDepartments.length === 0 ? (
                <TableRow>
                  <td
                    colSpan={5}
                    className="px-4 py-8 text-center text-gray-500 dark:text-gray-400"
                  >
                    Loading...
                  </td>
                </TableRow>
              ) : deletedDepartments.length > 0 ? (
                deletedDepartments.map((d) => (
                  <TableRow key={d.id}>
                    <TableCell className="px-4 py-4 font-normal text-gray-800 border border-gray-100 dark:border-white/[0.05] text-theme-sm dark:text-white/90 whitespace-nowrap">
                      <div>
                        <p className="block font-medium text-gray-800 text-theme-sm dark:text-white/90">
                          {d.name}
                        </p>
                        <span className="text-sm font-normal text-gray-500 dark:text-gray-400">
                          {d.description}
                        </span>
                      </div>
                    </TableCell>
                    <TableCell className="px-4 py-4 font-normal text-gray-800 border border-gray-100 dark:border-white/[0.05] text-theme-sm dark:text-white/90 whitespace-nowrap">
                      {d.code}
                    </TableCell>
                    <TableCell className="px-4 py-4 font-normal text-gray-800 border border-gray-100 dark:border-white/[0.05] text-theme-sm dark:text-white/90 whitespace-nowrap">
                      {d.floor_location}
                    </TableCell>
                    <TableCell className="px-4 py-4 font-normal text-gray-800 border border-gray-100 dark:border-white/[0.05] text-theme-sm dark:text-white/90 whitespace-nowrap">
                      {formatDateTime(d.deleted_at)}
                    </TableCell>
                    <TableCell className="px-4 py-4 font-normal text-gray-800 border border-gray-100 dark:border-white/[0.05] text-theme-sm dark:text-white/90 whitespace-nowrap">
                      <div className="flex items-center gap-2">
                        <button
                          type="button"
                          onClick={() => handleRestore(d)}
                          className="rounded-lg border border-brand-300 px-3 py-2 text-xs font-medium text-brand-700 transition-colors hover:bg-brand-50 dark:border-brand-700 dark:text-brand-300 dark:hover:bg-brand-500/10"
                        >
                          Restore
                        </button>
                        <button
                          type="button"
                          onClick={() => {
                            setSelectedDeleteId(String(d.id));
                            setSelectedDeleteName(d.name);
                          }}
                          className="rounded-lg border border-red-300 px-3 py-2 text-xs font-medium text-red-700 transition-colors hover:bg-red-50 dark:border-red-700 dark:text-red-300 dark:hover:bg-red-500/10"
                        >
                          Delete Permanently
                        </button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <td
                    colSpan={5}
                    className="px-4 py-8 text-center text-gray-500 dark:text-gray-400"
                  >
                    No deleted departments found
                  </td>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>

        <div className="border border-t-0 rounded-b-xl border-gray-100 py-4 pl-[18px] pr-4 dark:border-white/[0.05]">
          <div className="flex flex-col xl:flex-row xl:items-center xl:justify-between">
            <div className="pb-3 xl:pb-0">
              <p className="pb-3 text-sm font-medium text-center text-gray-500 border-b border-gray-100 dark:border-gray-800 dark:text-gray-400 xl:border-b-0 xl:pb-0 xl:text-left">
                Showing{" "}
                {meta.total_items === 0
                  ? 0
                  : (currentPage - 1) * rowsPerPage + 1}{" "}
                to {Math.min(currentPage * rowsPerPage, meta.total_items)} of{" "}
                {meta.total_items} entries
              </p>
            </div>
            <Pagination
              currentPage={currentPage}
              totalPages={meta.total_pages}
              onPageChange={handlePageChange}
            />
          </div>
        </div>
      </div>

      <DeleteModal
        isOpen={Boolean(selectedDeleteId)}
        itemName={selectedDeleteName}
        title="Delete Permanently"
        description="This data will be permanently deleted and cannot be recovered. Are you sure you want to delete"
        onClose={() => {
          setSelectedDeleteId(undefined);
          setSelectedDeleteName(undefined);
        }}
        onConfirm={handleHardDelete}
      />

      <SuccessModal
        title="Success"
        message={successMessage}
        buttonLabel="Close"
        isOpen={showSuccessModal}
        onButtonClick={() => setShowSuccessModal(false)}
      />
    </div>
  );
}

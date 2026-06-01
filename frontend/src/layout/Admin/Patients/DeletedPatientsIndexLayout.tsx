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
import { useAuth } from "../../../context/AuthContext";
import { patch, del } from "../../../services/api";
import {
  useDeletedPatients,
  type DeletedPatient,
} from "../../../hooks/Patients/useDeletedPatients";
import { getRolePatientsPath } from "../../../pages/Roles/shared/role-routing";

interface DeletedPatientsIndexLayoutProps {
  onBackToPatients?: () => void;
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

export default function DeletedPatientsIndexLayout({
  onBackToPatients,
}: DeletedPatientsIndexLayoutProps = {}) {
  const navigate = useNavigate();
  const { user } = useAuth();
  const { deletedPatients, loading, error, meta, fetchDeletedPatients } =
    useDeletedPatients();
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [search, setSearch] = useState("");
  const [selectedDeletePatientId, setSelectedDeletePatientId] = useState<
    string | undefined
  >(undefined);
  const [selectedDeletePatientName, setSelectedDeletePatientName] = useState<
    string | undefined
  >(undefined);
  const [showSuccessModal, setShowSuccessModal] = useState(false);
  const [successMessage, setSuccessMessage] = useState<string>("");
  const debounceRef = useRef<number | null>(null);
  const skipNextFetchRef = useRef(false);

  const triggerFetch = useCallback(
    (page: number, pageSize: number, searchValue: string) => {
      fetchDeletedPatients({
        page,
        page_size: pageSize,
        search: searchValue.trim() || undefined,
      });
    },
    [fetchDeletedPatients],
  );

  useEffect(() => {
    if (skipNextFetchRef.current) {
      skipNextFetchRef.current = false;
      return;
    }

    triggerFetch(currentPage, rowsPerPage, search);
  }, [currentPage, rowsPerPage, search, triggerFetch]);

  useEffect(() => {
    return () => {
      if (debounceRef.current) {
        window.clearTimeout(debounceRef.current);
      }
    };
  }, []);

  useEffect(() => {
    if (meta.total_pages > 0 && currentPage > meta.total_pages) {
      setCurrentPage(meta.total_pages);
    }
  }, [currentPage, meta.total_pages]);

  useEffect(() => {
    if (showSuccessModal) {
      const timer = setTimeout(() => {
        setShowSuccessModal(false);
        setSuccessMessage("");
      }, 3000);

      return () => clearTimeout(timer);
    }
    return undefined;
  }, [showSuccessModal]);

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    skipNextFetchRef.current = true;
    setSearch(value);
    setCurrentPage(1);

    if (debounceRef.current) {
      window.clearTimeout(debounceRef.current);
    }

    debounceRef.current = window.setTimeout(() => {
      triggerFetch(1, rowsPerPage, value);
      skipNextFetchRef.current = false;
    }, 150);
  };

  const handlePageChange = (page: number) => {
    if (page >= 1 && page <= meta.total_pages) {
      setCurrentPage(page);
    }
  };

  const handleRowsPerPageChange = (
    e: React.ChangeEvent<HTMLSelectElement>,
  ): void => {
    const nextRowsPerPage = Number.parseInt(e.target.value, 10);
    setRowsPerPage(nextRowsPerPage);
    setCurrentPage(1);
  };

  const handleRestorePatient = async (patient: DeletedPatient) => {
    try {
      await patch(`/patients/${patient.id}/restore`, {});
      setSuccessMessage(
        `Patient "${patient.full_name}" has been successfully restored.`,
      );
      setShowSuccessModal(true);
      triggerFetch(currentPage, rowsPerPage, search);
    } catch (err) {
      console.error("Failed to restore patient:", err);
    }
  };

  const handleHardDeletePatient = async () => {
    if (!selectedDeletePatientId) return;

    try {
      await del(`/patients/${selectedDeletePatientId}/hard-delete`);
      setSuccessMessage(
        `Patient "${selectedDeletePatientName ?? ""}" has been permanently deleted.`,
      );
      setShowSuccessModal(true);
      setSelectedDeletePatientId(undefined);
      setSelectedDeletePatientName(undefined);
      triggerFetch(currentPage, rowsPerPage, search);
    } catch (err) {
      console.error("Failed to hard delete patient:", err);
      throw err;
    }
  };

  return (
    <div className="space-y-4">
      <PageBreadcrumb pageTitle="Deleted Patients" />

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
          searchPlaceholder="Search deleted patients by name, code, phone, or email..."
          actionButtons={
            <Button
              variant="outline"
              size="sm"
              onClick={
                onBackToPatients ??
                (() => navigate(getRolePatientsPath(user?.role)))
              }
            >
              Back to Patients
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
                    Patient Code
                  </p>
                </TableCell>
                <TableCell
                  isHeader
                  className="px-4 py-3 border border-gray-100 dark:border-white/[0.05]"
                >
                  <p className="font-medium text-gray-700 text-theme-xs dark:text-gray-400">
                    Full Name
                  </p>
                </TableCell>
                <TableCell
                  isHeader
                  className="px-4 py-3 border border-gray-100 dark:border-white/[0.05]"
                >
                  <p className="font-medium text-gray-700 text-theme-xs dark:text-gray-400">
                    Age
                  </p>
                </TableCell>
                <TableCell
                  isHeader
                  className="px-4 py-3 border border-gray-100 dark:border-white/[0.05]"
                >
                  <p className="font-medium text-gray-700 text-theme-xs dark:text-gray-400">
                    Phone
                  </p>
                </TableCell>
                <TableCell
                  isHeader
                  className="px-4 py-3 border border-gray-100 dark:border-white/[0.05]"
                >
                  <p className="font-medium text-gray-700 text-theme-xs dark:text-gray-400">
                    Email
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
              {loading && deletedPatients.length === 0 ? (
                <TableRow>
                  <td
                    colSpan={7}
                    className="px-4 py-8 text-center text-gray-500 dark:text-gray-400"
                  >
                    Loading...
                  </td>
                </TableRow>
              ) : deletedPatients.length > 0 ? (
                deletedPatients.map((patient) => (
                  <TableRow key={patient.id}>
                    <TableCell className="px-4 py-4 font-normal text-gray-800 border border-gray-100 dark:border-white/[0.05] text-theme-sm dark:text-white/90 whitespace-nowrap">
                      {patient.patient_code}
                    </TableCell>
                    <TableCell className="px-4 py-4 font-normal text-gray-800 border border-gray-100 dark:border-white/[0.05] text-theme-sm dark:text-white/90 whitespace-nowrap">
                      {patient.full_name}
                    </TableCell>
                    <TableCell className="px-4 py-4 font-normal text-gray-800 border border-gray-100 dark:border-white/[0.05] text-theme-sm dark:text-white/90 whitespace-nowrap">
                      {patient.age}
                    </TableCell>
                    <TableCell className="px-4 py-4 font-normal text-gray-800 border border-gray-100 dark:border-white/[0.05] text-theme-sm dark:text-white/90 whitespace-nowrap">
                      {patient.phone || "-"}
                    </TableCell>
                    <TableCell className="px-4 py-4 font-normal text-gray-800 border border-gray-100 dark:border-white/[0.05] text-theme-sm dark:text-white/90 whitespace-nowrap">
                      {patient.email || "-"}
                    </TableCell>
                    <TableCell className="px-4 py-4 font-normal text-gray-800 border border-gray-100 dark:border-white/[0.05] text-theme-sm dark:text-white/90 whitespace-nowrap">
                      {formatDateTime(patient.deleted_at)}
                    </TableCell>
                    <TableCell className="px-4 py-4 font-normal text-gray-800 border border-gray-100 dark:border-white/[0.05] text-theme-sm dark:text-white/90 whitespace-nowrap">
                      <div className="flex items-center gap-2">
                        <button
                          type="button"
                          onClick={() => handleRestorePatient(patient)}
                          className="rounded-lg border border-brand-300 px-3 py-2 text-xs font-medium text-brand-700 transition-colors hover:bg-brand-50 dark:border-brand-700 dark:text-brand-300 dark:hover:bg-brand-500/10"
                        >
                          Restore
                        </button>
                        <button
                          type="button"
                          onClick={() => {
                            setSelectedDeletePatientId(String(patient.id));
                            setSelectedDeletePatientName(patient.full_name);
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
                    colSpan={7}
                    className="px-4 py-8 text-center text-gray-500 dark:text-gray-400"
                  >
                    No deleted patients found
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
        isOpen={Boolean(selectedDeletePatientId)}
        itemName={selectedDeletePatientName}
        title="Delete Permanently"
        description="This data will be permanently deleted and cannot be recovered. Are you sure you want to delete"
        onClose={() => {
          setSelectedDeletePatientId(undefined);
          setSelectedDeletePatientName(undefined);
        }}
        onConfirm={handleHardDeletePatient}
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

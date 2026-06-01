import { useCallback, useEffect, useRef, useState } from "react";
import PageBreadcrumb from "../../../components/common/PageBreadCrumb";
import BaseTable, {
  type ColumnDefinition,
} from "../../../components/tables/BaseTable";
import DeleteModal from "../../../components/modals/DeleteModal";
import SuccessModal from "../../../components/ui/notification/SuccessModal";
import { del } from "../../../services/api";
import { usePatients, type Patient } from "../../../hooks/Patients/usePatients";
import CreatePatientModal from "./CreatePatientModal";
import EditPatientModal from "./EditPatientModal";
import ShowPatientModal from "./ShowPatientModal";
import DeletedPatientsIndexLayout from "./DeletedPatientsIndexLayout";

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString("en-GB", {
    day: "numeric",
    month: "long",
    year: "numeric",
  });
};

const formatGender = (gender: string) =>
  gender ? gender.charAt(0).toUpperCase() + gender.slice(1) : "-";

export default function PatientsIndexLayout() {
  const [viewMode, setViewMode] = useState<"active" | "deleted">("active");
  const { patients, loading, error, meta, fetchPatients } = usePatients();
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [search, setSearch] = useState("");
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [selectedEditPatientId, setSelectedEditPatientId] = useState<
    string | undefined
  >(undefined);
  const [isShowModalOpen, setIsShowModalOpen] = useState(false);
  const [selectedPatientId, setSelectedPatientId] = useState<
    string | undefined
  >(undefined);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [selectedDeletePatientId, setSelectedDeletePatientId] = useState<
    string | undefined
  >(undefined);
  const [selectedDeletePatientName, setSelectedDeletePatientName] = useState<
    string | undefined
  >(undefined);
  const [successMessage, setSuccessMessage] = useState("");
  const [showSuccessModal, setShowSuccessModal] = useState(false);
  const debounceRef = useRef<number | null>(null);
  const skipNextFetchRef = useRef(false);

  const triggerFetch = useCallback(
    (page: number, pageSize: number, searchValue: string) => {
      fetchPatients({
        page,
        page_size: pageSize,
        search: searchValue.trim() || undefined,
      });
    },
    [fetchPatients],
  );

  useEffect(() => {
    if (viewMode === "deleted") {
      return;
    }

    if (skipNextFetchRef.current) {
      skipNextFetchRef.current = false;
      return;
    }

    triggerFetch(currentPage, rowsPerPage, search);
  }, [currentPage, rowsPerPage, search, triggerFetch, viewMode]);

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
    if (!showSuccessModal) {
      return undefined;
    }

    const timer = window.setTimeout(() => {
      setShowSuccessModal(false);
      setSuccessMessage("");
    }, 3000);

    return () => window.clearTimeout(timer);
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

  const handleCreatePatientSuccess = () => {
    setIsCreateModalOpen(false);
    triggerFetch(1, rowsPerPage, search);
  };

  const handleEditPatient = (patient: Patient) => {
    setSelectedEditPatientId(String(patient.id));
    setIsEditModalOpen(true);
  };

  const handleEditPatientClose = () => {
    setIsEditModalOpen(false);
    setSelectedEditPatientId(undefined);
  };

  const handleEditPatientSuccess = () => {
    handleEditPatientClose();
    triggerFetch(currentPage, rowsPerPage, search);
  };

  const handleDeletePatient = (patient: Patient) => {
    setSelectedDeletePatientId(String(patient.id));
    setSelectedDeletePatientName(patient.full_name);
    setIsDeleteModalOpen(true);
  };

  const handleViewPatient = (patient: Patient) => {
    setSelectedPatientId(String(patient.id));
    setIsShowModalOpen(true);
  };

  const patientColumns: ColumnDefinition<Patient>[] = [
    {
      key: "patient_code",
      header: "Patient Code",
      type: "text",
    },
    {
      key: "full_name",
      header: "Full Name",
      type: "text",
    },
    {
      key: "age",
      header: "Age",
      type: "text",
    },
    {
      key: "gender",
      header: "Gender",
      type: "custom",
      render: (value) => formatGender(String(value ?? "")),
    },
    {
      key: "phone",
      header: "Phone",
      type: "text",
    },
    {
      key: "email",
      header: "Email",
      type: "text",
    },
    {
      key: "insurance_provider",
      header: "Insurance Provider",
      type: "text",
    },
    {
      key: "created_at",
      header: "Created At",
      type: "custom",
      render: (value) => formatDate(String(value)),
    },
  ];

  return (
    <div className="space-y-4">
      {viewMode === "deleted" ? (
        <DeletedPatientsIndexLayout
          onBackToPatients={() => setViewMode("active")}
        />
      ) : (
        <>
          <PageBreadcrumb pageTitle="Patients" />
          <BaseTable
            data={patients}
            columns={patientColumns}
            loading={loading}
            error={error}
            currentPage={currentPage}
            rowsPerPage={rowsPerPage}
            search={search}
            totalItems={meta.total_items}
            totalPages={meta.total_pages}
            onSearchChange={handleSearchChange}
            onPageChange={handlePageChange}
            onRowsPerPageChange={handleRowsPerPageChange}
            onCreate={() => setIsCreateModalOpen(true)}
            onEdit={handleEditPatient}
            onDelete={handleDeletePatient}
            onView={handleViewPatient}
            searchPlaceholder="Search patients by name, code, phone, or email..."
            showDeleteButton={true}
            onDeleteAll={() => setViewMode("deleted")}
          />

          <CreatePatientModal
            isOpen={isCreateModalOpen}
            onClose={() => setIsCreateModalOpen(false)}
            onSuccess={handleCreatePatientSuccess}
          />

          <EditPatientModal
            isOpen={isEditModalOpen}
            id={selectedEditPatientId}
            onClose={handleEditPatientClose}
            onSuccess={handleEditPatientSuccess}
          />

          {isShowModalOpen && (
            <ShowPatientModal
              isOpen={isShowModalOpen}
              id={selectedPatientId}
              onClose={() => setIsShowModalOpen(false)}
            />
          )}

          <DeleteModal
            isOpen={isDeleteModalOpen}
            itemName={selectedDeletePatientName}
            onClose={() => {
              setIsDeleteModalOpen(false);
              setSelectedDeletePatientId(undefined);
              setSelectedDeletePatientName(undefined);
            }}
            onConfirm={async () => {
              if (!selectedDeletePatientId) throw new Error("Missing id");
              await del(`/patients/${selectedDeletePatientId}`);
              setIsDeleteModalOpen(false);
              setSelectedDeletePatientId(undefined);
              triggerFetch(currentPage, rowsPerPage, search);
              setSuccessMessage(
                `Patient "${selectedDeletePatientName ?? ""}" has been successfully deleted.`,
              );
              setShowSuccessModal(true);
            }}
          />

          <SuccessModal
            title="Success"
            message={successMessage}
            buttonLabel="Close"
            isOpen={showSuccessModal}
            onButtonClick={() => setShowSuccessModal(false)}
          />
        </>
      )}
    </div>
  );
}

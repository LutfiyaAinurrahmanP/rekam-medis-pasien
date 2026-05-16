import { useCallback, useEffect, useRef, useState } from "react";
import PageBreadcrumb from "../../../components/common/PageBreadCrumb";
import BaseTable, {
  type ColumnDefinition,
} from "../../../components/tables/BaseTable";

import DeleteModal from "../../../components/modals/DeleteModal";
import SuccessModal from "../../../components/ui/notification/SuccessModal";
import { del } from "../../../services/api";
import {
  useDepartments,
  type Department,
} from "../../../hooks/Departments/useDepartments";
import CreateDepartmentModal from "./CreateDepartmentModal";
import EditDepartmentModal from "./EditDepartmentModal";
import ShowDepartmentModal from "./ShowDepartmentModal";
import DeletedDepartmentsIndexLayout from "./DeletedDepartmentsIndexLayout";

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString("en-GB", {
    day: "numeric",
    month: "long",
    year: "numeric",
  });
};

export default function DepartmentsIndexLayout() {
  const { departments, loading, error, meta, fetchDepartments } =
    useDepartments();
  const [viewMode, setViewMode] = useState<"active" | "deleted">("active");
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [search, setSearch] = useState("");
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [selectedEditId, setSelectedEditId] = useState<string | undefined>(
    undefined,
  );
  const [isShowModalOpen, setIsShowModalOpen] = useState(false);
  const [selectedShowId, setSelectedShowId] = useState<string | undefined>(
    undefined,
  );
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [selectedDeleteId, setSelectedDeleteId] = useState<string | undefined>(
    undefined,
  );
  const [selectedDeleteName, setSelectedDeleteName] = useState<
    string | undefined
  >(undefined);
  const [successMessage, setSuccessMessage] = useState("");
  const [showSuccessModal, setShowSuccessModal] = useState(false);
  const debounceRef = useRef<number | null>(null);
  const skipNextFetchRef = useRef(false);

  const triggerFetch = useCallback(
    (page: number, pageSize: number, searchValue: string) => {
      fetchDepartments({
        page,
        page_size: pageSize,
        search: searchValue.trim() || undefined,
      });
    },
    [fetchDepartments],
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
    if (!showSuccessModal) return undefined;
    const t = window.setTimeout(() => {
      setShowSuccessModal(false);
      setSuccessMessage("");
    }, 3000);
    return () => window.clearTimeout(t);
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
    if (page >= 1 && page <= meta.total_pages) setCurrentPage(page);
  };

  const handleRowsPerPageChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const next = Number.parseInt(e.target.value, 10);
    setRowsPerPage(next);
    setCurrentPage(1);
  };

  const handleCreate = () => setIsCreateModalOpen(true);

  const handleCreateSuccess = () => {
    setIsCreateModalOpen(false);
    triggerFetch(1, rowsPerPage, search);
  };

  const handleEdit = (row: Department) => {
    setSelectedEditId(String(row.id));
    setIsEditModalOpen(true);
  };

  const handleEditSuccess = () => {
    setIsEditModalOpen(false);
    setSelectedEditId(undefined);
    triggerFetch(currentPage, rowsPerPage, search);
  };

  const handleView = (row: Department) => {
    setSelectedShowId(String(row.id));
    setIsShowModalOpen(true);
  };

  const handleDelete = (row: Department) => {
    setSelectedDeleteId(String(row.id));
    setSelectedDeleteName(row.name);
    setIsDeleteModalOpen(true);
  };

  const handleConfirmDelete = async () => {
    if (!selectedDeleteId) return;
    await del(`/departments/${selectedDeleteId}`);
    triggerFetch(currentPage, rowsPerPage, search);
    setSuccessMessage(
      `Department "${selectedDeleteName ?? ""}" has been deleted.`,
    );
    setShowSuccessModal(true);
    setIsDeleteModalOpen(false);
    setSelectedDeleteId(undefined);
    setSelectedDeleteName(undefined);
  };

  const columns: ColumnDefinition<Department>[] = [
    {
      key: "name",
      header: "Name",
      type: "text",
    },
    {
      key: "code",
      header: "Code",
      type: "text",
    },
    {
      key: "description",
      header: "Description",
      type: "custom",
      render: (value) => {
        const str = String(value ?? "");
        return str.length > 100 ? `${str.slice(0, 100)}...` : str;
      },
    },
    {
      key: "floor_location",
      header: "Floor/Location",
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
      <PageBreadcrumb pageTitle="Departments" />

      {viewMode === "deleted" ? (
        <DeletedDepartmentsIndexLayout
          onBackToDepartments={() => setViewMode("active")}
        />
      ) : (
        <BaseTable
          data={departments}
          columns={columns}
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
          onCreate={handleCreate}
          onEdit={handleEdit}
          onDelete={handleDelete}
          onView={handleView}
          searchPlaceholder="Search departments..."
          onDeleteAll={() => setViewMode("deleted")}
        />
      )}

      <CreateDepartmentModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onSuccess={handleCreateSuccess}
      />

      <EditDepartmentModal
        isOpen={isEditModalOpen}
        id={selectedEditId}
        onClose={() => {
          setIsEditModalOpen(false);
          setSelectedEditId(undefined);
        }}
        onSuccess={handleEditSuccess}
      />

      {isShowModalOpen && (
        <ShowDepartmentModal
          isOpen={isShowModalOpen}
          id={selectedShowId}
          onClose={() => setIsShowModalOpen(false)}
        />
      )}

      <DeleteModal
        isOpen={isDeleteModalOpen}
        itemName={selectedDeleteName}
        onClose={() => setIsDeleteModalOpen(false)}
        onConfirm={handleConfirmDelete}
      />

      <SuccessModal
        isOpen={showSuccessModal}
        title="Success"
        message={successMessage}
        buttonLabel="Close"
        onButtonClick={() => {
          setShowSuccessModal(false);
          setSuccessMessage("");
        }}
      />
    </div>
  );
}

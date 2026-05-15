import { useCallback, useEffect, useRef, useState } from "react";
import { useParams } from "react-router";
import PageBreadcrumb from "../../../components/common/PageBreadCrumb";
import CreateUserModal from "./CreateUserModal";
import DeletedUsersIndexLayout from "./DeletedUsersIndexLayout";
import EditUserModal from "./EditUserModal";
import BaseTable, {
  type ColumnDefinition,
} from "../../../components/tables/BaseTable";
import ShowUserModal from "./ShowUserModal";
import DeleteModal from "../../../components/modals/DeleteModal";
import SuccessModal from "../../../components/ui/notification/SuccessModal";
import StatusToggleModal from "../../../components/modals/StatusToggleModal";
import { del, patch } from "../../../services/api";
import { useUsers, type User } from "../../../hooks/Users/useUsers";

const getRoleColor = (
  role: string,
): "error" | "warning" | "success" | "info" | "light" => {
  switch (role) {
    case "admin":
      return "error";
    case "doctor":
      return "warning";
    case "patient":
      return "success";
    case "receptionist":
      return "info";
    case "super_admin":
      return "error";
    default:
      return "light";
  }
};

const getRoleLabel = (role: string) => {
  return role.charAt(0).toUpperCase() + role.slice(1).replace(/_/g, " ");
};

const getStatusButtonStyle = (isActive: boolean) =>
  isActive
    ? "border border-success-300 bg-success-50 text-success-700 hover:bg-success-100 dark:border-success-700 dark:bg-success-500/10 dark:text-success-300"
    : "border border-gray-300 bg-gray-50 text-gray-700 hover:bg-gray-100 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300";

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString("en-GB", {
    day: "numeric",
    month: "long",
    year: "numeric",
  });
};

export default function UsersIndexLayout() {
  const { role } = useParams();
  const [viewMode, setViewMode] = useState<"active" | "deleted">("active");
  const { users, loading, error, meta, fetchUsers } = useUsers();
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [search, setSearch] = useState("");
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [selectedEditUserId, setSelectedEditUserId] = useState<
    string | undefined
  >(undefined);
  const [isShowModalOpen, setIsShowModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [selectedUserId, setSelectedUserId] = useState<string | undefined>(
    undefined,
  );
  const [selectedDeleteUserId, setSelectedDeleteUserId] = useState<
    string | undefined
  >(undefined);
  const [selectedDeleteUserName, setSelectedDeleteUserName] = useState<
    string | undefined
  >(undefined);
  const [selectedStatusUser, setSelectedStatusUser] = useState<User | null>(
    null,
  );
  const [successMessage, setSuccessMessage] = useState("");
  const [showSuccessModal, setShowSuccessModal] = useState(false);
  const [showStatusModal, setShowStatusModal] = useState(false);
  const debounceRef = useRef<number | null>(null);
  const skipNextFetchRef = useRef(false);

  const triggerFetch = useCallback(
    (page: number, pageSize: number, searchValue: string) => {
      fetchUsers({
        page,
        page_size: pageSize,
        search: searchValue.trim() || undefined,
      });
    },
    [fetchUsers],
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

  const handleCreateUser = () => {
    setIsCreateModalOpen(true);
  };

  const handleCreateUserSuccess = () => {
    setIsCreateModalOpen(false);
    triggerFetch(1, rowsPerPage, search);
  };

  const handleEditUser = (user: User) => {
    setSelectedEditUserId(String(user.id));
    setIsEditModalOpen(true);
  };

  const handleEditUserClose = () => {
    setIsEditModalOpen(false);
    setSelectedEditUserId(undefined);
  };

  const handleEditUserSuccess = () => {
    handleEditUserClose();
    triggerFetch(currentPage, rowsPerPage, search);
  };

  const handleDeleteUser = (user: User) => {
    setSelectedDeleteUserId(String(user.id));
    setSelectedDeleteUserName(String(user.username ?? user.email ?? ""));
    setIsDeleteModalOpen(true);
  };

  const handleViewUser = (user: User) => {
    setSelectedUserId(String(user.id));
    setIsShowModalOpen(true);
  };

  const handleToggleUserStatus = (user: User) => {
    setSelectedStatusUser(user);
    setShowStatusModal(true);
  };

  const handleConfirmToggleUserStatus = async () => {
    if (!selectedStatusUser) {
      throw new Error("Missing user");
    }

    const nextAction = selectedStatusUser.is_active ? "deactivate" : "activate";
    await patch(`/users/${selectedStatusUser.id}/${nextAction}`, {});
    setSuccessMessage(
      `User "${selectedStatusUser.username}" has been successfully ${
        selectedStatusUser.is_active ? "deactivated" : "activated"
      }.`,
    );
    setShowStatusModal(false);
    setSelectedStatusUser(null);
    triggerFetch(currentPage, rowsPerPage, search);
    setShowSuccessModal(true);
  };

  const userColumns: ColumnDefinition<User>[] = [
    {
      key: "username",
      header: "User",
      type: "custom",
      render: (value, row) => (
        <div>
          <p className="block font-medium text-gray-800 text-theme-sm dark:text-white/90">
            {String(value)}
          </p>
          <span className="text-sm font-normal text-gray-500 dark:text-gray-400">
            {(row as User).email}
          </span>
        </div>
      ),
    },
    {
      key: "phone",
      header: "Phone",
      type: "text",
    },
    {
      key: "role",
      header: "Role",
      type: "badge",
      badgeColorMap: (value) => getRoleColor(String(value)),
      badgeLabel: (value) => getRoleLabel(String(value)),
    },
    {
      key: "is_active",
      header: "Status",
      type: "custom",
      render: (value, row) => {
        const isActive = Boolean(value);

        return (
          <button
            type="button"
            onClick={() => handleToggleUserStatus(row as User)}
            className={`inline-flex items-center rounded-full px-3 py-1 text-xs font-medium transition-colors ${getStatusButtonStyle(isActive)}`}
          >
            {isActive ? "Active" : "Inactive"}
          </button>
        );
      },
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
        <DeletedUsersIndexLayout onBackToUsers={() => setViewMode("active")} />
      ) : (
        <>
          <PageBreadcrumb pageTitle="Users" />
          <BaseTable
            data={users}
            columns={userColumns}
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
            onCreate={handleCreateUser}
            onEdit={handleEditUser}
            onDelete={handleDeleteUser}
            onView={handleViewUser}
            searchPlaceholder="Search users..."
            showDeleteButton={true}
            onDeleteAll={() => setViewMode("deleted")}
          />

          <CreateUserModal
            isOpen={isCreateModalOpen}
            onClose={() => setIsCreateModalOpen(false)}
            onSuccess={handleCreateUserSuccess}
            role={role}
          />

          <EditUserModal
            isOpen={isEditModalOpen}
            id={selectedEditUserId}
            role={role}
            onClose={handleEditUserClose}
            onSuccess={handleEditUserSuccess}
          />

          {isShowModalOpen && (
            <ShowUserModal
              isOpen={isShowModalOpen}
              id={selectedUserId}
              role={role}
              onClose={() => setIsShowModalOpen(false)}
            />
          )}

          <DeleteModal
            isOpen={isDeleteModalOpen}
            itemName={selectedDeleteUserName}
            onClose={() => {
              setIsDeleteModalOpen(false);
              setSelectedDeleteUserId(undefined);
              setSelectedDeleteUserName(undefined);
            }}
            onConfirm={async () => {
              if (!selectedDeleteUserId) throw new Error("Missing id");
              await del(`/users/${selectedDeleteUserId}`);
              triggerFetch(currentPage, rowsPerPage, search);
              setSuccessMessage(
                `User "${selectedDeleteUserName ?? ""}" has been successfully deleted.`,
              );
              setShowSuccessModal(true);
            }}
          />

          <StatusToggleModal
            isOpen={showStatusModal}
            actionLabel={
              selectedStatusUser?.is_active
                ? `Deactivate user "${selectedStatusUser?.username ?? ""}"?`
                : `Activate user "${selectedStatusUser?.username ?? ""}"?`
            }
            confirmLabel={
              selectedStatusUser?.is_active ? "Deactivate" : "Activate"
            }
            onClose={() => {
              setShowStatusModal(false);
              setSelectedStatusUser(null);
            }}
            onConfirm={handleConfirmToggleUserStatus}
          />

          <SuccessModal
            title="Success"
            message={successMessage}
            buttonLabel="Close"
            isOpen={showSuccessModal}
            onButtonClick={() => {
              setShowSuccessModal(false);
              setSuccessMessage("");
            }}
          />
        </>
      )}
    </div>
  );
}

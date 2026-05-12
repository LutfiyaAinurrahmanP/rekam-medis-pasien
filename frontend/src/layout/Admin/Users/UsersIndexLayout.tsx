import { useCallback, useEffect, useRef, useState } from "react";
import { useNavigate, useParams } from "react-router";
import PageBreadcrumb from "../../../components/common/PageBreadCrumb";
import CreateUserModal from "./CreateUserModal";
import BaseTable, {
  type ColumnDefinition,
} from "../../../components/tables/BaseTable";
import { useUsers, type User } from "../../../hooks/Users/useUsers";
import { getRoleUsersPath } from "../../../pages/Roles/shared/role-routing";

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

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString("en-GB", {
    day: "numeric",
    month: "long",
    year: "numeric",
  });
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
    type: "badge",
    badgeColorMap: (value) => (value ? "success" : "error"),
    badgeLabel: (value) => (value ? "Active" : "Inactive"),
  },
  {
    key: "created_at",
    header: "Created At",
    type: "custom",
    render: (value) => formatDate(String(value)),
  },
];

export default function UsersIndexLayout() {
  const navigate = useNavigate();
  const { role } = useParams();
  const { users, loading, error, meta, fetchUsers } = useUsers();
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [search, setSearch] = useState("");
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const debounceRef = useRef<number | null>(null);
  const skipNextFetchRef = useRef(false);

  const baseUsersPath = getRoleUsersPath(role ?? undefined);

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
    // Refetch users after successful creation
    triggerFetch(1, rowsPerPage, search);
  };

  const handleEditUser = (user: User) => {
    navigate(`${baseUsersPath}/${user.id}/edit`);
  };

  const handleDeleteUser = (user: User) => {
    console.log("Delete user:", user);
  };

  const handleDeleteAllUsers = () => {
    console.log("Delete all users");
  };

  return (
    <div className="space-y-4">
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
        // Optional header customizations
        searchPlaceholder="Search users..."
        showDeleteButton={true}
        onDeleteAll={handleDeleteAllUsers}
      />

      {/* Create User Modal */}
      <CreateUserModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onSuccess={handleCreateUserSuccess}
        role={role}
      />
    </div>
  );
}

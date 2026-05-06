import { useState, useEffect, useRef, useCallback } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
} from "../ui/table";
import { EyeIcon, PencilIcon, TrashBinIcon } from "../../icons";
import Badge from "../ui/badge/Badge";
import Pagination from "../tables/DataTables/TableThree/Pagination";
import Button from "../ui/button/Button";
import { useUsers, User } from "../../hooks/Users/useUsers";

interface UserTableProps {
  onEdit?: (user: User) => void;
  onDelete?: (user: User) => void;
}

const getRoleColor = (role: string) => {
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

export default function UserTable({ onEdit, onDelete }: UserTableProps) {
  const { users, loading, error, meta, fetchUsers } = useUsers();
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [search, setSearch] = useState("");
  const debounceRef = useRef<number | null>(null);
  const isInitialLoading = loading && users.length === 0;

  // ✅ Fungsi fetch terpusat
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

  // ✅ Fetch saat page atau rowsPerPage berubah (tanpa debounce)
  useEffect(() => {
    triggerFetch(currentPage, rowsPerPage, search);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [currentPage, rowsPerPage]);

  // ✅ Cleanup debounce saat unmount
  useEffect(() => {
    return () => {
      if (debounceRef.current) window.clearTimeout(debounceRef.current);
    };
  }, []);

  // ✅ Jaga currentPage tidak melebihi total_pages
  useEffect(() => {
    if (meta.total_pages > 0 && currentPage > meta.total_pages) {
      setCurrentPage(meta.total_pages);
    }
  }, [currentPage, meta.total_pages]);

  // ✅ Search dengan debounce 150ms langsung panggil triggerFetch
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
    if (page >= 1 && page <= meta.total_pages) {
      setCurrentPage(page);
    }
  };

  const handleRowsPerPageChange = (
    e: React.ChangeEvent<HTMLSelectElement>,
  ): void => {
    const newRowsPerPage = parseInt(e.target.value, 10);
    setRowsPerPage(newRowsPerPage);
    setCurrentPage(1);
  };

  const startIndex = (currentPage - 1) * rowsPerPage;
  const endIndex = Math.min(startIndex + rowsPerPage, meta.total_items);

  if (error) {
    return (
      <div className="flex items-center justify-center p-8">
        <p className="text-red-500">Error: {error}</p>
      </div>
    );
  }

  return (
    <div className="overflow-hidden rounded-xl bg-white dark:bg-white/[0.03]">
      {/* Header Section */}
      <div className="flex flex-col gap-2 px-4 py-4 border border-b-0 border-gray-100 dark:border-white/[0.05] rounded-t-xl sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3">
          <span className="text-gray-500 dark:text-gray-400">Show</span>
          <div className="relative z-20 bg-transparent">
            <select
              className="w-full py-2 pl-3 pr-8 text-sm text-gray-800 bg-transparent border border-gray-300 rounded-lg appearance-none dark:bg-dark-900 h-9 shadow-theme-xs placeholder:text-gray-400 focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90 dark:placeholder:text-white/30 dark:focus:border-brand-800"
              value={rowsPerPage}
              onChange={handleRowsPerPageChange}
            >
              <option value="10">10</option>
              <option value="20">20</option>
              <option value="50">50</option>
            </select>
            <span className="absolute z-30 text-gray-500 -translate-y-1/2 right-2 top-1/2 dark:text-gray-400">
              <svg
                className="stroke-current"
                width="16"
                height="16"
                viewBox="0 0 16 16"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M3.8335 5.9165L8.00016 10.0832L12.1668 5.9165"
                  stroke=""
                  strokeWidth="1.2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              </svg>
            </span>
          </div>
          <span className="text-gray-500 dark:text-gray-400">entries</span>
        </div>

        <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
          <div className="relative">
            <button className="absolute text-gray-500 -translate-y-1/2 left-4 top-1/2 dark:text-gray-400">
              <svg
                className="fill-current"
                width="20"
                height="20"
                viewBox="0 0 20 20"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  fillRule="evenodd"
                  clipRule="evenodd"
                  d="M3.04199 9.37363C3.04199 5.87693 5.87735 3.04199 9.37533 3.04199C12.8733 3.04199 15.7087 5.87693 15.7087 9.37363C15.7087 12.8703 12.8733 15.7053 9.37533 15.7053C5.87735 15.7053 3.04199 12.8703 3.04199 9.37363ZM9.37533 1.54199C5.04926 1.54199 1.54199 5.04817 1.54199 9.37363C1.54199 13.6991 5.04926 17.2053 9.37533 17.2053C11.2676 17.2053 13.0032 16.5344 14.3572 15.4176L17.1773 18.238C17.4702 18.5309 17.945 18.5309 18.2379 18.238C18.5308 17.9451 18.5309 17.4703 18.238 17.1773L15.4182 14.3573C16.5367 13.0033 17.2087 11.2669 17.2087 9.37363C17.2087 5.04817 13.7014 1.54199 9.37533 1.54199Z"
                  fill=""
                />
              </svg>
            </button>
            <input
              type="text"
              placeholder="Search..."
              value={search}
              onChange={handleSearchChange}
              className="dark:bg-dark-900 h-11 w-full rounded-lg border border-gray-300 bg-transparent py-2.5 pl-11 pr-4 text-sm text-gray-800 shadow-theme-xs placeholder:text-gray-400 focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90 dark:placeholder:text-white/30 dark:focus:border-brand-800 xl:w-[300px]"
            />
          </div>
          <Button variant="outline" size="sm">
            Create
          </Button>
          <Button variant="outline" size="sm">
            Delete Users
          </Button>
        </div>
      </div>

      {/* Table Section */}
      <div className="max-w-full overflow-x-auto custom-scrollbar">
        <Table>
          <TableHeader>
            <TableRow>
              <TableCell
                isHeader
                className="px-4 py-3 border border-gray-100 dark:border-white/[0.05]"
              >
                <div className="flex items-center gap-3">
                  <span className="font-medium text-gray-700 text-theme-xs dark:text-gray-400">
                    User
                  </span>
                </div>
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
                  Role
                </p>
              </TableCell>
              <TableCell
                isHeader
                className="px-4 py-3 border border-gray-100 dark:border-white/[0.05]"
              >
                <p className="font-medium text-gray-700 text-theme-xs dark:text-gray-400">
                  Status
                </p>
              </TableCell>
              <TableCell
                isHeader
                className="px-4 py-3 border border-gray-100 dark:border-white/[0.05]"
              >
                <p className="font-medium text-gray-700 text-theme-xs dark:text-gray-400">
                  Created At
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
            {isInitialLoading ? (
              <TableRow>
                <td
                  colSpan={6}
                  className="px-4 py-8 text-center text-gray-500 dark:text-gray-400"
                >
                  Loading users...
                </td>
              </TableRow>
            ) : users.length > 0 ? (
              users.map((user, index) => (
                <TableRow key={index}>
                  <TableCell className="px-4 py-4 border border-gray-100 dark:border-white/[0.05] dark:text-white/90 whitespace-nowrap">
                    <div>
                      <p className="block font-medium text-gray-800 text-theme-sm dark:text-white/90">
                        {user.username}
                      </p>
                      <span className="text-sm font-normal text-gray-500 dark:text-gray-400">
                        {user.email}
                      </span>
                    </div>
                  </TableCell>
                  <TableCell className="px-4 py-4 font-normal text-gray-800 border border-gray-100 dark:border-white/[0.05] text-theme-sm dark:text-gray-400 whitespace-nowrap">
                    {user.phone}
                  </TableCell>
                  <TableCell className="px-4 py-4 font-normal text-gray-800 border border-gray-100 dark:border-white/[0.05] text-theme-sm dark:text-white/90 whitespace-nowrap">
                    <Badge size="sm" color={getRoleColor(user.role)}>
                      {getRoleLabel(user.role)}
                    </Badge>
                  </TableCell>
                  <TableCell className="px-4 py-4 font-normal text-gray-800 border border-gray-100 dark:border-white/[0.05] text-theme-sm dark:text-white/90 whitespace-nowrap">
                    <Badge
                      size="sm"
                      color={user.is_active ? "success" : "error"}
                    >
                      {user.is_active ? "Active" : "Inactive"}
                    </Badge>
                  </TableCell>
                  <TableCell className="px-4 py-4 font-normal text-gray-800 border border-gray-100 dark:border-white/[0.05] text-theme-sm dark:text-white/90 whitespace-nowrap">
                    {new Date(user.created_at).toLocaleDateString()}
                  </TableCell>
                  <TableCell className="px-4 py-4 font-normal text-gray-800 border border-gray-100 dark:border-white/[0.05] text-theme-sm dark:text-white/90 whitespace-nowrap">
                    <div className="flex items-center w-full gap-2">
                      <button className="text-gray-500 hover:text-gray-800 dark:text-gray-400 dark:hover:text-white/90">
                        <EyeIcon className="fill-gray-500 dark:fill-gray-400 size-5" />
                      </button>
                      <button
                        onClick={() => onEdit?.(user)}
                        className="text-gray-500 hover:text-gray-800 dark:text-gray-400 dark:hover:text-white/90"
                      >
                        <PencilIcon className="size-5" />
                      </button>
                      <button
                        onClick={() => onDelete?.(user)}
                        className="text-gray-500 hover:text-error-500 dark:text-gray-400 dark:hover:text-error-500"
                      >
                        <TrashBinIcon className="size-5" />
                      </button>
                    </div>
                  </TableCell>
                </TableRow>
              ))
            ) : (
              <TableRow>
                <td
                  colSpan={6}
                  className="px-4 py-8 text-center text-gray-500 dark:text-gray-400"
                >
                  No users found
                </td>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      {/* Footer Section */}
      <div className="border border-t-0 rounded-b-xl border-gray-100 py-4 pl-[18px] pr-4 dark:border-white/[0.05]">
        <div className="flex flex-col xl:flex-row xl:items-center xl:justify-between">
          <div className="pb-3 xl:pb-0">
            <p className="pb-3 text-sm font-medium text-center text-gray-500 border-b border-gray-100 dark:border-gray-800 dark:text-gray-400 xl:border-b-0 xl:pb-0 xl:text-left">
              Showing {startIndex + 1} to {endIndex} of {meta.total_items}{" "}
              entries
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
  );
}

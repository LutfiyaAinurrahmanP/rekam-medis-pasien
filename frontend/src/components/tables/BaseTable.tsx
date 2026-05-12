import { ReactNode } from "react";
import { EyeIcon, PencilIcon, TrashBinIcon } from "../../icons";
import {
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
} from "../ui/table";
import Badge from "../ui/badge/Badge";
import Pagination from "./DataTables/TableThree/Pagination";
import TableHeaderComponent from "./TableHeader";

type BadgeColor = "error" | "success" | "warning" | "info" | "light";

export interface ColumnDefinition<T = object> {
  key: keyof T;
  header: string;
  type: "text" | "badge" | "custom";
  width?: string;
  badgeColorMap?: (value: unknown) => BadgeColor;
  badgeLabel?: (value: unknown) => string;
  render?: (value: unknown, row: T) => ReactNode;
}

interface BaseTableProps<T = object> {
  data: T[];
  columns: ColumnDefinition<T>[];
  loading: boolean;
  error: string | null;
  currentPage: number;
  rowsPerPage: number;
  search: string;
  totalItems: number;
  totalPages: number;
  onSearchChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onPageChange: (page: number) => void;
  onRowsPerPageChange: (e: React.ChangeEvent<HTMLSelectElement>) => void;
  onCreate: () => void;
  onEdit?: (row: T) => void;
  onDelete?: (row: T) => void;
  // Header customization (optional)
  searchPlaceholder?: string;
  showDeleteButton?: boolean;
  onDeleteAll?: () => void;
  actionButtons?: ReactNode;
}

export default function BaseTable<T extends object = object>({
  data,
  columns,
  loading,
  error,
  currentPage,
  rowsPerPage,
  search,
  totalItems,
  totalPages,
  onSearchChange,
  onPageChange,
  onRowsPerPageChange,
  onCreate,
  onEdit,
  onDelete,
  searchPlaceholder = "Search...",
  showDeleteButton = true,
  onDeleteAll,
  actionButtons,
}: BaseTableProps<T>) {
  const isInitialLoading = loading && data.length === 0;
  const startIndex = totalItems === 0 ? 0 : (currentPage - 1) * rowsPerPage;
  const endIndex = Math.min(startIndex + rowsPerPage, totalItems);

  const renderCellContent = (
    column: ColumnDefinition<T>,
    row: T,
  ): ReactNode => {
    const value = row[column.key];

    if (column.type === "custom" && column.render) {
      return column.render(value, row);
    }

    if (column.type === "badge") {
      const badgeColor = column.badgeColorMap
        ? column.badgeColorMap(value)
        : "light";
      const badgeLabel = column.badgeLabel
        ? column.badgeLabel(value)
        : String(value);
      return (
        <Badge size="sm" color={badgeColor}>
          {badgeLabel}
        </Badge>
      );
    }

    return String(value);
  };

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
      <TableHeaderComponent
        search={search}
        rowsPerPage={rowsPerPage}
        onSearchChange={onSearchChange}
        onRowsPerPageChange={onRowsPerPageChange}
        onCreate={onCreate}
        onDeleteAll={onDeleteAll}
        showDeleteButton={showDeleteButton}
        actionButtons={actionButtons}
        searchPlaceholder={searchPlaceholder}
      />

      {/* Table Section */}
      <div className="max-w-full overflow-x-auto custom-scrollbar">
        <Table>
          <TableHeader>
            <TableRow>
              {columns.map((column) => (
                <TableCell
                  key={String(column.key)}
                  isHeader
                  className="px-4 py-3 border border-gray-100 dark:border-white/[0.05]"
                >
                  <p className="font-medium text-gray-700 text-theme-xs dark:text-gray-400">
                    {column.header}
                  </p>
                </TableCell>
              ))}
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
                  colSpan={columns.length + 1}
                  className="px-4 py-8 text-center text-gray-500 dark:text-gray-400"
                >
                  Loading...
                </td>
              </TableRow>
            ) : data.length > 0 ? (
              data.map((row, index) => (
                <TableRow key={index}>
                  {columns.map((column) => (
                    <TableCell
                      key={String(column.key)}
                      className="px-4 py-4 font-normal text-gray-800 border border-gray-100 dark:border-white/[0.05] text-theme-sm dark:text-white/90 whitespace-nowrap"
                    >
                      {renderCellContent(column, row)}
                    </TableCell>
                  ))}
                  <TableCell className="px-4 py-4 font-normal text-gray-800 border border-gray-100 dark:border-white/[0.05] text-theme-sm dark:text-white/90 whitespace-nowrap">
                    <div className="flex items-center w-full gap-2">
                      <button className="text-gray-500 hover:text-gray-800 dark:text-gray-400 dark:hover:text-white/90">
                        <EyeIcon className="fill-gray-500 dark:fill-gray-400 size-5" />
                      </button>
                      <button
                        onClick={() => onEdit?.(row)}
                        className="text-gray-500 hover:text-gray-800 dark:text-gray-400 dark:hover:text-white/90"
                      >
                        <PencilIcon className="size-5" />
                      </button>
                      <button
                        onClick={() => onDelete?.(row)}
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
                  colSpan={columns.length + 1}
                  className="px-4 py-8 text-center text-gray-500 dark:text-gray-400"
                >
                  No data found
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
              Showing {startIndex + 1} to {endIndex} of {totalItems} entries
            </p>
          </div>
          <Pagination
            currentPage={currentPage}
            totalPages={totalPages}
            onPageChange={onPageChange}
          />
        </div>
      </div>
    </div>
  );
}

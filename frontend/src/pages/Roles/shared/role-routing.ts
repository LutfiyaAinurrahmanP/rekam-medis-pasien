export type RoleSlug = "admin" | "doctor" | "receptionist" | "patient" | "super-admin";

const roleSlugSet = new Set<RoleSlug>([
  "admin",
  "doctor",
  "receptionist",
  "patient",
  "super-admin",
]);

function normalizeRoleSegment(value: string) {
  return value.trim().toLowerCase().replace(/[_\s]+/g, "-");
}

export function normalizeRole(value?: string | null): RoleSlug | null {
  if (!value) {
    return null;
  }

  const normalized = normalizeRoleSegment(value);

  if (normalized === "superadmin") {
    return "super-admin";
  }

  if (roleSlugSet.has(normalized as RoleSlug)) {
    return normalized as RoleSlug;
  }

  return null;
}

export function getRoleDashboardPath(role?: string | null) {
  const normalizedRole = normalizeRole(role);

  if (!normalizedRole) {
    return "/admin/dashboard";
  }

  return `/${normalizedRole}/dashboard`;
}

export function getRoleReportsPath(role?: string | null) {
  const normalizedRole = normalizeRole(role);

  if (!normalizedRole) {
    return "/admin/reports";
  }

  return `/${normalizedRole}/reports`;
}

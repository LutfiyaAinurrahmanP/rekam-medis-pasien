/**
 * Role-based routing utilities
 */

const normalizeRolePathSegment = (role: string) =>
  role.trim().toLowerCase().replace(/_/g, "-");

/**
 * Get dashboard path based on user role
 * @param role User role (admin, doctor, patient, receptionist, super-admin)
 * @returns Dashboard path for the role
 */
export function getRoleDashboardPath(role?: string): string {
  if (!role) {
    return "/patient/dashboard"; // Default fallback
  }

  return `/${normalizeRolePathSegment(role)}/dashboard`;
}

/**
 * Get users management path based on user role
 * @param role User role
 * @returns Users management path for the role
 */
export function getRoleUsersPath(role?: string): string {
  if (!role) {
    return "/patient/users";
  }

  return `/${normalizeRolePathSegment(role)}/users`;
}

/**
 * Get reports path based on user role
 * @param role User role
 * @returns Reports path for the role
 */
export function getRoleReportsPath(role?: string): string {
  if (!role) {
    return "/patient/reports";
  }

  return `/${normalizeRolePathSegment(role)}/reports`;
}

/**
 * Get departments management path based on user role
 * @param role User role
 * @returns Departments management path for the role
 */
export function getRoleDepartmentsPath(role?: string): string {
  if (!role) {
    return "/admin/departments";
  }

  return `/${normalizeRolePathSegment(role)}/departments`;
}

/**
 * Get patients management path based on user role
 * @param role User role
 * @returns Patients management path for the role
 */
export function getRolePatientsPath(role?: string): string {
  if (!role) {
    return "/admin/patients";
  }

  return `/${normalizeRolePathSegment(role)}/patients`;
}

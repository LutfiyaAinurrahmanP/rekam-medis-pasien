/**
 * Role-based routing utilities
 */

/**
 * Get dashboard path based on user role
 * @param role User role (admin, doctor, patient, receptionist, super_admin)
 * @returns Dashboard path for the role
 */
export function getRoleDashboardPath(role?: string): string {
  if (!role) {
    return "/admin/dashboard"; // Default fallback
  }

  return `/${role}/dashboard`;
}

/**
 * Get users management path based on user role
 * @param role User role
 * @returns Users management path for the role
 */
export function getRoleUsersPath(role?: string): string {
  if (!role) {
    return "/admin/users";
  }

  return `/${role}/users`;
}

/**
 * Get reports path based on user role
 * @param role User role
 * @returns Reports path for the role
 */
export function getRoleReportsPath(role?: string): string {
  if (!role) {
    return "/admin/reports";
  }

  return `/${role}/reports`;
}

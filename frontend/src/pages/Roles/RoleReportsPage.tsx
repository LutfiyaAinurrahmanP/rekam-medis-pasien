import { Navigate, useParams } from "react-router";
import type { ComponentType } from "react";
import { useAuth } from "../../context/AuthContext";
import AdminReports from "./Admin/Reports";
import DoctorReports from "./Doctor/Reports";
import ReceptionistReports from "./Receptionist/Reports";
import SuperAdminReports from "./SuperAdmin/Reports";
import { getRoleDashboardPath, getRoleReportsPath, normalizeRole, type RoleSlug } from "./shared/role-routing";

const reports: Partial<Record<RoleSlug, ComponentType>> = {
  admin: AdminReports,
  doctor: DoctorReports,
  receptionist: ReceptionistReports,
  "super-admin": SuperAdminReports,
};

export default function RoleReportsPage() {
  const { role } = useParams();
  const { user, loading } = useAuth();

  if (loading) {
    return null;
  }

  const routeRole = normalizeRole(role);
  const userRole = normalizeRole(user?.role);

  if (!userRole) {
    return <Navigate to="/auth/login" replace />;
  }

  if (!routeRole) {
    return <Navigate to={getRoleDashboardPath(userRole)} replace />;
  }

  if (routeRole !== userRole) {
    return <Navigate to={getRoleReportsPath(userRole)} replace />;
  }

  if (routeRole === "patient") {
    return <Navigate to={getRoleDashboardPath(userRole)} replace />;
  }

  const Reports = reports[routeRole];

  if (!Reports) {
    return <Navigate to={getRoleDashboardPath(userRole)} replace />;
  }

  return <Reports />;
}

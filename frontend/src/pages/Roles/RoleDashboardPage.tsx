import { Navigate, useParams } from "react-router";
import type { ComponentType } from "react";
import { useAuth } from "../../context/AuthContext";
import AdminDashboard from "./Admin/Dashboard";
import DoctorDashboard from "./Doctor/Dashboard";
import PatientDashboard from "./Patient/Dashboard";
import ReceptionistDashboard from "./Receptionist/Dashboard";
import SuperAdminDashboard from "./SuperAdmin/Dashboard";
import { getRoleDashboardPath, normalizeRole, type RoleSlug } from "./shared/role-routing";

const dashboards: Record<RoleSlug, ComponentType> = {
  admin: AdminDashboard,
  doctor: DoctorDashboard,
  receptionist: ReceptionistDashboard,
  patient: PatientDashboard,
  "super-admin": SuperAdminDashboard,
};

export default function RoleDashboardPage() {
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
    return <Navigate to={getRoleDashboardPath(userRole)} replace />;
  }

  const Dashboard = dashboards[routeRole];

  return <Dashboard />;
}

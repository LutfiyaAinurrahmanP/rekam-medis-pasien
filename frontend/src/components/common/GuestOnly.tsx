import { Navigate, Outlet } from "react-router";
import authService from "../../services/auth";
import { useAuth } from "../../context/AuthContext";
import { getRoleDashboardPath } from "../../pages/Roles/shared/role-routing";

export default function GuestOnly() {
  const { user, loading } = useAuth();

  if (loading) {
    return null;
  }

  if (authService.isAuthenticated()) {
    return <Navigate to={getRoleDashboardPath(user?.role)} replace />;
  }

  return <Outlet />;
}

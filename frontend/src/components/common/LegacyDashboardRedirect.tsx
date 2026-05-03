import { Navigate } from "react-router";
import { useAuth } from "../../context/AuthContext";
import { getRoleDashboardPath } from "../../pages/Roles/shared/role-routing";

export default function LegacyDashboardRedirect() {
  const { user, loading } = useAuth();

  if (loading) {
    return null;
  }

  return <Navigate to={getRoleDashboardPath(user?.role)} replace />;
}

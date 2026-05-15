import { Navigate, Outlet, useLocation } from "react-router";
import authService from "../../services/auth";

const ROLE_PATH_SEGMENTS = new Set([
  "admin",
  "doctor",
  "patient",
  "receptionist",
  "super-admin",
]);

const normalizeRoleToPathSegment = (role: string) =>
  role.trim().toLowerCase().replace(/_/g, "-");

const getCurrentUserRole = (): string | undefined => {
  try {
    const storedUser = localStorage.getItem("currentUser");
    if (storedUser) {
      const parsed = JSON.parse(storedUser) as { role?: unknown };
      if (typeof parsed.role === "string" && parsed.role.trim() !== "") {
        return parsed.role;
      }
    }
  } catch (err) {
    console.warn("Failed to parse currentUser from localStorage:", err);
  }

  const payload = authService.getTokenPayload();
  const tokenRole = payload?.role;
  if (typeof tokenRole === "string" && tokenRole.trim() !== "") {
    return tokenRole;
  }

  return undefined;
};

export default function RequireAuth() {
  const location = useLocation();

  if (!authService.isAuthenticated()) {
    return <Navigate to="/auth/login" replace state={{ from: location }} />;
  }

  const firstPathSegment = location.pathname.split("/")[1]?.toLowerCase() ?? "";

  if (ROLE_PATH_SEGMENTS.has(firstPathSegment)) {
    const currentUserRole = getCurrentUserRole();
    const normalizedCurrentUserRole = currentUserRole
      ? normalizeRoleToPathSegment(currentUserRole)
      : undefined;

    if (
      !normalizedCurrentUserRole ||
      normalizedCurrentUserRole !== firstPathSegment
    ) {
      return <Navigate to="/unauthorized" replace />;
    }
  }

  return <Outlet />;
}

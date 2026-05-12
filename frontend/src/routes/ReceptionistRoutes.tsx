import { Fragment } from "react";
import { Route } from "react-router";
import RoleDashboardPage from "../pages/Roles/RoleDashboardPage";
import RoleReportsPage from "../pages/Roles/RoleReportsPage";

export default function ReceptionistRoutes() {
  return (
    <Fragment>
      <Route path="/receptionist/dashboard" element={<RoleDashboardPage />} />
      <Route path="/receptionist/reports" element={<RoleReportsPage />} />
    </Fragment>
  );
}

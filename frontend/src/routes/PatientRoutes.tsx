import { Fragment } from "react";
import { Route } from "react-router";
import RoleDashboardPage from "../pages/Roles/RoleDashboardPage";
import RoleReportsPage from "../pages/Roles/RoleReportsPage";

export default function PatientRoutes() {
  return (
    <Fragment>
      <Route path="/patient/dashboard" element={<RoleDashboardPage />} />
      <Route path="/patient/reports" element={<RoleReportsPage />} />
    </Fragment>
  );
}

import { Fragment } from "react";
import { Route } from "react-router";
import RoleDashboardPage from "../pages/Roles/RoleDashboardPage";
import RoleReportsPage from "../pages/Roles/RoleReportsPage";
import DepartmentsIndex from "../pages/Admin/Departments/Index";
// import DepartmentsShow from "../pages/Admin/Departments/Show";

export default function ReceptionistRoutes() {
  return (
    <Fragment>
      <Route path="/receptionist/dashboard" element={<RoleDashboardPage />} />
      <Route path="/receptionist/reports" element={<RoleReportsPage />} />

      <Route path="/receptionist/departments" element={<DepartmentsIndex />} />
    </Fragment>
  );
}

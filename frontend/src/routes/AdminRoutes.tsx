import { Fragment } from "react";
import { Route } from "react-router";
import RoleDashboardPage from "../pages/Roles/RoleDashboardPage";
import RoleReportsPage from "../pages/Roles/RoleReportsPage";
import UsersIndex from "../pages/Admin/Users/Index";
import DepartmentsIndex from "../pages/Admin/Departments/Index";
// import DepartmentsCreate from "../pages/Admin/Departments/Create";
// import DepartmentsEdit from "../pages/Admin/Departments/Edit";
// import DepartmentsShow from "../pages/Admin/Departments/Show";

export default function AdminRoutes() {
  return (
    <Fragment>
      <Route path="/admin/dashboard" element={<RoleDashboardPage />} />
      <Route path="/admin/reports" element={<RoleReportsPage />} />

      <Route path="/admin/users" element={<UsersIndex />} />
      <Route path="/admin/departments" element={<DepartmentsIndex />} />
    </Fragment>
  );
}

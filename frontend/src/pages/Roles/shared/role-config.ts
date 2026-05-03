export type RoleKey = "admin" | "doctor" | "receptionist" | "patient" | "superAdmin";

export type RoleDashboardConfig = {
  label: string;
  dashboardRoute: string;
  reportsRoute?: string;
  dashboardEndpoint: string;
  reportEndpoints: string[];
};

export const roleConfig: Record<RoleKey, RoleDashboardConfig> = {
  admin: {
    label: "Admin",
    dashboardRoute: "/dashboard/admin",
    reportsRoute: "/dashboard/admin/reports",
    dashboardEndpoint: "/api/v1/dashboard/admin",
    reportEndpoints: [
      "/api/v1/dashboard/overview",
      "/api/v1/dashboard/reports/appointments",
      "/api/v1/dashboard/reports/revenue",
      "/api/v1/dashboard/reports/patients",
    ],
  },
  doctor: {
    label: "Doctor",
    dashboardRoute: "/dashboard/doctor",
    reportsRoute: "/dashboard/doctor/reports",
    dashboardEndpoint: "/api/v1/dashboard/doctor",
    reportEndpoints: ["/api/v1/dashboard/reports/appointments"],
  },
  receptionist: {
    label: "Receptionist",
    dashboardRoute: "/dashboard/receptionist",
    reportsRoute: "/dashboard/receptionist/reports",
    dashboardEndpoint: "/api/v1/dashboard/receptionist",
    reportEndpoints: [
      "/api/v1/dashboard/reports/appointments",
      "/api/v1/dashboard/reports/patients",
    ],
  },
  patient: {
    label: "Patient",
    dashboardRoute: "/dashboard/patient",
    dashboardEndpoint: "/api/v1/dashboard/patient",
    reportEndpoints: [],
  },
  superAdmin: {
    label: "Super Admin",
    dashboardRoute: "/dashboard/super-admin",
    reportsRoute: "/dashboard/super-admin/reports",
    dashboardEndpoint: "/api/v1/dashboard/overview",
    reportEndpoints: [
      "/api/v1/dashboard/reports/appointments",
      "/api/v1/dashboard/reports/revenue",
      "/api/v1/dashboard/reports/patients",
    ],
  },
};
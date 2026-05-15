import { useEffect } from "react";
import PageBreadcrumb from "../components/common/PageBreadCrumb";
import UserInfoCard from "../components/UserProfile/UserInfoCard";
import PageMeta from "../components/common/PageMeta";
import { useAuth } from "../context/AuthContext";

export default function UserProfiles() {
  const { loadUserProfile } = useAuth();

  useEffect(() => {
    void loadUserProfile();
  }, [loadUserProfile]);

  return (
    <>
      <PageMeta
        title="Account Settings | Rekam Medis Pasien"
        description="Manage your account profile, password, and account actions."
      />
      <PageBreadcrumb pageTitle="Account Settings" />
      <div className="rounded-2xl border border-gray-200 bg-white p-5 dark:border-gray-800 dark:bg-white/[0.03] lg:p-6">
        <div className="space-y-6">
          <UserInfoCard />
        </div>
      </div>
    </>
  );
}

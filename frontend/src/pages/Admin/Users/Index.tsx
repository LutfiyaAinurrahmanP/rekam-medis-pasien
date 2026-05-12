import PageMeta from "../../../components/common/PageMeta";
import UsersIndexLayout from "../../../layout/Admin/Users/UsersIndexLayout";

export default function UsersIndex() {
  return (
    <>
      <PageMeta
        title="Users | Medical Records System - Admin Dashboard"
        description="Manage users in the medical records system"
      />
      <UsersIndexLayout />
    </>
  );
}

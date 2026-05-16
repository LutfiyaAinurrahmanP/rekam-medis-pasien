import PageMeta from "../../../components/common/PageMeta";
import DepartmentsIndexLayout from "../../../layout/Admin/Departments/DepartmentsIndexLayout";

export default function DepartmentsIndex() {
  return (
    <>
      <PageMeta title="Departments | Admin" description="Manage departments" />
      <DepartmentsIndexLayout />
    </>
  );
}

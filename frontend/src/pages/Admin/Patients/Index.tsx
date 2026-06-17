import PageMeta from "../../../components/common/PageMeta";
import PatientsIndexLayout from "../../../layout/Admin/Patients/PatientsIndexLayout";

export default function PatientsIndex() {
  return (
    <>
      <PageMeta title="Patients | Admin" description="Manage patients" />
      <PatientsIndexLayout />
    </>
  );
}

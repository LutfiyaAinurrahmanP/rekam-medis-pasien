import CreatePatientModal from "../../../layout/Admin/Patients/CreatePatientModal";
import { useNavigate } from "react-router";
import { useAuth } from "../../../context/AuthContext";
import { getRolePatientsPath } from "../../Roles/shared/role-routing";

export default function PatientsCreate() {
  const navigate = useNavigate();
  const { user } = useAuth();

  return (
    <CreatePatientModal
      isOpen={true}
      onClose={() => navigate(getRolePatientsPath(user?.role))}
    />
  );
}

import { useParams } from "react-router";
import { useNavigate } from "react-router";
import EditPatientModal from "../../../layout/Admin/Patients/EditPatientModal";
import { useAuth } from "../../../context/AuthContext";
import { getRolePatientsPath } from "../../Roles/shared/role-routing";

export default function PatientsEdit() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();

  return (
    <EditPatientModal
      isOpen={true}
      id={id}
      onClose={() => navigate(getRolePatientsPath(user?.role))}
    />
  );
}

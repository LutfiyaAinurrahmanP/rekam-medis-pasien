import { useParams } from "react-router";
import CreateUserModal from "../../../layout/Admin/Users/CreateUserModal";

export default function UsersCreate() {
  const { role } = useParams();

  return <CreateUserModal isOpen={true} role={role} />;
}

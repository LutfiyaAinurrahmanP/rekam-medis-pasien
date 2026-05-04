import PageMeta from "../../components/common/PageMeta";
import AuthLayout from "./AuthPageLayout";
import RegisterForm from "../../components/auth/RegisterForm";

export default function Register() {
  return (
    <>
      <PageMeta
        title="Register | Rekam Medis Pasien"
        description="Registration page to create a new patient account."
      />
      <AuthLayout>
        <RegisterForm />
      </AuthLayout>
    </>
  );
}

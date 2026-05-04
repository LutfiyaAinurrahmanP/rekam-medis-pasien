import PageMeta from "../../components/common/PageMeta";
import AuthLayout from "./AuthPageLayout";
import RegisterForm from "../../components/auth/RegisterForm";

export default function Register() {
  return (
    <>
      <PageMeta
        title="Daftar | Rekam Medis Pasien"
        description="Halaman pendaftaran untuk membuat akun pasien baru."
      />
      <AuthLayout>
        <RegisterForm />
      </AuthLayout>
    </>
  );
}

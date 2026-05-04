import PageMeta from "../../components/common/PageMeta";
import AuthLayout from "./AuthPageLayout";
import LoginForm from "../../components/auth/LoginForm";

export default function Login() {
  return (
    <>
      <PageMeta
        title="Masuk | Rekam Medis Pasien"
        description="Halaman masuk untuk mengakses dashboard sesuai peran pengguna."
      />
      <AuthLayout>
        <LoginForm />
      </AuthLayout>
    </>
  );
}

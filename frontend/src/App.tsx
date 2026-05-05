import { BrowserRouter as Router, Routes, Route } from "react-router";
import Ecommerce from "./pages/Dashboard/Ecommerce";
import Stocks from "./pages/Dashboard/Stocks";
import Crm from "./pages/Dashboard/Crm";
import Marketing from "./pages/Dashboard/Marketing";
import Analytics from "./pages/Dashboard/Analytics";
import Login from "./pages/AuthPages/Login";
import Register from "./pages/AuthPages/Register";
import NotFound from "./pages/OtherPage/NotFound";
import UserProfiles from "./pages/UserProfiles";
import Carousel from "./pages/UiElements/Carousel";
import Maintenance from "./pages/OtherPage/Maintenance";
import FiveZeroZero from "./pages/OtherPage/FiveZeroZero";
import FiveZeroThree from "./pages/OtherPage/FiveZeroThree";
import Videos from "./pages/UiElements/Videos";
import Images from "./pages/UiElements/Images";
import Alerts from "./pages/UiElements/Alerts";
import Badges from "./pages/UiElements/Badges";
import Pagination from "./pages/UiElements/Pagination";
import Avatars from "./pages/UiElements/Avatars";
import Buttons from "./pages/UiElements/Buttons";
import ButtonsGroup from "./pages/UiElements/ButtonsGroup";
import Notifications from "./pages/UiElements/Notifications";
import LineChart from "./pages/Charts/LineChart";
import BarChart from "./pages/Charts/BarChart";
import PieChart from "./pages/Charts/PieChart";
import Invoices from "./pages/Invoices";
import ComingSoon from "./pages/OtherPage/ComingSoon";
import FileManager from "./pages/FileManager";
import Calendar from "./pages/Calendar";
import BasicTables from "./pages/Tables/BasicTables";
import DataTables from "./pages/Tables/DataTables";
import PricingTables from "./pages/PricingTables";
import Faqs from "./pages/Faqs";
import Chats from "./pages/Chat/Chats";
import FormElements from "./pages/Forms/FormElements";
import FormLayout from "./pages/Forms/FormLayout";
import Blank from "./pages/Blank";
import EmailInbox from "./pages/Email/EmailInbox";
import EmailDetails from "./pages/Email/EmailDetails";

import TaskKanban from "./pages/Task/TaskKanban";
import BreadCrumb from "./pages/UiElements/BreadCrumb";
import Cards from "./pages/UiElements/Cards";
import Dropdowns from "./pages/UiElements/Dropdowns";
import Links from "./pages/UiElements/Links";
import Lists from "./pages/UiElements/Lists";
import Popovers from "./pages/UiElements/Popovers";
import Progressbar from "./pages/UiElements/Progressbar";
import Ribbons from "./pages/UiElements/Ribbons";
import Spinners from "./pages/UiElements/Spinners";
import Tabs from "./pages/UiElements/Tabs";
import Tooltips from "./pages/UiElements/Tooltips";
import Modals from "./pages/UiElements/Modals";
import ResetPassword from "./pages/AuthPages/ResetPassword";
import TwoStepVerification from "./pages/AuthPages/TwoStepVerification";
import Success from "./pages/OtherPage/Success";
import AppLayout from "./layout/AppLayout";
import AlternativeLayout from "./layout/AlternativeLayout";
import { ScrollToTop } from "./components/common/ScrollToTop";
import TaskList from "./pages/Task/TaskList";
import Saas from "./pages/Dashboard/Saas";
import Logistics from "./pages/Dashboard/Logistics";
import TextGeneratorPage from "./pages/Ai/TextGenerator";
import ImageGeneratorPage from "./pages/Ai/ImageGenerator";
import CodeGeneratorPage from "./pages/Ai/CodeGenerator";
import VideoGeneratorPage from "./pages/Ai/VideoGenerator";
import ProductList from "./pages/Ecommerce/ProductList";
import AddProduct from "./pages/Ecommerce/AddProduct";
import Billing from "./pages/Ecommerce/Billing";
import SingleInvoice from "./pages/Ecommerce/SingleInvoice";
import CreateInvoice from "./pages/Ecommerce/CreateInvoice";
import Transactions from "./pages/Ecommerce/Transactions";
import SingleTransaction from "./pages/Ecommerce/SingleTransaction";
import TicketList from "./pages/Support/TicketList";
import TicketReply from "./pages/Support/TicketReply";
import Integrations from "./pages/OtherPage/Integrations";
import ApiKeys from "./pages/OtherPage/ApiKeys";
import RequireAuth from "./components/common/RequireAuth";
import GuestOnly from "./components/common/GuestOnly";
import LegacyDashboardRedirect from "./components/common/LegacyDashboardRedirect";
import RoleDashboardPage from "./pages/Roles/RoleDashboardPage";
import UsersIndex from "./pages/Users/Index";
import UsersCreate from "./pages/Users/Create";
import UsersEdit from "./pages/Users/Edit";
import UsersShow from "./pages/Users/Show";
import RoleReportsPage from "./pages/Roles/RoleReportsPage";

export default function App() {
  return (
    <>
      <Router>
        <ScrollToTop />
        <Routes>
          {/* Auth */}
          <Route element={<GuestOnly />}>
            <Route path="/auth/login" element={<Login />} />
            <Route path="/auth/register" element={<Register />} />
          </Route>
          <Route path="/auth/reset-password" element={<ResetPassword />} />

          {/* Authorization Required */}
          <Route element={<RequireAuth />}>
            <Route element={<AppLayout />}>
              <Route path="/dashboard" element={<LegacyDashboardRedirect />} />
              <Route path="/:role/dashboard" element={<RoleDashboardPage />} />
              <Route path="/:role/reports" element={<RoleReportsPage />} />

              {/* Users Management */}
              <Route path="/:role/users" element={<UsersIndex />} />
              <Route path="/:role/users/create" element={<UsersCreate />} />
              <Route path="/:role/users/:id/edit" element={<UsersEdit />} />
              <Route path="/:role/users/:id" element={<UsersShow />} />

              {/* Template Sidebar */}
              {/* Dashboard Layout */}
              <Route index path="/template" element={<Ecommerce />} />
              <Route path="/template/analytics" element={<Analytics />} />
              <Route path="/template/marketing" element={<Marketing />} />
              <Route path="/template/crm" element={<Crm />} />
              <Route path="/template/stocks" element={<Stocks />} />
              <Route path="/template/saas" element={<Saas />} />
              <Route path="/template/logistics" element={<Logistics />} />

              <Route path="/template/calendar" element={<Calendar />} />
              <Route path="/template/invoice" element={<Invoices />} />
              <Route path="/template/invoices" element={<Invoices />} />
              <Route path="/template/chat" element={<Chats />} />
              <Route path="/template/file-manager" element={<FileManager />} />

              {/* E-commerce */}
              <Route path="/template/product-list" element={<ProductList />} />
              <Route path="/template/add-product" element={<AddProduct />} />
              <Route path="/template/billing" element={<Billing />} />
              <Route
                path="/template/single-invoice"
                element={<SingleInvoice />}
              />
              <Route
                path="/template/create-invoice"
                element={<CreateInvoice />}
              />
              <Route path="/template/transactions" element={<Transactions />} />
              <Route
                path="/template/single-transaction"
                element={<SingleTransaction />}
              />

              {/* Support */}
              <Route path="/template/ticket-list" element={<TicketList />} />
              <Route path="/template/ticket-reply" element={<TicketReply />} />

              {/* Others Page */}
              <Route path="/template/profile" element={<UserProfiles />} />
              <Route path="/template/faq" element={<Faqs />} />
              <Route
                path="/template/pricing-tables"
                element={<PricingTables />}
              />
              <Route path="/template/integrations" element={<Integrations />} />
              <Route path="/template/api-keys" element={<ApiKeys />} />
              <Route path="/template/blank" element={<Blank />} />

              {/* Forms */}
              <Route
                path="/template/form-elements"
                element={<FormElements />}
              />
              <Route path="/template/form-layout" element={<FormLayout />} />

              {/* Applications */}
              <Route path="/template/task-list" element={<TaskList />} />
              <Route path="/template/task-kanban" element={<TaskKanban />} />

              {/* Email */}
              <Route path="/template/inbox" element={<EmailInbox />} />
              <Route
                path="/template/inbox-details"
                element={<EmailDetails />}
              />

              {/* Tables */}
              <Route path="/template/basic-tables" element={<BasicTables />} />
              <Route path="/template/data-tables" element={<DataTables />} />

              {/* Ui Elements */}
              <Route path="/template/alerts" element={<Alerts />} />
              <Route path="/template/avatars" element={<Avatars />} />
              <Route path="/template/badge" element={<Badges />} />
              <Route path="/template/breadcrumb" element={<BreadCrumb />} />
              <Route path="/template/buttons" element={<Buttons />} />
              <Route
                path="/template/buttons-group"
                element={<ButtonsGroup />}
              />
              <Route path="/template/cards" element={<Cards />} />
              <Route path="/template/carousel" element={<Carousel />} />
              <Route path="/template/dropdowns" element={<Dropdowns />} />
              <Route path="/template/images" element={<Images />} />
              <Route path="/template/links" element={<Links />} />
              <Route path="/template/list" element={<Lists />} />
              <Route path="/template/modals" element={<Modals />} />
              <Route
                path="/template/notifications"
                element={<Notifications />}
              />
              <Route path="/template/pagination" element={<Pagination />} />
              <Route path="/template/popovers" element={<Popovers />} />
              <Route path="/template/progress-bar" element={<Progressbar />} />
              <Route path="/template/ribbons" element={<Ribbons />} />
              <Route path="/template/spinners" element={<Spinners />} />
              <Route path="/template/tabs" element={<Tabs />} />
              <Route path="/template/tooltips" element={<Tooltips />} />
              <Route path="/template/videos" element={<Videos />} />

              {/* Charts */}
              <Route path="/template/line-chart" element={<LineChart />} />
              <Route path="/template/bar-chart" element={<BarChart />} />
              <Route path="/template/pie-chart" element={<PieChart />} />
            </Route>

            {/* Alternative Layout - for special pages */}
            <Route element={<AlternativeLayout />}>
              {/* AI Generator */}
              <Route
                path="/template/text-generator"
                element={<TextGeneratorPage />}
              />
              <Route
                path="/template/image-generator"
                element={<ImageGeneratorPage />}
              />
              <Route
                path="/template/code-generator"
                element={<CodeGeneratorPage />}
              />
              <Route
                path="/template/video-generator"
                element={<VideoGeneratorPage />}
              />
            </Route>
          </Route>

          {/* Auth Layout */}
          <Route path="/signin" element={<Login />} />
          <Route path="/signup" element={<Register />} />
          <Route path="/reset-password" element={<ResetPassword />} />
          <Route
            path="/two-step-verification"
            element={<TwoStepVerification />}
          />

          {/* Fallback Route */}
          <Route path="*" element={<NotFound />} />
          <Route path="/maintenance" element={<Maintenance />} />
          <Route path="/success" element={<Success />} />
          <Route path="/five-zero-zero" element={<FiveZeroZero />} />
          <Route path="/five-zero-three" element={<FiveZeroThree />} />
          <Route path="/coming-soon" element={<ComingSoon />} />
        </Routes>
      </Router>
    </>
  );
}

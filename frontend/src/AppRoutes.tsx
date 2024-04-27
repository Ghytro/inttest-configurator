import React from "react";
import { Route, Routes } from "react-router-dom";
import AuthPage from "./pages/auth/Auth";
import ProjectsListPage from "./pages/projects/ProjectsList";
import { routesEnum } from "./routesEnum";
import UserEditPage from "./pages/users/UserEditor";
import ConfiguratorHeader from "./components/ConfiguratorHeader";

const AppRoutes = () => {
  return (
    <>
      <ConfiguratorHeader />
      <Routes>
        <Route path={routesEnum.auth} element={<AuthPage />} />
        <Route path={routesEnum.projects} element={<ProjectsListPage />} />
        <Route path={routesEnum.users} element={<UserEditPage />} />
        <Route path={routesEnum.projectPageMask element={<ProjectPage />}}
      </Routes>
    </>
  );
};

export default AppRoutes;

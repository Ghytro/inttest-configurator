import React from "react";
import { Route, Routes } from "react-router-dom";
import AuthPage from "./pages/auth/Auth";
import ProjectsPage from "./pages/projects/Projects";
import { routesEnum } from "./routesEnum";
import UserEditPage from "./pages/users/UserEditor";

const AppRoutes = () => {
  return (
    <Routes>
      <Route path={routesEnum.auth} element={<AuthPage />} />
      <Route path={routesEnum.projects} element={<ProjectsPage />} />
      <Route path={routesEnum.users} element={<UserEditPage />} />
    </Routes>
  );
};

export default AppRoutes;

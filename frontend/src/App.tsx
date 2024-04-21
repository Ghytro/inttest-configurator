import React from "react";

import { BrowserRouter } from "react-router-dom";

import "./App.css";
import AppRoutes from "./AppRoutes";
import { ConfigProvider } from "antd";
import { theme } from "./config/style";

function App() {
  return (
    <div className="App">
      <ConfigProvider theme={theme}>
        <BrowserRouter>
          <AppRoutes />
        </BrowserRouter>
      </ConfigProvider>
    </div>
  );
}

export default App;

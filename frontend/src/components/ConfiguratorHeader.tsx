import { Breadcrumb, Menu } from "antd";
import { Header } from "antd/es/layout/layout";
import Title from "antd/es/typography/Title";
import React, { Component } from "react";
import { routesEnum } from "../routesEnum";

class ConfiguratorHeader extends Component<any, any> {
  constructor(props: any) {
    super(props);
  }

  render(): React.ReactNode {
    return (
      <Header
        style={{
          display: "flex",
          flexDirection: "row",
          justifyContent: "center",
          marginBottom: "5vh",
        }}
      >
        <div style={{ alignSelf: "flex-start" }}>Configurator</div>
        <Menu
          mode="horizontal"
          defaultSelectedKeys={["2"]}
          inlineCollapsed={false}
          items={[
            {
              key: 0,
              label: "Список проектов",
              onClick: (info) => {
                window.location.href = routesEnum.projects;
              },
            },
            {
              key: 1,
              label: "Редактирование пользователей",
              onClick: (info) => {
                window.location.href = routesEnum.users;
              },
            },
          ]}
          style={{
            display: "flex",
            flex: 1,
            flexDirection: "row",
            justifyContent: "space-around",
            minWidth: 0,
          }}
        />
      </Header>
    );
  }
}

export default ConfiguratorHeader;

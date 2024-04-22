import React, { Component } from "react";

import "./UserEditor.css";
import { Button, Modal, Table } from "antd";
import { EditOutlined } from "@ant-design/icons";
import { InternalApiAuthListPermResponseItem, InternalApiAuthListRoleResponseItem, UsersApi } from "../../api/api";
import { $notify, ENotifyKind } from "../../notifier";

class UserEditPage extends Component<any, UserListModel> {
  private userApi: UsersApi;
  private roleById: Map<Number, InternalApiAuthListRoleResponseItem>;
  private permById: Map<Number, InternalApiAuthListPermResponseItem>;

  constructor(props) {
    super(props);
    this.userApi = new UsersApi();
    this.roleById = new Map<Number, InternalApiAuthListRoleResponseItem>,
    this.permById = new Map<Number, InternalApiAuthListPermResponseItem>,

    this.state = {
      loading: true,
      
      userListTableModel: undefined,
      editedUserId: undefined,
      editedUserName: undefined,
    };

    this.fetchUserListTableData();
  }

  fetchUserListTableData() {
    this.setState({loading: true});
    if (!this.permById.size) {
      this.fetchPermById(false);
    }
    if (!this.roleById.size) {
      this.fetchRoleById(false);
    }
    this.userApi.listUsers()
      .then(({ data }) => {
        this.setState({userListTableModel: data.map((item) => ({
          id: item.id,
          username: item.username,
          roleNames: item.role_ids
            ?.map((roleId) => this.roleById.get(roleId)?.name)
            .join(", "),
        }))});
      })
      .catch((e) => $notify(ENotifyKind.ERROR, e))
      .finally(() => {this.setState({loading: false})});
  }

  fetchRoleById(withLoading: boolean) {
    if (withLoading) {
      this.setState({loading: true});
    }
    this.userApi.listRoles()
      .then(({ data }) => {
        data.forEach((item) => this.roleById.set(item.id!, item))
      })
      .catch((e) => $notify(ENotifyKind.ERROR, e))
      .finally(() => { if (withLoading) {this.setState({loading: false}) }});
  }

  fetchPermById(withLoading: boolean) {
    if (withLoading) {
      this.setState({loading: true});
    }
    this.userApi.listPerms()
      .then(({ data }) => {
        data.forEach((item) => this.permById.set(item.id!, item))
      })
      .catch((e) => $notify(ENotifyKind.ERROR, e))
      .finally(() => { if (withLoading) {this.setState({loading: false}) }});
  }

  render() {
    console.log(this.state);
    const userListTableCols = [
      {
        title: "Логин",
        dataIndex: "username",
      },
      {
        title: "Роли",
        dataIndex: "roleNames",
      },
      {
        title: "Редактирование",
        dataIndex: "editedUser",
        render: (_, user) => (
          <Button
            shape="circle"
            icon={<EditOutlined />}
            onClick={() => {this.setState({editedUserId: user.id, editedUserName: user.username})}}
            size="large"
          />
        ),
      },
    ];

    return (
      <>
        <Table
          columns={userListTableCols}
          dataSource={this.state.userListTableModel}
          loading={this.state.loading}
        ></Table>
        <Modal
          title={"Редактирование пользователя "+this.state.editedUserName}
          open={this.state.editedUserId != undefined}
          onOk={() => console.log("OK pressed")}
          onCancel={() => this.setState({editedUserId: undefined})}
          >
          i'm here...
        </Modal>
      </>
    );
  }
}

interface UserListModel {
  loading: boolean;
  // roleById: Map<Number, InternalApiAuthListRoleResponseItem>;
  // permById: Map<Number, InternalApiAuthListPermResponseItem>;
  userListTableModel: any[] | undefined;
  editedUserId: number | undefined;
  editedUserName: string | undefined;
}

//   constructor(userApi: UsersApi) {
//     this.userApi = userApi;
//     this.loading = false;
//     this.roleById = new Map<Number, InternalApiAuthListRoleResponseItem>();
//     this.permById = new Map<Number, InternalApiAuthListPermResponseItem>();
//     this.editedUserId = undefined;
//   }


// }

export default UserEditPage;

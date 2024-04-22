import React, { Component } from "react";

import "./UserEditor.css";
import { Button, Checkbox, FloatButton, Form, Input, Modal, Table } from "antd";
import { DeleteOutlined, EditOutlined, UserAddOutlined, UsergroupAddOutlined } from "@ant-design/icons";
import { InternalApiAuthListPermResponseItem, InternalApiAuthListRoleResponseItem, UsersApi } from "../../api/api";
import { $notify, ENotifyKind } from "../../notifier";
import { routesEnum } from "../../routesEnum";

class UserEditPage extends Component<any, UserListModel> {
  private userApi: UsersApi;
  private roleById: Map<Number, InternalApiAuthListRoleResponseItem>;
  private permById: Map<Number, InternalApiAuthListPermResponseItem>;

  constructor(props) {
    super(props);
    this.userApi = new UsersApi();
    this.roleById = new Map<Number, InternalApiAuthListRoleResponseItem>;
    this.permById = new Map<Number, InternalApiAuthListPermResponseItem>;

    this.state = {
      loading: true,
      
      userListTableModel: undefined,
      editedUserId: undefined,
      editedUserName: undefined,

      roleListTableModel: undefined,

      editedRoleId: undefined,
      editedRoleName: undefined,
      editedRoleNameInputValue:  undefined,
      editedRoleDescInputValue:  undefined,
      editedRolePermsCheckboxValues: undefined,
      editedRolePermsCheckboxAllValues: undefined,

      permissionListTableModel: undefined,

      usernameEditInputValue: undefined,
      userRolesEditCheckboxValues: undefined,
      userRolesEditCheckboxAllValues: undefined,

      userCreateDialogOpened: false,
      userCreateLoginInputValue: undefined,
      userCreatePasswordInputValue: undefined,
      userCreateRolesCheckboxValues:  undefined,
      userCreateRolesCheckboxAllValues: undefined,

      roleCreateDialogOpened: false,
      roleCreateNameInputValue: undefined,
      roleCreateDescInputValue: undefined,
      roleCreatePermsCheckboxValues: undefined,
      roleCreatePermsCheckboxAllValues: undefined,
    };

    this.fetchUserListTableData();
  }

  fetchUserListTableData() {
    this.setState({loading: true});
    this.userApi.listRoles()
      .then((roleResp: any) => {
        this.userApi.listPerms()
          .then((permResp: any) => {
            permResp.data.forEach((item) => {
              this.permById.set(item.id!, item);
            });
            let permListTableModel: any[] = [];
            this.permById.forEach((v, k) => {
              permListTableModel.push({
                id: k,
                name: v.name,
                desc: v.desc
              });
            })
            this.setState({permissionListTableModel: permListTableModel});

            roleResp.data.forEach((item) => {
              this.roleById.set(item.id!, item);
            });
            let roleListTableModel: any[] = [];
            this.roleById.forEach((v, k) => {
              roleListTableModel.push({
                id: k,
                name: v.name,
                desc: v.desc,
                permissions: v.perm_ids?.map((permId) => this.permById.get(permId)?.name),
              });
            });
            this.setState({roleListTableModel: roleListTableModel});

            this.userApi.listUsers()
              .then(({ data }) => {
                this.setState({userListTableModel: data.map((item) => ({
                  id: item.id,
                  username: item.username,
                  roleNames: item.role_ids
                    ?.map((roleId) => this.roleById.get(roleId)?.name),
                }))});
              })
              .catch((e) => $notify(ENotifyKind.ERROR, e))
              .finally(() => { {this.setState({loading: false})} });
          })
          .catch((e) => $notify(ENotifyKind.ERROR, e));
      })
      .catch((e) => $notify(ENotifyKind.ERROR, e))
  }

  deleteRole(roleId: number) {
    this.userApi.deleteRole(roleId)
      .then(( {data} ) => {
        this.fetchUserListTableData();
      })
  }

  deleteUser(userId: number) {
    this.userApi.deleteUser(userId)
      .then(( {data} ) => {
        this.fetchUserListTableData();
      })
      .catch((e) => $notify(ENotifyKind.ERROR, e))
  }

  updateUser(id: number, newUsername: string, newRoleIds: number[]) {
    this.userApi.updateUser(id, {
      username: newUsername,
      role_ids: newRoleIds,
    })
    .then(( {data} ) => {
      this.fetchUserListTableData();
      this.setState({editedUserId: undefined});
    })
    .catch((e) => $notify(ENotifyKind.ERROR, e));
  }

  updateRole(id: number, newName: string, newDesc: string | undefined, newPermIds: number[]) {
    this.userApi.updateRole(id, {
      name: newName,
      desc: newDesc,
      perm_ids: newPermIds,
    })
    .then(( {data} ) => {
      this.setState({
        editedRoleId: undefined,
        editedRoleName: undefined,
        editedRoleNameInputValue: undefined,
        editedRoleDescInputValue: undefined,
        editedRolePermsCheckboxValues: undefined
      })
      this.fetchUserListTableData();
    })
    .catch((e) => $notify(ENotifyKind.ERROR, e));
  }

  createUser() {
    const roles = Array.from(this.roleById.values());
    this.userApi.createUser({
      username: this.state.userCreateLoginInputValue,
      password: this.state.userCreatePasswordInputValue,
      role_ids: this.state.userCreateRolesCheckboxValues?.map((roleName) => roles.find((role) => role.name == roleName)?.id!),
    })
    .then(({data}) => {
      this.setState({
        userCreateLoginInputValue: undefined,
        userCreatePasswordInputValue: undefined,
        userCreateRolesCheckboxValues: undefined,
        userCreateDialogOpened: false,
      });
      this.fetchUserListTableData();
    })
    .catch((e) => $notify(ENotifyKind.ERROR, e))
  }

  createRole() {
    const perms = Array.from(this.permById.values());
    this.userApi.createRole({
      name: this.state.roleCreateNameInputValue,
      desc: this.state.roleCreateDescInputValue,
      perm_ids: this.state.roleCreatePermsCheckboxValues?.map((permName) => perms.find((x) => x.name == permName)?.id!),
    })
    .then(({data}) => {
      this.setState({
        roleCreateNameInputValue: undefined,
        roleCreateDescInputValue: undefined,
        roleCreatePermsCheckboxValues: undefined,
        roleCreateDialogOpened: false,
      });
      this.fetchUserListTableData();
    })
    .catch((e) => $notify(ENotifyKind.ERROR, e))
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
        render: (_, user) => (
          user.roleNames.map((x) => <><span>{x}</span><br /></>)
        )
      },
      {
        title: "Редактировать",
        render: (_, user) => (
          <Button
            shape="circle"
            icon={<EditOutlined />}
            onClick={() => {
              const allCheckboxValues = this.state.roleListTableModel?.map((x) => x.name);
              this.setState({
                editedUserId: user.id,
                editedUserName: user.username,
                userRolesEditCheckboxValues: user.roleNames,
                userRolesEditCheckboxAllValues: allCheckboxValues,
                usernameEditInputValue: user.username,
              });
            }}
            size="large"
          />
        ),
      },
      {
        title: "Удалить",
        render: (_, user) => (
          <Button
            shape="circle"
            icon={<DeleteOutlined />}
            onClick={() => {this.deleteUser(user.id)}}
            size="large"
          />
        )
      }
    ];
    const roleListTableCols = [
      {
        title: "Наименование",
        dataIndex: "name",
      },
      {
        title: "Описание",
        dataIndex: "desc",
      },
      {
        title: "Разрешения",
        dataIndex: "permissions",
        render: (_, record) => (
          record.permissions.map((permName) => (<><span>{permName}</span><br /></>))
        )
      },
      {
        title: "Редактировать",
        render: (_, role) => (
          <Button
            shape="circle"
            icon={<EditOutlined />}
            onClick={() => {
              const allCheckboxValues = this.state.permissionListTableModel?.map((x) => x.name);
              this.setState({
                editedRoleId: role.id,
                editedRoleName: role.name,
                editedRoleNameInputValue: role.name,
                editedRoleDescInputValue: role.desc,
                editedRolePermsCheckboxValues: role.permissions,
                editedRolePermsCheckboxAllValues: allCheckboxValues,
              });
              console.log(this.state);
          }}
            size="large"
          />
        ),
      },
      {
        title: "Удалить",
        render: (_, role) => (
          <Button
            shape="circle"
            icon={<DeleteOutlined />}
            onClick={() => {this.deleteRole(role.id)}}
            size="large"
          />
        )
      }
    ];
    const permListTableCols = [
      {
        title: "Наименование",
        dataIndex: "name",
      },
      {
        title: "Описание",
        dataIndex: "desc"
      }
    ]

    return (
      <>
        <Button
          type="link"
          onClick={() => { window.location.href = routesEnum.projects }}
        >
          Назад
        </Button>
        <Table
          columns={userListTableCols}
          dataSource={this.state.userListTableModel}
          loading={this.state.loading}
        />
        <Table
          columns={roleListTableCols}
          dataSource={this.state.roleListTableModel}
          loading={this.state.loading}
        />
        <Table
          columns={permListTableCols}
          dataSource={this.state.permissionListTableModel}
          loading={this.state.loading}
        />

        <Modal
          title={"Редактирование пользователя "+this.state.editedUserName}
          open={this.state.editedUserId != undefined}
          onOk={() => {
            const allRoles = Array.from(this.roleById.values());
            this.updateUser(
              this.state.editedUserId!,
              this.state.usernameEditInputValue!,
              this.state.userRolesEditCheckboxValues!
              .map((item) => allRoles.find((v) => v.name == item)?.id!),
            );
          }}
          onCancel={() => this.setState({
            editedUserId: undefined,
            editedUserName: undefined,
            usernameEditInputValue: undefined,
            userRolesEditCheckboxValues: undefined,
          })}
        >
          <Form
            layout="vertical"
          >
            <Form.Item label="Логин">
              <Input
                value={this.state.usernameEditInputValue}
                onChange={(newValue) => this.setState({usernameEditInputValue: newValue.target.value})}
              />
            </Form.Item>
            <Form.Item label="Роли">
              <Checkbox.Group
                value={this.state.userRolesEditCheckboxValues}
                options={this.state.userRolesEditCheckboxAllValues}
                onChange={(newValue) => {console.log("onChange triggered"); this.setState({userRolesEditCheckboxValues: newValue.map((v) => v.toString())})}}
              />
            </Form.Item>
          </Form>
        </Modal>

        <Modal
          title={"Редактирование роли "+this.state.editedRoleName}
          open={this.state.editedRoleId != undefined}
          onOk={() => {
            const allPerms = Array.from(this.permById.values());
            const permIds = this.state.editedRolePermsCheckboxValues!.map((permName) => allPerms.find((p) => p.name == permName)?.id!);
            this.updateRole(this.state.editedRoleId!, this.state.editedRoleNameInputValue!, this.state.editedRoleDescInputValue, permIds);
          }}
          onCancel={() => this.setState({editedRoleId: undefined})}
        >
          <Form
            layout="vertical"
          >
            <Form.Item
              label="Наименование"
            >
              <Input 
                value={this.state.editedRoleName}
                onChange={(e) => this.setState({editedRoleNameInputValue: e.target.value})}
              />
            </Form.Item>
            <Form.Item
              label="Описание"
            >
              <Input 
                value={this.state.editedRoleDescInputValue}
                onChange={(e) => this.setState({editedRoleDescInputValue: e.target.value})}
              />
            </Form.Item>
            <Form.Item
              label="Разрешения"
              >
              <Checkbox.Group
                value={this.state.editedRolePermsCheckboxValues}
                options={this.state.editedRolePermsCheckboxAllValues}
                onChange={(newValue) => {console.log("onChange triggered"); this.setState({editedRolePermsCheckboxValues: newValue.map((v) => v.toString())})}}
              />
            </Form.Item>
          </Form>
        </Modal>

        <FloatButton
          icon={<UserAddOutlined />}
          type="primary"
          onClick={() => {
            const allRolesCheckboxes = this.state.roleListTableModel?.map((x) => x.name!);
            this.setState({userCreateDialogOpened: true, userCreateRolesCheckboxAllValues: allRolesCheckboxes});
          }}
          style={{right: 24}}
        />
        <FloatButton
          icon={<UsergroupAddOutlined />}
          type="default"
          onClick={() => {
            const allPermsCheckboxes = this.state.permissionListTableModel?.map((x) => x.name!);
            this.setState({roleCreateDialogOpened: true, roleCreatePermsCheckboxAllValues: allPermsCheckboxes});
          }}
          style={{right: 94}}
        />

        <Modal
          open={this.state.roleCreateDialogOpened}
          onCancel={() => this.setState({
            roleCreateDialogOpened: false,
            roleCreateNameInputValue: undefined,
            roleCreateDescInputValue: undefined,
            roleCreatePermsCheckboxValues: undefined,
          })}
          onOk={() => {
            this.createRole();
          }}
        >
          <Form
            layout="vertical"
          >
            <Form.Item label="Наименование">
              <Input
                onChange={(e) => this.setState({roleCreateNameInputValue: e.target.value})}
              />
            </Form.Item>
            <Form.Item label="Описание">
              <Input
                onChange={(e) => this.setState({roleCreateDescInputValue: e.target.value})}
              />
            </Form.Item>
            <Form.Item label="Разрешения">
              <Checkbox.Group
                options={this.state.roleCreatePermsCheckboxAllValues}
                onChange={(values) => this.setState({roleCreatePermsCheckboxValues: values.map((x) => x.toString())})}
              />
            </Form.Item>
          </Form>
        </Modal>
        <Modal
          open={this.state.userCreateDialogOpened}
          onCancel={() => this.setState({
            userCreateDialogOpened: false,
            userCreateLoginInputValue: undefined,
            userCreatePasswordInputValue: undefined,
            userCreateRolesCheckboxValues: undefined,
          })}
          onOk={() => {
            this.createUser();
          }}
        >
          <Form
            layout="vertical"
          >
            <Form.Item label="Логин">
              <Input
                onChange={(e) => this.setState({userCreateLoginInputValue: e.target.value})}
              />
            </Form.Item>
            <Form.Item label="Пароль">
              <Input
                type="password"
                onChange={(e) => this.setState({userCreatePasswordInputValue: e.target.value})}
              />
            </Form.Item>
            <Form.Item label="Роли">
              <Checkbox.Group
                options={this.state.userCreateRolesCheckboxAllValues}
                onChange={(values) => this.setState({userCreateRolesCheckboxValues: values.map((x) => x.toString())})}
              />
            </Form.Item>
          </Form>
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

  roleListTableModel: any[] | undefined;

  editedRoleId: number | undefined;
  editedRoleName: string | undefined;
  editedRoleNameInputValue: string | undefined;
  editedRoleDescInputValue: string | undefined;
  editedRolePermsCheckboxValues: string[] | undefined;
  editedRolePermsCheckboxAllValues: string[] | undefined;

  permissionListTableModel: any[] | undefined;

  usernameEditInputValue: string | undefined;
  userRolesEditCheckboxValues: string[] | undefined;
  userRolesEditCheckboxAllValues: string[] | undefined;

  userCreateDialogOpened: boolean;
  userCreateLoginInputValue: string | undefined;
  userCreatePasswordInputValue: string | undefined;
  userCreateRolesCheckboxValues: string[] | undefined;
  userCreateRolesCheckboxAllValues: string[] | undefined;

  roleCreateDialogOpened: boolean;
  roleCreateNameInputValue: string | undefined;
  roleCreateDescInputValue: string | undefined;
  roleCreatePermsCheckboxAllValues: string[] | undefined;
  roleCreatePermsCheckboxValues: string[] | undefined;
}

export default UserEditPage;

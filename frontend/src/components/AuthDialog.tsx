import React from "react";

import { Button, Col, Form, Input } from "antd";
import "./AuthDialog.css";
import { $notify, ENotifyKind } from "../notifier";
import { routesEnum } from "../routesEnum";
import { InternalApiAuthAuthRequest, UsersApi } from "../api/api";

const AuthDialog = () => {
  const [form] = Form.useForm();
  let pageModel = new AuthDialogModel(new UsersApi());

  return (
    <Col className="authDialog " span={6}>
      <Form
        form={form}
        layout="vertical"
        onFinish={
          (formData: any) => {
            pageModel.handleAuthSumbit(formData);
          }
        }>
        <Form.Item name="username" label="Логин">
          <Input />
        </Form.Item>
        <Form.Item name="password" label="Пароль">
          <Input type="password" />
        </Form.Item>
        <Form.Item>
          <Button type="primary" onClick={() => form.submit()}>
            Войти
          </Button>
        </Form.Item>
      </Form>
    </Col>
  );
};

class AuthDialogModel {
  private userApi: UsersApi;
  public loading: boolean;

  constructor(userApi: UsersApi) {
    this.handleAuthSumbit.bind(this);
    this.userApi = userApi;
    this.loading = false;
  }

  public handleAuthSumbit(formData: any) {
    this.loading = true;
    this.userApi.auth(formData as InternalApiAuthAuthRequest)
      .then((r: any) => { console.log(r.headers); window.location.href = routesEnum.projects })
      .catch((e) => { $notify(ENotifyKind.ERROR, e) })
      .finally(() => { this.loading = false });
  }
}

export default AuthDialog;

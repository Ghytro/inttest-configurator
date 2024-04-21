import React from "react";

import { Button, Col, Form, Input } from "antd";
import "./AuthDialog.css";
import UsersApi from "../api/src/api/UsersApi";
import { message } from "antd";
import { routesEnum } from "../routesEnum";

const AuthDialog = () => {
  const [form] = Form.useForm();

  // const userApi = new UsersApi();

  // const submitAuthData = (formData) => {
  //   userApi.auth(formData, (e, data, response) => {
  //     if (e) {
  //       message.error(response?.body?.message);
  //       return;
  //     }
  //     window.location.href = routesEnum.projects;
  //   });
  // };

  return (
    <Col className="authDialog " span={6}>
      <Form form={form} layout="vertical">
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

export default AuthDialog;

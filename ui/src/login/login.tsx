import { authLogin } from './api';
import { LockOutlined, UserOutlined } from '@ant-design/icons';
import { Button, Form, Input, message } from 'antd';
import React, { useEffect, useState } from 'react';
import { useHistory, useLocation } from 'react-router-dom';
import loginBg from '@/images/sgv/login-bg.svg';
import loginFormBg from '@/images/sgv/login-form-bg.svg';
import './login.less';

import { useTranslation } from 'react-i18next';
// import { getBaseUrl } from '@/utils/common';
export interface DisplayName {
    oidc: string;
    cas: string;
    oauth: string;
}
export default function Login() {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const history = useHistory();
  const location = useLocation();
  const redirect = location.search && new URLSearchParams(location.search).get('redirect');
  
  useEffect(() => {
     
  }, []);

  const handleSubmit = async () => {
      try {
          await form.validateFields();
          login();
      } catch {
          console.log(t('输入有误'));
      }
  };

  const login = async () => {
      let { username, password } = form.getFieldsValue();
      authLogin(username, password).then((res) => {
          const { dat, err } = res;
          const { access_token, refresh_token } = dat;
          localStorage.setItem('access_token', access_token);
          localStorage.setItem('refresh_token', refresh_token);
          if (!err) {
              window.location.href = redirect || getBaseUrl('/metric/explorer');
          }
      });
  };

  return (
    <>
      <img src={loginBg} alt="login" className="login-bg" />
      <div className="login-page-container">
        <div className="login-form-container">
          <div className="login-form-left">
            <img src={loginFormBg} />
          </div>
          <div className="login-form-right">
            <div className="login-form-content">
              <div className="form-title">欢迎使用 Gitee Ops 平台</div>
                <Form form={form} layout="vertical" requiredMark={true}>
                  <Form.Item
                      name="username"
                      rules={[
                          {
                              required: true,
                              message: t('请输入用户名'),
                          },
                      ]}
                  >
                      <Input
                          placeholder={t('请输入用户名')}
                          prefix={<UserOutlined className="site-form-item-icon" />}
                      />
                  </Form.Item>
                  <Form.Item
                      name="password"
                      rules={[
                          {
                              required: true,
                              message: t('请输入密码'),
                          },
                      ]}
                  >
                      <Input
                          type="password"
                          placeholder={t('请输入密码')}
                          onPressEnter={handleSubmit}
                          prefix={<LockOutlined className="site-form-item-icon" />}
                      />
                  </Form.Item>
                </Form>
                <div className="login-button" onClick={handleSubmit}>
                  {t('登录')}
                </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}

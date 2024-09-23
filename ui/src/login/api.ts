
import request from '@/utils/request';



export enum RequestMethod {
    Get = 'Get',
    Post = 'Post',
    Put = 'Put',
    Delete = 'Delete',
}


// 登录
export const authLogin = function (username: string, password: string) {
    return request(`/v1/auth/login`, {
        method: RequestMethod.Post,
        data: { username, password },
    });
};

// 刷新accessToken
export const UpdateAccessToken = function () {
    return request(`/v1/auth/refresh`, {
        method: RequestMethod.Post,
        data: {
            refresh_token: localStorage.getItem('refresh_token'),
        },
    });
};

// 更改密码
export const UpdatePwd = function (oldpass: string, newpass: string) {
    return request(`/v1/self/password`, {
        method: RequestMethod.Put,
        data: { oldpass, newpass },
    });
};

// 获取csrf token
export const GenCsrfToken = function () {
    return request(`/v1/csrf`, {
        method: RequestMethod.Get,
    });
};

// 退出
export const Logout = function () {
    return request(`/v1/auth/logout`, {
        method: RequestMethod.Post,
    });
};

export const getRedirectURL = function () {
    return request('/v1/auth/redirect', {
        method: RequestMethod.Get,
    });
};


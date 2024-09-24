
import request from '../utils/requests';


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




// 退出
export const Logout = function () {
    return request(`/v1/auth/logout`, {
        method: RequestMethod.Post,
    });
};



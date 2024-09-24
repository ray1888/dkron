/** Request 网络请求工具 更详细的 api 文档: https://github.com/umijs/umi-request */
import { extend } from 'umi-request';
import { notification } from 'antd';
import _ from 'lodash';
// import { UpdateAccessToken } from '@/login/api';

import { getBaseUrl } from './common';


/** 异常处理程序，所有的error都被这里处理，页面无法感知具体error */
const errorHandler = (error: Error): Response => {
    // 忽略掉 setting getter-only property "data" 的错误
    // 这是 umi-request 的一个 bug，当触发 abort 时 catch callback 里面不能 set data
    if (error.name !== 'AbortError' && error.message !== 'setting getter-only property "data"') {
        // @ts-ignore
        if (!error.silence) {
            notification.error({
                message: error?.message || '服务异常,请联系管理员。',
            });
        }
        // 暂时认定只有开启 silence 的时候才需要传递 error 详情以便更加精确的处理错误
        // @ts-ignore
        if (error.silence) {
            throw error;
        } else {
            throw new Error(error?.message);
        }
    }
    throw error;
};

/** 处理后端返回的错误信息 */
const processError = (res: any): string => {
    if (res?.error) {
        return _.isString(res?.error) ? res.error : JSON.stringify(res?.error);
    }
    if (res?.err) {
        return _.isString(res?.err) ? res.err : JSON.stringify(res?.err);
    }
    if (res?.errors) {
        const errors = res?.errors
        if (Array.isArray(errors) && errors[0]?.msg) {
            return errors[0]?.msg
        }
        return _.isString(res?.errors) ? res.errors : JSON.stringify(res?.errors);
    }
    if (res?.message) {
        return _.isString(res?.message) ? res.message : JSON.stringify(res?.message);
    }
    return JSON.stringify(res);
};

/** 配置request请求时的默认参数 */
const request = extend({
    errorHandler,
    credentials: 'include',
    prefix: '/v1'
});

request.interceptors.request.use((url, options) => {
    let headers = {
        ...options.headers,
    };
    
    return {
        url,
        options: { ...options, headers },
    };
});

/**
 * 响应拦截
 */
request.interceptors.response.use(
    async (response, options) => {
        const { status } = response;
        if (status === 200) {
            return response
                .clone()
                .json()
                .then((data) => {
                    const { url } = response;
                    // n9e 和 n9e-plus 大部分接口返回的数据结构是 { err: '', dat: {} }
                    if (data.err === '' || data.status === 'success' || data.error === '') {
                        return { ...data, success: true };
                    } else {
                        throw {
                            name: processError(data),
                            message: processError(data),
                            silence: options.silence,
                            data,
                            response,
                        };
                    }
                    window.location.href = "/"
                });
        } else if (status === 401) {
            // const loginUrl = getBaseUrl('/login')

            // if (window.location.pathname === loginUrl) return;
            // if (response.url.indexOf('/auth/refresh') > 0) {
            //     window.location.href = loginUrl
            // } else {
            //     localStorage.getItem('refresh_token')
            //         ? UpdateAccessToken().then((res) => {
            //             console.log('401 err', res);
            //             if (res.err) {
            //                 window.location.href = getBaseUrl('/login')
            //             } else {
            //                 const { access_token, refresh_token } = res.dat;
            //                 localStorage.setItem('access_token', access_token);
            //                 localStorage.setItem('refresh_token', refresh_token);
            //                 window.location.href = getBaseUrl('/login')
            //             }
            //         })
            //         : window.location.href = getBaseUrl('/login')
            // }
        } else {
            return response
                .clone()
                .text()
                .then((data) => {
                    let errObj = {};
                    try {
                        const parsed = JSON.parse(data);
                        const errMessage = processError(parsed);
                        errObj = {
                            name: errMessage,
                            message: errMessage,
                            data: parsed,
                        };
                    } catch (error) {
                        errObj = {
                            name: data,
                            message: data,
                        };
                    }
                    throw {
                        ...errObj,
                        silence: options.silence,
                    };
                });
        }
    },
    {
        global: false,
    },
);

export default request;

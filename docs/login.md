## 获取验证码

```
Type : HTTP
RequestUrl : https://api.creative-apartment.com.cn/auth/auth/sendCode
RequestType : POST
Body : 
{
    "type": 1,
    "username": "${手机号码}"
}

ResponseCode:200
Body : 
{
    "code": 200,
    "message": "操作成功",
    "content": "验证码已发送，请注意查收！"
}
```

## 登录

```
Type : HTTP
RequestUrl : https://api.creative-apartment.com.cn/auth/authentication/customer/phone/app
RequestType : POST
Body : 
{
    "code": "${手机验证码}",
    "password": "${AES_ECB_PKCS5, KEY: yBnulH9ODtonS5lj}",
    "registrationIds": "${极光推送id，应该可以随机生成}",
    "username": "${手机号码}"
}

ResponseCode:200
Body : 
{
    "code": 200,
    "message": "操作成功",
    "content": {
        "id": "000",
        "type": 2,
        "isRepairman": 2,
        "userName": "",
        "access_token": "xxx",
        "token_type": "bearer",
        "expires_in": 2591999,
        "refresh_token": "yyy",
        "scope": "[read]"
    }
}
```

# Golang OAuth 2.0 Client
<p align="center">
<img align="center" width="150px" src="https://www.oauth.com/wp-content/themes/oauthdotcom/images/oauth_logo@2x.png" />
</p>
<h3 align=center>Golang 实现的 OAuth2.0 客户端</h3>

[English](README.md) | 简体中文

## OAuth 2.0协议流程
     +--------+                               +---------------+
     |        |--(A)- Authorization Request ->|   Resource    |
     |        |                               |     Owner     |
     |        |<-(B)-- Authorization Grant ---|               |
     |        |                               +---------------+
     |        |
     |        |                               +---------------+
     |        |--(C)-- Authorization Grant -->| Authorization |
     | Client |                               |     Server    |
     |        |<-(D)----- Access Token -------|               |
     |        |                               +---------------+
     |        |
     |        |                               +---------------+
     |        |--(E)----- Access Token ------>|    Resource   |
     |        |                               |     Server    |
     |        |<-(F)--- Protected Resource ---|               |
     +--------+                               +---------------+

## 授权码流程
Web 和移动应用程序使用授权码授权类型。它与大多数其他授权类型不同，首先要求应用程序启动浏览器以开始流程。具有以下步骤：

- 应用程序打开浏览器请求发送到 OAuth 服务器
- 用户看到授权提示并批准应用程序的请求
- 授权成功后将用户重定向回应用程序并携带授权码 
- 应用程序携带访问令牌交换授权代码

### 获得用户的许可
OAuth 就是让用户能够授予对应用程序的有限访问权限。应用程序首先需要决定它请求的权限，然后将用户发送到浏览器以获得他们的权限。开始授权流程，应用程序构建如下所示的 URL 并打开浏览器访问该 URL。

```bash
https://authorization-server.com/auth
 ?response_type=code
 &client_id=29352915982374239857
 &redirect_uri=https%3A%2F%2Fexample-app.com%2Fcallback
 &scope=create+delete
 &state=xcoiv98y2kd22vusuye3kch
```
以下是对每个查询参数的解释：

response_type=code 这告诉授权服务器应用程序正在启动授权请求。
client_id- 应用程序的公共标识符，在开发人员首次注册应用程序时获得。
redirect_uri- 告诉授权服务器在用户批准请求后将用户重定向回何处。
scope- 一个或多个空格分隔的字符串，指示应用程序请求的权限。您使用的特定 OAuth API 将定义它支持的范围。
state- 应用程序生成一个随机字符串并将其包含在请求中。然后它应该检查在用户授权应用程序后是否返回相同的值。这用于防止CSRF 攻击。

当用户访问此 URL 时，授权服务器将向他们显示一个提示，询问他们是否愿意授权此应用程序的请求。

### 重定向回应用程序

如果用户批准请求，授权服务器会将浏览器重定向回redirect_uri应用程序指定的浏览器，并在查询字符串中添加code和state

例如，用户将被重定向回一个 URL，例如

```bash
https://example-app.com/redirect?code=g0ZGZmNjVmOWIjNTk2NTk4ZTYyZGI3&state=xcoiv98y2kd22vusuye3kch
```

该state值将与应用程序最初在请求中设置的值相同。应用程序应检查重定向中的状态是否与它最初设置的状态相匹配。这可以防止 CSRF 和其他相关攻击。

code是授权服务器生成的授权码。此代码的生命周期相对较短，通常会持续 1 到 10 分钟,有的 Oauth 服务只允许使用一次就会失效. 具体取决于 OAuth 服务。

## 使用授权码交换为访问令牌

我们即将结束流程。现在应用程序有了授权代码，它可以使用它来获取访问令牌。

应用程序使用以下参数向服务的令牌端点发出 POST 请求：

- grant_type=authorization_code 这告诉 Oauth 服务端应用程序正在使用授权代码授权类型。
- code 应用程序包含在重定向中提供的授权代码。
- redirect_uri- 请求代码时使用的相同重定向 URI。某些 API 不需要此参数，因此需要仔细检查您正在访问的特定 API 的文档,有的服务商可能需要。
- client_id- 应用程序的客户端 ID。
- client_secret- 应用程序的客户端机密。这确保获取访问令牌的请求仅来自应用程序，而不是来自可能拦截授权代码的潜在攻击者。

client_id,client_secret 的传参数需要参考 Oauth 服务商的约定.一般都回在header中传递 Basic 类型的 Authorization. 加密规则:base64_encode(client_id:client_secret)

令牌端点将验证请求中的所有参数，确保代码没有过期并且客户端 ID 和密码匹配。如果一切正常，它将生成一个访问令牌并在响应中返回它！

```http response
HTTP/1.1 200 OK
Content-Type: application/json
Cache-Control: no-store
Pragma: no-cache

{
  "access_token":"MTQ0NjJkZmQ5OTM2NDE1ZTZjNGZmZjI3",
  "token_type":"bearer",
  "expires_in":3600,
  "refresh_token":"IwOGYzYTlmM2YxOTQ5MGE3YmNmMDFkNTVk",
  "scope":"create delete"
}
```

授权码流程完成！该应用程序现在有一个访问令牌，它可以在发出 获取授权用户信息等相关 API 请求时使用。

## 何时使用授权代码流程

授权代码流程最适用于 Web 和移动应用程序。由于授权代码授予具有为访问令牌交换授权代码的额外步骤，因此它提供了隐式授权类型中不存在的附加安全层。

如果您在移动应用程序或无法存储客户端机密的任何其他类型的应用程序中使用授权代码流，那么您还应该使用 PKCE 扩展，它可以防止授权代码可能被攻击的其他攻击拦截。

代码交换步骤确保攻击者无法拦截访问令牌，因为访问令牌始终通过应用程序和 OAuth 服务器之间的安全反向通道发送。

## 具体使用请阅读 

[main.go ](main.go)

## Give a Star! ⭐
如果你喜欢或正在使用这个项目来学习或开始你的解决方案，请给它一颗星。谢谢！

## Buy me a coffee

<a href="https://www.buymeacoffee.com/demo007x" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 60px !important;width: 217px !important;" ></a>

## 问题讨论
<img src="https://user-images.githubusercontent.com/6418340/230716042-135d28f0-9912-4ba4-8adf-a8f14eb76b05.png" alt="discard with Me" style="height: 150px !important;width: 150px !important;" >



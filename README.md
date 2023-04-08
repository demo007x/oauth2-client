# Golang OAuth 2.0 Client

<p align="center">
<img align="center" width="150px" src="https://www.oauth.com/wp-content/themes/oauthdotcom/images/oauth_logo@2x.png" />
</p>
<h3 align=center>Oauth2 Client Package For Golang</h3>

English | [简体中文](README-CN.md) | [Oauth2 Flow](oauth-flow-en.md)

## OAuth 2.0 Protocol Flow
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




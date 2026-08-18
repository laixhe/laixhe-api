//! OpenAPI 文档聚合定义 (utoipa 注解, 对应 Go 端 `swag init` 生成的 docs)
//!
//! 生成方式: 运行 `make docs` 或 `webapi --gen-docs` 重新生成
//! `docs/swagger.json` / `docs/swagger.yaml`; 服务运行期直接由本模块的
//! `ApiDoc::openapi()` 序列化提供 (`/api/v1/swagger.json|yaml`), 无需读盘。

use utoipa::openapi::RefOr;
use utoipa::openapi::schema::Schema;
use utoipa::openapi::security::{HttpAuthScheme, HttpBuilder, SecurityRequirement, SecurityScheme};
use utoipa::{Modify, OpenApi};

use crate::app::controllers::auth;
use crate::app::controllers::health::{self, HealthResponse};
use crate::app::controllers::user;
use crate::app::entity::auth::{AuthLoginRequest, AuthRegisterRequest, AuthTokenResponse};
use crate::app::entity::user::{User, UserListResponse, UserUpdateRequest};

/// 统一错误响应体 (与 PHP/TS 端 Error 及 Go core.Error 对齐)
///
/// 字段仅供 utoipa 生成文档 schema 使用, 无需在业务中读写
#[derive(utoipa::ToSchema)]
#[allow(dead_code)]
pub struct Error {
    /// 错误码 (与 HTTP 状态码一致)
    pub code: u16,
    /// 错误描述
    pub message: String,
}

/// 向文档注入 BearerAuth 安全方案与根级默认安全 (与 PHP/TS 端一致),
/// 并为 User 的 sex/states/type_id 补充枚举取值 (对齐 Go 端 models.UserSex/UserState/UserType)
struct SecurityAddon;

impl Modify for SecurityAddon {
    fn modify(&self, openapi: &mut utoipa::openapi::OpenApi) {
        openapi.security = Some(vec![SecurityRequirement::new("BearerAuth", Vec::<String>::new())]);
        let Some(components) = openapi.components.as_mut() else {
            return;
        };
        components.add_security_scheme(
            "BearerAuth",
            SecurityScheme::Http(
                HttpBuilder::new()
                    .scheme(HttpAuthScheme::Bearer)
                    .bearer_format("JWT")
                    .description(Some("在请求头携带 Authorization: Bearer <token>"))
                    .build(),
            ),
        );
        // 枚举值: sex 0未填写/1男/2女, states 0禁用/1正常, type_id 1普通用户
        let enums: &[(&str, &[i64])] = &[
            ("sex", &[0, 1, 2]),
            ("states", &[0, 1]),
            ("type_id", &[1]),
        ];
        if let Some(RefOr::T(Schema::Object(user))) = components.schemas.get_mut("User") {
            for (field, values) in enums {
                if let Some(RefOr::T(Schema::Object(prop))) = user.properties.get_mut(*field) {
                    prop.enum_values = Some(
                        values
                            .iter()
                            .map(|v| serde_json::Value::from(*v))
                            .collect(),
                    );
                }
            }
        }
    }
}

/// 聚合的 OpenAPI 文档
#[derive(OpenApi)]
#[openapi(
    info(
        title = "API接口",
        description = "用户认证与用户管理 API 服务",
        version = "1.0"
    ),
    paths(
        auth::register,
        auth::login,
        auth::refresh,
        user::info,
        user::list,
        user::update,
        health::health
    ),
    components(schemas(
        Error,
        HealthResponse,
        AuthRegisterRequest,
        AuthLoginRequest,
        UserUpdateRequest,
        AuthTokenResponse,
        User,
        UserListResponse
    )),
    modifiers(&SecurityAddon)
)]
pub struct ApiDoc;

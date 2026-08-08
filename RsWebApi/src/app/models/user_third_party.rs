//! 用户第三方表实体 (对应 Go 项目 app/models/user_third_party.go)

use sea_orm::entity::prelude::*;
use sea_orm::sea_query::StringLen;

/// 用户第三方表名
#[allow(dead_code)]
pub const USER_THIRD_PARTY_TABLE: &str = "user_third_party";

/// 用户第三方
#[derive(Clone, Debug, PartialEq, Eq, DeriveEntityModel)]
#[sea_orm(table_name = "user_third_party")]
pub struct Model {
    /// 主键 (数据库为 INT UNSIGNED)
    #[sea_orm(primary_key)]
    pub id: u32,
    /// 用户UID (数据库为 INT UNSIGNED)
    #[sea_orm(indexed)]
    pub uid: u32,
    /// 微信unionid
    #[sea_orm(column_type = "String(StringLen::N(200))", default_value = "")]
    pub wechat_unionid: String,
    /// 微信openid
    #[sea_orm(column_type = "String(StringLen::N(200))", default_value = "", indexed)]
    pub wechat_openid: String,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {}

impl ActiveModelBehavior for ActiveModel {}

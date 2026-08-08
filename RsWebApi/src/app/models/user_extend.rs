//! 用户扩展表实体 (对应 Go 项目 app/models/user_extend.go)

use sea_orm::entity::prelude::*;

/// 用户扩展表名
#[allow(dead_code)]
pub const USER_EXTEND_TABLE: &str = "user_extend";

/// 用户扩展
#[derive(Clone, Debug, PartialEq, Eq, DeriveEntityModel)]
#[sea_orm(table_name = "user_extend")]
pub struct Model {
    /// 主键 (数据库为 INT UNSIGNED)
    #[sea_orm(primary_key)]
    pub id: u32,
    /// 用户UID (数据库为 INT UNSIGNED)
    #[sea_orm(indexed)]
    pub uid: u32,
    /// 生日(年月日)
    #[sea_orm(column_type = "Integer", default_value = 0)]
    pub birthday: i32,
    /// 身高(cm)
    #[sea_orm(column_type = "Integer", default_value = 0)]
    pub height: i32,
    /// 体重(kg)
    #[sea_orm(column_type = "Integer", default_value = 0)]
    pub weight: i32,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {}

impl ActiveModelBehavior for ActiveModel {}

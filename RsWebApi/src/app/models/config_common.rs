//! 通用配置表实体 (对应 Go 项目 app/models/config_common.go)

use sea_orm::entity::prelude::*;
use sea_orm::sea_query::StringLen;
use sea_orm::{EntityTrait, QueryFilter};

/// 通用配置
#[derive(Clone, Debug, PartialEq, Eq, DeriveEntityModel)]
#[sea_orm(table_name = "config_common")]
pub struct Model {
    /// 主键 (数据库为 INT)
    #[sea_orm(primary_key)]
    pub id: i32,
    #[sea_orm(column_type = "String(StringLen::N(50))", default_value = "")]
    pub key: String,
    #[sea_orm(column_type = "String(StringLen::N(500))", default_value = "")]
    pub value: String,
    /// 描述
    #[sea_orm(column_type = "String(StringLen::N(500))", default_value = "")]
    pub describe: String,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {}

impl ActiveModelBehavior for ActiveModel {}

/// 查询通用配置列表，可选按 key 过滤 (对应 ConfigCommon.List)
pub async fn list(db: &DatabaseConnection, keys: &[&str]) -> Result<Vec<Model>, DbErr> {
    let mut query = Entity::find();
    if !keys.is_empty() {
        query = query.filter(Column::Key.is_in(keys.iter().copied()));
    }
    query.all(db).await
}

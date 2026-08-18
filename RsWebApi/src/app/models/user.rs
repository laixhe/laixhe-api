//! 用户表实体 (对应 Go 项目 app/models/user.go)

use sea_orm::entity::prelude::*;
use sea_orm::sea_query::StringLen;
use sea_orm::ActiveValue::Set;
use sea_orm::{EntityTrait, QueryFilter, QueryOrder, QuerySelect, TransactionTrait};

use super::user_extend;
use super::user_third_party;

/// 用户
#[derive(Clone, Debug, PartialEq, Eq, DeriveEntityModel)]
#[sea_orm(table_name = "user")]
pub struct Model {
    /// 主键 (数据库为 INT; 与 Go 版有符号 int 建表保持一致, 兼容旧库)
    #[sea_orm(primary_key)]
    pub id: i32,
    /// 类型 1普通
    #[sea_orm(column_type = "Integer", default_value = 0)]
    pub type_id: i32,
    /// 账号 (全局唯一, 与 webapi.sql 的 UNIQUE KEY user_account_idx 保持一致)
    #[sea_orm(column_type = "String(StringLen::N(120))", default_value = "", unique)]
    pub account: String,
    /// 手机号
    #[sea_orm(column_type = "String(StringLen::N(120))", default_value = "", indexed)]
    pub mobile: String,
    /// 邮箱 (唯一索引, 与 webapi.sql 一致; 注册先查后插 + 数据库唯一约束双重防重)
    #[sea_orm(column_type = "String(StringLen::N(120))", default_value = "", unique)]
    pub email: String,
    /// 密码
    #[sea_orm(column_type = "String(StringLen::N(120))", default_value = "")]
    pub password: String,
    /// 昵称
    #[sea_orm(column_type = "String(StringLen::N(120))", default_value = "")]
    pub nickname: String,
    /// 头像地址
    #[sea_orm(column_type = "String(StringLen::N(255))", default_value = "")]
    pub avatar_url: String,
    /// 性别 0未填写 1男 2女
    #[sea_orm(column_type = "Integer", default_value = 0)]
    pub sex: i32,
    /// 状态 0封禁 1正常
    #[sea_orm(column_type = "Integer", default_value = 0)]
    pub states: i32,
    /// 创建时间
    /// 注: 类型为 chrono::NaiveDateTime 是 sea-orm 对 MySQL DATETIME 列的映射要求
    /// (with-chrono feature, 当前 sea-orm 不支持 jiff 映射); 业务时间统一由 jiff 生成
    pub created_at: chrono::NaiveDateTime,
    /// 更新时间
    pub updated_at: chrono::NaiveDateTime,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {}

impl ActiveModelBehavior for ActiveModel {}

/// 当前本地时间 (业务统一使用 jiff, 与 Go 的 time.Now() 语义一致)
///
/// 返回值类型为 chrono::NaiveDateTime, 仅为满足 sea-orm 实体字段类型。
pub fn now_local_naive() -> chrono::NaiveDateTime {
    let zoned = jiff::Zoned::now();
    let dt = zoned.datetime();
    // 当前时间必然构成合法的年月日/时分秒, *_opt 不会返回 None (unwrap 安全);
    // 用 expect 说明该不变式, 便于新手理解而非盲模仿 unwrap
    chrono::NaiveDate::from_ymd_opt(dt.year() as i32, dt.month() as u32, dt.day() as u32)
        .expect("当前年月日必然构成合法 NaiveDate")
        .and_hms_opt(dt.hour() as u32, dt.minute() as u32, dt.second() as u32)
        .expect("当前时分秒必然构成合法 NaiveTime")
}

/// 按主键查询用户 (对应 s.orm.GetById)
///
/// 说明: 不做列裁剪 (Go 端可用 UserColumnsNoPassword 排除 password, 那是 GORM 支持缺列映射;
/// sea-orm 的 into_model 要求查询包含实体的全部列, 裁剪需额外定义 PartialModel,
/// 教学规模下单行查询收益可忽略, 故保持全列查询, password 不会进入响应 DTO)
pub async fn find_by_id(db: &DatabaseConnection, uid: i32) -> Result<Option<Model>, DbErr> {
    Entity::find_by_id(uid).one(db).await
}

/// 按邮箱查询用户 (对应 s.orm.FirstByField)
pub async fn find_by_email(db: &DatabaseConnection, email: &str) -> Result<Option<Model>, DbErr> {
    Entity::find().filter(Column::Email.eq(email)).one(db).await
}

/// 仅查询邮箱对应的用户 ID (对应 Go: Select("id").First(...)，减少数据传输)
pub async fn find_id_by_email(db: &DatabaseConnection, email: &str) -> Result<Option<i32>, DbErr> {
    Entity::find()
        .select_only()
        .column(Column::Id)
        .filter(Column::Email.eq(email))
        .into_tuple::<(i32,)>()
        .one(db)
        .await
        .map(|row| row.map(|(id,)| id))
}

/// 在事务中创建用户，同时创建关联的扩展信息和第三方记录 (对应 CreateUser)
pub async fn create_user(db: &DatabaseConnection, user: ActiveModel) -> Result<i32, DbErr> {
    let txn = db.begin().await?;
    let result = Entity::insert(user).exec(&txn).await?;
    let uid = result.last_insert_id as i32;
    // 在同一事务中创建用户、扩展信息、第三方关联
    // INSERT INTO `user` (...)
    // INSERT INTO `user_extend` (...)
    // INSERT INTO `user_third_party` (...)
    let user_extend = user_extend::ActiveModel {
        uid: Set(uid),
        ..Default::default()
    };
    user_extend::Entity::insert(user_extend).exec(&txn).await?;
    let user_third_party = user_third_party::ActiveModel {
        uid: Set(uid),
        ..Default::default()
    };
    user_third_party::Entity::insert(user_third_party)
        .exec(&txn)
        .await?;
    txn.commit().await?;
    Ok(uid)
}

/// 用户可更新字段 (对应 Go UpdateUser 的 map[string]any 动态非零字段更新)
///
/// `None` 表示不更新该字段，`Some(v)` 表示更新为 v。
#[derive(Debug, Clone, Default)]
pub struct UserUpdateData {
    /// 类型 1普通
    pub type_id: Option<i32>,
    /// 手机号
    pub mobile: Option<String>,
    /// 邮箱
    pub email: Option<String>,
    /// 密码
    pub password: Option<String>,
    /// 昵称
    pub nickname: Option<String>,
    /// 头像地址
    pub avatar_url: Option<String>,
    /// 状态 0封禁 1正常
    pub states: Option<i32>,
}

/// 根据非零字段动态更新用户信息，同时更新 updated_at (对应 UpdateUser)
pub async fn update_user(
    db: &DatabaseConnection,
    uid: i32,
    data: &UserUpdateData,
) -> Result<(), DbErr> {
    if uid == 0 {
        return Err(DbErr::Custom("primary key required".to_string()));
    }
    let mut updates = ActiveModel {
        id: Set(uid),
        updated_at: Set(now_local_naive()),
        ..Default::default()
    };
    // 仅 Set 需要更新的字段，未 Set 的字段在 UPDATE 时保持不变
    if let Some(v) = data.type_id {
        updates.type_id = Set(v);
    }
    if let Some(v) = &data.mobile {
        updates.mobile = Set(v.clone());
    }
    if let Some(v) = &data.email {
        updates.email = Set(v.clone());
    }
    if let Some(v) = &data.password {
        updates.password = Set(v.clone());
    }
    if let Some(v) = &data.nickname {
        updates.nickname = Set(v.clone());
    }
    if let Some(v) = &data.avatar_url {
        updates.avatar_url = Set(v.clone());
    }
    if let Some(v) = data.states {
        updates.states = Set(v);
    }
    Entity::update(updates).exec(db).await?;
    Ok(())
}

/// 分页查询用户列表，按 ID 降序 (对应 ListUser)
pub async fn list_user(
    db: &DatabaseConnection,
    limit: u64,
    offset: u64,
) -> Result<(Vec<Model>, i32), DbErr> {
    // SELECT count(*) FROM `user` — InnoDB 下 count(*) 走全表扫描, 表数据量大时较慢;
    // 与 Go 版 ListUser 语义保持一致, 暂不引入计数缓存 (缓存会牺牲实时性, 属取舍)。
    // 若需优化可考虑: 独立计数表 / Redis 计数 / 分页接口改游标分页 (cursor-based)
    let total = Entity::find().count(db).await? as i32;
    if total == 0 {
        return Ok((Vec::new(), 0));
    }
    // SELECT * FROM `user` ORDER BY `id` DESC LIMIT ? OFFSET ?
    // 说明: 同 find_by_id, sea-orm 不支持缺列 into_model, 故全列查询; password 不进入响应 DTO
    let list = Entity::find()
        .order_by_desc(Column::Id)
        .limit(limit)
        .offset(offset)
        .all(db)
        .await?;
    Ok((list, total))
}
